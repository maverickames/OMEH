package main

import (
	"encoding/json"
	"net/http"

	"github.com/cbroglie/mustache"
	"github.com/maverickames/omeh"
)

// Render404notFound Returnes Mustache Template 404 not found
func (app *App) Render404notFound(w http.ResponseWriter, r *http.Request) *omeh.ErrResponse {

	renderdView, err := mustache.RenderFile("Error404.html.mustache")
	if err != nil {
		return app.em.ErrorHandler(err, omeh.StatusInternalServerError, r.RequestURI)
	}
	w.Write([]byte(renderdView))
	return nil
}

// RenderTemplate Render template
func (app *App) RenderTemplate(w http.ResponseWriter, r *http.Request, viewTemplate string, viewConstruct interface{}) *omeh.ErrResponse {

	renderdView, err := mustache.RenderFile(viewTemplate, viewConstruct)
	if err != nil {
		return app.em.ErrorHandler(err, omeh.StatusInternalServerError, r.RequestURI)
	}
	w.Write([]byte(renderdView))
	return nil
}

// RenderDefaultTemplate Render template inside of layout
func (app *App) RenderDefaultTemplate(w http.ResponseWriter, r *http.Request, viewTemplate string, viewConstruct interface{}) *omeh.ErrResponse {

	renderdView, err := mustache.RenderFileInLayout(viewTemplate, "defaultTemplate.html.mustache", viewConstruct)
	if err != nil {
		return app.em.ErrorHandler(err, omeh.StatusInternalServerError, r.RequestURI)
	}
	w.Write([]byte(renderdView))
	return nil
}

// RenderCustomTemplate Renders mustache template Pass the template file name and the viewConttruct that you want.
func (app *App) RenderCustomTemplate(w http.ResponseWriter, r *http.Request, defaultTemplate string, viewTemplate string, viewConstruct interface{}) *omeh.ErrResponse {

	renderdView, err := mustache.RenderFileInLayout(viewTemplate, defaultTemplate, viewConstruct)
	if err != nil {
		return app.em.ErrorHandler(err, omeh.StatusInternalServerError, r.RequestURI)
	}
	w.Write([]byte(renderdView))
	return nil
}

func (app *App) renderJSON(w http.ResponseWriter, r *http.Request, viewConstruct interface{}) *omeh.ErrResponse {

	jsonData, err := json.Marshal(viewConstruct)
	if err != nil {
		return app.em.ErrorHandler(err, nil, r.RequestURI)
	}
	w.Write(jsonData)
	return nil
}

// Referer returnes the page back to the refer page after processing POST request
func redirectToReferer(w http.ResponseWriter, r *http.Request) *omeh.ErrResponse {
	http.Redirect(w, r, r.Header.Get("Referer"), 302)
	return nil
}
