package demo

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

const (
	ivField1 = "firstString"
	ivField2 = "secondString"
	ivField3 = "separator"
	// Adding a new input field
	ivField4 = "useSeparator"
	ovResult = "result"
)

var activityLog = logger.GetLogger("tibco-activity-demo")

//ConcatActivity is a concat activity struct
type ConcatActivity struct {
	metadata *activity.Metadata
}

//NewActivity is a concat activity function
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &ConcatActivity{metadata: metadata}
}

//Metadata is a concat activity function
func (a *ConcatActivity) Metadata() *activity.Metadata {
	return a.metadata
}

//Eval is a concat activity function
func (a *ConcatActivity) Eval(context activity.Context) (done bool, err error) {
	activityLog.Info("Executing Concat activity")
	if context.GetInput(ivField1) == nil {
		return false, activity.NewError("First string is not configured", "CONCAT-4001", nil)
	}
	field1v := context.GetInput(ivField1).(string)

	if context.GetInput(ivField2) == nil {
		return false, activity.NewError("Second string is not configured", "CONCAT-4002", nil)
	}
	field2v := context.GetInput(ivField2).(string)

	// Get the new boolean value if we need to use a separator or not
	field4v := context.GetInput(ivField4).(bool)

	// The validation has changed slightly
	if field4v && context.GetInput(ivField3) == nil {
		return false, activity.NewError("Separator is not configured", "CONCAT-4003", nil)
	}
	field3v := context.GetInput(ivField3).(string)

	// Set the output value depending on whether we need a separator or not
	if field4v {
		//Use separator in concatenation
		context.SetOutput(ovResult, field1v+field3v+field2v)
	} else {
		//No separator in concatenation
		context.SetOutput(ovResult, field1v+field2v)
	}
	return true, nil
}
