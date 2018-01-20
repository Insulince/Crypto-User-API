package responses

type StatusError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
