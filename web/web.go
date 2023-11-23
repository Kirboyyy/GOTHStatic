package web

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin/render"
)

const (
	BlogName = "David Caudill"
	PageName = "David Caudill | Software Engineer"
)

type TemplRender struct {
	Code int
	Data templ.Component
}

func (t TemplRender) Render(w http.ResponseWriter) error {
	t.WriteContentType(w)
	w.WriteHeader(t.Code)
	if t.Data != nil {
		return t.Data.Render(context.Background(), w)
	}
	return nil
}

func (t TemplRender) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (t *TemplRender) Instance(name string, data interface{}) render.Render {
	if templData, ok := data.(templ.Component); ok {
		return &TemplRender{
			Code: http.StatusOK,
			Data: templData,
		}
	}
	return nil
}

// big thanks to https://github.com/stephenafamo / https://github.com/a-h/templ/issues/175
func Raw(s string, errs ...error) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := fmt.Fprint(w, s)
		return errors.Join(append(errs, err)...)
	})
}
