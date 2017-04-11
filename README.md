# metro
adelaide metro realtime transport timetable go api

usage:
`metro.FetchStop(stopID string, previewTimeSeconds string)` returns `RealTimeResponse`

example:
```
response := metro.FetchStop("12624", "120")
for _, stop := range response.StopMonitoringDelivery {
  for _, visted := range stop.MonitoredStopVisit {
	  fmt.Println(visted.MonitoredVehicleJourney.Lineref.Value, "-", visted.MonitoredVehicleJourney.MonitoredCall.LatestExpectedArrivalTime)
  }
}
```

stopID lookup coming
