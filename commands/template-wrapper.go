package commands

import (
	"fmt"
	htmlTemplate "html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

type templateWrapper struct {
	tmpl             *template.Template
	htmlTmpl         *htmlTemplate.Template
	parseFilesFn     func(filenames ...string) error
	parseGlobFn      func(pattern string) error
	executeFn        func(wr io.Writer, data interface{}) error
	createTextTmplFn func(templateText string) error
}

func (w templateWrapper) newTextTemplate(templateText string) error {
	return w.createTextTmplFn(templateText)
}
func (w templateWrapper) parseFiles(filenames ...string) error {
	return w.parseFilesFn(filenames...)
}
func (w templateWrapper) parseGlob(pattern string) error {
	return w.parseGlobFn(pattern)
}
func (w templateWrapper) execute(wr io.Writer, data interface{}) error {
	return w.executeFn(wr, data)
}

func newTemplateWrapper(format string) (w templateWrapper) {

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

func fileOrDirectoryExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func isDirectory(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// getEnv - utility function available in templates as "getenv"
func getEnv(key string) string {
	value, found := os.LookupEnv(key)

	if found {
		return value
	}
	return ""
}

// getFirstMatchedFile - from the given pattern, it turns the filename (without dir) of the first matching file
func getFirstMatchedFile(pattern string) (string, error) {
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(filenames) == 0 {
		return "", fmt.Errorf("No files matched for pattern: %s", pattern)
	}

	_, file := path.Split(filenames[0])
	return file, nil
}
