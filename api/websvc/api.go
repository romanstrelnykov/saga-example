package websvc

type request struct {
	RequestID string `json:"requestId"`
	Action    string `json:"type"`
}
