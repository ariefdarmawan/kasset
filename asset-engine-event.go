package kasset

import (
	"fmt"
	"io"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/kaos"
	"github.com/eaciit/toolkit"
	"github.com/h2non/filetype"
)

type AssetEngine struct {
	fs          AssetFS
	topicPrefix string
}

type AssetData struct {
	Asset   *Asset `json:"asset"`
	Content []byte `json:"content"`
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
	h, e := ctx.DefaultHub()
	if e != nil {
		return nil, e
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
	asset.ContentType, ext, _ = getFileType(attachReq.Content[:512])
	if ext != "" && ext[0] != '.' {
		ext = "." + ext
	}
	if asset.ContentType == "" {
		return nil, fmt.Errorf("unknown file type")
	}
	newFileName := asset.NewFileName
	if newFileName == "" {
		newFileName = fmt.Sprintf("%s_%s%s",
			asset.ID, toolkit.GenerateRandomString("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 6),
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

func (a *AssetEngine) Read(ctx *kaos.Context, id string) (*Asset, error) {
	h, e := ctx.DefaultHub()
	if e != nil {
		return nil, e
	}

	ast := new(Asset)
	ast.ID = id
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
	h, e := ctx.DefaultHub()
	if e != nil {
		return 0, e
	}
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
	h, e := ctx.DefaultHub()
	if e != nil {
		return "", e
	}
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

func getFileType(buffer []byte) (string, string, error) {
	kind, err := filetype.Match(buffer)
	if err != nil {
		return "", "", err
	}
	return kind.MIME.Value, kind.Extension, nil
}
