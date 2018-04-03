package context

import (
	"github.com/labstack/echo"
)

// ZContext 自定义中间件扩展
type ZContext struct {
	echo.Context
}

// InitZContext 注册中间件
func InitZContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("X-Powered-By", "PHP/7.1.15")
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			cc := &ZContext{c}
			return next(cc)
		}
	}
}
