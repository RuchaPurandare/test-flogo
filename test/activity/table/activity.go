package table

import (
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

const (
	ivTableField = "My Table"
	ovJSONField  = "tableoutput"
)

var activityLog = logger.GetLogger("tibco-activity-concat")

type TypedValue struct {
	Name  string `json:"fieldName"`
	Value bool   `json:"selected"`
}

type TableActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &TableActivity{metadata: metadata}
}

func (a *TableActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func (a *TableActivity) Eval(context activity.Context) (done bool, err error) {
	activityLog.Info("Executing Concat activity")

	//table field

	ivTableField := context.GetOutput(ivTableField).([]interface{})
	var terms []string
	var objectArray []TypedValue
	for _, obj := range ivTableField {
		//fmt.Println(obj)
		//obj["FieldName"]
		var typeObj TypedValue
		for k, v := range obj.(map[string]string) {
			//fmt.Println("Key is", k, " and Value is", v)
			if k == "FieldName" {
				typeObj.Name = v
				terms = append(terms, v)
			} else {
				if v == "true" {
					typeObj.Value = true
				} else {
					typeObj.Value = false
				}
			}
		}
		objectArray = append(objectArray, typeObj)
		//keys = append(keys, key)
	}

	fmt.Println(objectArray)

	ovJSONField := context.GetOutput(ovJSONField).(*data.ComplexObject).Value
	//context.SetOutput(ovResult, terms[0])
	fmt.Println(ovJSONField)
	return true, nil
}
