package converter

type Generator interface {
	NumChannels() int
	// Returns "" indicate not generating a feature.
	Feature(record *Record, channel int) string
}
