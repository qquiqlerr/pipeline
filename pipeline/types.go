package pipeline
type Block struct {
	Name     string                 `json:"name"`
	Function string                 `json:"function"`
	Input    map[string]interface{} `json:"input"`
	Output   []string               `json:"output"`
}

type Function func(input map[string]any) any

type blockEnv struct {
	Block Block
	In    map[string]interface{}
}