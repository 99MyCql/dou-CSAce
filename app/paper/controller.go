package paper

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"douCSAce/app/paper/model"
	"douCSAce/pkg"
)

// @Summary Count
// @Tags Paper
// @Accept json
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/paper/count [get]
func Count(c *gin.Context) {
	count, err := model.Count()
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("服务端错误"))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("success", count))
}

// @Summary Get
// @Tags Paper
// @Accept json
// @Param key query string true "唯一标识"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/paper/get [get]
func Get(c *gin.Context) {
	var key string
	if key = c.DefaultQuery("key", ""); key == "" {
		pkg.Log.Error("need key")
		c.JSON(http.StatusOK, pkg.ClientErr("need key"))
		return
	}
	pkg.Log.Info(key)
	paper, err := model.FindByKey(key)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("find by key error:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("", paper))
}

type ListAuthorReq struct {
	Key    string `json:"key" binding:"required,excludes= "`
	Offset uint64 `json:"offset" binding:""`
	Count  uint64 `json:"count" binding:""`
}

// @Summary ListAuthor
// @Tags Paper
// @Accept json
// @Param data body ListAuthorReq true "ListAuthorReq"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/paper/listAuthor [post]
func ListAuthor(c *gin.Context) {
	var req ListAuthorReq
	if err := c.ShouldBind(&req); err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	paper, err := model.FindByKey(req.Key)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("find by key error:"+err.Error()))
		return
	}
	authors, err := paper.ListAuthor(req.Offset, req.Count)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("list author error:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("", authors))
}

// @Summary GetPublishVenue
// @Tags Paper
// @Accept json
// @Param key query string true "唯一标识"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/paper/getPublishVenue [get]
func GetPublishVenue(c *gin.Context) {
	var key string
	if key = c.DefaultQuery("key", ""); key == "" {
		pkg.Log.Error("need key")
		c.JSON(http.StatusOK, pkg.ClientErr("need key"))
		return
	}
	pkg.Log.Info(key)
	paper, err := model.FindByKey(key)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("find by key error:"+err.Error()))
		return
	}
	venue, err := paper.GetPublishVenue()
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("get venue error:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("", venue))
}
