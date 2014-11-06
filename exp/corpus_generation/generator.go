package converter

type Generator interface {
	NumChannels() int
	// Returns nil to indicate that no feature generated.
	Feature(records []*Record, year, channel int) []string
}
