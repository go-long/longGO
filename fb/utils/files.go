package utils

import (
	"os"
	"path/filepath"
	"strings"
	"crypto/md5"
	"io"
	"encoding/hex"
)

//获取程序路径
func GetAppDirectory(joinPath ...string)string{
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	path= strings.Replace(path, "\\", "/", -1)

	if len(joinPath)>0 {
		var paths []string
		paths=append(paths,path)
		paths=append(paths,joinPath...)
		path=filepath.Join( paths...)
	}
	return path
}

//获取工作路径
func GetWokingDirectory(joinPath ...string) string{
	dir,_:=os.Getwd()
	if len(joinPath)>0 {
		var dirs []string
		dirs=append(dirs,dir)
		dirs=append(dirs,joinPath...)
		dir=filepath.Join( dirs...)
	}
	return dir
}

//获取指定路径文件列表
func GetDirectoryFiles(root string) (files []os.FileInfo, err error){
	err=filepath.Walk(root,
		func(path string,f os.FileInfo, err error) error {
			if (f == nil) {
				return err
			}
			if f.IsDir() {
				return nil
			}
			files=append(files,f)
//			println(path)
			return nil
		})
	return files,err
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//获取文件MD5
func FileMD5(filename string)(string,error ) {
	file, inerr := os.Open(filename)
	if inerr == nil {
		md5h := md5.New()
		io.Copy(md5h, file)
		return hex.EncodeToString(md5h.Sum(nil)) ,nil //md5
	}
	return "",inerr
}