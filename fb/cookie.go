package fb

import (
"github.com/labstack/echo"
"jex/cn/longGo/fb/middleware/session"
	"fmt"
)

type ICookie struct {
	ctx echo.Context
}

func GetCookie(c echo.Context)*ICookie{
	return &ICookie{c}
}

func (this *ICookie)SetCookie(key interface{}, val interface{}){
	session := session.Default(this.ctx)
	session.Set(key,val)
	session.Save()
}

func (this *ICookie)Cookie(key interface{}) interface{}{
	session := session.Default(this.ctx)
	return session.Get(key)
}

func (this *ICookie)DelCookie(key interface{}){
	fmt.Println("ctx:",this.ctx)
	session := session.Default(this.ctx)
        session.Delete(key)

}

func (this *ICookie)ClearCookie(){
	session := session.Default(this.ctx)
	session.Clear()
}