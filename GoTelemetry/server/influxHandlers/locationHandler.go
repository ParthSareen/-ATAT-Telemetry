package influxHandlers

import (
	telemetry "GoTelemetry/server/pb"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

func LocationData(tel *telemetry.TelemetryEvent, writeAPI api.WriteAPI) {
	measurement := "Location"
	eventUuid := uuid.NewV4().String()

	data := tel.TelLoc.Data
	rows := int(tel.TelLoc.Rows)
	cols := int(tel.TelLoc.Cols)
	total := rows * cols

	// TODO: Finish impl
	for index := 0; index < total; index++ {
		i := index / 6
		name := fmt.Sprintf("Location: Row %d, Col %d", i)
		p := influxdb2.NewPointWithMeasurement(measurement).
			AddField("Event ID", eventUuid).
			AddField(name, data[index]).
			SetTime(time.Now())

		writeAPI.WritePoint(p)
		writeAPI.Flush()
		log.Println("Uploaded Location Data")
	}

}
