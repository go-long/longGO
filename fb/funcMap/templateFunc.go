package funcMap

import (
	"html/template"
"github.com/go-long/i18n"
)

var (
	AppFuncMaps    []template.FuncMap
	tplFuncMaps  template.FuncMap
)

func init() {

	tplFuncMaps = make(template.FuncMap)
	//tplFuncMaps["title"]= GetAppTitle
	//tplFuncMaps["getCnfValue"]=  GetCnfValue
	//tplFuncMaps["substr"]= utils.Substr
	//tplFuncMaps["byteUnitStr"]= utils.ByteUnitStr
	//tplFuncMaps["byteUnitStr_uint"]= utils.ByteUnitStr_uint
	//tplFuncMaps["Html2str"]=  Html2str
	tplFuncMaps["Flashes"]=  func()string{return ""}
	tplFuncMaps["T"]=  func()string{return ""}
        tplFuncMaps["LanguageName"]= i18n.LanguageName
	AppFuncMaps=append(AppFuncMaps,tplFuncMaps)

}

func AddFuncMap(key string, funname interface{})   {
	tplFuncMaps[key]=funname
}

//func GetAppTitle() string{
//	return JexHttp.Cnf.AppName
//}

