package request

type (
	Message struct {
		Body     string `json:"body"`
		UserName string `json:"user_name"`
	}
)
