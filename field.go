/*
@Author: Weny Xu
@Date: 2021/06/02 22:51
*/

package odm

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
	"time"
)

type IDField struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" odm:"primaryID"`
}

// PrepareID method prepares the ID value to ObjectID
func (f *IDField) PrepareID(id interface{}) (interface{}, error) {
	if idStr, ok := id.(string); ok {
		return primitive.ObjectIDFromHex(idStr)
	}
	//TODO
	return id, nil
}

// GetID method returns a model's ID
func (f *IDField) GetID() interface{} {
	return f.ID
}

// SetID sets the value of a model's ID field.
func (f *IDField) SetID(id interface{}) {
	f.ID = id.(primitive.ObjectID)
}

const (
	tag            = "odm"
	autoCreateTime = "autoCreateTime"
	autoUpdateTime = "autoUpdateTime"
	deleteTime     = "deleteTime"
	primaryID      = "primaryID"
)

type TimestampFields struct {
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty" odm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty" odm:"autoUpdateTime"`
}

type DeletedAtField struct {
	DeletedAt time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty" odm:"deleteTime"`
}

type FieldsConfig struct {
	AutoCreateTimeField string
	AutoUpdateTimeField string
	DeleteTimeField     string
	PrimaryIDField      string
}

func (info FieldsConfig) SoftDeletable() bool {
	return info.DeleteTimeField != ""
}

type FieldsInfo []FieldInfo

func (info FieldsInfo) genModelConfig() *FieldsConfig {
	mc := &FieldsConfig{}
	for _, f := range info {
		if strings.Contains(f.RawTag, autoCreateTime) {
			mc.AutoCreateTimeField = f.Name
		} else if strings.Contains(f.RawTag, autoUpdateTime) {
			mc.AutoUpdateTimeField = f.Name
		} else if strings.Contains(f.RawTag, deleteTime) {
			mc.DeleteTimeField = f.Name
		} else if strings.Contains(f.RawTag, primaryID) {
			mc.PrimaryIDField = f.Name
		}
	}
	return mc
}

type FieldInfo struct {
	Name   string
	RawTag string
	Tags   []string
	Index  []int
}

func genModelFieldsInfo(data interface{}) *FieldsInfo {
	rt := reflect.TypeOf(data)
	for rt.Kind() == reflect.Ptr ||
		rt.Kind() == reflect.Interface {
		rt = rt.Elem()
	}
	var i []int
	var fi FieldsInfo
	getModelFieldsInfo(rt, &fi, i)
	return &fi
}

func getModelFieldsInfo(t reflect.Type, fieldsInfo *FieldsInfo, index []int) {
	for i := 0; i < t.NumField(); i++ {
		// increase
		field := t.Field(i)
		tag := field.Tag.Get(tag)
		if tag != "" {
			*fieldsInfo = append(*fieldsInfo, FieldInfo{
				RawTag: tag,
				Tags:   strings.Split(tag, ","),
				Index:  append(index, i),
				Name:   field.Name,
			})
		}
		if field.Type.Kind() == reflect.Struct {
			getModelFieldsInfo(field.Type, fieldsInfo, append(index, i))
		}
	}
}
