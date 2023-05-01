package service

import (
	"net/http"
	"strconv"

	"github.com/Web-Engineering-XDU/Project-Backend/app/models"
	agentsystem "github.com/Web-Engineering-XDU/Project-Backend/app/service/agent_system"
	"github.com/gin-gonic/gin"
	_ "github.com/Web-Engineering-XDU/Project-Backend/docs/swaggo"

)

type getListParams struct {
	ID     int `form:"id"`
	Number int `form:"number"`
	Page   int `form:"page"`
}

// @Summary      List agents
// @Description  get agents
// @Tags         agents
// @Accept       json
// @Produce      json
// @Param        id			query	int    false    "agent id. Don't include in request if you don't specify a id"  
// @Param        number		query	int    true		"number of agents in a page"
// @Param        page   	query	int    true		"page sequence number"
// @Success      200  {object}   swaggo.GetAgentListResponse
// @Router       /agent [get]
func GetAgentList(c *gin.Context) {
	params := &getListParams{}
	c.ShouldBind(params)

	results := make([]models.AgentDetail, 0)

	if params.ID == 0 {
		results = append(
			results,
			models.SelectAgentDetailList(
				params.Number,
				(params.Page-1)*params.Number,
			)...,
		)
	} else {
		agent, ok := models.SelectAgentDetailByID(params.ID)
		if ok {
			results = append(results, agent)
		}
	}

	c.JSON(http.StatusOK, makeRespBody(200, "ok", makeCountContent(len(results), results)))
}

// @Summary      New agents
// @Description  new agents
// @Tags         agents
// @Accept       json
// @Produce      json
// @Param        enable			formData	bool		true	"enable the agent"
// @Param		 type_id		formData	int			true	"agent type id"
// @Param		 name			formData	string		true	"name of the agent"
// @Param		 description	formData	string		true	"description"
// @Param		 event_forever	formData	bool		true	"whether keep the event forever"
// @Param		 event_max_age	formData	int			true	"event max age in timestamp"
// @Param		 prop_json_str	formData	string		true	"props used by specific agent type in json"
// @Success      200  {object}   swaggo.NewAgentResponse
// @Router       /agent [put]
func NewAgent(c *gin.Context) {
	params := &models.Agent{}
	c.ShouldBind(params)
	err := models.InsertAgent(params)
	if err != nil {
		c.JSON(http.StatusOK, makeRespBody(400, err.Error(), makeCountContent(0, nil)))
		return
	}
	acInt, ok := c.Get("agents")
	if !ok {
		panic("gin need a agent collection to work")
	}
	ac := acInt.(*agentsystem.AgentCollection)
	ac.AddAgent(*params)

	c.JSON(http.StatusOK, makeRespBody(200, "ok", map[string]int{"id": params.ID}))
}

// @Summary      Delete agents
// @Description  delete agents
// @Tags         agents
// @Accept       json
// @Produce      json
// @Param        id		query	int    true    "agent id"  
// @Success      200  {object}   swaggo.StateInfo
// @Router       /agent [delete]
func DeleteAgent(c *gin.Context) {
	IdStr, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusOK, makeRespBody(400, "missing param: id", nil))
		return
	}

	acInt, ok := c.Get("agents")
	if !ok {
		panic("gin need a agent collection to work")
	}
	ac := acInt.(*agentsystem.AgentCollection)

	ID, err := strconv.Atoi(IdStr)

	if err != nil {
		c.JSON(http.StatusOK, makeRespBody(400, "invalid param: id", nil))
		return
	}

	if ac.DeleteAgent(ID) {
		c.JSON(http.StatusOK, makeRespBody(200, "ok", nil))
		return
	}

	c.JSON(http.StatusOK, makeRespBody(400, "agent with this id doesn't exist", nil))
}

// @Summary      Update agents
// @Description  update agents
// @Tags         agents
// @Accept       json
// @Produce      json
// @Param        enable			formData	bool		true	"enable the agent"
// @Param		 name			formData	string		true	"name of the agent"
// @Param		 description	formData	string		true	"description"
// @Param		 event_forever	formData	bool		true	"whether keep the event forever"
// @Param		 event_max_age	formData	int			true	"event max age in timestamp"
// @Param		 prop_json_str	formData	string		true	"props used by specific agent type in json"
// @Success      200  {object}   swaggo.StateInfo
// @Router       /agent [post]
func UpdateAgent(c *gin.Context) {
	params := &models.Agent{}
	c.ShouldBind(params)
	acInt, ok := c.Get("agents")
	if !ok {
		panic("gin need a agent collection to work")
	}
	ac := acInt.(*agentsystem.AgentCollection)
	ok = ac.UpdateAgent(*params)
	if !ok {
		c.JSON(http.StatusOK, makeRespBody(400, "agent with this id doesn't exist", nil))
		return
	}
	c.JSON(http.StatusOK, makeRespBody(200, "ok", nil))
}

// @Summary      List Events
// @Description  get events
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        id			query	int    false    "src agent id. Don't include in request if you don't specify it"  
// @Param        number		query	int    true		"number of events in a page"
// @Param        page   	query	int    true		"page sequence number"
// @Success      200  {object}   swaggo.GetEventListResponse
// @Router       /event [get]
func GetEventList(c *gin.Context) {
	params := &getListParams{}
	c.ShouldBind(params)

	results := make([]models.Event, 0)

	if params.ID == 0 {
		results = append(
			results,
			models.SelectEventList(
				params.Number,
				(params.Page-1)*params.Number,
			)...,
		)
	} else {
		agent := models.SelectEventListByAgentID(
			params.ID,
			params.Number,
			(params.Page-1)*params.Number,
		)
		results = append(results, agent...)
	}

	c.JSON(http.StatusOK, makeRespBody(200, "ok", makeCountContent(len(results), results)))
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
