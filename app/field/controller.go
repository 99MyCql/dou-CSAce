package field

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"douCSAce/app/field/model"
	"douCSAce/pkg"
)

// @Summary Get
// @Tags Field
// @Accept json
// @Param key query string true "唯一标识"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/field/get [get]
func Get(c *gin.Context) {
	var key string
	if key = c.DefaultQuery("key", ""); key == "" {
		pkg.Log.Error("need key")
		c.JSON(http.StatusOK, pkg.ClientErr("need key"))
		return
	}
	pkg.Log.Info(key)
	field, err := model.FindByKey(key)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("find by key error:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("", field))
}

// @Summary List
// @Tags Field
// @Accept json
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/field/list [get]
func List(c *gin.Context) {
	var key string
	if key = c.DefaultQuery("key", ""); key == "" {
		pkg.Log.Error("need key")
		c.JSON(http.StatusOK, pkg.ClientErr("need key"))
		return
	}
	pkg.Log.Info(key)
	fields, err := model.List()
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("list error:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("", fields))
}

type ListReq struct {
	Key      string `json:"key" binding:"required,excludes= "`
	Offset   uint64 `json:"offset" binding:""`
	Count    uint64 `json:"count" binding:""`
	SortAttr string `json:"sortAttr" binding:""`
	SortType string `json:"sortType" binding:""`
}

// @Summary ListPaper
// @Tags Field
// @Accept json
// @Param data body ListReq true "ListReq"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/field/listPaper [post]
func ListPaper(c *gin.Context) {
	var req ListReq
	if err := c.ShouldBind(&req); err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	field, err := model.FindByKey(req.Key)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("find by key error:"+err.Error()))
		return
	}
	papers, err := field.ListPaper(req.Offset, req.Count, req.SortAttr, req.SortType)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("list paper error:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("", papers))
}

// @Summary ListAuthor
// @Tags Field
// @Accept json
// @Param data body ListReq true "ListReq"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/field/listAuthor [post]
func ListAuthor(c *gin.Context) {
	var req ListReq
	if err := c.ShouldBind(&req); err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	field, err := model.FindByKey(req.Key)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("find by key error:"+err.Error()))
		return
	}
	authors, err := field.ListAuthor(req.Offset, req.Count, req.SortAttr, req.SortType)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("list author error:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("", authors))
}

// @Summary ListVenue
// @Tags Field
// @Accept json
// @Param data body ListReq true "ListReq"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/field/listVenue [post]
func ListVenue(c *gin.Context) {
	var req ListReq
	if err := c.ShouldBind(&req); err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	field, err := model.FindByKey(req.Key)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("find by key error:"+err.Error()))
		return
	}
	venues, err := field.ListVenue(req.Offset, req.Count, req.SortAttr, req.SortType)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("list venue error:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("", venues))
}
