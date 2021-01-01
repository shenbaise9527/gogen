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
	dir string
	pkg string
}

func NewGenerator(dir string) (*generator, error) {
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
		dir: dir,
		pkg: pkg,
	}

	return gen, nil
}

func (g *generator) Start(matchTables []*schemas.Table) error {
	for _, table := range matchTables {
		modelFileName := table.LowerStartCamelObject() + "model"
		name := modelFileName + ".go"
		filename := filepath.Join(g.dir, name)
		_, err := os.Stat(filename)
		if err == nil {
			fmt.Printf("%s already exists, ignored.\n", name)

			continue
		}

		output, err := genModel(g.pkg, table)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filename, []byte(output), os.ModePerm)
		if err != nil {
			return err
		}
	}

	filename := filepath.Join(g.dir, "builder.go")
	output, err := genBuilder(g.pkg)
	if err != nil {
		return err
	}

	_, err = os.Stat(filename)
	if err == nil {
		fmt.Println("builder.go already exists, ignored.")

		return nil
	}

	err = ioutil.WriteFile(filename, []byte(output), os.ModePerm)

	return err
}
