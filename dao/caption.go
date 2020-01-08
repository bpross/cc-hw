package dao

// Captioner defines the interface for generating captions
type Captioner interface {
	Generate(string) ([]string, error)
}
