package core

type Embedder interface {
	EmbedText(text string) ([]float32, error)
}

type Retriever interface {
	Search(query string, topK int) ([]Document, error)
}

type LLMClient interface {
	Generate(prompt string) (string, error)
}

type Cache interface {
	Get(key string) (string, bool)
	Set(key string, value string) error
}

type Document struct {
	ID       string
	Content  string
	Metadata map[string]string
	Vector   []float32
}
