package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/shenbaise9527/gogen/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli/v2"
)

var (
	BuildVersion = "0.0.1"
	commands     = []*cli.Command{
		{
			Name:  "model",
			Usage: "generate model code, only support mysql",
			Subcommands: []*cli.Command{
				{
					Name:  "datasource",
					Usage: "generate model from datasource",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "url",
							Usage: `data soucre of database, mysql: "root:password@tcp(127.0.0.1:3306)/database"`,
						},
						&cli.StringFlag{
							Name:  "table, t",
							Usage: `the tables in the database,support for comma separation`,
						},
						&cli.StringFlag{
							Name:  "dir, d",
							Usage: "the target dir",
						},
						&cli.StringFlag{
							Name:  "cache, c",
							Usage: "generate code with cache [optional]",
						},
						&cli.StringFlag{
							Name:  "tracing",
							Usage: "generate code with tracing [optional]",
						},
					},
					Action: model.SQLDataSource,
				},
			},
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Usage = "a cli tool to generate code"
	app.Version = fmt.Sprintf("%s %s/%s", BuildVersion, runtime.GOOS, runtime.GOARCH)
	app.Commands = commands

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
