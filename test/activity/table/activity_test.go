package table

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/stretchr/testify/assert"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {
	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}
		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}
	return activityMetadata
}

func TestActivityRegistration(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(act.Metadata())
	//setup attrs

	var m1 = make(map[string]string)
	m1["FieldName"] = "Option1"
	m1["Selected"] = "true"

	var m2 = make(map[string]string)
	m2["FieldName"] = "Option2"
	m2["Selected"] = "false"

	var v = []interface{}{m1, m2}

	tc.SetInput("My Table", v)

	//tc.SetInput("headers", "&{{\"type\":\"object\",\"properties\":{\"sdsd\":{\"required\":\"true\",\"type\":\"array\",\"items\":{\"type\":\"string\"}},\"Accept\":{\"required\":\"false\",\"type\":\"string\"},\"Accept- Charset\":{\"required\":\"false\",\"type\":\"string\"}},\"required\":[\"sdsd\"]} {}}")
	_, err := act.Eval(tc)
	assert.Nil(t, err)
	result := tc.GetOutput("result")
	assert.Equal(t, result, "Option1")
}
