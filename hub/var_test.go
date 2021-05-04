package main

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_chunk2map(t *testing.T) {
	m := chunk2map("./services.md", "```", "```", "=", "$")
	spew.Dump(m)
}
