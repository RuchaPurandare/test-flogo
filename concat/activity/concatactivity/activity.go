package concatactivity

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"git.tibco.com/git/product/ipaas/wi-contrib.git/engine/conversion"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

const (
	ivField1            = "firstString"
	ivField2            = "secondString"
	ivPasswordField     = "password"
	ivdropDownField     = "separator"
	ivFileSelectorField = "fileSelector"
	ivParamField        = "headers"
	ivJSONField         = "responseBody"
	ovResult            = "result"
)

var activityLog = logger.GetLogger("tibco-activity-concat")

type Param struct {
	Name      string
	Type      string
	Repeating string
	Required  string
}
type TypedValue struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type Parameters struct {
	Headers []*TypedValue `json:"headers"`
}

type ConcatActivity struct {
	metadata *activity.Metadata
}

func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &ConcatActivity{metadata: metadata}
}

func (a *ConcatActivity) Metadata() *activity.Metadata {
	return a.metadata
}

func convertToMap(data interface{}) (map[string]interface{}, error) {
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
}

func GetComplexValueAsMap(context activity.Context, name string) (map[string]interface{}, error) {
	valueIN := context.GetInput(name)
	if valueIN != nil {
		complex := valueIN.(*data.ComplexObject)
		if complex != nil {
			switch t := complex.Value.(type) {
			case string:
				m := map[string]interface{}{}
				err := json.Unmarshal([]byte(t), &m)
				if err != nil {
					return nil, err
				}
				return m, nil
			default:
				return convertToMap(complex.Value)

			}
		}
	}
	return nil, nil
}

func GetParameter(context activity.Context) (params *Parameters, err error) {
	params = &Parameters{}
	//Headers
	headersMap, _ := LoadJsonSchemaFromMetadata(context.GetInput("headers"))
	fmt.Println("HeaderMap", headersMap)
	if headersMap != nil {
		headers, err := ParseParams(headersMap)
		//fmt.Println("Headers", headers)
		if err != nil {
			return params, err
		}

		if headers != nil {
			inputHeaders, err := GetComplexValueAsMap(context, "headers")
			//fmt.Println("inputHeaders", inputHeaders)
			if err != nil {
				return params, err
			}
			var typeValuesHeaders []*TypedValue
			for _, hParam := range headers {
				isRequired := hParam.Required
				paramName := hParam.Name
				if isRequired == "true" && inputHeaders[paramName] == nil {
					return nil, fmt.Errorf("Required header parameter [%s] is not configured.", paramName)
				}
				if inputHeaders[paramName] != nil {
					if hParam.Repeating == "true" {
						val := inputHeaders[paramName]
						switch reflect.TypeOf(val).Kind() {
						case reflect.Slice:
							s := reflect.ValueOf(val)
							for i := 0; i < s.Len(); i++ {
								typeValue := &TypedValue{}
								typeValue.Name = paramName
								typeValue.Value = s.Index(i).Interface()
								typeValue.Type = hParam.Type
								typeValuesHeaders = append(typeValuesHeaders, typeValue)
							}
						}
					} else {
						typeValue := &TypedValue{}
						typeValue.Name = paramName
						typeValue.Value = inputHeaders[paramName]
						typeValue.Type = hParam.Type
						typeValuesHeaders = append(typeValuesHeaders, typeValue)
					}
					params.Headers = typeValuesHeaders
				}
			}
		}
	}
	return params, err
}
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

	//Parse json file and get file name
	var result map[string]interface{}
	json.Unmarshal([]byte(ivFileSelectorField), &result)
	var fileName = ""
	if result != nil {
		fmt.Println(result["filename"])
		fileName = result["filename"].(string)
	}

	//Param field

	//ivParamField := context.GetInput(ivParamField).(*data.ComplexObject)
	//fmt.Println(ivParamField)

	parameters, err2 := GetParameter(context)
	if err2 != nil {
		fmt.Println("Err is:", err2)
		return false, err2
	}
	//fmt.Println("Params are:", parameters.Headers)
	for _, value := range parameters.Headers {
		fmt.Println("Param Name is:", value.Name)
		fmt.Println("Param value is:", value.Value)
		fmt.Println("Param Type is:", value.Type)
	}

	//Complex object- JSON
	ivJSONField := context.GetInput(ivJSONField).(*data.ComplexObject).Value

	//fmt.Println("JSON output:", ivJSONField.Value)
	/*Sample json schema:{"key":["string"]}
	Sample input: array.create("abc", "abc", "67")
	*/
	for key, value := range ivJSONField.(map[string]interface{}) {
		activityLog.Infof("Key Value of JSON param -- %s : %s", key, value)
		var buffer bytes.Buffer
		fmt.Println("TypeOfObject:", reflect.TypeOf(value))
		for k1, v1 := range value.([]interface{}) {
			fmt.Printf("key[%s] value[%s]\n", k1, v1)
			buffer.WriteString(v1.(string))
			buffer.WriteString(",")
		}
		fmt.Println(strings.TrimRight(buffer.String(), ","))
	}
	//Set output
	context.SetOutput(ovResult, field1v+field2v+ivPasswordField+ivdropDownField+fileName)
	return true, nil
}
