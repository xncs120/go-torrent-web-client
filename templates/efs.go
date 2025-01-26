package templates

import (
	_ "embed"
)

//go:embed "index.html"
var IndexHTML []byte

//go:embed "player.html"
var PlayerHTML []byte
