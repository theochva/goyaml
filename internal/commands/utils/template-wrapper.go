package utils

import (
	htmlTemplate "html/template"
	"io"
	"path"
	"text/template"
)

// TemplateWrapper - adapter type around either html/template or text/template packages
type TemplateWrapper struct {
	tmpl             *template.Template
	htmlTmpl         *htmlTemplate.Template
	parseFilesFn     func(filenames ...string) error
	parseGlobFn      func(pattern string) error
	executeFn        func(wr io.Writer, data interface{}) error
	createTextTmplFn func(templateText string) error
}

// NewTextTemplate - create new text/template
func (w TemplateWrapper) NewTextTemplate(templateText string) error {
	return w.createTextTmplFn(templateText)
}

// ParseFiles - parse files
func (w TemplateWrapper) ParseFiles(filenames ...string) error {
	return w.parseFilesFn(filenames...)
}

// ParseGlob - parse glob
func (w TemplateWrapper) ParseGlob(pattern string) error {
	return w.parseGlobFn(pattern)
}

// Execute - execute
func (w TemplateWrapper) Execute(wr io.Writer, data interface{}) error {
	return w.executeFn(wr, data)
}

// NewTemplateWrapper - create new template wrapper
func NewTemplateWrapper(format string) (w TemplateWrapper) {

	if format == "text" {
		newTmpl := func(name string) *template.Template {
			return template.New(name).Funcs(template.FuncMap{
				"getenv": getEnv,
			})
		}
		w.createTextTmplFn = func(templateText string) (err error) {
			w.tmpl, err = newTmpl("main").Parse(templateText)
			return
		}
		w.parseFilesFn = func(filenames ...string) (err error) {
			if w.tmpl == nil {
				// We need to name the template.
				// So we will use the file name (without dir)
				_, file := path.Split(filenames[0])
				w.tmpl = newTmpl(file)
			}
			w.tmpl, err = w.tmpl.ParseFiles(filenames...)

			return
		}
		w.parseGlobFn = func(pattern string) (err error) {
			if w.tmpl == nil {
				// We need to name the template.  So, we
				// will use the first matching file
				var file string
				if file, err = getFirstMatchedFile(pattern); err != nil {
					return err
				}
				w.tmpl = newTmpl(file)
			}
			w.tmpl, err = w.tmpl.ParseGlob(pattern)

			return
		}
		w.executeFn = func(wr io.Writer, data interface{}) error {
			return w.tmpl.Execute(wr, data)
		}
	} else {
		newTmpl := func(name string) *htmlTemplate.Template {
			return htmlTemplate.New(name).Funcs(htmlTemplate.FuncMap{
				"getenv": getEnv,
			})
		}

		w.createTextTmplFn = func(templateText string) (err error) {
			w.htmlTmpl, err = newTmpl("main").Parse(templateText)
			return
		}
		w.parseFilesFn = func(filenames ...string) (err error) {
			if w.htmlTmpl == nil {
				// We need to name the template.
				// So we will use the file name (without dir)
				_, file := path.Split(filenames[0])
				w.htmlTmpl = newTmpl(file)
			}
			w.htmlTmpl, err = w.htmlTmpl.ParseFiles(filenames...)

			return
		}
		w.parseGlobFn = func(pattern string) (err error) {
			if w.htmlTmpl == nil {
				// We need to name the template.  So, we
				// will use the first matching file
				var file string
				if file, err = getFirstMatchedFile(pattern); err != nil {
					return err
				}
				w.htmlTmpl = newTmpl(file)
			}
			w.htmlTmpl, err = w.htmlTmpl.ParseGlob(pattern)

			return
		}
		w.executeFn = func(wr io.Writer, data interface{}) error {
			return w.htmlTmpl.Execute(wr, data)
		}
	}

	return w
}
