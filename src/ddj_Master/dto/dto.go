//This package contains structures that are data transfer objects (DTO). These structures are send between node and master
//and  between packages in master.
package dto

//This is main interface in project. Structures which implements this interface could be send or received from node
type Dto interface {
	Encoder
	Decoder
	String() string
	Size() int
}

//This interface is responsible for encoding structure to slice of bytes
type Encoder interface {
	Encode() ([]byte, error)
}

//This interface is responsible for decoding slice of bytes to object.
type Decoder interface {
	Decode(buf []byte) error
}
