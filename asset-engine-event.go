package kasset

import (
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/kaos"
	"github.com/sebarcode/codekit"
)

type AssetEngine struct {
	fs          AssetFS
	topicPrefix string
}

type AssetData struct {
	Asset   *Asset
	Content []byte
}

type AssetDataBase64 struct {
	Asset   *Asset
	Content string
}

func NewAssetData() *AssetData {
	ad := new(AssetData)
	ad.Asset = new(Asset)
	return ad
}

func NewAssetEngine(fs AssetFS, topicPrefix string) *AssetEngine {
	a := new(AssetEngine)
	a.fs = fs
	a.topicPrefix = topicPrefix
	return a
}

func (a *AssetEngine) Write(ctx *kaos.Context, attachReq *AssetData) (*Asset, error) {
	var e error
	h := GetTenantDBFromContext(ctx)
	if h == nil {
		return nil, fmt.Errorf("missing: db")
	}

	asset := attachReq.Asset
	if asset == nil {
		return nil, fmt.Errorf("asset is invalid")
	}
	if asset.ID == "" {
		asset.PreSave(nil)
	}
	if len(attachReq.Content) == 0 {
		return nil, fmt.Errorf("content is blank")
	}
	if asset.NewFileName != "" {
		// if newFileName is not blank, replace the asset
		other := new(Asset)
		if e = h.GetByParm(other, dbflex.NewQueryParam().SetWhere(dbflex.Eq("uri", asset.NewFileName))); e == nil {
			other.NewFileName = asset.NewFileName
			asset = other
		}
	}

	// save the file
	ext := ""
	ext1 := ""
	fileNameParts := strings.Split(attachReq.Asset.OriginalFileName, ".")
	if len(fileNameParts) > 1 {
		ext = "." + fileNameParts[len(fileNameParts)-1]
	} else {
		ext = ".txt"
	}

	contentType := ""
	if asset.ContentType != "" {
		contentType, ext1, _ = getFileType(attachReq.Content[:512])
		if ext == "" {
			if ext1 != "" && ext1[0] != '.' {
				ext1 = "." + ext
			}
			ext = ext1
		}

		if asset.ContentType == "" {
			asset.ContentType = contentType
		}

		if asset.ContentType == "" {
			return nil, fmt.Errorf("unknown file type")
		}
	}

	newFileName := asset.NewFileName
	if newFileName == "" {
		newFileName = fmt.Sprintf("%s_%s%s",
			asset.ID, codekit.GenerateRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 6),
			ext)
	}
	asset.URI = newFileName
	asset.Size = len(attachReq.Content)

	if e = a.fs.Save(newFileName, attachReq.Content); e != nil {
		return nil, fmt.Errorf("fail to write file %s: %s", newFileName, e.Error())
	}

	if e = h.Save(asset); e != nil {
		// rollback, delete the file
		a.fs.Delete(newFileName)
		return nil, fmt.Errorf("unable to save file metadata. %s", e.Error())
	}

	return asset, nil
}

func (a *AssetEngine) WriteWithContent(ctx *kaos.Context, attachReq *AssetDataBase64) (*Asset, error) {
	bs, e := base64.StdEncoding.DecodeString(attachReq.Content)
	if e != nil {
		return nil, fmt.Errorf("fail to decode content. %s", e.Error())
	}

	req := new(AssetData)
	req.Asset = attachReq.Asset
	req.Content = bs

	return a.Write(ctx, req)
}

func (a *AssetEngine) Read(ctx *kaos.Context, id string) (*Asset, error) {
	h := GetTenantDBFromContext(ctx)
	if h == nil {
		return nil, fmt.Errorf("missing: db")
	}

	ast := new(Asset)
	ast.ID = id
	var e error
	if e = h.Get(ast); e != nil {
		return nil, e
	}

	/*
		bs, e := a.fs.Read(ast.URI)
		if e != nil {
			return nil, fmt.Errorf("error reading file. %s", e.Error())
		}
	*/

	return ast, nil
}

func (ae *AssetEngine) Delete(ctx *kaos.Context, id string) (int, error) {
	h := GetTenantDBFromContext(ctx)
	if h == nil {
		return 0, fmt.Errorf("missing: db")
	}

	var e error
	a := new(Asset)
	a.ID = id
	if e = h.Get(a); e != nil {
		if e == io.EOF {
			return 0, nil
		}
		return 0, e
	}
	if e = ae.fs.Delete(a.URI); e != nil {
		if e != io.EOF {
			return 0, e
		}
	}
	h.DeleteQuery(new(AssetReference), dbflex.Eq("assetid", id))
	if e = h.Delete(a); e != nil {
		return 0, e
	}
	return a.Size, nil
}

type SaveAttrRequest struct {
	ID   string                 `json:"_id"`
	Data map[string]interface{} `json:"data"`
}

func (ae *AssetEngine) SaveAttr(ctx *kaos.Context, req *SaveAttrRequest) (string, error) {
	h := GetTenantDBFromContext(ctx)
	if h == nil {
		return "", fmt.Errorf("missing: db")
	}

	var e error
	a := new(Asset)
	a.ID = req.ID
	if e = h.Get(a); e != nil {
		return "", e
	}
	fields := []string{}
	for k := range req.Data {
		fields = append(fields, k)
	}
	cmd := dbflex.From(a.TableName()).Where(dbflex.Eq("_id", a.ID)).Update(fields...)
	if _, e = h.Execute(cmd, req.Data); e != nil {
		return "", e
	}
	return req.ID, nil
}
