package agentsystem

import (
	"context"
	"reflect"

	"github.com/osteele/liquid"
	"github.com/valyala/fasthttp"
)

type Selector struct {
	VarName         string `json:"var_name"`
	SelectorType    string `json:"selector_type"`
	SelectorContent string `json:"selector_content"`
}

type httpAgentCore struct {
	Url    []string          `json:"urls"`
	Method string            `json:"method"`
	Header map[string]string `json:"header"`
	Body   string            `json:"body"`

	DocType   string            `json:"doc_type"`
	Selectors []Selector        `json:"selectors"`
	Template  map[string]string `json:"template"`

	Mode string `json:"mode"`
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

	template := make(map[string]string)
	for k, v := range hac.Template {
		template[k] = v
	}
	hac.Template = nil

	err := renderAllString(hac, event.Msg)
	if err != nil {
		newEvent.MetError = true
		//TODO
	}

	//Http Request
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.Header.SetMethod(hac.Method)
	for key, val := range hac.Header {
		req.Header.Set(key, val)
	}
	req.SetBodyString(hac.Body)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = fasthttp.Do(req, resp)
	if err != nil {
		newEvent.MetError = true
		//TODO
	}

	//extract data from doc
	resultMap := make(map[string]string)
	err = nil
	switch hac.DocType {
	case "html":
		err = selectHtml(resp.Body(), hac.Selectors, resultMap)
	case "json":
		err = selectJson(resp.Body(), hac.Selectors, resultMap)
	case "text":
		err = selectText(resp.Body(), hac.Selectors, resultMap)
	default:
		newEvent.MetError = true
		//TODO
	}

	if err != nil {
		newEvent.MetError = true
		//TODO
	}

}

func (hac *httpAgentCore) Stop() {
	//TODO
}

func renderAllString(v interface{}, bindings map[string]interface{}) error {
	rv := reflect.ValueOf(v).Elem()
	rt := rv.Type()
	n := rt.NumField()
	var err error

	for i := 0; i < n; i++ {
		field := rv.Field(i)

		switch field.Kind() {
		case reflect.Struct:
			err = renderAllString(field.Addr().Interface(), bindings)
			if err != nil {
				return err
			}
		case reflect.Map:
			keys := field.MapKeys()
			for _, key := range keys {
				elem := field.MapIndex(key)
				if elem.Kind() == reflect.String {
					err = render(&elem, bindings)
					if err != nil {
						return nil
					}
				}
			}
		case reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				if elem.Kind() == reflect.String {
					err = render(&elem, bindings)
					if err != nil {
						return nil
					}
				}
			}
		case reflect.String:
			err = render(&field, bindings)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}

var engine = liquid.NewEngine()

func render(elem *reflect.Value, bindings map[string]interface{}) error {
	out, err := engine.ParseAndRenderString(elem.String(), bindings)
	if err != nil {
		return err
	}
	elem.SetString(out)
	return nil
}

func selectHtml(doc []byte, selectors []Selector, result map[string]string) error {
	return nil
}

func selectJson(doc []byte, selectors []Selector, result map[string]string) error {
	return nil
}

func selectText(doc []byte, selectors []Selector, result map[string]string) error {
	return nil
}
