package author

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"douCSAce/app/author/service"
	"douCSAce/pkg"
)

// @Summary Count
// @Tags Author
// @Accept json
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/author/count [get]
func Count(c *gin.Context) {
	count, err := service.Count()
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("服务端错误"))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("success", count))
}

// @Summary Get
// @Tags Author
// @Accept json
// @Param key query string true "author唯一标识"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/author/get [get]
func Get(c *gin.Context) {
	var key string
	if key = c.DefaultQuery("key", ""); key == "" {
		pkg.Log.Error("need key")
		c.JSON(http.StatusOK, pkg.ClientErr("need key"))
		return
	}
	pkg.Log.Info(key)
	author, err := service.FindByKey(key)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("find by key error:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("", author))
}

type ListPaperReq struct {
	Key      string `json:"key" binding:"required,excludes= "`
	Offset   uint64 `json:"offset" binding:""`
	Count    uint64 `json:"count" binding:""`
	SortAttr string `json:"sortAttr" binding:""`
	SortType string `json:"sortType" binding:""`
}

// @Summary ListPaper
// @Tags Author
// @Accept json
// @Param ListPaperReq body ListPaperReq true "ListPaperReq"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/author/listPaper [post]
func ListPaper(c *gin.Context) {
	var req ListPaperReq
	if err := c.ShouldBind(&req); err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ClientErr(err.Error()))
		return
	}
	author, err := service.FindByKey(req.Key)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("find by key error:"+err.Error()))
		return
	}
	papers, err := author.ListPaper(req.Offset, req.Count, req.SortAttr, req.SortType)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("list paper error:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("", papers))
}

// @Summary ListCoAuthor
// @Tags Author
// @Accept json
// @Param key query string true "author唯一标识"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/author/listCoAuthor [get]
func ListCoAuthor(c *gin.Context) {
	var key string
	if key = c.DefaultQuery("key", ""); key == "" {
		pkg.Log.Error("need key")
		c.JSON(http.StatusOK, pkg.ClientErr("need key"))
		return
	}
	pkg.Log.Info(key)
	author, err := service.FindByKey(key)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("find by key error:"+err.Error()))
		return
	}
	coauthors, err := author.ListCoAuthor()
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("list coauthor error:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, pkg.SucWithData("", coauthors))
}
