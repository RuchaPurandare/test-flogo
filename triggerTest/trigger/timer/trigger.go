package sample

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// Create a logger for the Sample Trigger
// Logger Name : <category>-trigger-<type>
var triggerLog = logger.GetLogger("demo-trigger-sample")

// Trigger must define a struct
type SampleTrigger struct {
	// Trigger Metadata
	metadata *trigger.Metadata
	// Flow action that would create a flow instance
	runner action.Runner
	// Trigger configuration
	config *trigger.Config
}

// Sample Trigger factory
// Trigger must define a factory
type SampleTriggerFactory struct {
	// Trigger Metadata
	metadata *trigger.Metadata
}

// NewFactory create a new Trigger factory
// Trigger must define this function
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &SampleTriggerFactory{metadata: md}
}

// Creates a new trigger instance for a given id
// Trigger must define this method
func (t *SampleTriggerFactory) New(config *trigger.Config) trigger.Trigger {
	return &SampleTrigger{metadata: t.metadata, config: config}
}

// Returns trigger metadata
// Trigger must define this method
func (t *SampleTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Set flow runner
// Trigger must define this method
func (t *SampleTrigger) Init(runner action.Runner) {
	// Flow runner must be stored
	t.runner = runner
}

// Start trigger. Start will be called once engine is started successfully.
func (t *SampleTrigger) Start() error {
	return nil
}

// Stop trigger. Stop will be called in case engine is gracefully stopped.
func (t *SampleTrigger) Stop() error {
	return nil
}
