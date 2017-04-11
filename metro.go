package metro

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const baseURL = "http://realtime.adelaidemetro.com.au/SiriWebServiceSAVM/SiriStopMonitoring.svc/json/SM?"

// RealTimeResponse is the root of the json dict
type RealTimeResponse struct {
	StopMonitoringDelivery []StopMonitoringDelivery
}

// StopMonitoringDelivery is your requests responses
type StopMonitoringDelivery struct {
	ResponseTimeStamp  Timestamp
	MonitoredStopVisit []MonitoredStopVisit
}

// MonitoredStopVisit is a monitoring struct for a bus/train
type MonitoredStopVisit struct {
	MonitoringRef           Ref
	MonitoredVehicleJourney MonitoredVehicleJourney
}

// MonitoredVehicleJourney is the vechicle infomation. LineRef is bus number
type MonitoredVehicleJourney struct {
	Lineref       Ref
	DirectionRef  Ref
	MonitoredCall MonitoredCall
}

// MonitoredCall is the estimated and aimed arrival time data for a ourney
type MonitoredCall struct {
	StopPointRef              Ref
	AimedArrivalTime          Timestamp
	LatestExpectedArrivalTime Timestamp
}

// Ref holds a value as string
type Ref struct {
	Value string
}

// Timestamp wraps timer.Time for json Unmarshal
// because the time received by the timer is in a C# timestamp format
type Timestamp struct {
	time.Time
}

// UnmarshalJSON unmarshals a C# string timestamp into a time.Time
func (t *Timestamp) UnmarshalJSON(value []byte) error {

	// remove quotes
	s := strings.Trim(string(value), "\"")
	// check null json value
	if s == "null" {
		t.Time = time.Time{}
		return nil
	}

	// extract timestamp and timezone
	var ts, tz int64
	_, err := fmt.Sscanf(s, "\\/Date(%d+%d)\\/", &ts, &tz)
	if err != nil {
		// only timestamp present
		_, err = fmt.Sscanf(s, "\\/Date(%d)\\/", &ts)
		if err != nil {
			// invalid format
			t.Time = time.Time{}
			return nil
		}
	}

	// convert time to seconds from milliseconds
	// go time.Time automatically uses current timezone
	t.Time = time.Unix(ts/1000, 0)

	return nil
}

// FetchStop ...
func FetchStop(stopID string, future string) RealTimeResponse {
	// fetch json from siri server
	response, err := http.Get(fmt.Sprintf("%sMonitoringRef=%s&PreviewInterval=%s", baseURL, stopID, future))
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Unmarshal json string to our struct schema
	var realtime RealTimeResponse
	err = json.Unmarshal(body, &realtime)
	if err != nil {
		log.Fatalln(err)
	}

	return realtime
}
