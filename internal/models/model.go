package models

type InputOutput struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}
