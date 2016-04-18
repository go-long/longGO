package token

import (
	"time"
	//"encoding/json"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
"io"
	"encoding/base64"
	"crypto/md5"
	"encoding/hex"

	"strings"
	"fmt"
)

type Token struct {
	Id         UUID
	Data       map[string]string
	Expiration int64
}

type TokenError struct {
	Message string
	Details error
}

const HeaderTAG="JusTin"


//新建token
// duration 有效时长(小时)
//var hash = crypto.SHA256

//例子: NewToken(time.Hour * 72)
func NewToken(duration time.Duration) (*Token) {
	return &Token{
		Id:NewUUID(10),
		Data:make(map[string]string),
		Expiration: time.Now().Add(duration).Unix(),
	}
}

func ParseToken(tokenString string) (*Token, *TokenError) {
	if len(tokenString)<16 {
		return nil, &TokenError{"Token is not complete", nil}
	}

	ps:=strings.Split(tokenString,".")
	if len(ps)<2 {
		return nil, &TokenError{"Token Was a fake", nil}
	}

	//获取时间戳
	expirstr := ps[1]
	expiration,err:=Decode(expirstr)
	if err != nil {
		return nil, &TokenError{"Token is not complete", err}
	}

	//获取uuid
	id,err:= Decode(ps[0][:10])
	if err!=nil{
		return nil, &TokenError{"Token Was a fake", nil}
	}
	datastr:=ps[0][10:]

	key:=GetMD5Hash(ps[1])
        data:=DecryptString(key,datastr)
	fmt.Println("data",data)
	ls:=strings.Split(data,";")
	fmt.Println("ls:",ls)
	if len(ls)==0 || ls[0]!=HeaderTAG{
		return nil, &TokenError{"Token Was a fake", nil}
	}


	tk:= &Token{
		Id: id,
		Data: make(map[string]string),
		Expiration:int64(expiration),
	}
	for n,l:=range ls {
		if n>0 {
			s := strings.Split(l, ":")
			tk.Data[s[0]]=s[1]
		}
	}

        if err:= tk.IsValid();err!=nil{
		return nil,err
	}

	return tk, nil
}

func (t *Token) AddData(key string,value string) {
	  t.Data[key]=value

}


func (t *Token) String() string {
	ext:=UUID(t.Expiration).Encode()
	key:=GetMD5Hash(ext)
	strs:=HeaderTAG+";"
	for k,v:=range t.Data{
		strs=strs+k+":"+v+";"
	}
	strs=strings.TrimSuffix(strs,";")
	//bytes,_:=json.Marshal(t)

	mstr:= EncryptString(key,strs)
	return t.Id.Encode()+mstr+"."+ext
}

func (t *Token) IsValid() *TokenError {
	now := time.Now().Unix()

	if t.Expiration < now {
		return &TokenError{"Token is expired!", nil}
	}
	return nil
}

//字符串加密/解密

var iv = []byte{34, 35, 35, 57, 68, 4, 35, 36, 7, 8, 35, 23, 35, 86, 35, 23}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func EncryptString(key string, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	fmt.Println("eny:", ciphertext)
	// convert to base64
	b4data:=base64.URLEncoding.EncodeToString(ciphertext)
	return  strings.Replace(b4data,"=", "~",-1)
}

// decrypt from base64 to decrypted string
func DecryptString(key string, cryptoText string) string {
	cryptoText= strings.Replace(cryptoText,"~", "=",-1)
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)
	fmt.Println("dec:", ciphertext)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
