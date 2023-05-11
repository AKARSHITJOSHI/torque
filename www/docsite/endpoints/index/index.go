package index

import (
	_ "embed"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tylermmorton/tmpl"
	"github.com/tylermmorton/torque"
	"github.com/tylermmorton/torque/www/docsite/domain/content"
	"github.com/tylermmorton/torque/www/docsite/model"
	"net/http"
)

// TODO(tylermorton) update this when tmpl is refactored to use viper
//go:generate tmplbind

// DotContext is the dot context of the index page template.
//
//tmpl:bind index.tmpl.html --watch
type DotContext struct {
	Document *model.Document
}

var Template = tmpl.MustCompile(&DotContext{})

// RouteModule is the torque route module to be registered with the torque app.
type RouteModule struct {
	ContentSvc content.Service
}

var _ interface {
	torque.Loader
	torque.Renderer
	torque.ErrorBoundary
} = &RouteModule{}

func (rm *RouteModule) Load(req *http.Request) (any, error) {
	pageName, ok := mux.Vars(req)["pageName"]
	if !ok {
		return nil, fmt.Errorf("fail to get document name in route vars '%s'", pageName)
	}

	doc, err := rm.ContentSvc.Get(req.Context(), pageName)
	if err != nil {
		return nil, fmt.Errorf("failed to get document by name '%s': %+v", pageName, err)
	}

	return doc, nil
}

func (rm *RouteModule) Render(wr http.ResponseWriter, req *http.Request, loaderData any) error {
	if document, ok := loaderData.(*model.Document); !ok {
		return fmt.Errorf("invalid loader data type. expected *model.Document, got %T", loaderData)
	} else {
		return Template.Render(wr, &DotContext{
			Document: document,
		})
	}
}

func (rm *RouteModule) ErrorBoundary(wr http.ResponseWriter, req *http.Request, err error) http.HandlerFunc {
	panic(err)
}
