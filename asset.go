package kasset

import (
	"git.kanosolution.net/kano/appkit"
	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
)

type Asset struct {
	orm.DataModelBase `json:"-" bson:"-"`
	ID                string `json:"_id" bson:"_id"`
	Title             string
	OriginalFileName  string
	NewFileName       string
	URI               string
	ContentType       string
	Size              int
	Tags              []string
	Kind              string
	RefID             string
	Data              map[string]string
}

func (a *Asset) TableName() string {
	return "AppAssets"
}

func (a *Asset) GetID(c dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{a.ID}
}

func (a *Asset) SetID(keys ...interface{}) {
	if len(keys) > 0 {
		a.ID = keys[0].(string)
	}
}

func (a *Asset) PreSave(c dbflex.IConnection) error {
	if a.ID == "" {
		a.ID = appkit.MakeID("", 32)
	}
	return nil
}
