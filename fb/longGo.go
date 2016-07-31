package fb

import (
	"fmt"
	"path"
//"regexp"
	"path/filepath"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo/engine/standard"
	"github.com/go-long/longGO/fb/render"

	"github.com/go-long/longGO/fb/db"
	"github.com/labstack/gommon/color"
	mdw "github.com/go-long/longGO/fb/middleware"
	"github.com/go-long/longGO/fb/middleware/session"
	"net/http"
	"github.com/go-long/i18n"
	"github.com/go-long/longGO/fb/funcMap"
	"github.com/labstack/gommon/log"

)

type Long struct {
	Echo     *echo.Echo
	// 模块列表
	Modules  map[string]*Module
	// 插件列表
	// Addons *Addons
	// 模板引擎
	//*Template
	// 配置信息
	Config
	// 框架信息
	Author   string
	Version  string
	renderer *render.Renderer
	sessionStore session.Store
	FirstMiddleware  []echo.MiddlewareFunc
        LastMiddleware []echo.MiddlewareFunc
}

// 重要配置，涉及项目架构，请勿修改
const (
// 模块应用目录名
	MODULES_PACKAGE = "modules"
// 视图文件目录名
	VIEW_PACKAGE = "views"
// 公共目录
	APP_PACKAGE = "application"
// 主题目录名
	THEME_PACKAGE = "themes"
// 资源文件目录名
	PUBLIC_PACKAGE = "public"
// 资源文件静态目录名
	STATIC_PACKAGE = "static"
// 上传根目录名
	UPLOADS_PACKAGE = "uploads"


//运行环境
	ENV_DEVEL = "development"
	ENV_production = "production"
	ENV_TEST = "test"
)

// 全局运行实例
var (
	LongGo = newLongGo()
	DB = new(db.LongDB)
	Log = log.New("LGO")
	db_Tables  []interface{}
)

func newLongGo() *Long{

	t := &Long{
		// 业务数据
		Echo:    echo.New(),
		Modules: Modules,
		// Addons:  newAddons(),
		Config: getConfig(),
		// 框架信息
		Author:  AUTHOR,
		Version: VERSION,

	}



	//t.Echo.Blackfile(".html")

	t.Echo.SetDebug(t.Config.Env == ENV_DEVEL)

	t.renderer = render.NewRenderer(render.Options{
		Funcs:funcMap.AppFuncMaps,
	})
	t.renderer.SetDebug(true)
	//t.Echo.SetRenderer(t.render)
	//t.htmlPrepare()

	//t.Hook()
	// t.Echo.SetBinder(b)
	// t.Echo.SetHTTPErrorHandler(HTTPErrorHandler)
	// t.Echo.SetLogOutput(w io.Writer)
	// t.Echo.SetHTTPErrorHandler(h HTTPErrorHandler)

	Log.SetLevel(t.Config.LogLevel)
	Log.EnableColor()
	Log.SetHeader("${prefix}|${level}|${message}\n")
        t.Echo.SetLogger(Log)
	t.Echo.SetHTTPErrorHandler(t.httpErrorHandler)

	t.Echo.SetBinder(&binder{})
	i18n.Init("i18n","zh-cn");
	funcMap.AddFuncMap("Languages",  i18n.Translations)
	//funcMap.AddFuncMap("LanguageTemplate",i18n.LanguageTemplate)
	return t
}



func (this *Long)SessionStore(store session.Store){
	this.sessionStore=store
}

