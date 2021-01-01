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
)

func SQLDataSource(ctx *cli.Context) error {
	url := strings.TrimSpace(ctx.String(flagURL))
	dir := strings.TrimSpace(ctx.String(flagDir))
	tablenames := strings.Split(strings.TrimSpace(ctx.String(flagTable)), ",")

	// 获取表信息.
	tables, err := schemas.GetTableInfos(tablenames, url)
	if err != nil {
		return err
	}

	g, err := gen.NewGenerator(dir)
	if err != nil {
		return err
	}

	err = g.Start(tables)

	return err
}
