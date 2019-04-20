package main

import "github.com/armon/circbuf"

//the serial backbuffer which holds the data displayed in the serial console

var scrollbackBuffer *circbuf.Buffer

func resetScrollback() {
	scrollbackBuffer, _ = circbuf.NewBuffer(int64(globalSettings.ScrollbackBufferSize))
}
