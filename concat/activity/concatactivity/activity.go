package concat

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

const (
	ivField1            = "firstString"
	ivField2            = "secondString"
	ivPasswordField     = "password"
	ivdropDownField     = "separator"
	ivFileSelectorField = "fileSelector"
	ivParamField        = "headers"
	ovResult            = "result"
)

var activityLog = logger.GetLogger("tibco-activity-concat")

type ConcatActivity struct {
	metadata *activity.Metadata
}

/*type Param struct {
	Name      string
	Type      string
	Repeating string
	Required  string
}*/

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &ConcatActivity{metadata: metadata}
}

func (a *ConcatActivity) Metadata() *activity.Metadata {
	return a.metadata
}

/*func convertToMap(data interface{}) (map[string]interface{}, error) {
	switch t := data.(type) {
	case string:
		if t != "" {
			m := map[string]interface{}{}
			err := json.Unmarshal([]byte(t), &m)
			if err != nil {
				return nil, err
			}
			return m, nil
		}
	case map[string]interface{}:
		return t, nil
	case interface{}:
		b, err := json.Marshal(t)
		if err != nil {
			return nil, err
		}
		m := map[string]interface{}{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return nil, err
		}
		return m, nil
	}

	return nil, nil
}
func LoadJsonSchemaFromMetadata(valueIN interface{}) (map[string]interface{}, error) {
	if valueIN != nil {
		complex := valueIN.(*data.ComplexObject)
		if complex != nil {
			params, err := convertToMap(complex.Metadata)
			if err != nil {
				return nil, err
			}
			return params, nil
		}
	}
	return nil, nil
}

func ParseParams(paramSchema map[string]interface{}) ([]Param, error) {

	if paramSchema == nil {
		return nil, nil
	}

	var parameter []Param

	//Structure expected to be JSON schema like
	props := paramSchema["properties"].(map[string]interface{})
	for k, v := range props {
		param := &Param{}
		param.Name = k
		propValue := v.(map[string]interface{})
		for k1, v1 := range propValue {
			if k1 == "required" {
				param.Required = v1.(string)
			} else if k1 == "type" {
				if v1 != "array" {
					param.Repeating = "false"
				}
				param.Type = v1.(string)
			} else if k1 == "items" {
				param.Repeating = "true"
				items := v1.(map[string]interface{})
				s, err := conversion.ConvertToString(items["type"])
				if err != nil {
					return nil, err
				}
				param.Type = s
			}
		}
		parameter = append(parameter, *param)
	}

	return parameter, nil
}*/

func (a *ConcatActivity) Eval(context activity.Context) (done bool, err error) {
	activityLog.Info("Executing Concat activity")
	//Read Inputs
	if context.GetInput(ivField1) == nil {
		// First string is not configured
		// return error to the engine
		return false, activity.NewError("First string is not configured", "CONCAT-4001", nil)
	}
	field1v := context.GetInput(ivField1).(string)

	if context.GetInput(ivField2) == nil {
		// Second string is not configured
		// return error to the engine
		return false, activity.NewError("Second string is not configured", "CONCAT-4002", nil)
	}
	field2v := context.GetInput(ivField2).(string)

	//Password field
	if context.GetInput(ivPasswordField) == nil {
		return false, activity.NewError("Passsword not available", "CONCAT-4003", nil)
	}
	ivPasswordField := context.GetInput(ivPasswordField).(string)
	if strings.HasSuffix(ivPasswordField, "=") {
		data, err := base64.StdEncoding.DecodeString(ivPasswordField)
		if err == nil {
			ivPasswordField = string(data)
		}
	}

	//Dropdown field
	if context.GetInput(ivdropDownField) == nil {
		return false, activity.NewError("Separator not available", "CONCAT-4004", nil)
	}
	ivdropDownField := context.GetInput(ivdropDownField).(string)

	//File selector field

	if context.GetInput(ivFileSelectorField) == nil {
		return false, activity.NewError("No file selected", "CONCAT-4005", nil)
	}
	ivFileSelectorField := context.GetInput(ivFileSelectorField).(string)

	data1, err1 := ioutil.ReadFile(ivFileSelectorField)
	if err1 != nil {
		fmt.Println("Can't read file:", ivFileSelectorField)
		panic(err1)
	}
	fmt.Println("File content is:")
	fmt.Println(string(data1))

	//Param field
	/*queryParamsMap, _ := LoadJsonSchemaFromMetadata(context.GetInput(ivParamField))
	if queryParamsMap != nil {
		queryParams, err2 := ParseParams(queryParamsMap)
		if err != nil {
			return params, err2
		}
	}*/
	//Set output
	context.SetOutput(ovResult, field1v+field2v+ivPasswordField+ivdropDownField+ivFileSelectorField)
	return true, nil
}
