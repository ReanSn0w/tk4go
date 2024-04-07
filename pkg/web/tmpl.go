package web

import (
	"bytes"
	"context"
	"html/template"
	"net/http"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

// NewTemplate создает новый экземпляр шаблона
func NewTemplate(log tools.Logger, dir, errorTMPL string) (*Template, error) {
	tmpl := template.New("")
	return ParseTemplate(log, tmpl, dir, errorTMPL)
}

// ParseTemplate парсит шаблоны из директории
func ParseTemplate(log tools.Logger, tmpl *template.Template, dir string, errorTMPL string) (*Template, error) {
	tmpl, err := tmpl.ParseGlob(dir + "/*.html")
	return &Template{tmpl: tmpl, log: log, errorTMPL: errorTMPL}, err
}

type (
	Template struct {
		tmpl      *template.Template
		log       tools.Logger
		errorTMPL string
	}

	HTMLResponse interface {
		Context() context.Context
		Template() string
		Data() any
	}
)

// Write записывает шаблон в ответ
func (t *Template) Write(w http.ResponseWriter, code int, data HTMLResponse) {
	buf := new(bytes.Buffer)

	err := t.tmpl.ExecuteTemplate(buf, data.Template(), data)
	if err != nil {
		if _, ok := data.Data().(error); ok || t.errorTMPL == "" {
			t.log.Logf("[ERROR] template error: %v", err)
			return
		}

		data := NewHTMLData(data.Context(), t.errorTMPL, err)
		t.Write(w, http.StatusInternalServerError, data)
		return
	}

	w.WriteHeader(code)
	buf.WriteTo(w)
}

// MARK: - Базовая реализация структуры для ответа в шаблоне

type (
	HTMLData struct {
		template string
		context  context.Context
		data     any
	}
)

func NewHTMLData(ctx context.Context, tmpl string, data any) *HTMLData {
	return &HTMLData{context: ctx, data: data, template: tmpl}
}

func (h *HTMLData) Context() context.Context {
	return h.context
}

func (h *HTMLData) Template() string {
	return h.template
}

func (h *HTMLData) Data() any {
	return h.data
}
