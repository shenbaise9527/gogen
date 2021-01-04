package template

import (
	"bytes"
	goformat "go/format"
	gotemplate "text/template"
)

type defaultTemplate struct {
	name  string
	text  string
	goFmt bool
}

func With(name string) *defaultTemplate {
	return &defaultTemplate{
		name: name,
	}
}

func (t *defaultTemplate) Parse(text string) *defaultTemplate {
	t.text = text

	return t
}

func (t *defaultTemplate) GoFmt(format bool) *defaultTemplate {
	t.goFmt = format

	return t
}

func (t *defaultTemplate) Execute(data interface{}) (*bytes.Buffer, error) {
	tpl, err := gotemplate.New(t.name).Parse(t.text)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	if err != nil {
		return nil, err
	}

	if !t.goFmt {
		return buf, nil
	}

	formatOutput, err := goformat.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}

	buf.Reset()
	buf.Write(formatOutput)

	return buf, nil
}
