package service

import (
	"net/http"
	"strconv"

	"github.com/Web-Engineering-XDU/Project-Backend/app/models"
	agentsystem "github.com/Web-Engineering-XDU/Project-Backend/app/service/agent_system"
	"github.com/gin-gonic/gin"
)

type getAgentListParams struct {
	ID     int `form:"id"`
	Number int `form:"number"`
	Page   int `form:"page"`
}

func GetAgentList(c *gin.Context) {
	params := &getAgentListParams{}
	c.ShouldBind(params)

	results := models.SelectAgentDetailList(
		params.Number,
		(params.Page-1)*params.Number,
	)

	c.JSON(http.StatusOK, makeRespBody(200, "ok", makeCountContent(len(results), results)))
}

func NewAgent(c *gin.Context) {
	params := &models.Agent{}
	c.ShouldBind(params)
	err := models.InsertAgent(params)
	if err != nil {
		c.JSON(http.StatusOK, makeRespBody(400, err.Error(), makeCountContent(0, nil)))
	}
	acInt, ok := c.Get("agents")
	if !ok {
		panic("gin need a agent collection to work")
	}
	ac := acInt.(*agentsystem.AgentCollection)
	ac.AddAgent(*params)

	c.JSON(http.StatusOK, makeRespBody(200, "ok", map[string]int{"id": params.ID}))
}

func DeleteAgent(c *gin.Context) {
	IdStr, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusOK, makeRespBody(400, "missing param: id", nil))
	}

	acInt, ok := c.Get("agents")
	if !ok {
		panic("gin need a agent collection to work")
	}
	ac := acInt.(*agentsystem.AgentCollection)

	ID, err := strconv.Atoi(IdStr)

	if err != nil {
		c.JSON(http.StatusOK, makeRespBody(400, "invalid param: id", nil))
	}

	if ac.DeleteAgent(ID) {
		c.JSON(http.StatusOK, makeRespBody(200, "ok", nil))
	}

	c.JSON(http.StatusOK, makeRespBody(400, "agent with this id doesn't exist", nil))
}

func UpdateAgent(c *gin.Context) {
	
}

func makeCountContent(count int, content any) gin.H {
	if content == nil {
		content = []struct{}{}
	}
	return gin.H{
		"count":   count,
		"content": content,
	}
}

func makeRespBody(code int, msg string, result any) gin.H {
	return gin.H{
		"code":   code,
		"msg":    msg,
		"result": result,
	}
}
