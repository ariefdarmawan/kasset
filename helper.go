package kasset

import (
	"git.kanosolution.net/kano/kaos"
	"github.com/ariefdarmawan/datahub"
	"github.com/h2non/filetype"
	"github.com/sebarcode/codekit"
)

func getFileType(buffer []byte) (string, string, error) {
	kind, err := filetype.Match(buffer)
	if err != nil {
		return "", "", err
	}
	return kind.MIME.Value, kind.Extension, nil
}

func GetFileType(buffer []byte) (string, string, error) {
	data := buffer
	if len(data) > 512 {
		data = data[:512]
	}
	return getFileType(data)
}

func GetTenantDBFromContext(ctx *kaos.Context) *datahub.Hub {
	jwtdata := ctx.Data().Get("jwt_data", codekit.M{}).(codekit.M)
	tenantID := jwtdata.GetString("TenantID")
	if tenantID == "" {
		tenantID = "Demo"
	}
	h, _ := ctx.GetHub(tenantID, "tenant")
	return h
}
