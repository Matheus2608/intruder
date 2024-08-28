package structs

type Responses struct {
	List []ResponseData
	URL  string
}

func NewResponses(size int, url string) Responses {
	return Responses{
		List: make([]ResponseData, size),
		URL:  url,
	}
}

func (r *Responses) AddResponse(response ResponseData) {
	r.List[response.RequestId-1] = response
}
