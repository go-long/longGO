package fb

import (
"github.com/go-long/echo"
"github.com/go-long/longGO/fb/middleware/session"
	"fmt"
)

type ICookie struct {
	ctx echo.Context
}

func GetCookie(c echo.Context)*ICookie{
	return &ICookie{c}
}

func (this *ICookie)Set(key interface{}, val interface{}){
	session := session.Default(this.ctx)
	session.Set(key,val)
	session.Save()
}

func (this *ICookie)Get(key interface{}) interface{}{
	session := session.Default(this.ctx)
	return session.Get(key)
}

func (this *ICookie)Del(key interface{}){
	fmt.Println("ctx:",this.ctx)
	session := session.Default(this.ctx)
        session.Delete(key)
	session.Save()
}

func (this *ICookie)Clear(){
	session := session.Default(this.ctx)
	session.Clear()
	session.Save()
}