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
	measurement := "Motor_Test_Encoder"
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
	measurement := "Motor_Test_Speed"
	eventUuid := uuid.NewV4().String()

	// TODO Motor speed might need to be refactored with motor speed left and right
	p := influxdb2.NewPointWithMeasurement(measurement).
		AddField("Event ID", eventUuid).
		AddField("Motor Speed", tel.TelMotorSpeed.MotorSpeed).
		SetTime(time.Now())

	writeAPI.WritePoint(p)
	writeAPI.Flush()
	log.Println("Uploaded Motor Speed Data")
}
