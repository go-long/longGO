package fb

import (
	"jex/cn/longGo/fb/middleware/session"
	"encoding/gob"
)



type (
	Flash struct {
		session session.Session
	}

	FlashMessage struct {
		Type    string
		Message string
	}
)

func (this *Flash)Success(value string) {
	data := &FlashMessage{Type:"Success", Message:value}
	this.session.AddFlash(data)
	this.session.Save()
}

func (this *Flash)Error(value string) {
	data := &FlashMessage{Type:"Error", Message:value}
	this.session.AddFlash(data)
	this.session.Save()
}

func (this *Flash)Warning(value string) {
	data := &FlashMessage{Type:"Warning", Message:value}
	this.session.AddFlash(data)
	this.session.Save()
}

func (this *Flash)Info(value string) {
	data := &FlashMessage{Type:"Info", Message:value}
	this.session.AddFlash(data)
	this.session.Save()
}

func (this *Flash)Flashes() []FlashMessage {
	var datas []FlashMessage
	for _, v := range this.session.Flashes() {
		datas = append(datas, v.(FlashMessage))
	}
	this.session.Save()
	return datas
}

func (f *FlashMessage)IsSuccess() bool {
	return f.Type == "Success"
}

func (f *FlashMessage)IsError() bool {
	return f.Type == "Error"
}

func (f *FlashMessage)IsWarning() bool {
	return f.Type == "Warning"
}

func (f *FlashMessage)IsInfo() bool {
	return f.Type == "Info"
}

func init() {
	gob.Register(FlashMessage{})
}