package influxHandlers

import (
	telemetry "GoTelemetry/server/pb"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

func MotorDataEncoder(tel *telemetry.TelemetryEvent, writeAPI api.WriteAPI) {
	measurement := "Motor_Encoder"
	eventUuid := uuid.NewV4().String()

	p := influxdb2.NewPointWithMeasurement(measurement).
		AddField("Event ID", eventUuid).
		AddField("Left Motor Enc", tel.TelEnc.LeftMotor).
		AddField("Right Motor Enc", tel.TelEnc.RightMotor).
		SetTime(time.Now())

	writeAPI.WritePoint(p)
	writeAPI.Flush()
	log.Println("Uploaded Motor Encoder Data")
}

func MotorDataSpeed(tel *telemetry.TelemetryEvent, writeAPI api.WriteAPI) {
	measurement := "Motor_Speed_New2"
	eventUuid := uuid.NewV4().String()
	stuff := float32(tel.TelMotorSpeed.MotorSpeed)
	log.Println(stuff, measurement, eventUuid)
	// TODO Motor speed might need to be refactored with motor speed left and right
	p := influxdb2.NewPointWithMeasurement(measurement).
		AddField("Event ID", eventUuid).
		AddField("Motor Speed", stuff).
		SetTime(time.Now())

	writeAPI.WritePoint(p)
	writeAPI.Flush()
	log.Println("Uploaded Motor Speed Data")
}
