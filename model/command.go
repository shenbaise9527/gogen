package model

import (
	"strings"

	"github.com/shenbaise9527/gogen/model/gen"

	"github.com/shenbaise9527/gogen/model/schemas"

	"github.com/urfave/cli/v2"
)

const (
	flagURL   = "url"
	flagTable = "table"
	flagDir   = "dir"
	flagCache = "cache"
)

func SQLDataSource(ctx *cli.Context) error {
	url := strings.TrimSpace(ctx.String(flagURL))
	dir := strings.TrimSpace(ctx.String(flagDir))
	cache := ctx.Bool(flagCache)
	tablenames := strings.Split(strings.TrimSpace(ctx.String(flagTable)), ",")
	tables, err := schemas.GetTableInfos(tablenames, url, cache)
	if err != nil {
		return err
	}

	g, err := gen.NewGenerator(dir, cache)
	if err != nil {
		return err
	}

	err = g.Start(tables)

	return err
}
