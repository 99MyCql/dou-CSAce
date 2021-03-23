package confInstance

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"douCSAce/app/confInstance/model"
	"douCSAce/pkg"
)

// @Summary Count
// @Tags Conference Instance
// @Accept json
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/confInstance/count [get]
func Count(c *gin.Context) {
	count, err := model.Count()
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("服务端错误"))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("success", count))
}
