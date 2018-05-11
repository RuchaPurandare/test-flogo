package concatactivity

import (
	"encoding/base64"
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
	tc.SetInput("firstString", "Hello")
	tc.SetInput("secondString", "World!")
	bytedata := []byte("test")
	data := base64.StdEncoding.EncodeToString(bytedata)
	tc.SetInput("password", data)
	tc.SetInput("separator", "#")
	tc.SetInput("fileSelector", "{\"content\":\"data:text/plain;base64,MS5XaGVyZSB3aWxsIHdlIHNldCB1cCBkb2NrZXI/DQoyLkFyZSB3ZSBkZXZlbG9waW5nIGNvbW11bml0eSBlZGl0aW9uIGNvbnRyaWJ1dGlvbj8NCjMuV2hpY2ggcGx1Z2luIHdlIGFyZSBjb252ZXJ0aW5nIHRvIGNvbnRpYnV0aW9uPw0KIA==\",\"filename\":\"questions.txt\"}")
	//tc.SetInput("headers", "&{{\"type\":\"object\",\"properties\":{\"sdsd\":{\"required\":\"true\",\"type\":\"array\",\"items\":{\"type\":\"string\"}},\"Accept\":{\"required\":\"false\",\"type\":\"string\"},\"Accept- Charset\":{\"required\":\"false\",\"type\":\"string\"}},\"required\":[\"sdsd\"]} {}}")
	_, err := act.Eval(tc)
	assert.Nil(t, err)
	result := tc.GetOutput("result")
	assert.Equal(t, result, "HelloWorld!test#questions.txt")
}
