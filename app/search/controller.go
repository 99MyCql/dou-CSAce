package search

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"douCSAce/pkg"
)

// @Summary Search
// @Tags Search
// @Accept json
// @Param query query string true "查询内容"
// @Success 200 {string} json "{"code":0,"data":{},"msg":""}"
// @Failure 200 {string} json "{"code":!0,"data":{},"msg":""}"
// @Router /api/v1/search [get]
func Search(c *gin.Context) {
	var query string
	if query = c.DefaultQuery("query", ""); query == "" {
		pkg.Log.Error("need key")
		c.JSON(http.StatusOK, pkg.ClientErr("need search key"))
		return
	}
	pkg.Log.Info(query)

	u := fmt.Sprintf("http://dblp.org/search/author/api?q=%s&format=json&h=10&f=0&c=0", url.QueryEscape(query))
	rsp, err := http.Get(u)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("request dblp/search error"))
		return
	}
	if rsp.StatusCode != http.StatusOK {
		pkg.Log.Error(rsp.Status + ":" + u)
		c.JSON(http.StatusOK, pkg.ServerErr("request dblp/search http error:"+rsp.Status))
		return
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	var authorJson map[string]interface{}
	if err := json.Unmarshal(body, &authorJson); err != nil {
		pkg.Log.Fatal(err)
	}

	u = fmt.Sprintf("http://dblp.org/search/venue/api?q=%s&format=json&h=10&f=0&c=0", url.QueryEscape(query))
	rsp, err = http.Get(u)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("request dblp/search error"))
		return
	}
	if rsp.StatusCode != http.StatusOK {
		pkg.Log.Error(rsp.Status + ":" + u)
		c.JSON(http.StatusOK, pkg.ServerErr("request dblp/search http error:"+rsp.Status))
		return
	}
	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	var venueJson map[string]interface{}
	if err := json.Unmarshal(body, &venueJson); err != nil {
		pkg.Log.Fatal(err)
	}

	u = fmt.Sprintf("http://dblp.org/search/publ/api?q=%s&format=json&h=5&f=0&c=0", url.QueryEscape(query))
	rsp, err = http.Get(u)
	if err != nil {
		pkg.Log.Error(err)
		c.JSON(http.StatusOK, pkg.ServerErr("request dblp/search error"))
		return
	}
	if rsp.StatusCode != http.StatusOK {
		pkg.Log.Error(rsp.Status + ":" + u)
		c.JSON(http.StatusOK, pkg.ServerErr("request dblp/search http error:"+rsp.Status))
		return
	}
	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		pkg.Log.Fatal(err)
	}
	var paperJson map[string]interface{}
	if err := json.Unmarshal(body, &paperJson); err != nil {
		pkg.Log.Fatal(err)
	}

	c.JSON(http.StatusOK, pkg.SucWithData("", map[string]interface{}{
		"query":  query,
		"author": authorJson,
		"venue":  venueJson,
		"paper":  paperJson,
	}))
}
