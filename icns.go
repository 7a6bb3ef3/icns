package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type ICNS struct {
	Length	int64
	rawData	[]byte
	Icons	[]*Icon
}

func (i *ICNS) parseHeader() error {
	if !(i.rawData[0] == 'i' && i.rawData[1] == 'c' && i.rawData[2] == 'n' && i.rawData[3] == 's') {
		return fmt.Errorf("magic literal: %x" ,i.rawData[:4])
	}
	i.Length = int64(i.rawData[4]) << 24 + int64(i.rawData[5]) << 16 + int64(i.rawData[6]) << 8 + int64(i.rawData[7])
	return nil
}

func (i *ICNS) parseData(offset int64) (*Icon ,error) {
	if len(i.rawData) < int(offset) + 8 {
		return nil ,errors.New("size of icon header < 8")
	}
	ico := &Icon{
		offset: offset,
		Type:   string(i.rawData[offset:offset+4]),
		Length: int64(i.rawData[offset + 4]) << 24 + int64(i.rawData[offset + 5]) << 16 + int64(i.rawData[offset + 6]) << 8 + int64(i.rawData[offset + 7]) - 8,
	}
	end := int(offset) + int(ico.Length)
	if len(i.rawData) < end {
		return nil ,fmt.Errorf("can not read icon data: %d ,%d" ,len(i.rawData) ,end)
	}
	ico.Bytes = i.rawData[offset+8:end]
	ico.data = bytes.NewBuffer(ico.Bytes)
	return ico ,nil
}



func Decode(r io.Reader) (*ICNS ,error) {
	buf := make([]byte ,1024)
	dst := &bytes.Buffer{}
	for true{
		n ,e := r.Read(buf)
		if errors.Is(e ,io.EOF){
			break
		}else if e != nil {
			return nil ,e
		}
		dst.Write(buf[:n])
	}
	img := &ICNS{
		rawData: dst.Bytes(),
	}
	if e := img.parseHeader();e != nil {
		return nil ,fmt.Errorf("parseHeader: %w" ,e)
	}

	var offset int64 = 8
	for true{
		if offset >= img.Length{
			break
		}
		ico ,e := img.parseData(offset)
		if e != nil {
			return nil ,e
		}
		offset += 8 + ico.Length
		img.Icons = append(img.Icons ,ico)
	}
	return img ,nil
}



