package middle

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/utils"
	"net/http"
)

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对login接口放行
		if len(c.Request.URL.String()) >= 10 && c.Request.URL.String()[0:10] == "/api/login" {
			c.Next()
		} else {
			//获取Header中的Authorization
			token := c.Request.Header.Get("Authorization")
			if token == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"msg":  "请求未携带token，无权限访问",
					"data": nil,
				})
				c.Abort()
				return
			}

			// parseToken 解析token包含的信息
			claims, err := utils.JWTToken.ParseToken(token)
			if err != nil {
				// token延期错误
				if err.Error() == "TokenExpired" {
					c.JSON(http.StatusBadRequest, gin.H{"msg": "授权已过期", "data": nil})
					c.Abort()
					return
				}
				//其他解析错误
				c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error(), "data": nil})
				c.Abort()
				return
			}

			// 继续交给下一个路由处理,并将解析出的信息传递下去
			c.Set("claims", claims)
			c.Next()
		}
	}
}
