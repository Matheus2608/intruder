package handlers

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"intruder/strategy/matrix"
	strategy "intruder/strategy/request"
	"intruder/structs"
	"log"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Helper to extract method and path safely
func getMethodAndPath(requestLine string) (string, string, error) {
	parts := strings.SplitN(requestLine, " ", 3)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid request line: %q", requestLine)
	}
	// Basic validation: Method should be uppercase, path should start with /
	method := strings.ToUpper(parts[0])
	path := parts[1]
	if path == "" || path[0] != '/' {
		// Could be more robust, e.g., check for valid characters
		// log.Printf("Warning: Path %q does not start with '/'. Assuming it's valid.", path)
	}
	// Consider validating method against known HTTP methods if needed
	return method, path, nil
}

func AttackHandler(res http.ResponseWriter, req *http.Request) {
	startTimeOverall := time.Now()
	log.Println("Received attack request")

	// Parse the form data
	if err := req.ParseForm(); err != nil {
		log.Printf("ERROR: Unable to parse form: %v", err)
		http.Error(res, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// --- Configuration ---
	reqString := req.FormValue("requestData")
	attackType := req.FormValue("typeOfAttack")
	concurrencyStr := req.FormValue("concurrency")
	delayStr := req.FormValue("delay")
	skipVerifyStr := req.FormValue("skipVerify")

	if reqString == "" {
		log.Println("ERROR: Request data is empty")
		http.Error(res, "Request data cannot be empty", http.StatusBadRequest)
		return
	}

	concurrency, err := strconv.Atoi(concurrencyStr)
	if err != nil || concurrency <= 0 {
		defaultConcurrency := runtime.NumCPU() * 2 // Default based on CPU cores
		log.Printf("Warning: Invalid concurrency value '%s', defaulting to %d. Error: %v", concurrencyStr, defaultConcurrency, err)
		concurrency = defaultConcurrency
	}
	if concurrency > 500 { // Safety limit
		log.Printf("Warning: Concurrency %d exceeds limit of 500, capping at 500", concurrency)
		concurrency = 500
	}


	delayMs, err := strconv.Atoi(delayStr)
	if err != nil || delayMs < 0 {
		log.Printf("Warning: Invalid delay value '%s', defaulting to 0ms. Error: %v", delayStr, err)
		delayMs = 0
	}
	delayDuration := time.Duration(delayMs) * time.Millisecond

	skipVerify := skipVerifyStr == "true"

	log.Printf("Configuration: AttackType=%s, Concurrency=%d, Delay=%v, SkipVerify=%t",
		attackType, concurrency, delayDuration, skipVerify)


	// --- Request Parsing ---
	reqStringList := strings.SplitN(reqString, "\r\n", 2) // Split only the first line
	if len(reqStringList) == 0 {
		log.Println("ERROR: Request data is empty or malformed")
		http.Error(res, "Request data is malformed (no lines)", http.StatusBadRequest)
		return
	}
	requestLine := reqStringList[0]
	method, path, err := getMethodAndPath(requestLine)
	if err != nil {
		log.Printf("ERROR: Cannot parse request line: %v", err)
		http.Error(res, fmt.Sprintf("Cannot parse request line: %v", err), http.StatusBadRequest)
		return
	}
	log.Printf("Parsed Request: Method=%s, Path=%s", method, path)

	// --- Payload Matrix Generation ---
	matrixStrategy := matrix.GetMatrixStrategy(attackType) // Panics on invalid type, consider returning error
	payloads, err := matrixStrategy.BuildMatrix(req)       // Modify BuildMatrix to return error
	if err != nil {
		log.Printf("ERROR: Failed to build payload matrix: %v", err)
		http.Error(res, fmt.Sprintf("Failed to build payload matrix: %v", err), http.StatusBadRequest)
		return
	}
	lenPayloads := len(payloads)
	if lenPayloads == 0 {
		log.Println("Warning: No payloads generated based on input.")
		// Decide whether to error out or return an empty result page
		http.Error(res, "No payloads were generated for the attack.", http.StatusBadRequest)
		return
	}
	log.Printf("Generated %d payload combinations for attack", lenPayloads)

	// --- Strategy Cloning (Prepare request templates) ---
	// Clone Variables
	strategyClones := make([]strategy.RequestStrategy, lenPayloads)
	var clonesWG sync.WaitGroup
	cloneErrors := make(chan error, lenPayloads) // Channel to collect errors during cloning

	originalStrategy := strategy.ChooseStrategy(method) // Panics on invalid method, consider returning error

	clonesWG.Add(lenPayloads)
	for idx, payload := range payloads {
		go func(currentIndex int, currentPayload []string) {
			defer clonesWG.Done()
			err := originalStrategy.CloneWithDifferentPayload(currentIndex, reqString, currentPayload, &strategyClones)
			if err != nil {
				errMsg := fmt.Sprintf("Error cloning strategy for payload index %d: %v", currentIndex, err)
				log.Printf("ERROR: %s", errMsg)
				cloneErrors <- fmt.Errorf(errMsg) // Send error to channel
			}
		}(idx, payload)
	}

	clonesWG.Wait()
	close(cloneErrors)

	// Check for cloning errors before proceeding
	var firstCloneError error
	for err := range cloneErrors {
		if firstCloneError == nil {
			firstCloneError = err
		}
		// Log all errors, but maybe only report the first one to the user
	}
	if firstCloneError != nil {
		http.Error(res, fmt.Sprintf("Failed to prepare requests: %v", firstCloneError), http.StatusInternalServerError)
		return
	}

	log.Println("Successfully prepared all request strategies")


	// --- HTTP Client Setup ---
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second, // Connection timeout
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          concurrency + 10, // Adjust based on concurrency
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second, // Handshake timeout
			ExpectContinueTimeout: 1 * time.Second,
			ForceAttemptHTTP2:     true, // Try HTTP/2
		},
		Timeout: 60 * time.Second, // Overall request timeout (adjust as needed)
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Do not follow redirects automatically
		},
	}

	// --- Worker Pool Execution ---
	resultsChan := make(chan structs.ResponseData, lenPayloads)
	jobs := make(chan int, lenPayloads)
	var workerWG sync.WaitGroup
	var targetURL string
	var urlOnce sync.Once

	responseList := structs.NewResponses(lenPayloads) // Pre-allocate result storage

	log.Printf("Starting %d workers for %d requests...", concurrency, lenPayloads)
	startTimeAttack := time.Now()

	// Start workers
	for w := 1; w <= concurrency; w++ {
		workerWG.Add(1)
		go func(workerID int) {
			defer workerWG.Done()
			// log.Printf("Worker %d started", workerID)
			for idx := range jobs {
				// Optional delay
				if delayDuration > 0 {
					time.Sleep(delayDuration)
				}

				clone := strategyClones[idx]
				if clone == nil { // Check if cloning failed for this index
					log.Printf("Worker %d: Skipping index %d due to cloning error", workerID, idx)
					resultsChan <- structs.NewErrorResponse(idx, "Failed to prepare request strategy")
					continue
				}

				// Create the actual HTTP request
				httpReq, err := clone.CreateRequest(path)
				if err != nil {
					errMsg := fmt.Sprintf("Create Request Error (Index %d): %v", idx, err)
					log.Printf("Worker %d: ERROR %s", workerID, errMsg)
					resultsChan <- structs.NewErrorResponse(idx, errMsg)
					continue
				}

				// Capture target URL safely once
				urlOnce.Do(func() {
					if httpReq != nil && httpReq.URL != nil {
						targetURL = httpReq.URL.String() // Store the first successfully created URL
						log.Printf("Target URL identified: %s", targetURL)
					}
				})

				// Send the request
				httpRes, elapsedTime, err := strategy.SendRequest(client, httpReq) // SendRequest handles timing

				// Process result/error
				var response structs.ResponseData
				if err != nil {
					errMsg := fmt.Sprintf("Send Request Error (Index %d): %v", idx, err)
					log.Printf("Worker %d: ERROR %s", workerID, errMsg)
					response = structs.NewErrorResponse(idx, errMsg)
					// Ensure body is closed even on error if response exists
					if httpRes != nil && httpRes.Body != nil {
						httpRes.Body.Close()
					}
				} else {
					// Successful request (even if 4xx/5xx)
					cloneHttpReqStr, clonePayloadStr := clone.ToString() // Get string representations
					response = structs.NewResponse(
						httpRes,
						elapsedTime,
						cloneHttpReqStr,
						clonePayloadStr,
						idx, // Use the original index
						true) // Success = true (request sent without network error)

					// IMPORTANT: Ensure the response body is closed after processing
					// NewResponse reads the body, so it should handle closing or pass it back.
					// If NewResponse doesn't close it, close it here:
					// defer httpRes.Body.Close() // This might be too late if NewResponse errors

					// Let's modify NewResponse/makeHttpResString to ensure closure. (See structs/response_data.go)

				}
				resultsChan <- response // Send result (or error wrapper) to channel
			}
			// log.Printf("Worker %d finished", workerID)
		}(w)
	}

	// Send jobs to the workers
	for idx := range strategyClones {
		jobs <- idx
	}
	close(jobs) // Signal that no more jobs will be sent

	// Wait for all workers to complete processing
	workerWG.Wait()
	close(resultsChan) // Close results channel *after* workers are done

	attackDuration := time.Since(startTimeAttack)
	log.Printf("Workers finished. Attack duration: %s. Collecting results...", attackDuration)


	// --- Collect Results ---
	for response := range resultsChan {
		responseList.AddResponse(response) // AddResponse uses RequestId for index
	}

	// Assign collected URL and elapsed time
	responseList.URL = targetURL
	responseList.ElapsedTime = attackDuration.String() // Use the actual attack duration

	log.Printf("Collected %d results. Total time: %s", len(responseList.List), time.Since(startTimeOverall))

	// --- Render Response ---
	tmpl, err := template.ParseFiles("templates/attack.html")
	if err != nil {
		log.Printf("ERROR: Error parsing template 'attack.html': %v", err)
		http.Error(res, "Error rendering results page", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(res, responseList)
	if err != nil {
		log.Printf("ERROR: Error executing template: %v", err)
		// Don't try sending another error if headers are already sent
	}
	log.Println("Attack response sent to client.")
}