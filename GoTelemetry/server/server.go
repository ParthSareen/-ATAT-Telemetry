package main

import (
	telemetry "GoTelemetry/server/pb"
	"context"
	"flag"
	//"fmt"
	"github.com/golang/protobuf/proto"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	uuid "github.com/satori/go.uuid"
	"log"
	"net"
	"os"
	"time"
)

var (
	addr, network string
	dbAddr        string
	authToken     string
)

func main() {
	network = "tcp"
	flag.StringVar(&addr, "e", ":10101", "service endpoint")
	flag.StringVar(&dbAddr, "r", "http://localhost:8086", "influxDB endpoint")
	//flag.StringVar(&dbUname, "u", "admin", "influxDB username")
	//flag.StringVar(&dbPwd, "p", "admin", "influxDB password")
	flag.Parse()
	authToken = os.Getenv("INFLUX_TOKEN")

	ln, err := net.Listen(network, addr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {

		}
	}(ln)

	log.Printf("Telemetry Service Initialized: (%s) %s\n", network, addr)

	for {
		client := influxdb2.NewClient(dbAddr, authToken)

		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			err := conn.Close()
			if err != nil {
				return
			}
			continue
		}
		log.Println("Connected to ", conn.RemoteAddr())
		// TODO look into how locking will work if needed + pass down
		handleConnection(conn, client)
		// TODO probably should close this at some point lol
		//client.Close()
	}
}

func handleConnection(conn net.Conn, client influxdb2.Client) {
	defer func() {
		log.Println("INFO: closing connection")
		if err := conn.Close(); err != nil {
			log.Println("error closing connection:", err)
		}
	}()

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}
	if n <= 0 {
		log.Println("no data received")
		return
	}

	var tel telemetry.TelemetryEvent
	if err := proto.Unmarshal(buf[:n], &tel); err != nil {
		log.Println("failed to unmarshal:", err)
		return
	}
	writeAPI := client.WriteAPI("380", "380")
	switch tel.TelCmd {
	case telemetry.TelemetryEvent_CMD_READ_DATA:
		handleReadBackup(conn, &tel)
	case telemetry.TelemetryEvent_CMD_ULTRASONIC:
		uploadUltrasonicData(&tel, writeAPI)
	case telemetry.TelemetryEvent_CMD_ACCELERATION:
		uploadImuDataAccel(&tel, writeAPI)
	}
	//log.Println("Tel proto read")
	// Do not run read in goroutine as it will not have enough time to read
	//handleReadBackup(conn, tel)
	return

}

func uploadUltrasonicData(tel *telemetry.TelemetryEvent, writeAPI api.WriteAPI) {
	measurement := "Ultrasonic_Test"
	eventUuid := uuid.NewV4().String()

	p := influxdb2.NewPointWithMeasurement(measurement).
		AddField("Event ID", eventUuid).
		AddField("US_Left", tel.TelUs.UsLeft).
		AddField("US_Front", tel.TelUs.UsFront).
		AddField("US_Back", tel.TelUs.UsBack).
		SetTime(time.Now())

	writeAPI.WritePoint(p)
	writeAPI.Flush()
	log.Println("Uploaded Ultrasonic Data")
}

func uploadImuDataAccel(tel *telemetry.TelemetryEvent, writeAPI api.WriteAPI) {
	measurement := "IMU_Test"
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

func handleReadBackup(conn net.Conn, event *telemetry.TelemetryEvent) {
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
