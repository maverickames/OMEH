package main

import (
	"errors"
	"net/http"

	"github.com/maverickames/omeh"
)

func (app *App) root(w http.ResponseWriter, r *http.Request) *omeh.ErrResponse {

	RenderdView := struct {
		Title       string
		ContentData string
	}{
		Title:       "Title of the current page.",
		ContentData: "This is the content data",
	}

	return app.RenderDefaultTemplate(w, r, "dashboardPlan.html.mustache", RenderdView)
}

func (app *App) brokenRoute(w http.ResponseWriter, r *http.Request) *omeh.ErrResponse {

	RenderdView := struct {
		Title       string
		ContentData string
	}{
		Title:       "Page 2 broken route",
		ContentData: "This is the content data of the broke route",
	}

	// faking function return of non nil error
	err := errors.New("did something and it broke")
	if err != nil {
		return app.em.ErrorHandler(err, omeh.StatusInternalServerError, r.RequestURI)
	}

	return app.RenderDefaultTemplate(w, r, "dashboardPlan.html.mustache", RenderdView)
}
