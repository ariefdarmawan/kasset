package kasset

import (
	"errors"
	"net/http"
	"strings"

	"git.kanosolution.net/kano/kaos"
	"github.com/sebarcode/codekit"
)

var (
	Event kaos.EventHub
	Topic string
)

func (ae *AssetEngine) View(ctx *kaos.Context, assetid string) ([]byte, error) {
	r, rOK := ctx.Data().Get("http_request", nil).(*http.Request)
	w, wOK := ctx.Data().Get("http_writer", nil).(http.ResponseWriter)

	if !rOK {
		return nil, errors.New("not a http compliant request")
	}

	if !wOK {
		return nil, errors.New("not a http compliant writer")
	}

	assetID := r.URL.Query().Get("id")
	dl := r.URL.Query().Get("t") == "dl"

	ast := new(Asset)
	h, _ := ctx.DefaultHub()
	if e := h.GetByID(ast, assetID); e != nil {
		return nil, e
	}

	content, e := ae.fs.Read(ast.URI)
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(e.Error()))
	}

	if dl {
		if ast.OriginalFileName != "" {
			w.Header().Set("Content-Disposition", "attachment; filename=\""+ast.OriginalFileName+"\"")
		} else {
			newFileName := codekit.GenerateRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 16)
			uriParts := strings.Split(ast.URI, ".")
			newFileName += "." + uriParts[len(uriParts)-1]
			w.Header().Set("Content-Disposition", "attachment; filename=\""+newFileName+"\"")
		}
	}
	w.Header().Set("Content-Type", ast.ContentType)
	w.Write(content)

	ctx.Data().Set("kaos-command-1", "stop")
	return content, nil
}

// its nothing