func (this *Long) Run(onBeforeRuning func(),mw... echo.MiddlewareFunc ) {

	DB.Init(this.DBConfig)
	//this.DB.DB.SetLogger(this.log)
	DB.LogMode(this.Echo.Debug())
	DB.AutoMigrate()

	//使用session中间件
	if this.sessionStore==nil {
		this.sessionStore = session.NewCookieStore([]byte("secret-key"))
	}
	//默认中间件
	var tmw []echo.MiddlewareFunc
	tmw=append(tmw, mdw.Log(), session.Sessions("ESESSION",  this.sessionStore))
	this.FirstMiddleware= append(tmw,this.FirstMiddleware...)
	this.FirstMiddleware= append(tmw,mw...)


	//this.Echo.Use()
	this.LastMiddleware=append(this.LastMiddleware,middleware.Recover())

	//根据模块设置静态路由
	this.setModuesStaticRoutes()


	//查看模版路径列表
	for _, t := range this.renderer.Template().Templates() {
		Log.Infof("%s%s", color.Blue("template|"), color.Green(t.Name()))
	}

	//查看路由列表
	Log.Infof(color.CyanBg(color.Yellow("Routes:")))
	for _, r := range this.Echo.Routes() {
		Log.Infof("%s %s %s", color.Blue(r.Method), color.Green(r.Path), r.Handler)
	}

	//启动服务
	serverEng := standard.New(fmt.Sprintf("%s:%d", this.Config.HttpAddr, this.Config.HttpPort))
	Log.Debugf("%s %s:%d", color.Yellow("Start httpService"), color.Green(this.Config.HttpAddr), this.Config.HttpPort)


	if onBeforeRuning!=nil {
		onBeforeRuning()
	}
	this.Echo.Run(serverEng)
}


// 定义中间件
func (this *Long) FirstUse(m ...echo.MiddlewareFunc)  {
   this.FirstMiddleware = append(this.FirstMiddleware, m...)

}


//根据模块设置静态路由
func (this *Long)setModuesStaticRoutes() {

	this.Echo.Static("/favicon",filepath.Join(MODULES_PACKAGE, APP_PACKAGE, PUBLIC_PACKAGE, "favicon"))
	//this.Echo.Favicon(filepath.Join(MODULES_PACKAGE, APP_PACKAGE, PUBLIC_PACKAGE, "favicon", "favicon.ico"))
	this.Echo.Static("/uploads", UPLOADS_PACKAGE)
	this.Echo.Static("/static", "static")
	if this.Config.Env == ENV_DEVEL {
		this.Echo.Static("/swagger", "swagger")
	}



	for name, m := range Modules {
		this.Echo.Static(path.Join("/", name, "public"), filepath.Join(MODULES_PACKAGE, name, PUBLIC_PACKAGE, THEME_PACKAGE, m.Themes.cur))
		this.Echo.Static(path.Join("/", name, "static"), filepath.Join(MODULES_PACKAGE, name, PUBLIC_PACKAGE, STATIC_PACKAGE))
		//layout
	 	this.renderer.AddLayoutDirectory(name, filepath.Join(MODULES_PACKAGE, name, PUBLIC_PACKAGE, VIEW_PACKAGE))
	}
}

func (this *Long) Group(prefix string, m ...echo.MiddlewareFunc) *echo.Group {
	return this.Echo.Group(prefix, m...)
}


func (this *Long) httpErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := http.StatusText(code)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	}

	//debug模式显示错误信息
	if this.Echo.Debug() {
		msg = err.Error()
	}
	this.Echo.Logger().Debugf("%s %s %s %s", color.Blue(c.Request().Method()), color.Green(c.Request().URI()), color.Red(err),  color.Yellow(code))
	switch code {
	case http.StatusNotFound:
	     //404
		//fmt.Println("url:",c.Request().URI())
		 c.File(filepath.Join(MODULES_PACKAGE,APP_PACKAGE, PUBLIC_PACKAGE,VIEW_PACKAGE,"404.html"))
	default:
		if !c.Response().Committed() {
			c.String(code, msg)
		}
	}

	//this.Echo.Logger().Debug(err)

}


//func (this *Long) Hook() {
//	this.Echo.Hook(func(r engine.Request, w engine.Response) {
//		//fs := this.Echo.fileSystem.path
//		//if fs != "" && strings.HasPrefix(r.URL.Path, fs) {
//		//	return
//		//}
//		//p := strings.Trim(r.URL.Path, "/")
//		fmt.Println("rrr:",r.URL().Path())
//	})
//}