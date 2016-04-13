package utils

import (
	"os"
	"path/filepath"
)

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