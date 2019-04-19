package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"text/template"

	"github.com/cenkalti/backoff"
	"github.com/gernest/sydent-go/embed"
	"github.com/gernest/sydent-go/logger"

	"github.com/hashicorp/hcl"
	"github.com/rodaine/hclencoder"
	"go.uber.org/zap"

	"github.com/gernest/signedjson"
	// we only support postgres
	_ "github.com/lib/pq"
)

// ApplicationName the name of the application.
const ApplicationName = "sydent-go"

// AgentName the name used for User-Agent header. All http clients making
// requests between matrix serves include this header.
const AgentName = "sydent-go"

// MaxRetries is the maximum number of retries for http requests made by SYDENT
// clients before backing off.
const MaxRetries = 16

// ThreepidSessionValidationTimeout duration in MS of the time a user will have
// to wait before the session he/she created is validates.
const ThreepidSessionValidationTimeout = 24 * 60 * 60 * 1000

// ThreepidSessionValidationLifetime duration in ms of the retaining period for
// sessions.
const ThreepidSessionValidationLifetime = 24 * 60 * 60 * 1000

// Matrix is the central configuration for the whole server. This stores
// settings for all services api and evertything needed to run the server.
type Matrix struct {
	Mode      string     `hcl:"mode"`
	Server    Server     `hcl:"server"`
	DB        DB         `hcl:"db"`
	Email     Email      `hcl:"email"`
	Templates []Template `hcl:"templates"`
	Peers     []Peer     `hcl:"peer"`
	tpl       *template.Template
}

func (m *Matrix) LoadTemplates() error {
	fs := embed.New()
	tpl := template.New("email")
	for k, v := range MergeTemplates(EmbedTemplates(), m.Templates) {
		t := tpl.New(k)
		b, err := readFile(fs, v.Path)
		if err != nil {
			return err
		}
		_, err = t.Parse(string(b))
		if err != nil {
			return err
		}
	}
	m.tpl = tpl
	return nil
}

func (m *Matrix) GetTemplate() *template.Template {
	return m.tpl
}

// Validate validate m and logs the error message in case m is not valid to lg.
func (m *Matrix) Validate(lg logger.Logger) (ok bool) {
	lg.Info("validating configuration ...")
	if v := m.Valid(); !v.IsValid() {
		v.log(lg)
		return false
	}
	lg.Info("validating configuration ... OK")
	return true
}

// Valid performs validation of m. Call Validation.IsValid() to see if
// validation passed.
func (m *Matrix) Valid() *Validation {
	v := &Validation{Namespace: "matrix"}
	v.add(m.Server)
	v.add(m.DB)
	err := m.LoadTemplates()
	if err != nil {
		v.Fields = append(v.Fields, Field{
			Name: "templates", Value: err.Error(),
		})
	}
	e := m.Email.Valid(m.GetTemplate())
	if !e.IsValid() {
		v.Children = append(v.Children, e)
	}

	return v
}

// Server defines setting for the main server. Unlike the reference
// implementation, all services are provided under the same server, service
// specific configuration is used to enable/disable services. So, this is the
// one and only server that runs on boot up.
//
// Only TLS server is supported. Failure to provide tls certificates will result
// in validation error.
//
// TODO: Add let's encrypt support with auto-tls.
type Server struct {
	Name           string `hcl:"name"`
	Port           string `hcl:"port"`
	ClientHTTPBase string `hcl:"client_http_base"`
	Crypto         Crypto `hcl:"crypto"`
}

// Valid validates s settings.
func (s Server) Valid() *Validation {
	v := &Validation{Namespace: "server"}
	if s.Name == "" {
		v.Set("name", missingField)
	}
	if s.Port == "" {
		v.Set("port", missingField)
	} else {
		_, err := strconv.Atoi(s.Port)
		if err != nil {
			v.Set("port", err.Error())
		}
	}
	if s.ClientHTTPBase == "" {
		v.Set("client_http_base", missingField)
	}
	if cv := s.Crypto.Valid(); !cv.IsValid() {
		v.Children = append(v.Children, cv)
	}
	return v
}

// Address returns the server address to bind to.
func (s Server) Address() string {
	return fmt.Sprintf(":%s", s.Port)
}

// Peer defines settings for remote replication.
type Peer struct {
	Name               string `hcl:",key"`
	BaseReplicationURL string `hcl:"base_replication_url"`
}

type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Validation struct {
	Namespace string        `json:"namespace"`
	Fields    []Field       `json:"fields,omitempty"`
	Children  []*Validation `json:",omitempty"`
}

func (v *Validation) Set(name, value string) {
	v.Fields = append(v.Fields, Field{Name: name, Value: value})
}

func (v Validation) String() string {
	b, _ := json.Marshal(v)
	return string(b)
}

func (v *Validation) add(vd Validator) {
	if val := vd.Valid(); !val.IsValid() {
		v.Children = append(v.Children, val)
	}
}

