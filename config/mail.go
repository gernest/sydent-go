package config

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io/ioutil"
	"text/template"
	"time"

	"github.com/gernest/sydent-go/embed"
	"github.com/gernest/sydent-go/models"
)

var _ EmailProvider = (*Kmail)(nil)
var _ EmailProvider = NoopMail{}

// Mail is an interface for sending transactional email for the matrixid
// home server.
type Mail interface {
	SendMail(ctx context.Context, tmpl, from string, to []string, data map[string]string) error
}

// NoopMail implements Mail intreface but does not actually send the email.
type NoopMail struct{}

// SendMail does nothing and always returns nil.
func (n NoopMail) SendMail(ctx context.Context, tmpl, from string, to []string, data map[string]string) error {
	return nil
}

func (n NoopMail) Valid() *Validation {
	return &Validation{Namespace: "noop"}
}

const dateFormat = "Mon, 02 Jan 2006 15:04:05 -0700"
const dateFormatGMT = "Mon, 02 Jan 2006 15:04:05 MST"

// Kmail provide methods for sending emails.
type Kmail struct {
	tpl    *template.Template
	emebd  embed.Embed
	client Client
}

func readFile(vfs embed.Embed, tpl string) ([]byte, error) {
	f, err := vfs.Open(tpl)
	if err != nil {
		return ioutil.ReadFile(tpl)
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func New(client Client, tpl *template.Template) (*Kmail, error) {
	return &Kmail{tpl: tpl, client: client}, nil
}

// SendMail uses tmpl template to send and email, data is a context object which
// is passed to the cached template and rendered to generate the email message.
func (k *Kmail) SendMail(ctx context.Context, tmpl, from string, to []string, data map[string]string) error {
	if ctx.Err() != nil {
		return nil
	}
	if data == nil {
		data = make(map[string]string)
	}
	midRand := models.RandomString(16)
	now := time.Now()
	data["messageid"] = fmt.Sprintf("<%d%s%s>", models.MS(&now), midRand, k.client.Host())
	data["date"] = FormatDate(now, false, false)
	data["from"] = from
	data["to"] = to[0]
	for k, v := range data {
		data[k+"_forhtml"] = html.EscapeString(v)
		data[k+"_forurl"] = fmt.Sprintf("%q", v)
	}
	var buf bytes.Buffer
	err := k.tpl.ExecuteTemplate(&buf, tmpl, data)
	if err != nil {
		return err
	}
	return k.client.Send(from, to, buf.Bytes())
}

func (k *Kmail) Valid() *Validation {
	return k.client.Valid()
}

// FormatDate similar to python email.utils.formatdate
func FormatDate(ts time.Time, local, useGMT bool) string {
	s := dateFormat
	if local && useGMT {
		s = dateFormatGMT
	}
	if !local {
		ts = ts.UTC()
	}
	return ts.Format(s)
}
