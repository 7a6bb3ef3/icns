package icns

import (
	"bytes"
)

type Icon struct {
	offset	int64
	Type	string
	Length	int64
	data	*bytes.Buffer
	Bytes	[]byte
}

func (i *Icon) Read(p []byte) (int ,error){
	return i.data.Read(p)
}
