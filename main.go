package main

import "github.com/adamroach/webrd/darwin"

func main() {
	c := darwin.NewCapturer()
	c.Start()
}
