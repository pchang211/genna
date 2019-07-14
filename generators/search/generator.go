package search

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/dizzyfool/genna/lib"
	"github.com/dizzyfool/genna/util"

	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

// Generator represents search generator
type Generator struct {
	genna.Genna
}

// New creates generator
func New(url string, logger *zap.Logger) Generator {
	return Generator{
		Genna: genna.New(url, logger),
	}
}

// Generate runs whole generation process
func (g Generator) Generate(options Options) error {
	options.def()

	entities, err := g.Read(options.Tables, options.FollowFKs, false)
	if err != nil {
		return xerrors.Errorf("read database error: %w", err)
	}

	parsed, err := template.New("search").Parse(templateSearch)
	if err != nil {
		return xerrors.Errorf("parsing template error: %w", err)
	}

	pack := NewTemplatePackage(entities, options)

	var buffer bytes.Buffer
	if err := parsed.ExecuteTemplate(&buffer, "search", pack); err != nil {
		return xerrors.Errorf("processing model template error: %w", err)
	}

	saved, err := util.FmtAndSave(buffer.Bytes(), options.Output)
	if err != nil {
		if !saved {
			return xerrors.Errorf("saving file error: %w", err)
		}
		g.Logger.Error("formatting file error", zap.Error(err), zap.String("file", options.Output))
	}

	g.Logger.Info(fmt.Sprintf("succesfully generated %d models\n", len(entities)))

	return nil
}
