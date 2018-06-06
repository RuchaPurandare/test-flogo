package timer

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"io/ioutil"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

var jsonMetadata = getJsonMetadata()

func getJsonMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

const testConfig3 string = `{
  "name": "tibco-wi-timer",
  "settings": {
  },
  "endpoints": [
    {
      "flowURI": "local://testFlow2",
      "settings": {
        "repeating": "false"
      }
    }
  ]
}`

const testConfig string = `{
  "name": "tibco-wi-timer",
  "settings": {
  },
  "endpoints": [
    {
      "flowURI": "local://testFlow2",
      "settings": {
        "repeating": "false",
        "startDate" : "2016-05-03T19:25:00Z-04:00"
      }
    }
  ]
}`

const testConfig2 string = `{
  "name": "tibco-wi-timer",
  "settings": {
  },
  "endpoints": [
    {
      "flowURI": "local://testFlow2",
      "settings": {
      	"notImmediate": "false",
        "repeating": "true",
        "seconds": "5"
      }
    }
  ]
}`

const testConfig4 string = `{
  "name": "tibco-wi-timer",
  "settings": {
  },
  "endpoints": [
    {
      "flowURI": "local://testFlow",
      "settings": {
        "repeating": "false",
        "startDate" : "05/01/2016, 12:25:01"
      }
    },
    {
      "flowURI": "local://testFlow2",
      "settings": {
        "repeating": "true",
        "startDate" : "05/01/2016, 12:25:01",
        "hours": "24"
      }
    },
    {
      "flowURI": "local://testFlow3",
      "settings": {
        "repeating": "true",
        "startDate" : "05/01/2016, 12:25:01",
        "minutes": "60"
      }
    },
    {
      "flowURI": "local://testFlow3",
      "settings": {
        "repeating": "true",
        "startDate" : "05/01/2016, 12:25:01",
        "seconds": "30"
      }
    }
  ]
}`

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Debugf("Ran Action: %v", uri)
	return 0, nil, nil
}

func (tr *TestRunner) RunAction(ctx context.Context, act action.Action, options map[string]interface{}) (results map[string]*data.Attribute, err error) {
	return nil, nil
}

func (tr *TestRunner) Execute(ctx context.Context, act action.Action, inputs map[string]*data.Attribute) (results map[string]*data.Attribute, err error) {
	return nil, nil
}

func TestRegistered(t *testing.T) {
	time := &TimerFactory{trigger.NewMetadata(jsonMetadata)}
	config := &trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	timer := time.New(config)
	if timer == nil {
		t.Error("Timer Trigger Not Registered")
		t.Fail()
		return
	}
}

func TestInit(t *testing.T) {
	time := &TimerFactory{trigger.NewMetadata(jsonMetadata)}
	config := &trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	timer := time.New(config)

	_, isNew := timer.(trigger.Initializable)

	if !isNew {
		runner := &TestRunner{}
		tgr, isOld := timer.(trigger.InitOld)
		if isOld {
			tgr.Init(runner)

		}
	}
}

func TestTimer(t *testing.T) {

	log.Debugf("TestTimer")
	timeFactory := &TimerFactory{trigger.NewMetadata(jsonMetadata)}
	config := &trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	timer := timeFactory.New(config)

	timer.Start()
	<-time.After(time.Millisecond * 2000)
	defer timer.Stop()

	log.Debug("Test timer done")
}
