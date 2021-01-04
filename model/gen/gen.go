package gen

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/shenbaise9527/gogen/model/schemas"
)

type generator struct {
	dir       string
	pkg       string
	withCache bool
}

func NewGenerator(dir string, cache bool) (*generator, error) {
	if dir == "" {
		return nil, errors.New("the target dir is empty")
	}

	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	dir = dirAbs
	pkg := filepath.Base(dirAbs)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	}

	if err != nil {
		return nil, err
	}

	gen := &generator{
		dir:       dir,
		pkg:       pkg,
		withCache: cache,
	}

	return gen, nil
}

func (g *generator) Start(matchTables []*schemas.Table) error {
	for _, table := range matchTables {
		modelFileName := table.LowerStartCamelObject() + "model"
		name := modelFileName + ".go"
		err := g.buildFile(name, func() (string, error) {
			return genModel(g.pkg, g.withCache, table)
		})

		if err != nil {
			return err
		}
	}

	err := g.buildFile("builder.go", func() (string, error) {
		return genBuilder(g.pkg)
	})

	if err != nil {
		return err
	}

	err = g.buildFile("dbconn.go", func() (string, error) {
		return genDBConn(g.pkg)
	})

	return err
}

func (g *generator) buildFile(name string, f func() (string, error)) error {
	filename := filepath.Join(g.dir, name)
	output, err := f()
	if err != nil {
		return err
	}

	_, err = os.Stat(filename)
	if err == nil {
		fmt.Printf("%s already exists, ignored.\n", name)

		return nil
	}

	err = ioutil.WriteFile(filename, []byte(output), os.ModePerm)

	return err
}
