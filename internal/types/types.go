package types

type Result struct {
	Source string
	Data   map[string]interface{}
	Err    error
}

type ConsultCep interface {
	Fetch(cep string, ch chan<- Result)
	ToString(result Result) string
}
