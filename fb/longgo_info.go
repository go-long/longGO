// Copyright 2016 henrylee2cn.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package fb

// 框架信息
const (
	NAME    = "LongGo"
	VERSION = "Ver 0.0.1"
	AUTHOR  = "jex.h"
)

func LongGoInfo() string {
	var s string
	s += "[THINKGO-INFO] @NAME:		" + NAME + "\n"
	s += "[THINKGO-INFO] @VERSION:	" + VERSION + "\n"
	s += "[THINKGO-INFO] @AUTHOR:		" + AUTHOR + "\n"
	return s
}
