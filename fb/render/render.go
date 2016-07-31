// Package render is a middleware for Martini that provides easy JSON serialization and HTML template rendering.
//
//  package main
//
//  import (
//    "encoding/xml"
//
//    "github.com/go-martini/martini"
//    "github.com/martini-contrib/render"
//  )
//
//  type Greeting struct {
//    XMLName xml.Name `xml:"greeting"`
//    One     string   `xml:"one,attr"`
//    Two     string   `xml:"two,attr"`
//  }
//
//  func main() {
//    m := martini.Classic()
//    m.Use(render.Renderer()) // reads "templates" directory by default
//
//    m.Get("/html", func(r render.Render) {
//      r.HTML(200, "mytemplate", nil)
//    })
//
//    m.Get("/json", func(r render.Render) {
//      r.JSON(200, "hello world")
//    })
//
//    m.Get("/xml", func(r render.Render) {
//      r.XML(200, Greeting{One: "hello", Two: "world"})
//    })
//
//    m.Get("/file/:filename", func(r render.Render, params martini.Params) {
//      r.Download(params[filename])
//    })
//    m.Run()
//  }
package render

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
  "io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/oxtoacart/bpool"

	//"github.com/go-martini/martini"

	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/standard"
)

const (
	ContentType    = "Content-Type"
	ContentLength  = "Content-Length"
	ContentBinary  = "application/octet-stream"
	ContentText    = "text/plain"
	ContentJSON    = "application/json"
	ContentHTML    = "text/html"
	ContentXHTML   = "application/xhtml+xml"
	ContentXML     = "text/xml"
	defaultCharset = "UTF-8"
)

// Provides a temporary buffer to execute templates into and catch errors.
var bufpool *bpool.BufferPool

// Included helper functions for use when rendering html
var helperFuncs = template.FuncMap{
	"partial": func() (string, error) {
	   return "", fmt.Errorf("partial called with no layout defined")
	},
	"yield": func() (string, error) {
		return "", fmt.Errorf("yield called with no layout defined")
	},
	"current": func() (string, error) {
		return "", nil
	},
	"defaultSubLayout": func() (string, error) {
		return "", fmt.Errorf("default called with no sublayout defined")
	},
}

// Render is a service that can be injected into a Martini handler. Render provides functions for easily writing JSON and
// HTML templates out to a http Response.
type Render interface {
	// JSON writes the given status and JSON serialized version of the given value to the http.ResponseWriter.
	JSON(status int, v interface{})
	// HTML renders a html template specified by the name and writes the result and given status to the http.ResponseWriter.
	HTML(status int, name string, v interface{}, htmlOpt ...HTMLOptions)
	// XML writes the given status and XML serialized version of the given value to the http.ResponseWriter.
	XML(status int, v interface{})
	// Data writes the raw byte array to the http.ResponseWriter.
	Data(status int, v []byte)
	// Text writes the given status and plain text to the http.ResponseWriter.
	Text(status int, v string)
	// Error is a convenience function that writes an http status to the http.ResponseWriter.
	Error(status int)
	// Status is an alias for Error (writes an http status to the http.ResponseWriter)
	Status(status int)
	// Redirect is a convienience function that sends an HTTP redirect. If status is omitted, uses 302 (Found)
	Redirect(location string, status ...int)
	// Template returns the internal *template.Template used to render the HTML
	Template() *template.Template
	// Header exposes the header struct from http.ResponseWriter.
	Header() http.Header
	// Download forces response for download file, it prepares the download response header automatically.
	Download(file string, filename ...string)
}

// Delims represents a set of Left and Right delimiters for HTML template rendering
type Delims struct {
	// Left delimiter, defaults to {{
	Left string
	// Right delimiter, defaults to }}
	Right string
}

// Options is a struct for specifying configuration options for the render.Renderer middleware
type Options struct {
	// Directory to load templates. Default is "templates"
	Directory string
	// Layout template name. Will not render a layout if "". Defaults to "".
	Layout string
	// Extensions to parse template files from. Defaults to [".tmpl"]
	Extensions []string
	// Funcs is a slice of FuncMaps to apply to the template upon compilation. This is useful for helper functions. Defaults to [].
	Funcs []template.FuncMap
	// Delims sets the action delimiters to the specified strings in the Delims struct.
	Delims Delims
	// Appends the given charset to the Content-Type header. Default is "UTF-8".
	Charset string
	// Outputs human readable JSON
	IndentJSON bool
	// Outputs human readable XML
	IndentXML bool
	// Prefixes the JSON output with the given bytes.
	PrefixJSON []byte
	// Prefixes the XML output with the given bytes.
	PrefixXML []byte
	// Allows changing of output to XHTML instead of HTML. Default is "text/html"
	HTMLContentType string
	Debug bool
}

// HTMLOptions is a struct for overriding some rendering Options for specific HTML call
type HTMLOptions struct {
	// Layout template name. Overrides Options.Layout.
	Layout string
}

