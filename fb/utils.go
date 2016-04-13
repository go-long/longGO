package fb

import (
	"path/filepath"
	"os"
       "log"
"strings"
	"crypto/md5"
	"encoding/hex"
)

type PublicViewData map[string]interface{}

func (h PublicViewData) Set(k string, v interface{}) PublicViewData {
	if _, ok := h[k]; !ok {
		h[k] = v
	}
	return h
}

// 遍历文件，可指定后缀，返回相对路径
func WalkRelFiles(targpath string, suffixes ...string) (filelist []string) {
	if !filepath.IsAbs(targpath) {
		targpath, _ = filepath.Abs(targpath)
	}
	err := filepath.Walk(targpath, func(retpath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if len(suffixes) == 0 {
			filelist = append(filelist, RelPath(retpath))
			return nil
		}
		_retpath := RelPath(retpath)
		for _, suffix := range suffixes {
			if strings.HasSuffix(_retpath, suffix) {
				filelist = append(filelist, _retpath)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("utils.WalkRelFiles: %v\n", err)
		return
	}

	return
}

// 遍历目录，可指定后缀，返回相对路径
func WalkRelDirs(targpath string, suffixes ...string) (dirlist []string) {
	if !filepath.IsAbs(targpath) {
		targpath, _ = filepath.Abs(targpath)
	}
	err := filepath.Walk(targpath, func(retpath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			return nil
		}
		if len(suffixes) == 0 {
			dirlist = append(dirlist, RelPath(retpath))
			return nil
		}
		_retpath := RelPath(retpath)
		for _, suffix := range suffixes {
			if strings.HasSuffix(_retpath, suffix) {
				dirlist = append(dirlist, _retpath)
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("utils.WalkRelDirs: %v\n", err)
		return
	}

	return
}

// 转相对路径
func RelPath(targpath string) string {
	basepath, _ := filepath.Abs("./")
	rel, _ := filepath.Rel(basepath, targpath)
	return strings.Replace(rel, `\`, `/`, -1)
}

// snake string, XxYy to xx_yy
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func Md5String(str string)string{
	h := md5.New()
	h.Write([]byte(str))
	return string(hex.EncodeToString(h.Sum(nil)))
}