package influxHandlers

import (
	telemetry "GoTelemetry/server/pb"
	"context"
	"github.com/golang/protobuf/proto"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log"
	"net"
)

func HandleReadBackup(conn net.Conn, event *telemetry.TelemetryEvent) {
	//defer func(conn net.Conn) {
	//	err := conn.Close()
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//}(conn)

	//results := getInfluxData(backup, client)

	//for resultType, result := range results {
	//	log.Println(resultType, result)
	//	_, err := conn.Write([]byte("return ok"))
	//	if err != nil {
	//		return
	//	}
	//}
	//telE := telemetry.TelemetryEvent{
	//	Timestamp: 0,
	//	UsFront:   0,
	//	UsLeft:    44,
	//	UsBack:    0,
	//	AccelX:    0,
	//	AccelY:    0,
	//	AccelZ:    0,
	//	GyroX:     0,
	//	GyroY:     0,
	//	GyroZ:     0,
	//}
	tempProto, err := proto.Marshal(event)
	if err != nil {
		panic(err)
	}
	log.Println("writing")
	_, err = conn.Write(tempProto)
	if err != nil {
		panic(err)
	}
	log.Println("written")

}

func getInfluxData(client influxdb2.Client) map[string]float64 {
	// TODO return a collection of data from here
	//s := string(backup.ReadFrom)
	// Get query client
	queryAPI := client.QueryAPI("380")
	queries := make(map[string]string)

	queries["US_Back"] = `from(bucket:"380")|> range(start: -1m) |> filter(fn: (r) => r._measurement == "Influx_Test_Event") |> filter(fn: (r) => r._field == "US_Back") |> last()`
	queries["US_Front"] = `from(bucket:"380")|> range(start: -1m) |> filter(fn: (r) => r._measurement == "Influx_Test_Event") |> filter(fn: (r) => r._field == "US_Front") |> last()`
	queries["US_Left"] = `from(bucket:"380")|> range(start: -1m) |> filter(fn: (r) => r._measurement == "Influx_Test_Event") |> filter(fn: (r) => r._field == "US_Left") |> last()`

	results := make(map[string]float64)
	// get QueryTableResult
	for queryType, query := range queries {
		log.Println(queryType)
		result, err := queryAPI.Query(context.Background(), query)
		if err != nil {
			panic(err)
		}

		// Iterate over query response TODO cleanup this as there is always only one response
		for result.Next() {
			// Notice when group key has changed
			//if result.TableChanged() {
			//	log.Printf("table: %s\n", result.TableMetadata().String())
			//}
			// Access data
			//log.Printf("value: %v\n", result.Record().Value())
		}
		log.Printf("%s: %v\n", queryType, result.Record().Value())
		// TODO refactor
		switch i := result.Record().Value().(type) {
		case float64:
			results[queryType] = i
		default:
			// Type panic
			panic(i)
		}
	}
	return results

	// check for an error
	//if result.Err() != nil {
	//	fmt.Printf("query parsing error: %\n", result.Err().Error())
	//}

}
