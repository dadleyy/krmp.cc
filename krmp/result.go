package krmp

import "bytes"

type Result struct {
	*bytes.Buffer
	ContentType string
}
