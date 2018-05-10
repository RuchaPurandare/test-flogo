package concat

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
	tc.SetInput("fileSelector", "C:/Users/rpuranda/Desktop/questions.txt")
	//tc.SetInput("headers", "[{\"parameterName\":\"sddd\",\"type\":\"number\",\"repeating\":\"false\",\"required\":\"false\",\"visible\":true},{\"parameterName\":\"Accept- Charset\",\"type\":\"string\",\"repeating\":\"false\",\"required\":\"false\",\"visible\":false}]")
	_, err := act.Eval(tc)
	assert.Nil(t, err)
	result := tc.GetOutput("result")
	assert.Equal(t, result, "HelloWorld!test#"+"C:/Users/rpuranda/Desktop/questions.txt")
}