func (v *Validation) log(lg logger.Logger) {
	if !v.IsValid() {
		b, _ := json.Marshal(v)
		lg.Error("failed validating configuration",
			zap.String("validation_object", string(b)),
		)
	}
}

type Validator interface {
	Valid() *Validation
}

func (v Validation) IsValid() bool {
	return v.Fields == nil && v.Children == nil
}

const missingField = "missing"
const notValidFile = " not a valid file path"
const algorithmNotSupported = "not supported"

type DB struct {
	Driver string `hcl:"driver"`
	Conn   string `hcl:"conn"`
}

func (db DB) Valid() *Validation {
	v := &Validation{Namespace: "db"}
	if db.Driver == "" {
		v.Set("driver", missingField)
	} else {
		if db.Conn == "" {
			v.Set("conn", missingField)
		} else {
			q, err := EnsureConnected(db.Driver, db.Conn)
			if err != nil {
				v.Set("conn", err.Error())
			} else {
				defer q.Close()
				err = q.Ping()
				if err != nil {
					v.Set("conn", err.Error())
				}
			}
		}
	}
	return v
}

// EnsureConnected for automated deployments there are some cases when the
// database might be booting up so connections here will fall, this does
// exponential retries to connect to the database.
func EnsureConnected(driver, conn string) (db *sql.DB, err error) {
	b := backoff.NewExponentialBackOff()
	var xerr error
	err = backoff.Retry(func() error {
		db, xerr = sql.Open(driver, conn)
		if xerr != nil {
			return xerr
		}
		return nil
	}, b)
	return
}

// Crypto cryptographic keys used for signing and verifying messages.
type Crypto struct {
	Algorithm  string `hcl:"algorithm"`
	Version    string `hcl:"version"`
	SingingKey string `hcl:"signing_key"`
	VerifyKey  string `hcl:"verify_key"`
}

// Valid returns nil if c is a valid configuration for identity service crypto.
func (c Crypto) Valid() *Validation {
	v := &Validation{Namespace: "crypto"}
	var hasAlgorithm bool
	for _, v := range signedjson.SupportedAlgorithms {
		if v == c.Algorithm {
			hasAlgorithm = true
		}
	}
	if !hasAlgorithm {
		if c.Algorithm == "" {
			v.Set("algorithm", missingField)
		} else {
			v.Set("algorithm", algorithmNotSupported)
		}
	}
	if c.Version == "" {
		v.Set("version", missingField)
	}
	if c.SingingKey == "" {
		v.Set("signing_key", missingField)
	}
	if c.VerifyKey == "" {
		v.Set("verify_key", missingField)
	}
	return v
}

// Email defines all possible configuration for email client.
type Email struct {
	Providers    []Provider   `hcl:"provider"`
	Invite       Invite       `hcl:"invite"`
	Verification Verification `hcl:"verification"`
}

// embedded templates
var (
	InviteTpl                = "/email/invite_template.eml"
	InviteVector             = "/email/invite_template_vector.eml"
	VerificationTpl          = "/email/verification_template.eml"
	VerificationVector       = "/email/verification_template_vector.eml"
	VerifyResponsePage       = "/email/verify_response_page_template"
	VerifyResponsePageVector = "/email/verify_response_page_template_vector_im"
)

type Template struct {
	Name string `hcl:",key"`
	Path string `hcl:"path"`
}

func EmbedTemplates() []Template {
	return []Template{
		{
			Name: "invite", Path: InviteTpl,
		},
		{
			Name: "invite_vector", Path: InviteVector,
		},
		{
			Name: "verification", Path: VerificationTpl,
		},
		{
			Name: "verification_vector", Path: VerificationVector,
		},
		{
			Name: "verify_response", Path: VerifyResponsePage,
		},
		{
			Name: "verify_response_vector", Path: VerifyResponsePage,
		},
	}
}

func MergeTemplates(a, b []Template) map[string]Template {
	m := make(map[string]Template)
	if a == nil {
		for _, v := range b {
			m[v.Name] = v
		}
		return m
	}
	if b == nil {
		for _, v := range a {
			m[v.Name] = v
		}
		return m
	}
	for _, v := range a {
		m[v.Name] = v
	}
	for _, v := range b {
		m[v.Name] = v
	}
	return m

}
func (e Email) Provider(templates *template.Template) (EmailProvider, error) {
	for _, v := range e.Providers {
		if v.State == "enabled" {
			if v.Name == "smtp" {
				var s SMTPEmail
				s.Enabled = true
				if x, ok := v.Settings["username"]; ok {
					s.Username = x.(string)
				}
				if x, ok := v.Settings["password"]; ok {
					s.Password = x.(string)
				}
				if x, ok := v.Settings["host"]; ok {
					s.Host = x.(string)
				}
				if x, ok := v.Settings["port"]; ok {
					xs := x.(string)
					if xs != "" {
						m, err := strconv.Atoi(xs)
						if err != nil {
							log.Fatalf("incorrect port value :%v", err)
						}
						s.Port = int64(m)
					}
				}
				c := NewSMAPCLient(s)
				e, err := New(c, templates)
				if err != nil {
					return nil, err
				}
				return e, nil
			}
		}

	}
	return NoopMail{}, nil
}

