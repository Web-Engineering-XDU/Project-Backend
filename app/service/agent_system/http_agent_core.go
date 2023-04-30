package agentsystem

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/antchfx/jsonquery"
	jsoniter "github.com/json-iterator/go"
	"github.com/mohae/deepcopy"
	"github.com/osteele/liquid"
	"github.com/valyala/fasthttp"
)

type Selector struct {
	VarName         string `json:"var_name"`
	SelectorType    string `json:"selector_type"`
	SelectorContent string `json:"selector_content"`
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
	Mode       string `json:"mode"`
	MergeEvent bool   `json:"merge_event"`

	httpRequstsTemplate
	Template map[string]string `json:"template"`

	DocType   string     `json:"doc_type"`
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

func (hac *httpAgentCore) Run(ctx context.Context, agent *Agent, event *Event) {
	newEvent := &Event{}

	agent.Mutex.RLock()
	httpReqTemp := deepcopy.Copy(hac.httpRequstsTemplate).(httpRequstsTemplate)
	temp := deepcopy.Copy(hac.Template).(map[string]string)
	docType := hac.DocType
	selectors := deepcopy.Copy(hac.Selectors).([]Selector)
	mergeEvent := hac.MergeEvent
	agent.Mutex.RUnlock()

	err := httpReqTemp.render(event.Msg)
	if err != nil {
		newEvent.MetError = true
		//TODO
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
			newEvent.MetError = true
			//TODO
		}

		//extract data from doc
		resultMap := make([]map[string]string, 0, 20)
		err = nil
		switch docType {
		case "html":
			err = selectHtml(resp.Body(), selectors, &resultMap)
		case "json":
			err = selectJson(resp.Body(), selectors, &resultMap)
		case "text":
			err = selectText(resp.Body(), selectors, &resultMap)
		default:
			newEvent.MetError = true
			//TODO
		}

		if err != nil {
			newEvent.MetError = true
			//TODO
		}

		for _, v := range resultMap {
			bindings := deepcopy.Copy(v).(map[string]string)
			mergeMap(bindings, event.Msg)
			err = renderTemplate(temp, bindings)
			if err != nil {
				newEvent.MetError = true
				//TODO
			}
			mergeMap(temp, v)
			if mergeEvent {
				mergeMap(temp, event.Msg)
			}
			agent.ac.eventHdl.PushEvent(&Event{
				SrcAgent:      agent,
				CreateTime:    time.Now(),
				DeleteTime:    time.Now().Add(agent.EventMaxAge),
				Msg:           temp,
				ToBeDelivered: true,
			})
		}

	}
}

func (hac *httpAgentCore) Stop() {
	//TODO
}

// func renderAllStringExceptTemplate(v interface{}, bindings map[string]interface{}) error {

// }

var engine = liquid.NewEngine()

func selectHtml(doc []byte, selectors []Selector, result *[]map[string]string) error {
	//TODO
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
			return errors.New("unsupported json selector: " + v.SelectorType)
		}
		if err != nil {
			return err
		}
	}
	for eventInd := 0; eventInd < nodesMaxMum; eventInd++ {
		newMap := make(map[string]string)
		for selectorInd, v := range selectors {
			if eventInd <= len(nodesList[selectorInd]) {
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
