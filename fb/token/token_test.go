package token

import (
	"testing"
	"time"
)

var (
   tokenString string
)

const datastr="oopp88765m,#2"

func TestCreateToken(t *testing.T) {
	token := NewToken(time.Second * 3)
	token.AddData("url",datastr)

	if err := token.IsValid(); err != nil {
		t.Error(err)
		return
	}


	tokenString = token.String()
	t.Log("tokenStr-->"+tokenString+"\n")

}

func TestIsValid(t *testing.T) {
	parseToken, err := ParseToken(tokenString)
	 t.Log(parseToken)
	if err != nil {
		t.Error(err)
		return
	}
	//
	if err := parseToken.IsValid(); err != nil {
		t.Error(err)
		return
	} else {
		t.Log("token is Valid")
	}
	if parseToken.Data["url"]!=datastr{
		t.Error("Data Err:",parseToken.Data["url"],"<>",datastr)
	}
	//
	t.Log(parseToken.Id)
	t.Log(parseToken.String())
	t.Log(parseToken.Expiration)
	t.Log(parseToken.Data)

	time.Sleep(time.Second*5)
	if err := parseToken.IsValid(); err != nil {
		t.Error(err)
		return
	}
}
