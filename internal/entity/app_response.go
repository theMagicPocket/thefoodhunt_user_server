package entity

type envelope struct {
	Ok      bool        `json:"ok"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Ok(message string, data interface{}) envelope {
	res := envelope{
		Ok:      true,
		Message: message,
	}
	if data != nil {
		res.Data = data
	}
	return res
}

func NotOk(message string) envelope {
	res := envelope{
		Ok:      false,
		Message: message,
	}

	return res
}
