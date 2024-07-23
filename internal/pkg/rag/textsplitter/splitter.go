package textsplitter

// TextSplitter is the standard interface for splitting texts.
type TextSplitter interface {
	Split(text string) ([]string, error)
}
