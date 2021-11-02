package main

import "io"

type Converter interface {
	Copy(ico Icon ,w io.Writer) error
	Format() string
}


type PNGConverter struct {

}