package common

type ContentInterface interface {
	Marshal()
	Unmarshal()
}

type BaseContent struct {
	Content string
	// Images  []string
}
