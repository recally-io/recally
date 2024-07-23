package document

// Document is the struct for interacting with a document.
type Document struct {
	Content  string
	Metadata map[string]any
	Score    float32
}

