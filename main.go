package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gernest/sydent-go/matrixid"
	"github.com/gernest/sydent-go/store"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gernest/matrixid/embed"
	"github.com/gernest/sydent-go/core"
	"github.com/gernest/sydent-go/logger"
	"github.com/gernest/sydent-go/store/query"
	"github.com/gernest/sydent-go/store/schema"
	"go.uber.org/zap"

	"github.com/gernest/sydent-go/config"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "sydent-go"
	app.Usage = "matrix identity service in Go"
	app.Commands = []cli.Command{id()}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func id() cli.Command {
	return cli.Command{
		Name:  "id",
		Usage: "high performance matrix identity service",
		Action: func(ctx *cli.Context) error {
			file := ctx.Args().First()
			var c *config.Matrix
			var err error
			if file != "" {
				b, err := ioutil.ReadFile(file)
				if err != nil {
					return err
				}
				c, err = config.LoadFile(config.ProcessFile(b))
				if err != nil {
					return err
				}
			} else {
				c = config.LoadFromEnv()
			}

			if config.Empty(c) {
				return errors.New(missingConfigNotice)
			}

			lg, err := logger.New()
			if err != nil {
				return err
			}
			defer lg.Sync()
			if !c.Validate(lg) {
				return nil
			}
			db, err := sql.Open(c.DB.Driver, c.DB.Conn)
			if err != nil {
				return err
			}
			defer db.Close()
			opts := core.Ctx{
				Config: c,
				Log:    lg.With(zap.Namespace("matrix")),
			}
			vfs := embed.New()
			sq := query.New(db)
			driver, err := store.NewDriver(c.DB.Driver)
			if err != nil {
				return err
			}
			storeMetrics := store.NewMetric(prometheus.Opts{
				Namespace: "matrix",
				Subsystem: "storage",
			})
			storage := store.NewStore(sq, driver, storeMetrics)
			err = schema.IdentityUp(context.Background(), vfs, storage.DB())
			if err != nil {
				return err
			}
			mail, err := c.Email.Provider(c.GetTemplate())
			if err != nil {
				return err
			}
			opts.Email = mail
			opts.Store = storage
			m := matrixid.NewMetric()
			e := matrixid.Service(opts.Namespace("matrixid"), m)
			lg.Info("statring matrix identity service", zap.String("address", c.Server.Address()))
			return http.ListenAndServe(c.Server.Address(), e)
		},
	}
}

const missingConfigNotice = `Can't find configuration for starting this service

You can pass configuration file path as the first argument fot matrixid command like
	matrixid path/to/config.hcl

Or you can use environment variables eg
	MX_MODE=prod matrixid

Note that you can't mix, you can only apply one configuration i.e either file or 
environment variable but never both. If you use file, environment variables will
be ignored.
`
