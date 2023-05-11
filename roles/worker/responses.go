package worker

type Response struct {
	Status  int    `json:"status"`
	Payload string `json:"payload"`
}
