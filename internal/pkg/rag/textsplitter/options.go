package textsplitter

const (
	_defaultTokenChunkSize    = 1000
	_defaultTokenChunkOverlap = 200
)

// Options is a struct that contains options for a text splitter.
type Options struct {
	ChunkSize     int
	ChunkOverlap  int
	Separators    []string
	LenFunc       func(string) int
	KeepSeparator bool
}

// DefaultOptions returns the default options for all text splitter.
func DefaultOptions() Options {
	return Options{
		ChunkSize:     _defaultTokenChunkSize,
		ChunkOverlap:  _defaultTokenChunkOverlap,
		Separators:    []string{"\n\n", "\n", " ", ""},
		LenFunc:       func(s string) int { return len(s) },
		KeepSeparator: true,
	}
}

type Option func(*Options)

// WithChunkSize sets the chunk size for a text splitter.
func WithChunkSize(chunkSize int) Option {
	return func(o *Options) {
		o.ChunkSize = chunkSize
	}
}

// WithChunkOverlap sets the chunk overlap for a text splitter.
func WithChunkOverlap(chunkOverlap int) Option {
	return func(o *Options) {
		o.ChunkOverlap = chunkOverlap
	}
}

// WithSeparators sets the separators for a text splitter.
func WithSeparators(separators []string) Option {
	return func(o *Options) {
		o.Separators = separators
	}
}

// WithLenFunc sets the lenfunc for a text splitter.
func WithLenFunc(lenFunc func(string) int) Option {
	return func(o *Options) {
		o.LenFunc = lenFunc
	}
}

// WithKeepSeparator sets the keep separator for a text splitter.
func WithKeepSeparator(keepSeparator bool) Option {
	return func(o *Options) {
		o.KeepSeparator = keepSeparator
	}
}
