// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package fb

import (
//"path"
	"path/filepath"
//"reflect"
	"regexp"
	"runtime"
	"strings"
	"strconv"
	"github.com/labstack/echo"
//	"github.com/labstack/echo/middleware"
	"reflect"
	"path"
)

type(
// 应用模块
	Module struct {
		id          string
		Name        string
		Description string
		*Themes
		//*echo.Group
		mw          []echo.MiddlewareFunc
		Routes      []LongRoute
	}

	LongRoute struct {
		Method  string
		Path    string
		Handler string
		Tags    string
	}
)

var (
	Modules = map[string]*Module{}
	re = regexp.MustCompile("^[/]?([a-zA-Z0-9_]+)([\\./\\?])?")
)

const (
	PERMISSION_PUBLIC = iota
	PERMISSION_USER
	PERMISSION_ADMIN
)

func (lr *LongRoute)Desc()string{
	return lr.Tag("desc")
}

func (lr *LongRoute)Tag(key string)string{
	tag:=lr.Tags
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}


		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		if key == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value
		}
	}
	return ""
}



// 创建模块
// 自动设置default主题
// 文件名作为id，且文件名应与模块目录名、包名保存一致
func NewModule(description string, mw ...echo.MiddlewareFunc) *Module {
	m := &Module{
		Description: description,
		Themes:      &Themes{},
	}

	// 设置默认主题
	m.Themes.Set(&Theme{
		Name:        "default",
		Description: "default",
		Src:         map[string]string{},
	})

	// 设置id
	_, file, _, _ := runtime.Caller(1)
	m.id = strings.TrimSuffix(filepath.Base(file), ".go")

	// 设置Name
	m.Name = strings.Title(m.id)

	m.Use(mw...)

	//// 生成url前缀
	//prefix := "/" + m.id
	//
	//// 创建分组并修改请求路径c.path "/[模块]/[控制器]/[操作]"为"/[模块]/[主题]/[控制器]/[操作]"
	//m.Group = LongGo.Echo.Group(
	//	prefix,
	//
	//		  func(h echo.HandlerFunc) echo.HandlerFunc {
	//			return func(c echo.Context) error {
	//				fmt.Println("mmmmmmmmmmmmmm")
	//				// 补全主题字段
	//				//p := strings.Split(c.Path(), "/:")[0]
	//			//p = path.Join(prefix, m.Themes.Cur().Name, strings.TrimPrefix(p, prefix))
	//				//c.SetPath(p)
	//				// 静态文件前缀
	//				//c.Set("__PUBLIC__", path.Join("/public", prefix, m.Themes.Cur().Name))
	//				return nil
	//			}
	//
	//	},
	//	middleware.Recover(),
	//	middleware.Log(),
	//)
	//m.Group.Use(mw...)

	// 模块登记
	Modules[m.id] = m

	return m
}

// 获取Id
func (this *Module) GetId() string {
	return this.id
}

// 设置Id
func (this *Module) SetId(id string) *Module {
	this.id = id
	return this
}

// 获取Name
func (this *Module) GetName() string {
	return this.Name
}

// 获取Description
func (this *Module) GetDescription() string {
	return this.Description
}

// 设置主题，并默认设置传入的第1个主题为当前主题
func (this *Module) SetThemes(themes ...*Theme) *Module {
	this.Themes.Set(themes...)
	return this
}

// 设置当前主题
func (this *Module) UseTheme(name string) *Module {
	this.Themes.Use(name)
	return this
}

// 定义中间件
func (this *Module) Use(m ...echo.MiddlewareFunc) *Module {
	//this.Group.Use(m...)
	for _, h := range m {
		this.mw = append(this.mw, h)
	}
	return this
}

func (this *Module)addRoute(c Controller, route echo.Route, m ...echo.MiddlewareFunc) *Module {
	t := reflect.TypeOf(c)
	e := t.Elem()

	//cname := SnakeString(strings.TrimSuffix(e.Name(), "Controller"))
	//group := this.Group.Group(cname, m...)
	//group.Match(strings.Split( route.Method,"|"),route.Path,echo.HandlerFunc(func(ctx echo.Context) error {

	//func(next Handler) Handler {
	//	return HandlerFunc(func(c Context) error {
	//		if err := h.Handle(c); err != nil {
	//			return err
	//		}
	//		return next.Handle(c)
	//	})


	h := echo.HandlerFunc(func(ctx echo.Context) error {

			/////////
			var v = reflect.New(e)
			//控制器默认layout布局文件
			v.Interface().(Controller).SetLayout(strings.ToLower(this.Name + "/layouts/default"))
			v.Interface().(Controller).Object().module = this

			//初始化Context
			v.Interface().(Controller).autoInit(ctx,route.Handler)
		        v.Interface().(Controller).Init()

			//运行页面处理函数
			rets := v.MethodByName(route.Handler).Call([]reflect.Value{})
			if len(rets) > 0 {
				if err, ok := rets[0].Interface().(error); ok {
					return err
				}
			}
			return nil
		})

	LongGo.Echo.Add(route.Method, route.Path,echo.HandlerFunc(func(ctx echo.Context) error {
			//调用框架默认中间件(首部)
			for i := len(LongGo.FirstMiddleware) - 1; i >= 0; i-- {
				h = LongGo.FirstMiddleware[i](h)
			}
			//调用模块中间件
			for i := len(this.mw) - 1; i >= 0; i-- {
				h = this.mw[i](h)
			}
			//调用路由指定中间件
			// Chain middleware with handler in the end
			for i := len(m) - 1; i >= 0; i-- {
				h = m[i](h)
			}

			//调用框架默认中间件(尾部)
			for i := len(LongGo.LastMiddleware) - 1; i >= 0; i-- {
				h = LongGo.LastMiddleware[i](h)
			}

			// Execute chain
			if err := h(ctx); err != nil {
				LongGo.Echo.DefaultHTTPErrorHandler(err, ctx)
			}
			return nil
		}))
	return this
}


// 翻译路由
func (this *Module) Router(rootpath string,c Controller,routes []LongRoute, m ...echo.MiddlewareFunc) *Module {
	if rootpath==LongGo.Config.DefaultModule{
		rootpath=""
	}
	for _, route := range routes {
		methods := strings.Split(route.Method, "|")
		for _, method := range methods {
			this.addRoute(c, echo.Route{method,path.Join("/",rootpath,route.Path), route.Handler}, m...)
		}
		//this.addRoute(c, route,m...)
	}
        this.Routes=append(this.Routes,routes...)
	return this
}

//func (this *Module) Router(c  Controller, m ...echo.Middleware) *Module {
//t := reflect.TypeOf(c)
//e := t.Elem()
//cname := SnakeString(strings.TrimSuffix(e.Name(), "Controller"))
//group := this.Group.Group(cname, m...)
//for i := t.NumMethod() - 1; i >= 0; i-- {
//	fname := t.Method(i).Name
//	idx := strings.LastIndex(fname, "_")
//	if idx == -1 {
//		continue
//	}
//	pattern := SnakeString(fname[:idx])
//	method := strings.ToUpper(fname[idx+1:])
//	switch method {
//	case "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT", "TRACE", "SOCKET":
//		group.Match([]string{method}, pattern, func(ctx *echo.Context) error {
//			var v = reflect.New(e)
//			v.Interface().(Controller).AutoInit(ctx)
//			rets := v.MethodByName(fname).Call([]reflect.Value{})
//			if len(rets) > 0 {
//				if err, ok := rets[0].Interface().(error); ok {
//					return err
//				}
//			}
//			return nil
//		})
//	}
//}
//	return this
//}



