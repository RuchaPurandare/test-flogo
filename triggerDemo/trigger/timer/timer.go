package timer

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/carlescere/scheduler"
)

// log is the default package logger
var log = logger.GetLogger("general-trigger-timer")

type TimerTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
	timers   map[string]*scheduler.Job
	myId     string
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &TimerFactory{metadata: md}
}

// TimerFactory Timer Trigger factory
type TimerFactory struct {
	metadata *trigger.Metadata
}

// Metadata implements trigger.Trigger.Metadata
func (t *TimerTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

//New Creates a new trigger instance for a given id
func (t *TimerFactory) New(config *trigger.Config) trigger.Trigger {
	return &TimerTrigger{metadata: t.metadata, config: config}
}

// Init implements ext.Trigger.Init
func (t *TimerTrigger) Init(runner action.Runner) {
	log.Debugf("In init, id: '%s', Metadata: '%+v', Config: '%+v'", t.myId, t.metadata, t.config)
	t.runner = runner
}

// Start implements ext.Trigger.Start
func (t *TimerTrigger) Start() error {

	log.Debug("Start")
	t.timers = make(map[string]*scheduler.Job)
	endpoints := t.config.Handlers

	log.Debug("Processing endpoints")
	for _, endpoint := range endpoints {

		repeatingIn := endpoint.Settings["Repeating"]
		var repeating bool
		if repeatingIn != nil {
			switch repeatingIn.(type) {
			case bool:
				repeating = repeatingIn.(bool)
			default:
				repeat, err := strconv.ParseBool(repeatingIn.(string))
				if err != nil {
					log.Errorf("Parser %s to bool error %s", repeatingIn, err.Error())
					repeating = false
				} else {
					repeating = repeat
				}
			}
		}

		log.Debug("Repeating: ", repeating)
		if repeating {
			interval := endpoint.Settings["Time Interval"]
			if interval == nil {
				log.Error("Time Interval must be specified for Timer trigger")
				return errors.New("Time Interval must be specified for Timer trigger")
			}

			intervalUnit := endpoint.Settings["Interval Unit"]
			if intervalUnit == nil {
				log.Error("Time Interval must be specified for Timer trigger")
				return errors.New("Interval unit must be selected for Timer trigger")
			}
			err := t.scheduleRepeating(endpoint)
			if err != nil {
				return err
			}
		} else {
			err := t.scheduleOnce(endpoint)
			if err != nil {
				return err
			}
		}
		log.Debug("Settings repeating: ", endpoint.Settings["Repeating"])
		log.Debugf("Processing Handler: %s", endpoint.ActionId)
	}

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *TimerTrigger) Stop() error {

	log.Debug("Stopping endpoints")
	for k, v := range t.timers {
		if t.timers[k].IsRunning() {
			log.Debug("Stopping timer for : ", k)
			v.Quit <- true
		} else {
			log.Debugf("Timer: %s is not running", k)
		}
	}

	return nil
}

func (t *TimerTrigger) scheduleOnce(endpoint *trigger.HandlerConfig) error {
	log.Info("Scheduling a run once job")

	seconds := getInitialStartInSeconds(endpoint)
	log.Debug("Seconds till trigger fires: ", seconds)

	var interval int = 0
	var err error
	tInterval := endpoint.Settings["Time Interval"]
	if tInterval != nil {
		tInterval, err = data.CoerceToValue(tInterval, data.TypeInteger)
		if err != nil {
			log.Errorf("Invalid Time Interval [%s]. ", tInterval)
			return err
		} else {
			interval = tInterval.(int)
			if interval < 0 {
				err = fmt.Errorf("Invalid Time Interval [%d]. Must not be a negative value. ", tInterval)
				log.Error(err.Error())
				return err
			}
		}
	}

	intervalUnit := endpoint.Settings["Interval Unit"]

	switch intervalUnit {
	case "Minute":
		interval = seconds + interval*60
	case "Hour":
		interval = seconds + interval*60*60
	case "Day":
		interval = seconds + interval*60*60*24
	case "Week":
		interval = seconds + interval*60*60*24*7
	}

	log.Debug("Seconds: ", interval)

	timerJob := scheduler.Every(int(interval))

	if timerJob == nil {
		err = fmt.Errorf("Failed to create timer job")
		log.Error(err.Error())
		return err
	}

	fn := func() {
		log.Debug("-- Starting \"Once\" timer process")

		action := action.Get(endpoint.ActionId)
		log.Debugf("Found action: '%+x'", action)
		log.Debugf("ActionID: '%s'", endpoint.ActionId)
		_, _, err := t.runner.Run(context.Background(), action, endpoint.ActionId, nil)
		if err != nil {
			log.Error("Error starting action: ", err.Error())
		}
		timerJob.Quit <- true
	}

	if interval == 0 {
		//Run now
		go fn()
	} else {
		timerJob, err = timerJob.Seconds().NotImmediately().Run(fn)
		if err != nil {
			log.Error("Error scheduleOnce flow err: ", err.Error())
			return err
		}
	}

	return nil
}

func (t *TimerTrigger) scheduleRepeating(endpoint *trigger.HandlerConfig) error {
	log.Info("Scheduling a repeating job")

	//	seconds := getInitialStartInSeconds(endpoint)

	fn2 := func() {
		log.Debug("-- Starting \"Repeating\" (repeat) timer action")

		action := action.Get(endpoint.ActionId)
		log.Debugf("Found action: '%+x'", action)
		log.Debugf("ActionID: '%s'", endpoint.ActionId)
		_, _, err := t.runner.Run(context.Background(), action, endpoint.ActionId, nil)
		if err != nil {
			log.Error("Error starting flow: ", err.Error())
		}
	}

	timerJob, err := t.scheduleJobEverySecond(endpoint, fn2)

	if timerJob == nil {
		err = fmt.Errorf("Failed to create timer job")
		log.Error(err.Error())
		return err
	}

	if err != nil {
		log.Error("Failed to schedule a job:" + err.Error())
		return err
	}

	t.timers["r:"+endpoint.ActionId] = timerJob
	return nil

}

func getInitialStartInSeconds(endpoint *trigger.HandlerConfig) int {

	startDate, _ := data.CoerceToString(endpoint.Settings["Start Date"])
	if startDate == "" {
		return 0
	}

	layout := time.RFC3339
	idx := strings.LastIndex(startDate, "Z")
	timeZone := startDate[idx+1 :]
	log.Debug("Time Zone: ", timeZone)
	startDate = strings.TrimSuffix(startDate, timeZone)
	log.Debug("Start Date: ", startDate)

	// is timezone negative
	var isNegative bool
	isNegative = strings.HasPrefix(timeZone, "-")
	// remove sign
	timeZone = strings.TrimPrefix(timeZone, "-")

	triggerDate, err := time.Parse(layout, startDate)
	if err != nil {
		log.Error("Error parsing time err: ", err.Error())
	}
	log.Debug("Time parsed from settings: ", triggerDate)

	var hour int
	var minutes int

	sliceArray := strings.Split(timeZone, ":")
	if len(sliceArray) != 2 {
		log.Error("Time zone has wrong format: ", timeZone)
	} else {
		hour, _ = strconv.Atoi(sliceArray[0])
		minutes, _ = strconv.Atoi(sliceArray[1])

		log.Debug("Duration hour: ", time.Duration(hour)*time.Hour)
		log.Debug("Duration minutes: ", time.Duration(minutes)*time.Minute)
	}

	hours, _ := strconv.Atoi(timeZone)
	log.Debug("Hours: ", hours)
	if isNegative {
		log.Debug("Adding to triggerDate")
		triggerDate = triggerDate.Add(time.Duration(hour) * time.Hour)
		triggerDate = triggerDate.Add(time.Duration(minutes) * time.Minute)
	} else {
		log.Debug("Subtracting to triggerDate")
		triggerDate = triggerDate.Add(time.Duration(hour * -1))
		triggerDate = triggerDate.Add(time.Duration(minutes))
	}

	currentTime := time.Now().UTC()
	log.Debug("Current time: ", currentTime)
	log.Debug("Setting start time: ", triggerDate)
	duration := time.Since(triggerDate)

	return int(math.Abs(duration.Seconds()))
}

type PrintJob struct {
	Msg string
}

func (j *PrintJob) Run() error {
	log.Debug(j.Msg)
	return nil
}

func (t *TimerTrigger) scheduleJobEverySecond(endpoint *trigger.HandlerConfig, fn func()) (*scheduler.Job, error) {

	var interval int = 0
	var err error
	tInterval := endpoint.Settings["Time Interval"]
	if tInterval != nil {
		tInterval, err = data.CoerceToValue(tInterval, data.TypeInteger)
		if err != nil {
			log.Errorf("Invalid Time Interval [%s]. ", tInterval)
			return nil, err
		} else {
			interval = tInterval.(int)
			if interval < 0 {
				err = fmt.Errorf("Invalid Time Interval [%d]. Must not be a negative value. ", tInterval)
				log.Error(err.Error())
				return nil, err
			}
		}
	}

	intervalUnit := endpoint.Settings["Interval Unit"]

	switch intervalUnit {
	case "Minute":
		interval = interval * 60
	case "Hour":
		interval = interval * 60 * 60
	case "Day":
		interval = interval * 60 * 60 * 24
	case "Week":
		interval = interval * 60 * 60 * 24 * 7
	}

	log.Debug("Repeating seconds: ", interval)
	// schedule repeating
	timerJob, err := scheduler.Every(interval).Seconds().Run(fn)

	return timerJob, err
}
