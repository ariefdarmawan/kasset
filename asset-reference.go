package kasset

import (
	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
)

type AssetReference struct {
	orm.DataModelBase `json:"-" bson:"-"`
	ID                string `json:"_id" bson:"_id"`
	AssetID           string `json:"assetid"`
	RefType           string `json:"reftype"`
	RefID             string `json:"refid"`
	Feature           string `json:"feature"`
}

func (ar *AssetReference) TableName() string {
	return "assetreferences"
}

func (ar *AssetReference) GetID(_ dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{ar.ID}
}

func (ar *AssetReference) SetID(keys ...interface{}) {
	ar.ID = keys[0].(string)
}
