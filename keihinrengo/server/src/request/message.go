package request

type (
	Message struct {
		Id       int    `json:"id"`
		Body     string `json:"body"`
		UserName string `json:"user_name"`
	}
)
