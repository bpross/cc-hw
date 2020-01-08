package caption

// Generator defines the interface for generating captions
type Generator interface {
	Create(string, int) ([]string, error)
}
