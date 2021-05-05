package journal

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"douCSAce/app/journal/model"
	"douCSAce/pkg"
)

// @Summary Count
// @Tags Journal
// @Accept json
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/journal/count [get]
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
// @Tags Journal
// @Accept json
// @Param key query string true "唯一标识"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/journal/get [get]
func Get(c *gin.Context) {
	var key string
	if key = c.DefaultQuery("key", ""); key == "" {
		pkg.Log.Error("need key")
		c.JSON(http.StatusOK, pkg.ClientErr("need key"))
		return
	}
	pkg.Log.Info(key)
	journal, err := model.FindByKey(key)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("find by key error:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("", journal))
}

type ListReq struct {
	Key      string `json:"key" binding:"required,excludes= "`
	Offset   uint64 `json:"offset" binding:""`
	Count    uint64 `json:"count" binding:""`
	SortAttr string `json:"sortAttr" binding:""`
	SortType string `json:"sortType" binding:""`
}

// @Summary ListPaper
// @Tags Journal
// @Accept json
// @Param data body ListReq true "ListReq"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/journal/listPaper [post]
func ListPaper(c *gin.Context) {
	var req ListReq
	if err := c.ShouldBind(&req); err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	journal, err := model.FindByKey(req.Key)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("find by key error:"+err.Error()))
		return
	}
	papers, err := journal.ListPaper(req.Offset, req.Count, req.SortAttr, req.SortType)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("list paper error:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("", papers))
}

// @Summary ListAuthor
// @Tags Journal
// @Accept json
// @Param data body ListReq true "ListReq"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/journal/listAuthor [post]
func ListAuthor(c *gin.Context) {
	var req ListReq
	if err := c.ShouldBind(&req); err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	journal, err := model.FindByKey(req.Key)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("find by key error:"+err.Error()))
		return
	}
	authors, err := journal.ListAuthor(req.Offset, req.Count, req.SortAttr, req.SortType)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("list author error:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("", authors))
}
