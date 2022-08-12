package kasset

import (
	"fmt"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/kaos"
	"github.com/google/uuid"
)

type AssetAPIEngine struct {
}

type ReferenceRequest struct {
	RefType string
	RefID   string
}

func (ae *AssetAPIEngine) FindByRefID(ctx *kaos.Context, req *ReferenceRequest) ([]*Asset, error) {
	res := []*Asset{}
	h, e := ctx.DefaultHub()
	if e != nil {
		return res, e
	}

	ars := []*AssetReference{}
	if e = h.Gets(new(AssetReference),
		dbflex.NewQueryParam().SetWhere(dbflex.And(dbflex.Eq("RefType", req.RefType), dbflex.Eq("RefID", req.RefID))),
		&ars); e != nil {
		return res, fmt.Errorf("unable to get reference. %s", e.Error())
	}

	for _, ar := range ars {
		a := new(Asset)
		a.ID = ar.AssetID
		if e = h.Get(a); e != nil {
			return res, fmt.Errorf("unable to get asset %s: %s", ar.AssetID, e.Error())
		}
		res = append(res, a)
	}

	return res, nil
}

func (ae *AssetAPIEngine) MakeRef(ctx *kaos.Context, req *AssetReference) (string, error) {
	h, e := ctx.DefaultHub()
	if e != nil {
		return "", e
	}

	// validate asset
	a := new(Asset)
	a.ID = req.AssetID
	if e = h.Get(a); e != nil {
		return "", e
	}

	// write the ref
	if req.ID == "" {
		req.ID = uuid.New().String()
	}
	if e = h.Save(req); e != nil {
		return "", e
	}
	return req.ID, nil
}
