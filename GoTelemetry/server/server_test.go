package main

import (
	telemetry "GoTelemetry/server/pb"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"os"
	"testing"
)

func TestInfluxUpload(t *testing.T) {
	tel := telemetry.TelemetryEvent{
		Timestamp: 0,
		UsFront:   5,
		UsLeft:    20,
		UsBack:    30,
		AccelX:    0,
		AccelY:    0,
		AccelZ:    0,
		GyroX:     0,
		GyroY:     0,
		GyroZ:     0,
	}
	dbAddr = "http://localhost:8086"
	authToken = os.Getenv("INFLUX_TOKEN")
	client := influxdb2.NewClient(dbAddr, authToken)
	writeAPI := client.WriteAPI("380", "380")
	handleTelemetryData(&tel, writeAPI, "Influx_Test_Event")
	t.Log("Successful upload")
	defer client.Close()
}

func TestInfluxDataGet(t *testing.T) {
	dbAddr = "http://localhost:8086"
	client := influxdb2.NewClient(dbAddr, "kq-DUuEgrNFdBcTbb_AaAzJ5oBSWU5cOCYcLAnHU_oZV3T-4fTIzpk6salOSvQrKnz_de0rUXZcvfs3MGtmOlw==")
	backup := telemetry.ReadBackup{ReadFrom: 2}
	getInfluxData(&backup, client)
	defer client.Close()

}
