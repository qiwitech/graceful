// template.go
package main

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/descriptor"
)

var (
	funcLower = template.FuncMap{
		"ToLower": strings.ToLower,
	}
)

func applyTemplate(p headerParams) (string, error) {
	w := bytes.NewBuffer(nil)
	if err := headerTemplate.Execute(w, p); err != nil {
		return "", err
	}
	for _, svc := range p.Services {
		glog.V(1).Infof("Processing service %s", svc.GetName())

		for _, m := range svc.Methods {
			glog.V(1).Infof("\tmethod %s.%s", svc.GetName(), m.GetName())
			// Приведение имен к виду
			methName := strings.Title(*m.Name)
			m.Name = &methName
			// trim package name
			inputType := (*m.InputType)[len(p.GetPackage())+2:]
			m.InputType = &inputType
			outputType := (*m.OutputType)[len(p.GetPackage())+2:]
			m.OutputType = &outputType
		}

		if err := serviceTemplate.Execute(w, map[string]interface{}{"Service": svc}); err != nil {
			return "", err
		}

		if err := interfaceTemplate.Execute(w, bindingParams{Service: svc, Methods: svc.Methods}); err != nil {
			return "", err
		}
	}

	return w.String(), nil
}

type headerParams struct {
	*descriptor.File
	Imports []descriptor.GoPackage
}

var (
	headerTemplate = template.Must(template.New("header").Parse(`// Code generated by protoc-gen-httpgw
// source: {{.GetName}}
// DO NOT EDIT!

/*
Package {{.GoPkg.Name}} is a http proxy.
*/

package {{.GoPkg.Name}}

import (
{{range $i := .Imports}}{{if $i.Standard}}{{$i | printf "\t%s\n"}}{{end}}{{end}}

{{range $i := .Imports}}{{if not $i.Standard}}{{$i | printf "\t%s\n"}}{{end}}{{end}}
)
`))
)

// type descriptor.Service

var (
	serviceTemplate = template.Must(template.New("service").Parse(`
func New{{.Service.Name}}Handler(srv {{.Service.Name}}Interface, c graceful.Codec) graceful.Handlerer {
	return Add{{.Service.Name}}Handlers(nil, srv, c)
}
func Add{{.Service.Name}}Handlers(mux graceful.Handlerer, srv {{.Service.Name}}Interface, c graceful.Codec) graceful.Handlerer {
	if mux == nil {
		mux = chi.NewMux()
	}
{{range $m := .Service.Methods}}
	mux.Handle("/{{$m.Name}}", graceful.NewHandler(
		c,
		func() interface{} { return &{{$m.InputType}}{} },
		func(ctx context.Context, args interface{}) (interface{}, error) { return srv.{{$m.Name}}(ctx, args.(*{{$m.InputType}})) }))
{{end}}
	return mux
}

type {{.Service.Name}}HTTPClient struct {
	*graceful.Client
}
func New{{.Service.Name}}HTTPClient(cl *graceful.Client) {{.Service.Name}}HTTPClient {
	return {{.Service.Name}}HTTPClient{
		Client: cl,
	}
}

{{range $m := .Service.Methods}}
func (cl {{$.Service.Name}}HTTPClient) {{$m.Name}}(ctx context.Context, args *{{$m.InputType}}) (*{{$m.OutputType}}, error) {
	var resp {{$m.OutputType}}
	err := cl.Client.Call(ctx, "{{$m.Name}}", args, &resp)
	return &resp, err
}
{{end}}
`))
)

type bindingParams struct {
	Service *descriptor.Service
	Methods []*descriptor.Method
}

var (
	interfaceTemplate = template.Must(template.New("interface").Parse(`
type {{.Service.Name}}Interface interface {
{{range $m := $.Methods}}
	{{$m.Name}}(context.Context, *{{$m.InputType}}) (*{{$m.OutputType}}, error)
{{end}}
}
`))
)
