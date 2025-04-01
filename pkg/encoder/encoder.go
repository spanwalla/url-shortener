package encoder

type Encoder interface {
	Encode(text string, length int) string
}
