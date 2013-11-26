package dto

type Dto interface {
	Encoder
	Decoder
	String() string
	Size() int
}

type InDto interface {
	Decoder
	String() string
	Size() int
}

type OutDto interface {
	Encoder
	String() string
	Size() int
}

type Encoder interface {
	Encode() ([]byte, error)
}

type Decoder interface {
	Decode(buf []byte) error
}
