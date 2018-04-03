package faceserver

import (
	"github.com/shiguanghuxian/face-login/faceserver/v1"
)

// Route 根路由
func (fs *FaceServer) Route() {
	g := fs.e.Group("/v1")
	v1.Route(g)
	// g.OPTIONS("/*", func(c echo.Context) error {
	// 	c.Response().Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// 	c.Response().Header().Add("Access-Control-Allow-Headers", "Authorization,Content-Type")
	// 	return nil
	// })
}
