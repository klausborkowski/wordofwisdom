package quotes

import (
	_ "embed"
)

var Quotes []string

//go:embed quotes.txt
var quotes []byte
