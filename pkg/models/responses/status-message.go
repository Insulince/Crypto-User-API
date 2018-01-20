package responses

type StatusMessage struct {
	Valid  bool `json:"valid"`
	Message string `json:"message"`
}
