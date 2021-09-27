package kasset

import (
	"errors"
	"net/http"

	"git.kanosolution.net/kano/kaos"
)

var (
	Event kaos.EventHub
	Topic string

	e error
)

func (ae *AssetEngine) View(ctx *kaos.Context, assetid string) ([]byte, error) {
	r := ctx.Data().Get("http-request", &http.Request{}).(*http.Request)
	w := ctx.Data().Get("http-response", &http.Response{}).(http.ResponseWriter)

	assetID := r.URL.Query().Get("id")
	dl := r.URL.Query().Get("t") == "dl"

	if Event == nil {
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte("EventHub is not initialized"))
		return nil, errors.New("eventhub is not initialized")
	}

	ast := new(Asset)
	if e = Event.Publish(Topic, assetID, ast); e != nil {
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte(e.Error()))
		return nil, e
	}

	content, e := ae.fs.Read(ast.URI)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(e.Error()))
	}

	if dl {
		w.Header().Set("Content-Disposition", "attachment; filename=\""+ast.OriginalFileName+"\"")
	}
	w.Header().Set("Content-Type", ast.ContentType)
	w.Write(content)

	ctx.Data().Set("kaos-command-1", "stop")

	return content, nil
}
