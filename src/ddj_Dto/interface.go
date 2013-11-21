package ddj_Dto

type Dto interface {
	String() string
	Encode() ([]byte, error)
	Decode(buf []byte) error
	Size() int
}
