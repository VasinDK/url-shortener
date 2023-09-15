package responce

type Responce struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "Ok"
	StatusError = "Error"
)

func Ok() Responce {
	return Responce{
		Status: StatusOk,
	}
}

func Error(msg string) Responce {
	return Responce{
		Status: StatusError,
		Error:  msg,
	}
}
