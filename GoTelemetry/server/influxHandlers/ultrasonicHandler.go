package influxHandlers

import (
	telemetry "GoTelemetry/server/pb"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

func UltrasonicData(tel *telemetry.TelemetryEvent, writeAPI api.WriteAPI) {
	measurement := "Ultrasonic"
	eventUuid := uuid.NewV4().String()

	p := influxdb2.NewPointWithMeasurement(measurement).
		AddField("Event ID", eventUuid).
		AddField("US_Left", tel.TelUs.UsLeft).
		AddField("US_Front", tel.TelUs.UsFront).
		AddField("US_Back", tel.TelUs.UsBack).
		SetTime(time.Now())

	//writeAPI.WritePoint(p)
	//writeAPI.Flush()
	log.Println("Uploaded Ultrasonic Data", p)

	go ImuDataAccel(tel, writeAPI)
	go ImuDataGyro(tel, writeAPI)
	go RobotOrientation(tel, writeAPI)
}
