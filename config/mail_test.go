package config

import (
	"testing"

	"github.com/gernest/sydent-go/embed"
)

func TestEmbedEmailtemplates(t *testing.T) {
	e := embed.New()
	for _, v := range EmbedTemplates() {
		f, err := e.Open(v.Path)
		if err != nil {
			t.Errorf("%s :%v", v, err)
			continue
		}
		f.Close()
	}
}

func TestGetSample(t *testing.T) {
	c := Sample()
	c.WriteToFile("config.out")
}