func NewRenderer (options ...Options) *Renderer {
	render:=new(Renderer)
	render.moduleDirs=make(map[string]string)
	render.opt = prepareOptions(options)
	render.compiledCharset = prepareCharset(render.opt.Charset)
	render.t=template.New("")
	template.Must(render.t.Parse("longgo"))
	bufpool = bpool.NewBufferPool(64)
	return render
}
/**
  装载模版指定路径到缓存
 */
func (lr *Renderer) AddLayoutDirectory(moduleName,dir string) {
	lr.moduleDirs[moduleName]=dir
	t := lr.t.New(dir)
	t.Delims(lr.opt.Delims.Left, lr.opt.Delims.Right)
	// parse an initial template in case we don't have any
	//template.Must(t.Parse("Martini"))

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		r, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		ext := getExt(r)

		for _, extension := range lr.opt.Extensions {
			if ext == extension {

				buf, err := ioutil.ReadFile(path)
				if err != nil {
					panic(err)
				}

				name := (r[0 : len(r) - len(ext)])
				tmpl := t.New(moduleName+"/"+filepath.ToSlash(name))

				// add our funcmaps
				for _, funcs := range lr.opt.Funcs {
					tmpl.Funcs(funcs)
				}

				// Bomb out if parse fails. We don't want any silent server starts.
				template.Must(tmpl.Funcs(helperFuncs).Parse(string(buf)))
				break
			}
		}

		return nil
	})

	//return
}

func  (lr *Renderer)PrepareTemplate(){
	if lr.opt.Debug {
		lr.t=template.New("")
		template.Must(lr.t.Parse("longgo"))
		for name,dir:=range lr.moduleDirs{
			lr.AddLayoutDirectory(name,dir)
		}
	}
}

func (lr *Renderer)SetDebug(b bool){
	lr.opt.Debug=b
}
// Renderer is a Middleware that maps a render.Render service into the Martini handler chain. An single variadic render.Options
// struct can be optionally provided to configure HTML rendering. The default directory for templates is "templates" and the default
// file extension is ".tmpl".
//
// If MARTINI_ENV is set to "" or "development" then templates will be recompiled on every request. For more performance, set the
// MARTINI_ENV environment variable to "production"
//func Renderer(options ...Options) martini.Handler {
//	opt := prepareOptions(options)
//	cs := prepareCharset(opt.Charset)
//	t := compile(opt)
//	bufpool = bpool.NewBufferPool(64)
//	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
//		var tc *template.Template
//		if martini.Env == martini.Dev {
//			// recompile for easy development
//			tc = compile(opt)
//		} else {
//			// use a clone of the initial template
//			tc, _ = t.Clone()
//		}
//		c.MapTo(&renderer{res, req, tc, opt, cs}, (*Render)(nil))
//	}
//}

func prepareCharset(charset string) string {
	if len(charset) != 0 {
		return "; charset=" + charset
	}

	return "; charset=" + defaultCharset
}

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults
	if len(opt.Directory) == 0 {
		opt.Directory = "templates"
	}
	if len(opt.Extensions) == 0 {
		opt.Extensions = []string{".tmpl"}
	}
	if len(opt.HTMLContentType) == 0 {
		opt.HTMLContentType = ContentHTML
	}

	return opt
}

func compile(options Options) *template.Template {
	dir := options.Directory
	t := template.New(dir)
	t.Delims(options.Delims.Left, options.Delims.Right)
	// parse an initial template in case we don't have any
	template.Must(t.Parse("longgo"))

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		r, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		ext := getExt(r)

		for _, extension := range options.Extensions {
			if ext == extension {

				buf, err := ioutil.ReadFile(path)
				if err != nil {
					panic(err)
				}

				name := (r[0 : len(r)-len(ext)])
				tmpl := t.New(filepath.ToSlash(name))

				// add our funcmaps
				for _, funcs := range options.Funcs {
					tmpl.Funcs(funcs)
				}

				// Bomb out if parse fails. We don't want any silent server starts.
				template.Must(tmpl.Funcs(helperFuncs).Parse(string(buf)))
				break
			}
		}

		return nil
	})

	return t
}

func getExt(s string) string {
	if strings.Index(s, ".") == -1 {
		return ""
	}
	return "." + strings.Join(strings.Split(s, ".")[1:], ".")
}

type Renderer struct {
	engine.Response
	Req             engine.Request
	t               *template.Template
	opt             Options
	compiledCharset string
	moduleDirs map[string]string
}

func (r *Renderer) JSON(status int, v interface{}) {
	var result []byte
	var err error
	if r.opt.IndentJSON {
		result, err = json.MarshalIndent(v, "", "  ")
	} else {
		result, err = json.Marshal(v)
	}
	if err != nil {
		r.Error(500)
		//http.Error(r.Response., err.Error(), 500)
		return
	}

	// json rendered fine, write out the result
	r.Header().Set(ContentType, ContentJSON+r.compiledCharset)
	r.WriteHeader(status)
	if len(r.opt.PrefixJSON) > 0 {
		r.Write(r.opt.PrefixJSON)
	}
	r.Write(result)
}

