package influxHandlers

import (
	telemetry "GoTelemetry/server/pb"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

func ShutdownData(tel *telemetry.TelemetryEvent, writeAPI api.WriteAPI) {
	measurement := "General_Shutdown"
	eventUuid := uuid.NewV4().String()

	p := influxdb2.NewPointWithMeasurement(measurement).
		AddField("Event ID", eventUuid).
		AddField("Shutdown", tel.ImproperShutdown).
		SetTime(time.Now())

	writeAPI.WritePoint(p)
	writeAPI.Flush()
	log.Println("Uploaded Shutdown Data")
}

func RobotOrientation(tel *telemetry.TelemetryEvent, writeAPI api.WriteAPI) {
	measurement := "General_Orientation"
	eventUuid := uuid.NewV4().String()
	num := int(tel.TelOrientation.Orientation.Number())
	p := influxdb2.NewPointWithMeasurement(measurement).
		AddField("Event ID", eventUuid).
		AddField("Orientation", num).
		SetTime(time.Now())

	writeAPI.WritePoint(p)
	writeAPI.Flush()
	log.Println("Uploaded Orientation Data")
}
