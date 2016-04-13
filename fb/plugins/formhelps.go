package plugins

import (
	"jex/cn/longGo/fb"
	"reflect"
	"github.com/coscms/forms"
	"github.com/coscms/forms/common"
)

type (
	FormEditorConfig struct  {
	  Controller fb.Controller
	  DBModel interface{}
	  ID int
	  BeforeValidData FuncHandle //验证数据前
        }

     FuncHandle  func (*FormEditorConfig)
)

func LongFormEditor(conf *FormEditorConfig) *forms.Form{
	controller:=conf.Controller.Object()

	if conf.ID>0 {
		fb.DB.Find(conf.DBModel, "id = ?", conf.ID)
	}
	if controller.Request().Method()=="POST" {
		if err := controller.Bind(conf.DBModel); err != nil {
			fb.Log.Error(err)
			controller.Flash.Error("数据错误"+err.Error())
		}
		if conf.ID>0{
			reflect.ValueOf(conf.DBModel).Elem().FieldByName("ID").SetUint(uint64(conf.ID))

		}

	}

	form := forms.NewFormFromModel(conf.DBModel, formcommon.BOOTSTRAP, "POST",controller.Request().URI())
	form.SetId("agentSvrEditForm")
	form.Field("submit").SetText("提交更改")
	form.Field("reset").SetText("关闭")
	controller.Data["form"] = form

	if controller.Request().Method() == "POST" {
		if conf.BeforeValidData!=nil{
			conf.BeforeValidData(conf)
		}
		valid, passed := form.Valid()
		if !passed {
			// validation does not pass
			controller.Data["valid"] = valid;
			fb.Log.Debug(valid)
			controller.Flash.Error("提交失败")
		} else {
			//			fmt.Println("tmp_user:",tmp_user)
			var err error
			if conf.ID>0{
				err= fb.DB.Save(conf.DBModel).Error
			}else{
				err=fb.DB.Create(conf.DBModel).Error
			}
			if err!=nil{
				fb.Log.Error(err)
				controller.Flash.Error("保存失败："+err.Error())
			}else{
				controller.Flash.Success("修改成功")
			}
		}
	}

	return form
}