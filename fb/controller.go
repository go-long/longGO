// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package fb

import (
//"bytes"
//"io/ioutil"
//"net/http"

	"github.com/go-long/echo"

	"github.com/go-long/longGO/fb/render"
	"net/http"

	"runtime"
	"strings"
	"path/filepath"
"github.com/go-long/longGO/fb/middleware/session"

	"github.com/go-long/longGO/fb/funcMap"
	"github.com/go-long/i18n"
	"path"
	"errors"

)

type (
// 控制器接口
	Controller interface {
		autoInit(ctx echo.Context,routeHandle string) Controller
		Init()
		SetLayout(tmplName... string)
		Layout() string
		Set(key string, v interface{})
		Get(key string) interface{}
		SetPath(path string)
		Path() string
		Object() *BaseController
		SetLanguage(languageTag string)

}
// 基础控制器
	BaseController struct {
		echo.Context    // 请求上下文
		layout   string //模版
		path     string //页面
		Data     map[string]interface{}
		Renderer *render.Renderer
		module   *Module
		Flash    *Flash
		Cookie   *ICookie
		Language *i18n.Language
                T      i18n.TranslateFunc
	}


)

func (this *BaseController) Object() *BaseController {
	return this
}

func (this *BaseController) Init() {

}

// 自动初始化
func (this *BaseController) autoInit(ctx echo.Context,routeHandle string) Controller {

	//if ctx.Sections == nil {
	//	ctx.Sections = map[string]string{}
	//}
	//this.Context.Set("")
	this.Context = ctx
	this.Cookie=&ICookie{ctx}
	this.Flash=&Flash{session.Default(ctx)}
        this.Language,_=this._Language()

	this.Renderer = LongGo.renderer
	this.Renderer.Response = ctx.Response()
	this.Renderer.Req = ctx.Request()
	this.Set("_ROUTEHANDLE",routeHandle)
	//设置Application模块public路径变量
	this.Set("_APP_PUBLIC_",path.Join("/",APP_PACKAGE,PUBLIC_PACKAGE))
	//设置当前模块public路径变量
	this.Set("_PUBLIC_", "/"+strings.ToLower(this.module.Name)+"/public")
        //当前语言
	this.Set("Language",this.Language)


	funcMap.AddFuncMap("Flashes",  this.Flash.Flashes)
	funcMap.AddFuncMap("T",  this.T)


	//语言设置
	q := this.QueryParam("lang")
	if q != "" {
		this.SetLanguage(q)
	}

	return this
}

//设置模版
func (this *BaseController)SetLayout(tmplName... string) {
	if len(tmplName)>0{
		this.layout = tmplName[0]
	}else{
		this.layout = ""
	}
}

//返回模版
func (this *BaseController)Layout() string {
	return this.layout
}

//设置参数
func (this *BaseController)Set(key string, v interface{}) {
	if this.Data == nil {
		this.Data = make(map[string]interface{})
	}
	this.Data[key] = v
}

//获取参数
func (this *BaseController)Get(key string) interface{} {
	if this.Data == nil {
		this.Data = make(map[string]interface{})
	}
	return this.Data[key]
}

func (this *BaseController)SetPath(path string) {
	this.path = path
}

func (this *BaseController)Path() string {
	return this.path
}

func (this *BaseController)_Language() (*i18n.Language,error){

	cookieLang:=this.Cookie.Get("Language")
	if cookieLang==nil{
		cookieLang=""
	}
	acceptLang:=this.Request().Header().Get("Accept-Language")

	translation:= i18n.TranslationMatch(cookieLang.(string), acceptLang)

	if translation==nil{
		e:=errors.New("setLanguage error:"+cookieLang.(string)+","+ acceptLang)
		this.Logger().Debug(e)
		return nil,e
	}
        this.T=translation.Tr
	return &translation.Language,nil
}

func (this *BaseController)SetLanguage(languageTag string){
	this.Cookie.Set("Language",languageTag)
}

func (this *BaseController) Render(status ...int) {
	if len(status) == 0 {
		status = append(status, http.StatusOK)
	}

	//fmt.Println("ttrrrrr:",this.layout)
	////自动补齐布局路径(模块名)
	//if len(this.layout)>0 &&  strings.HasSuffix(this.layout, strings.ToLower(this.module.Name+"/"))==false{
	//	this.layout=strings.ToLower(this.module.Name+"/"+this.layout)
	//	fmt.Println("ttrrrrr:",this.layout)
	//}
	////自动补齐页面路径(模块名)
	//fmt.Println("pppll:",this.path,this.module.Name)
	//if len(this.path)>0 &&  strings.HasSuffix(this.path, strings.ToLower(this.module.Name+"/"))==false{
	//	this.path=strings.ToLower(this.module.Name+"/"+this.path)
	//}


	//无指定页面文件, 自动生成默认页面路径
	//指定了页面文件名称,但没有指定路径,使用自动生成路径
	if len(this.path) == 0 || strings.Index(this.path,"/")<0 {

		var old_file string
		var old_pc uintptr
		i := 0
		for (i < 10) {
			pc, file, _, _ := runtime.Caller(i)
			if strings.HasSuffix(file, ".s") {
				break
			}
			old_file = file
			old_pc = pc
			i++
		}
		if len(old_file) > 0 {
			fileName := strings.TrimSuffix(filepath.Base(old_file), filepath.Ext(old_file))
			f := runtime.FuncForPC(old_pc)
			funName:=this.path
			if len(this.path)== 0{
			   funName = SnakeString(strings.TrimLeft(filepath.Ext(f.Name()), "."))
			}
			this.SetPath(strings.ToLower(filepath.Join(this.module.Name, fileName, funName)))
		}
	}

	this.Renderer.HTML(status[0], this.Path(), this.Data, render.HTMLOptions{Layout:this.layout})
}


//func (this *BaseController) Render(code ...int) error {
//	if len(code) == 0 {
//		code = append(code, http.StatusOK)
//	}

//
//if this.Context.Layout != "" {
//	render := this.Echo().Render
//	for k, v := range this.Context.Sections {
//		if v == "" {
//			this.Set(k, "")
//			continue
//		}
//		sectionBytes := bytes.NewBufferString("")
//		render(sectionBytes, v, this.Context.GetAll())
//		sectionContent, _ := ioutil.ReadAll(sectionBytes)
//		this.Set(k, template.HTML(sectionContent))
//	}
//} else {
//	this.Context.Layout = this.Context.Path()
//}
//return this.Context.Render(code[0], this.Context.Layout, this.Context.GetAll())
//return  this.Context.Render(code[0],this.layout,this.data)
//}

