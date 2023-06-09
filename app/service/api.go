package service

import (
	"net/http"
	"strconv"

	"github.com/Web-Engineering-XDU/Project-Backend/app/models"
	agentsystem "github.com/Web-Engineering-XDU/Project-Backend/app/service/agent_system"
	_ "github.com/Web-Engineering-XDU/Project-Backend/docs/swaggo"
	"github.com/gin-gonic/gin"
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

	var results []models.AgentDetail
	var totalCount int64

	if params.ID == 0 {
		results, totalCount = models.SelectAgentDetailList(
			params.Number,
			(params.Page-1)*params.Number,
		)
	} else {
		agent, ok := models.SelectAgentDetailByID(params.ID)
		if ok {
			results = []models.AgentDetail{agent}
		}
	}

	c.JSON(http.StatusOK, makeRespBody(200, "ok", makeCountContent(len(results), int(totalCount), results)))
}

type getRelationableAgentsParams struct {
	ID      int    `form:"id"`
	Keyword string `form:"keyword"`
	Type    string `form:"type"`
	Number  int    `form:"number"`
	Page    int    `form:"page"`
}

// @Summary      List agents available for relation op
// @Description  List agents available for relation op
// @Tags         agents
// @Accept       json
// @Produce      json
// @Param        id			query	int    true    "agent id"
// @Param        keyword	query	string false   "keyword (%val%)"
// @Param        type		query	string true    "prev or next" Enums(prev, next)
// @Param        number		query	int    true		"number of agents in a page"
// @Param        page   	query	int    true		"page sequence number"
// @Success      200  {object}   swaggo.GetRelationableAgentsResp
// @Router       /agent/relationable [get]
func GetRelationableAgents(c *gin.Context) {
	params := &getRelationableAgentsParams{}
	c.ShouldBind(params)

	var results []models.AgentIdAndName
	var totalCount int64
	var err error

	switch params.Type {
	case "prev":
		results, totalCount, err = models.SelectPrevableAgents(
			params.ID,
			params.Keyword,
			params.Number,
			(params.Page-1)*params.Number,
		)
	case "next":
		results, totalCount, err = models.SelectNextableAgents(
			params.ID,
			params.Keyword,
			params.Number,
			(params.Page-1)*params.Number,
		)
	default:
		c.JSON(http.StatusOK, makeRespBody(400, "unknown relation type", nil))
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, makeRespBody(400, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, makeRespBody(200, "ok", makeCountContent(len(results), int(totalCount), results)))
}

// @Summary      New agents
// @Description  Create new agents. Also support application/json
// @Tags         agents
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        enable			formData	bool		true	"enable the agent"
// @Param		 typeId		    formData	int			true	"agent type id"
// @Param		 name			formData	string		true	"name of the agent"
// @Param		 description	formData	string		true	"description"
// @Param		 eventForever	formData	bool		true	"whether keep the event forever"
// @Param		 eventMaxAge	formData	int			true	"event max age in nanosecond count"
// @Param		 propJsonStr	formData	string		true	"props used by specific agent type in json"
// @Success      200  {object}   swaggo.NewAgentResponse
// @Router       /agent [put]
func NewAgent(c *gin.Context) {
	params := &models.Agent{}
	c.ShouldBind(params)

	acInt, ok := c.Get("agents")
	if !ok {
		panic("gin need a agent collection to work")
	}
	ac := acInt.(*agentsystem.AgentCollection)

	err := ac.AddAgent(params)

	if err != nil {
		c.JSON(http.StatusOK, makeRespBody(400, err.Error(), nil))
		return
	}

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

	ok, err = ac.DeleteAgent(ID)

	if err != nil {
		c.JSON(http.StatusOK, makeRespBody(400, err.Error(), nil))
		return
	}
	if !ok {
		c.JSON(http.StatusOK, makeRespBody(200, "agent with this id does not exist", nil))
		return
	}
	c.JSON(http.StatusOK, makeRespBody(200, "ok", nil))
}

// @Summary      Update agents
// @Description  Update agents with specific id. Also support application/json
// @Tags         agents
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        id			    formData	int         true    "agent id"
// @Param        enable			formData	bool		true	"enable the agent"
// @Param		 name			formData	string		true	"name of the agent"
// @Param		 description	formData	string		true	"description"
// @Param		 eventForever	formData	bool		true	"whether keep the event forever"
// @Param		 eventMaxAge	formData	int			true	"event max age in timestamp"
// @Param		 propJsonStr	formData	string		true	"props used by specific agent type in json"
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
	err := ac.UpdateAgent(*params)
	if err != nil {
		c.JSON(http.StatusOK, makeRespBody(400, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, makeRespBody(200, "ok", nil))
}

type getEventListParams struct {
	getListParams
	EventId int `form:"eventId"`
}

// @Summary      List Events
// @Description  get events
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        id			query	int    false    "src agent id. Don't include in request if you don't specify it"
// @Param        eventId	query	int    false    "event id. Don't include in request if you don't specify it"
// @Param        number		query	int    true		"number of events in a page"
// @Param        page   	query	int    true		"page sequence number"
// @Success      200  {object}   swaggo.GetEventListResponse
// @Router       /event [get]
func GetEventList(c *gin.Context) {
	params := &getEventListParams{}
	c.ShouldBind(params)

	var results []models.EventAndSrcAgentName
	var totalCount int64

	results, totalCount = models.SelectEventList(
		params.ID,
		params.EventId,
		params.Number,
		(params.Page-1)*params.Number,
	)
	results = append([]models.EventAndSrcAgentName{}, results...)

	c.JSON(http.StatusOK, makeRespBody(200, "ok", makeCountContent(len(results), int(totalCount), results)))
}

type setAgentRelationParams struct {
	AgentId int   `json:"agentId"`
	Srcs    []int `json:"srcs"`
	Dsts    []int `json:"dsts"`
}

// @Summary      Set Agent Relations
// @Description  Set relations of an agent with specific id. Also support x-www-form-urlencoded
// @Tags         relations
// @Accept       json
// @Produce      json
// @Param        request    body    setAgentRelationParams    true    "all relations of an agent"
// @Success      200  {object}   swaggo.GetEventListResponse
// @Router       /agent-relation [post]
func SetAgentRelation(c *gin.Context) {
	params := &setAgentRelationParams{}
	c.ShouldBind(params)
	acInt, ok := c.Get("agents")
	if !ok {
		panic("gin need a agent collection to work")
	}
	ac := acInt.(*agentsystem.AgentCollection)
	err := ac.SetAgentRelation(params.AgentId, params.Srcs, params.Dsts)
	if err != nil {
		c.JSON(http.StatusOK, makeRespBody(400, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, makeRespBody(200, "ok", nil))
}

// @Summary      List All Relations
// @Description  get all relations
// @Tags         relations
// @Accept       json
// @Produce      json
// @Success      200  {object}   swaggo.GetRelationsResponse
// @Router       /agent-relation [get]
func GetAllAgentRelations(c *gin.Context) {
	result := models.SelectAllAgentRelations()
	c.JSON(http.StatusOK, makeRespBody(200, "ok", makeCountContent(len(result), len(result), result)))
}

// @Summary      Get Relations By Agent ID
// @Description  Get Relations By Agent ID
// @Tags         relations
// @Accept       json
// @Produce      json
// @Param        id			query	int    false    "src agent id. Don't include in request if you don't specify it"
// @Success      200  {object}   swaggo.GetRelationsForEditResp
// @Router       /agent-relation/for-edit [get]
func GetAgentRelationsForEdit(c *gin.Context) {
	IdStr, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusOK, makeRespBody(400, "missing param: id", nil))
		return
	}
	ID, err := strconv.Atoi(IdStr)
	if err != nil {
		c.JSON(http.StatusOK, makeRespBody(400, "invalid param: id", nil))
		return
	}
	relations := models.SelectAgentRelationsAbout(ID)
	prevIds := make([]int, 0, 3)
	nextIds := make([]int, 0, 3)
	for _, v := range relations {
		if v.SrcAgentId != ID {
			prevIds = append(prevIds, v.SrcAgentId)
		} else {
			nextIds = append(nextIds, v.DstAgentId)
		}
	}
	prevs := models.SelectAgentIdAndNames(prevIds)
	nexts := models.SelectAgentIdAndNames(nextIds)
	c.JSON(http.StatusOK, makeRespBody(200, "ok", map[string]any{
		"srcs": prevs,
		"dsts": nexts,
	}))
}

type dryRunParams struct {
	AgentTypeId      int                 `form:"agentTypeId" json:"agentTypeId"`
	AgentPropJsonStr string              `form:"agentPropJsonStr" json:"agentPropJsonStr"`
	Msg              agentsystem.Message `form:"event" json:"event"`
}

// @Summary      Dry run
// @Description  Dry run
// @Tags         agents
// @Accept       json
// @Produce      json
// @Param        request    body    dryRunParams    true    "agent and event"
// @Success      200  {object}   swaggo.DryRunResponse
// @Router       /agent/dry-run [post]
func DryRun(c *gin.Context) {
	params := &dryRunParams{}
	c.ShouldBind(params)

	msg, err := agentsystem.DryRunAgent(&models.Agent{
		AgentBasic: models.AgentBasic{
			TypeId: params.AgentTypeId,
		},
		PropJsonStr: params.AgentPropJsonStr,
	}, params.Msg)
	if err != nil {
		c.JSON(http.StatusOK, makeRespBody(400, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, makeRespBody(200, "ok", msg))
}

func makeCountContent(count int, totalCount int, content any) gin.H {
	if content == nil {
		content = []struct{}{}
	}
	return gin.H{
		"count":      count,
		"totalCount": totalCount,
		"content":    content,
	}
}

func makeRespBody(code int, msg string, result any) gin.H {
	return gin.H{
		"code":   code,
		"msg":    msg,
		"result": result,
	}
}
