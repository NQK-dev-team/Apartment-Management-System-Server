package structs

import "net/textproto"

type CustomFileStruct struct {
	Filename string
	Header   textproto.MIMEHeader
	Size     int64
	Content  []byte
}
