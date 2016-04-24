package middleware

import (
	"net"
	"time"
"strings"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/color"
	"github.com/labstack/echo/engine/standard"
)

var LogBlacklist =&TBlacklist{}

type (
	TBlacklist []blackNote

	blackNote struct {
	  Method string
          Prefix string
        }
)

func Log() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return  func(c echo.Context) error {

			req := c.Request()
			res := c.Response()

			logger := c.Logger()
			method := req.Method()
			path := req.(*standard.Request).RequestURI
			if path == "" {
				path = "/"
			}
                        //黑名单不打印日志
			if LogBlacklist.Exists(method,path) {
				//return nil
				return h(c)
			}

			remoteAddr := req.RemoteAddress()
			if ip := req.Header().Get(echo.HeaderXRealIP); ip != "" {
				remoteAddr = ip
			} else if ip = req.Header().Get(echo.HeaderXForwardedFor); ip != "" {
				remoteAddr = ip
			} else {
				remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
			}

			start := time.Now()
			if err := h(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()


			size := res.Size()

			n := res.Status()
			code := color.Green(n)
			switch {
			case n >= 500:
				code = color.Red(n)
			case n >= 400:
				code = color.Yellow(n)
			case n >= 300:
				code = color.Cyan(n)
			}

			logger.Infof("%s|%s %s %s %s %s %d",stop.Format("15:04:05.999"), remoteAddr, method, path, code, stop.Sub(start), size)
			return nil
		}
	}
}

func (b *TBlacklist)Add(method,url string){
	*b=append(*b,blackNote{method,url})
}

func (b *TBlacklist)Exists(method,url string) bool{
	for _,v:=range *b{
		if v.Method==method && strings.HasPrefix(url, v.Prefix){
			return true
		}
	}
	return false
}