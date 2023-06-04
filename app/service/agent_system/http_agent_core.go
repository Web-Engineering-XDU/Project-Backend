package agentsystem

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/jsonquery"
	"github.com/antchfx/xmlquery"
	jsoniter "github.com/json-iterator/go"
	"github.com/mohae/deepcopy"
	"github.com/osteele/liquid"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/html"
)

type Selector struct {
	VarName         string `json:"varName"`
	SelectorType    string `json:"selectorType"`
	SelectorContent string `json:"selectorContent"`
}

type httpRequstsTemplate struct {
	Urls   []string          `json:"urls"`
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`
}

func convertMap(m map[string]string) map[string]interface{} {
	res := make(map[string]interface{}, len(m))
	for k, v := range m {
		res[k] = v
	}
	return res
}

func (t *httpRequstsTemplate) render(bindings map[string]string) error {
	out, err := engine.ParseAndRenderString(t.Method, convertMap(bindings))
	if err != nil {
		return err
	}
	t.Method = out

	out, err = engine.ParseAndRenderString(t.Body, convertMap(bindings))
	if err != nil {
		return err
	}
	t.Body = out

	for k, v := range t.Header {
		out, err = engine.ParseAndRenderString(v, convertMap(bindings))
		if err != nil {
			return err
		}
		t.Header[k] = out
	}

	for i, v := range t.Urls {
		out, err = engine.ParseAndRenderString(v, convertMap(bindings))
		if err != nil {
			return err
		}
		t.Urls[i] = out
	}

	return nil
}

type httpAgentCore struct {
	OnUpdate   bool `json:"onUpdate"`
	MergeEvent bool `json:"mergeEvent"`

	httpRequstsTemplate
	Template map[string]string `json:"template"`

	DocType   string     `json:"docType"`
	Selectors []Selector `json:"selectors"`
}

func renderTemplate(template, bindings map[string]string) error {
	var out string
	var err error
	for k, v := range template {
		out, err = engine.ParseAndRenderString(v, convertMap(bindings))
		if err != nil {
			return err
		}
		template[k] = out
	}
	return nil
}

func (a *Agent) loadHttpAgentCore() error {
	core := &httpAgentCore{}
	err := json.UnmarshalFromString(a.AgentCoreJsonStr, core)
	if err != nil {
		return err
	}
	a.AgentCore = core
	return nil
}

func (hac *httpAgentCore) Run(ctx context.Context, agent *Agent, event *Event, callBack func(e []*Event)) {
	deleteTime := time.Now().Add(agent.EventMaxAge)
	if agent.EventForever {
		deleteTime = time.Now().AddDate(100, 0, 0)
	}
	newEvent := &Event{
		SrcAgent:   agent,
		CreateTime: time.Now(),
		DeleteTime: deleteTime,
	}

	agent.Mutex.RLock()
	httpReqTemp := deepcopy.Copy(hac.httpRequstsTemplate).(httpRequstsTemplate)
	docType := hac.DocType
	selectors := deepcopy.Copy(hac.Selectors).([]Selector)
	mergeEvent := hac.MergeEvent
	agent.Mutex.RUnlock()

	err := httpReqTemp.render(event.Msg)
	if err != nil {
		handleRunError(newEvent, err, callBack)
		return
	}

	//Http Request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(httpReqTemp.Method)
	for key, val := range httpReqTemp.Header {
		req.Header.Set(key, val)
	}
	req.SetBodyString(httpReqTemp.Body)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	for _, v := range httpReqTemp.Urls {
		req.SetRequestURI(v)
		err = fasthttp.Do(req, resp)
		if err != nil {
			handleRunError(newEvent, err, callBack)
			return
		}

		//extract data from doc
		resultMap := make([]map[string]string, 0, 20)
		err = nil
		switch docType {
		case "html":
			err = selectHtml(resp.Body(), selectors, &resultMap)
		case "xml":
			err = selectXml(resp.Body(), selectors, &resultMap)
		case "json":
			err = selectJson(resp.Body(), selectors, &resultMap)
		case "text":
			err = selectText(resp.Body(), selectors, &resultMap)
		default:
			handleRunError(newEvent, errors.New("unsupported doc type: "+docType), callBack)
			return
		}

		if err != nil {
			handleRunError(newEvent, err, callBack)
			return
		}
		results := make([]*Event, len(resultMap))
		for i, v := range resultMap {
			temp := deepcopy.Copy(hac.Template).(map[string]string)
			bindings := deepcopy.Copy(v).(map[string]string)
			mergeMap(bindings, event.Msg)
			err = renderTemplate(temp, bindings)
			if err != nil {
				handleRunError(newEvent, err, callBack)
				return
			}
			mergeMap(temp, v)
			if mergeEvent {
				mergeMap(temp, event.Msg)
			}
			results[i] = &Event{
				SrcAgent:      newEvent.SrcAgent,
				CreateTime:    newEvent.CreateTime,
				DeleteTime:    newEvent.DeleteTime,
				Msg:           temp,
				ToBeDelivered: true,
			}
		}
		callBack(results)
	}
}

func handleRunError(newEvent *Event, err error, callBack func(e []*Event)) {
	if newEvent == nil {
		return
	}
	newEvent.MetError = true
	newEvent.ToBeDelivered = false
	newEvent.Log = err.Error()
	callBack([]*Event{newEvent})
}

func (hac *httpAgentCore) Stop() {
	//TODO
}

func (hac *httpAgentCore) IgnoreDuplicateEvent() bool {
	return hac.OnUpdate
}

func (hac *httpAgentCore) ValidCheck() error {
	if hac.DocType != "text" && hac.DocType != "json" && hac.DocType != "html" && hac.DocType != "xml"{
		return errors.New("unsupported doc type: " + hac.DocType)
	}
	for _, v := range hac.Selectors {
		if v.SelectorType != "xpath" {
			return errors.New("unsupported selector type: " + v.SelectorType)
		}
	}
	return nil
}

var engine = liquid.NewEngine()

func selectHtml(doc []byte, selectors []Selector, result *[]map[string]string) error {
	docNode, err := htmlquery.Parse(bytes.NewReader(doc))
	if err != nil {
		return err
	}
	selectorNum := len(selectors)
	nodesList := make([][]*html.Node, 0, selectorNum)
	nodesMaxMum := 0
	for i, v := range selectors {
		switch v.SelectorType {
		case "xpath":
			var nodes []*html.Node
			nodes, err = htmlquery.QueryAll(docNode, v.SelectorContent)
			if nodes == nil {
				nodes = []*html.Node{}
			}
			nodesList = append(nodesList, nodes)
			if len(nodesList[i]) > nodesMaxMum {
				nodesMaxMum = len(nodesList[i])
			}
		default:
			return errors.New("unsupported selector: " + v.SelectorType)
		}
		if err != nil {
			return err
		}
	}
	for eventInd := 0; eventInd < nodesMaxMum; eventInd++ {
		newMap := make(map[string]string)
		for selectorInd, v := range selectors {
			if eventInd <= len(nodesList[selectorInd]) - 1 {
				attr := false
				Secs := strings.Split(v.SelectorContent, "/")
				LastSec := Secs[len(Secs)-1]
				if LastSec[0] == '@' {
					attr = true
				}
				if attr {
					newMap[v.VarName] = htmlquery.SelectAttr(nodesList[selectorInd][eventInd], LastSec[1:])
				} else {
					newMap[v.VarName] = htmlquery.InnerText(nodesList[selectorInd][eventInd])
				}
			} else {
				newMap[v.VarName] = ""
			}
			if err != nil {
				return err
			}
		}
		*result = append(*result, newMap)
	}
	return nil
}

func selectXml(doc []byte, selectors []Selector, result *[]map[string]string) error {
	docNode, err := xmlquery.Parse(bytes.NewReader(doc))
	if err != nil {
		return err
	}
	selectorNum := len(selectors)
	nodesList := make([][]*xmlquery.Node, 0, selectorNum)
	nodesMaxMum := 0
	for i, v := range selectors {
		switch v.SelectorType {
		case "xpath":
			var nodes []*xmlquery.Node
			nodes, err = xmlquery.QueryAll(docNode, v.SelectorContent)
			if nodes == nil {
				nodes = []*xmlquery.Node{}
			}
			nodesList = append(nodesList, nodes)
			if len(nodesList[i]) > nodesMaxMum {
				nodesMaxMum = len(nodesList[i])
			}
		default:
			return errors.New("unsupported selector: " + v.SelectorType)
		}
		if err != nil {
			return err
		}
	}
	for eventInd := 0; eventInd < nodesMaxMum; eventInd++ {
		newMap := make(map[string]string)
		for selectorInd, v := range selectors {
			if eventInd <= len(nodesList[selectorInd]) - 1 {
				attr := false
				Secs := strings.Split(v.SelectorContent, "/")
				LastSec := Secs[len(Secs)-1]
				if LastSec[0] == '@' {
					attr = true
				}
				if attr {
					newMap[v.VarName] = nodesList[selectorInd][eventInd].SelectAttr(LastSec[1:])
				} else {
					newMap[v.VarName] = nodesList[selectorInd][eventInd].InnerText()
				}
			} else {
				newMap[v.VarName] = ""
			}
			if err != nil {
				return err
			}
		}
		*result = append(*result, newMap)
	}
	return nil
}

func selectJson(doc []byte, selectors []Selector, result *[]map[string]string) error {
	docNode, err := jsonquery.Parse(bytes.NewReader(doc))
	if err != nil {
		return err
	}
	selectorNum := len(selectors)
	nodesList := make([][]*jsonquery.Node, 0, selectorNum)
	nodesMaxMum := 0
	for i, v := range selectors {
		switch v.SelectorType {
		case "xpath":
			var nodes []*jsonquery.Node
			nodes, err = jsonquery.QueryAll(docNode, v.SelectorContent)
			if nodes == nil {
				nodes = []*jsonquery.Node{}
			}
			nodesList = append(nodesList, nodes)
			if len(nodesList[i]) > nodesMaxMum {
				nodesMaxMum = len(nodesList[i])
			}
		default:
			return errors.New("unsupported selector: " + v.SelectorType)
		}
		if err != nil {
			return err
		}
	}
	for eventInd := 0; eventInd < nodesMaxMum; eventInd++ {
		newMap := make(map[string]string)
		for selectorInd, v := range selectors {
			if eventInd <= len(nodesList[selectorInd]) - 1 {
				newMap[v.VarName], err = jsonNodeToStr(nodesList[selectorInd][eventInd])
			} else {
				newMap[v.VarName] = ""
			}
			if err != nil {
				return err
			}
		}
		*result = append(*result, newMap)
	}

	return nil
}

func selectText(doc []byte, selectors []Selector, result *[]map[string]string) error {
	//TODO
	return nil
}

func jsonNodeToStr(node *jsonquery.Node) (string, error) {
	v := node.Value()
	if reflect.ValueOf(v).Kind() == reflect.String {
		return v.(string), nil
	}
	return jsoniter.MarshalToString(v)
}

// if has same key, won't overwrite
func mergeMap(dst, src map[string]string) {
	for k, v := range src {
		_, ok := dst[k]
		if !ok {
			dst[k] = v
		}
	}
}
