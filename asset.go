package kasset

import (
	"time"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
	"github.com/sebarcode/codekit"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Data              codekit.M
	LastUpdated       time.Time
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
		a.ID = primitive.NewObjectID().Hex()
	}
	a.LastUpdated = time.Now()
	return nil
}

func (a *Asset) Indexes() []dbflex.DbIndex {
	return []dbflex.DbIndex{
		{Name: "KindRefIndex", IsUnique: false, Fields: []string{"Kind", "RefID"}},
		{Name: "OriginalFileNameIndex", IsUnique: false, Fields: []string{"OriginalFileName"}},
		{Name: "NewFileName", IsUnique: false, Fields: []string{"NewFileName"}},
		{Name: "TagIndex", IsUnique: false, Fields: []string{"Tags"}},
	}
}
