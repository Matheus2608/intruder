package structs

type Responses struct {
	List        []ResponseData
	URL         string
	ElapsedTime string
}

func NewResponses(size int) Responses {
	return Responses{
		List: make([]ResponseData, size),
	}
}

func (r *Responses) AddResponse(response ResponseData) {
	r.List[response.RequestId-1] = response
}