func (e Email) Valid(templates *template.Template) *Validation {
	v := &Validation{Namespace: "email"}
	p, err := e.Provider(templates)
	if err != nil {
		v.Fields = append(v.Fields, Field{
			Name: "provider", Value: err.Error(),
		})
	} else {
		v.add(p)
	}
	v.add(e.Invite)
	return v
}

type EmailProvider interface {
	Validator
	Mail
}

// credit  https://github.com/asaskevich/govalidator
const email = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"

var emailRexex = regexp.MustCompile(email)

func IsValidEmail(email string) bool {
	return emailRexex.MatchString(email)
}

type Invite struct {
	From     string `hcl:"from"`
	Template string `hcl:"template"`
}

func (i Invite) Valid() *Validation {
	v := &Validation{Namespace: "invite"}
	if !emailRexex.MatchString(i.From) {
		v.Set("from", "not a valid email address")
	}
	return v
}

// Verification stores details used when sending token validation emails.
type Verification struct {
	From         string `hcl:"from"`
	Template     string `hcl:template"`
	ResponsePage string `hcl:"response_page_template"`
}

func (i Verification) Valid() *Validation {
	v := &Validation{Namespace: "invite"}
	if !emailRexex.MatchString(i.From) {
		v.Set("from", "not a valid email address")
	}
	return v
}

// Provider email provider settings.
type Provider Container

// Container abstract struct for storing namespaced settings.
type Container struct {
	Name     string                 `hcl:",key"`
	State    string                 `hcl:"state"`
	Settings map[string]interface{} `hcl:"settings"`
}

// SMTPEmail smtp plain auth configurations.
type SMTPEmail struct {
	Enabled  bool
	Username string
	Password string
	Host     string
	Port     int64
}

func (s SMTPEmail) Valid() *Validation {
	v := &Validation{Namespace: "smtp_provider"}
	if s.Username == "" {
		v.Set("username", missingField)
	}
	if s.Password == "" {
		v.Set("password", missingField)
	}
	if s.Password == "" {
		v.Set("host", missingField)
	}
	if s.Port == 0 {
		v.Set("port", missingField)
	}
	return v
}

func (c Crypto) Key() *signedjson.Key {
	pub, _ := signedjson.DecodeBase64(c.VerifyKey)
	priv, _ := signedjson.DecodeBase64(c.SingingKey)
	return &signedjson.Key{
		PrivateKey: priv,
		PublicKey:  pub,
		Alg:        c.Algorithm,
		Version:    c.Version,
	}
}

type Key struct {
	ID    string `hcl:"key"`
	Value string `hcl:"value"`
}

// ProcessFile expands environment variables in src. This means it will replace
// all $VAR or ${VAR} declarations with values found in environment variables.
func ProcessFile(src []byte) []byte {
	return []byte(os.ExpandEnv(string(src)))
}

// LoadFile decodes Matrix object from src. src is either json/hcl configuration
// in raw bytes.
func LoadFile(src []byte) (*Matrix, error) {
	var m Matrix
	err := hcl.Decode(&m, string(src))
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// WriteToFile marshals m and writes the hcl configuration to the the file
// filename.
func (m *Matrix) WriteToFile(filename string) error {
	b, err := hclencoder.Encode(*m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0600)
}

// Sample sample home server configuration.
func Sample() *Matrix {
	return &Matrix{
		Mode: "prod",
		Server: Server{
			Name:           ApplicationName,
			Port:           "9891",
			ClientHTTPBase: "https://localhost:9891",
			Crypto: Crypto{
				Algorithm:  "ed25519",
				Version:    "0",
				SingingKey: "${SYDENT_PRIVATE_KEY}",
				VerifyKey:  "${SYDENT_PUBLIC_KEY}",
			},
		},
		DB: DB{
			Driver: "postgres",
			Conn:   "${SYDENT_DB_CONN}",
		},
		Email: Email{
			Providers: []Provider{
				{
					Name:  "smtp",
					State: "enabled",
					Settings: map[string]interface{}{
						"host":     "$SYDENT_SMTP_HOST",
						"port":     "$SYDENT_SMTP_PORT",
						"username": "$SYDENT_SMTP_USERNAME",
						"password": "$SYDENT_SMTP_PASSWORD",
					},
				},
				{
					Name:  "sendgrid",
					State: "disabled",
					Settings: map[string]interface{}{
						"api_key": "$SYDENT_SENDGRID_APIKEY",
					},
				},
			},
			Invite: Invite{
				From:     "$SYDENT_SMTP_USERNAME",
				Template: "invite",
			},
			Verification: Verification{
				From:         "$SYDENT_SMTP_USERNAME",
				Template:     "verification",
				ResponsePage: "verify_response",
			},
		},
	}
}
