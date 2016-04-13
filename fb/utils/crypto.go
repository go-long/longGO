package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"
	"crypto/rand"
)

/**
 * 对一个字符串进行MD5加密,不可解密
 */
func MD5Encode(s string) string {
	h := md5.New()
	h.Write([]byte(s)) //使用zhifeiya名字做散列值，设定后不要变
	return hex.EncodeToString(h.Sum(nil))
}

func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func Base64Decode(s string) string {
	bstr, _:=base64.StdEncoding.DecodeString(s)
	return string(bstr)
}

/**
 * 获取一个Guid值
 */
func GetGuid() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return MD5Encode(base64.URLEncoding.EncodeToString(b))
}


//xxtea 加密
func XXEncrypt(str, key string) string {
	encrypt_data:= Encrypt([]byte(str), []byte(key))
	b4data:= base64.StdEncoding.EncodeToString(encrypt_data)
	b4data= strings.Replace(b4data,"+", "-",-1)
	b4data= strings.Replace(b4data,"/", "_",-1)
	b4data= strings.Replace(b4data,"=", "~",-1)
	return b4data
}

//xxtea 解密
func XXDecrypt(str, key string) (string,error) {
	str= strings.Replace(str,"-", "+",-1)
	str= strings.Replace(str,"_", "/",-1)
	str= strings.Replace(str,"~", "=",-1)
	encrypt_data,err:=base64.StdEncoding.DecodeString(str)
	if err!=nil{
		return "",err
	}
	return  string(Decrypt(encrypt_data, []byte(key))),nil
}

