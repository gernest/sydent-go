package embed

import (
	"log"
	"net/http"

	_ "github.com/gernest/sydent-go/embed/statik"
	"github.com/rakyll/statik/fs"
)

//go:generate statik -src=files/   -f

// Embed is an alias for http.FileSystem that contains all embeded files used by
// matrixid.
type Embed = http.FileSystem

var vfs Embed

func init() {
	h, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	vfs = h
}

func New() Embed {
	return vfs
}
