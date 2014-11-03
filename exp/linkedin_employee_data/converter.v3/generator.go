package converter

type Generator interface {
	NumChannels() int
	// Returns "" indicate not generating a feature.
	Feature(record *Record, year, channel int) []string
}
