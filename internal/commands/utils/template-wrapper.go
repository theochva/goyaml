package utils

import (
	"bytes"
	"fmt"
	htmlTemplate "html/template"
	"io"
	"path"
	"text/template"

	"github.com/Masterminds/sprig/v3"
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

func getFuncs() map[string]interface{} {
	funcs := map[string]interface{}{
		"mustToToml":   mustToTOML,
		"toToml":       toTOML,
		"mustToYaml":   mustToYAML,
		"toYaml":       toYAML,
		"mustFromYaml": mustFromYAML,
		"fromYaml":     fromYAML,
		"mustToJson":   mustToJSON,
		"toJson":       toJSON,
		"mustFromJson": mustFromJSON,
		"fromJson":     fromJSON,
		"getenv":       getEnv,
	}

	return funcs
}

// NewTemplateWrapper - create new template wrapper
func NewTemplateWrapper(format string) (w TemplateWrapper) {

	if format == "text" {
		newTmpl := func(name string) *template.Template {
			// return template.New(name).Funcs(template.FuncMap{
			// 	"getenv": getEnv,
			// })
			theTmpl := template.New(name)
			// funcMap := template.FuncMap{}
			funcMap := template.FuncMap(sprig.TxtFuncMap())
			for key, funcValue := range getFuncs() {
				// If not already in map, then add it.  If already in map, ignore it
				if _, exists := funcMap[key]; !exists {
					funcMap[key] = funcValue
				}
			}
			funcMap["include"] = func(name string, data interface{}) (string, error) {
				var buf = bytes.NewBuffer([]byte{})

				err := theTmpl.ExecuteTemplate(buf, name, data)
				return buf.String(), err
			}
			// For text templates, add the "expand" as well
			funcMap["expand"] = func(name string, data interface{}) (result string, err error) {
				var (
					buf          = bytes.NewBuffer([]byte{})
					templateText = fmt.Sprintf(`{{ template "%s" .}}`, name)
					newTmpl2     *template.Template
				)

				if newTmpl2, err = theTmpl.Parse(templateText); err != nil {
					return "", err
				}
				newTmpl2.Execute(buf, data)
				return buf.String(), nil
			}
			return theTmpl.Funcs(funcMap)
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
			// return htmlTemplate.New(name).Funcs(htmlTemplate.FuncMap{
			// 	"getenv": getEnv,
			// })
			theTmpl := htmlTemplate.New(name)
			// funcMap := htmlTemplate.FuncMap{}
			funcMap := htmlTemplate.FuncMap(sprig.HtmlFuncMap())
			for key, funcValue := range getFuncs() {
				// If not already in map, then add it.  If already in map, ignore it
				if _, exists := funcMap[key]; !exists {
					funcMap[key] = funcValue
				}
			}
			funcMap["include"] = func(name string, data interface{}) (string, error) {
				var buf = bytes.NewBuffer([]byte{})

				err := theTmpl.ExecuteTemplate(buf, name, data)
				return buf.String(), err
			}
			return theTmpl.Funcs(funcMap)
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