func (r *Renderer) HTML(status int, name string, binding interface{}, htmlOpt ...HTMLOptions) {
	r.PrepareTemplate()
	opt := r.prepareHTMLOptions(htmlOpt)
	// assign a layout if there is one
	if len(opt.Layout) > 0 {
		r.addYield(name, binding)
		name = opt.Layout
	}


	buf, err := r.execute(name, binding)
	if err != nil {
		r.Error(http.StatusInternalServerError,err.Error())
		//http.Error(r, err.Error(), http.StatusInternalServerError)
		return
	}

	// template rendered fine, write out the result
	r.Header().Set(ContentType, r.opt.HTMLContentType+r.compiledCharset)
	r.WriteHeader(status)
	io.Copy(r, buf)
	bufpool.Put(buf)
}

func (r *Renderer) XML(status int, v interface{}) {
	var result []byte
	var err error
	if r.opt.IndentXML {
		result, err = xml.MarshalIndent(v, "", "  ")
	} else {
		result, err = xml.Marshal(v)
	}
	if err != nil {
		r.Error(500)
		//http.Error(r, err.Error(), 500)
		return
	}

	// XML rendered fine, write out the result
	r.Header().Set(ContentType, ContentXML+r.compiledCharset)
	r.WriteHeader(status)
	if len(r.opt.PrefixXML) > 0 {
		r.Write(r.opt.PrefixXML)
	}
	r.Write(result)
}

func (r *Renderer) Data(status int, v []byte) {
	if r.Header().Get(ContentType) == "" {
		r.Header().Set(ContentType, ContentBinary)
	}
	r.WriteHeader(status)
	r.Write(v)
}

func (r *Renderer) Download(file string,filename ...string ) {
	r.Header().Set("Content-Description","File Transfer")
	r.Header().Set("Content-Type","application/octet-stream")
	if len(filename) > 0 && filename[0] != "" {
		r.Header().Set("Content-Disposition", "attachment; filename="+filename[0])
	} else {
		r.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(file))
	}
	r.Header().Set("Content-Transfer-Encoding", "binary")
	r.Header().Set("Expires", "0")
	r.Header().Set("Cache-Control", "must-revalidate")
	r.Header().Set("Pragma", "public")

 	http.ServeFile(r.Response.(* standard.Response).ResponseWriter,r.Req.(*standard.Request).Request,file)


}

func (r *Renderer) Text(status int, v string) {
	if r.Header().Get(ContentType) == "" {
		r.Header().Set(ContentType, ContentText+r.compiledCharset)
	}
	r.WriteHeader(status)
	r.Write([]byte(v))
}

// Error writes the given HTTP status to the current ResponseWriter
func (r *Renderer) Error(status int,msg ...string) {
	r.WriteHeader(status)
	if len(msg)>0 {
		r.Write([]byte(msg[0]))
      }
}

func (r *Renderer) Status(status int) {
	r.WriteHeader(status)
}

func (r *Renderer) Redirect(location string, status ...int) {
	code := http.StatusFound
	if len(status) == 1 {
		code = status[0]
	}
       // r.Redirect(location, code)

	http.Redirect(r.Response.(*standard.Response).ResponseWriter, r.Req.(*standard.Request).Request, location, code)

}

func (r *Renderer) Template() *template.Template {
	return r.t
}

func (r *Renderer) execute(name string, binding interface{}) (*bytes.Buffer, error) {
	buf := bufpool.Get()
	return buf, r.t.ExecuteTemplate(buf, name, binding)
}



func (r *Renderer) addYield(name string, binding interface{}) {
	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			buf, err := r.execute(name, binding)
			// return safe html here since we are rendering our own template
			return template.HTML(buf.String()), err
		},
		"current": func() (string, error) {
			return name, nil
		},
		"defaultSubLayout": func()(template.HTML, error) {
			ls:=strings.Split(name,"/")
			ls[len(ls)-1]="default"
			defaultSubLayout:=strings.Join(ls,"/")

			if r.t.Lookup(defaultSubLayout)!=nil{
				buf, err := r.execute(defaultSubLayout, binding)
				return template.HTML(buf.String()), err
			}else{
				buf, err := r.execute(name, binding)
				return template.HTML(buf.String()), err
			}
		},

	}
	r.t.Funcs(funcs)
}

func (r *Renderer) prepareHTMLOptions(htmlOpt []HTMLOptions) HTMLOptions {
	if len(htmlOpt) > 0 {
		return htmlOpt[0]
	}

	return HTMLOptions{
		Layout: r.opt.Layout,
	}
}
