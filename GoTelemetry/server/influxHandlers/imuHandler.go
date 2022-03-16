package influxHandlers

import (
	telemetry "GoTelemetry/server/pb"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

func ImuDataAccel(tel *telemetry.TelemetryEvent, writeAPI api.WriteAPI) {
	measurement := "IMU_Test_Accel"
	eventUuid := uuid.NewV4().String()

	p := influxdb2.NewPointWithMeasurement(measurement).
		AddField("Event ID", eventUuid).
		AddField("Accel_X", tel.TelAcc.AccelX).
		AddField("Accel_Y", tel.TelAcc.AccelY).
		AddField("Accel_Z", tel.TelAcc.AccelZ).
		SetTime(time.Now())

	writeAPI.WritePoint(p)
	writeAPI.Flush()
	log.Println("Uploaded IMU Acceleration Data")
}

func ImuDataGyro(tel *telemetry.TelemetryEvent, writeAPI api.WriteAPI) {
	measurement := "IMU_Test_Gyro"
	eventUuid := uuid.NewV4().String()

	p := influxdb2.NewPointWithMeasurement(measurement).
		AddField("Event ID", eventUuid).
		AddField("Gyro_X", tel.TelGyro.GyroX).
		AddField("Gyro_Y", tel.TelGyro.GyroY).
		AddField("Gyro_Z", tel.TelGyro.GyroZ).
		SetTime(time.Now())

	writeAPI.WritePoint(p)
	writeAPI.Flush()
	log.Println("Uploaded IMU Gyro Data")
}
