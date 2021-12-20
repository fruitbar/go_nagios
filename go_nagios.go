package nagios

import (
	"fmt"
	"os"
	"strings"
)

type NagiosStatusVal int

// The values with which a Nagios check can exit
const (
	NAGIOS_OK NagiosStatusVal = iota
	NAGIOS_WARNING
	NAGIOS_CRITICAL
	NAGIOS_UNKNOWN
)

// Maps the NagiosStatusVal entries to output strings
var (
	valMessages = []string{
		"OK:",
		"WARNING:",
		"CRITICAL:",
		"UNKNOWN:",
	}
)

//--------------------------------------------------------------
// A type representing a Nagios check status. The Value is a the exit code
// expected for the check and the Message is the specific output string.
type NagiosStatus struct {
	Message  string
	Value    NagiosStatusVal
	Perfdata NagiosPerformanceVal
}

// Take a bunch of NagiosStatus pointers and find the highest value, then
// combine all the messages. Things win in the order of highest to lowest.
// Combines messages as well
func (status *NagiosStatus) Aggregate(otherStatuses []*NagiosStatus) {
	perfFormat := "'%s'=%s%s;%s;%s;%s;%s"
	msgFormat := "%s %s"
	perfDataStringArr := []string{}
	perfDataString := ""
	longMessageArr := []string{}
	longMessage := ""
	msg := ""
	for _, s := range otherStatuses {
		fmt.Printf("Status: %d Message: %s PerfData: "+perfFormat+"\n", s.Value, s.Message,
			s.Perfdata.Label,
			s.Perfdata.Value,
			s.Perfdata.Uom,
			s.Perfdata.WarnThreshold,
			s.Perfdata.CritThreshold,
			s.Perfdata.MinValue,
			s.Perfdata.MaxValue)
		if status.Value < s.Value {
			status.Value = s.Value
		}
		if status.Message == "" {
			msg = fmt.Sprintf(msgFormat+" | "+perfFormat,
				valMessages[s.Value],
				s.Message,
				s.Perfdata.Label,
				s.Perfdata.Value,
				s.Perfdata.Uom,
				s.Perfdata.WarnThreshold,
				s.Perfdata.CritThreshold,
				s.Perfdata.MinValue,
				s.Perfdata.MaxValue)
		} else {
			longMessageArr = append(longMessageArr, fmt.Sprintf(msgFormat,
				valMessages[s.Value],
				s.Message))
			perfDataStringArr = append(perfDataStringArr, fmt.Sprintf(perfFormat,
				s.Perfdata.Label,
				s.Perfdata.Value,
				s.Perfdata.Uom,
				s.Perfdata.WarnThreshold,
				s.Perfdata.CritThreshold,
				s.Perfdata.MinValue,
				s.Perfdata.MaxValue))
		}
	}
	if len(longMessageArr) > 0 {
		longMessage = strings.Join(longMessageArr, "\n")
		msg += "\n" + longMessage
	}
	if len(perfDataStringArr) > 0 {
		perfDataString = strings.Join(perfDataStringArr, "\n")
		msg += " | " + perfDataString
	}
	status.Message = msg

}

// Construct the Nagios message
func (status *NagiosStatus) constructedNagiosMessage() string {
	return valMessages[status.Value] + " " + status.Message
}

// NagiosStatus: Issue a Nagios message to stdout and exit with appropriate Nagios code
func (status *NagiosStatus) NagiosExit() {
	fmt.Fprintln(os.Stdout, status.constructedNagiosMessage())
	os.Exit(int(status.Value))
}

//--------------------------------------------------------------
// A type representing a Nagios performance data value.
// https://nagios-plugins.org/doc/guidelines.html#AEN200
// http://docs.pnp4nagios.org/pnp-0.6/about#system_requirements
type NagiosPerformanceVal struct {
	Label         string
	Value         string
	Uom           string
	WarnThreshold string
	CritThreshold string
	MinValue      string
	MaxValue      string
}

//--------------------------------------------------------------
// A type representing a Nagios check status and performance data.
type NagiosStatusWithPerformanceData struct {
	*NagiosStatus
	Perfdata NagiosPerformanceVal
}

// Construct the Nagios message with performance data
func (status *NagiosStatusWithPerformanceData) constructedNagiosMessage() string {
	msg := fmt.Sprintf("%s %s | '%s'=%s%s;%s;%s;%s;%s",
		valMessages[status.Value],
		status.Message,
		status.Perfdata.Label,
		status.Perfdata.Value,
		status.Perfdata.Uom,
		status.Perfdata.WarnThreshold,
		status.Perfdata.CritThreshold,
		status.Perfdata.MinValue,
		status.Perfdata.MaxValue)
	return msg
}

// Issue a Nagios message (with performance data) to stdout and exit with appropriate Nagios code
func (status *NagiosStatusWithPerformanceData) NagiosExit() {
	fmt.Fprintln(os.Stdout, status.constructedNagiosMessage())
	os.Exit(int(status.Value))
}

//--------------------------------------------------------------

// Exit with an UNKNOWN status and appropriate message
func Unknown(output string) {
	ExitWithStatus(&NagiosStatus{output, NAGIOS_UNKNOWN, NagiosPerformanceVal{}})
}

// Exit with an CRITICAL status and appropriate message
func Critical(err error) {
	ExitWithStatus(&NagiosStatus{err.Error(), NAGIOS_CRITICAL, NagiosPerformanceVal{}})
}

// Exit with an WARNING status and appropriate message
func Warning(output string) {
	ExitWithStatus(&NagiosStatus{output, NAGIOS_WARNING, NagiosPerformanceVal{}})
}

// Exit with an OK status and appropriate message
func Ok(output string) {
	ExitWithStatus(&NagiosStatus{output, NAGIOS_OK, NagiosPerformanceVal{}})
}

// Exit with a particular NagiosStatus
func ExitWithStatus(status *NagiosStatus) {
	status.NagiosExit()
}
