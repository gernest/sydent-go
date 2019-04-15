package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/gernest/sydent-go/config"

	"github.com/gernest/sydent-go/core"
	"github.com/gernest/matrixid/embed"
	"github.com/gernest/sydent-go/logger"
	"github.com/gernest/sydent-go/store"
	"github.com/gernest/sydent-go/store/query"
	"github.com/gernest/sydent-go/store/schema"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"go.uber.org/zap"
)

var _ prometheus.Counter = (*TestCounter)(nil)
var _ Metric = (*TestMetric)(nil)
var _ logger.Logger = (*TestLogger)(nil)

type TestMetric struct {
	counters map[string]*TestCounter
}

func (tm *TestMetric) CountError(h string) prometheus.Counter {
	c := &TestCounter{}
	if tm.counters == nil {
		tm.counters = make(map[string]*TestCounter)
	}
	tm.counters[h] = c
	return c
}

type TestCounter struct {
	emptyCollector
	emptyMetric
	n float64
}

func (tc *TestCounter) Inc() {
	tc.n++
}

func (tc *TestCounter) Add(v float64) {
	tc.n += v
}

type emptyCollector struct{}

func (emptyCollector) Describe(chan<- *prometheus.Desc) {}
func (emptyCollector) Collect(chan<- prometheus.Metric) {}

type emptyMetric struct{}

func (emptyMetric) Desc() *prometheus.Desc  { return nil }
func (emptyMetric) Write(*dto.Metric) error { return nil }

type logEntry struct {
	With   []zap.Field
	Level  string
	Msg    string
	Fields []zap.Field
}

type TestLogger struct {
	with    []zap.Field
	entries []logEntry
}

func (lg *TestLogger) Info(msg string, fields ...zap.Field) {
	lg.add("info", msg, fields...)
}

func (lg *TestLogger) Error(msg string, fields ...zap.Field) {
	lg.add("error", msg, fields...)
}

func (lg *TestLogger) Sync() error {
	return nil
}

func (lg *TestLogger) add(level, msg string, fields ...zap.Field) {
	lg.entries = append(lg.entries, logEntry{
		With:   lg.with,
		Level:  level,
		Msg:    msg,
		Fields: fields,
	})
}

func (lg *TestLogger) String() string {
	b, _ := json.Marshal(lg.entries)
	return string(b)
}

func (lg *TestLogger) With(fields ...zap.Field) logger.Logger {
	return &TestLogger{
		with: append(lg.with, fields...),
	}
}

var sqlConn *sql.DB
var mainContext *core.Ctx

func TestMain(m *testing.M) {
	driverName := "postgres"
	conn := os.Getenv("MATRIXID_DB_CONN")
	db, err := sql.Open(driverName, conn)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlConn.Close()
	sqlConn = db
	q := query.New(db)
	fs := embed.New()
	ctx := context.Background()
	err = schema.IdentityDown(ctx, fs, q)
	if err != nil {
		log.Fatal(err)
	}
	err = schema.IdentityUp(ctx, fs, q)
	if err != nil {
		log.Fatal(err)
	}
	mctx, err := NewTestCtx(store.Metric{})
	if err != nil {
		log.Fatal(err)
	}
	mainContext = mctx
	code := m.Run()
	err = schema.IdentityUp(ctx, fs, q)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func NewTestCtx(m store.Metric) (*core.Ctx, error) {
	b, err := ioutil.ReadFile("fixture/config.hcl")
	if err != nil {
		return nil, err
	}
	c, err := config.LoadFile(b)
	if err != nil {
		return nil, err
	}
	err = c.LoadTemplates()
	if err != nil {
		return nil, err
	}
	drv, err := store.NewDriver("postgres")
	if err != nil {
		return nil, err
	}
	db := store.NewStore(query.New(sqlConn), drv, m)
	return &core.Ctx{
		Config: c,
		Store:  db,
		Log:    &TestLogger{},
	}, nil
}

type TestEmailClient struct {
	send func(from string, to []string, msg []byte) error
	host string
}

// SendMail does nothing and always returns nil.
func (n TestEmailClient) Send(from string, to []string, msg []byte) error {
	if n.send != nil {
		return n.send(from, to, msg)
	}
	return nil
}

func (n TestEmailClient) Host() string {
	return n.host
}

func (n TestEmailClient) Valid() *config.Validation {
	return &config.Validation{Namespace: "noop"}
}
