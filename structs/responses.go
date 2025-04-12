package structs

type Responses struct {
	List        []ResponseData // Slice to hold results
	URL         string         // Target URL (first successful request)
	ElapsedTime string         // Total attack duration string
}

// NewResponses creates a Responses struct with a pre-allocated slice
func NewResponses(size int) Responses {
	if size <= 0 {
		size = 0 // Handle edge case of zero payloads
	}
	return Responses{
		// Pre-allocate the slice to avoid append overhead and potential race conditions
		// Note: Indices will be 0 to size-1, corresponding to RequestId 1 to size.
		List: make([]ResponseData, size),
	}
}

// AddResponse adds a response to the correct index in the pre-allocated slice.
// It's crucial that RequestId is correctly set (index + 1).
func (r *Responses) AddResponse(response ResponseData) {
	index := int(response.RequestId - 1)
	// Basic bounds check
	if index >= 0 && index < len(r.List) {
		r.List[index] = response
	} else {
		// This should not happen if NewResponses and RequestId are correct
		// Log this error, as it indicates a logic flaw
		// Consider appending if List wasn't pre-allocated, but with pre-allocation, this is an error.
		// log.Printf("ERROR: Attempted to add response with invalid RequestId %d (index %d) to list of size %d",
		// 	response.RequestId, index, len(r.List))
		// Fallback: Append if list wasn't sized correctly (though it should be)
		if index == len(r.List) { // Only allow appending if it's the very next item sequentially
			r.List = append(r.List, response)
		} else {
			// Log error - data might be lost or incorrect
		}
	}
}