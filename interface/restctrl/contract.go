package restctrl

// RestRequest ...
type RestRequest struct {
	Body []byte
}

// RestResponse ...
type RestResponse struct {
	Body       []byte
	StatusCode int
}
