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
		// TODO: Look into cleaning client setup code
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
		go handleConnection(conn, client)
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
	// TODO add case to send data back
	//if n ...

	var tel telemetry.TelemetryEvent
	if err := proto.Unmarshal(buf[:n], &tel); err != nil {
		//log.Println("failed to unmarshal:", err)
		//return
		// TODO cleanup not returning error here
	}
	if err == nil {
		writeAPI := client.WriteAPI("380", "380")
		go handleTelemetryData(&tel, writeAPI, "Test_Event")

	}

	var backup telemetry.ReadBackup
	if err := proto.Unmarshal(buf[:n], &backup); err != nil {
		log.Println("failed to unmarshal:", err)
		return
	}

	go getInfluxData(&backup, client)

}

func handleTelemetryData(tel *telemetry.TelemetryEvent, writeAPI api.WriteAPI, measurement string) {
	eventUuid := uuid.NewV4().String()
	log.Println("Telemetry data:")
	log.Printf("{Timestamp:%d, EventID:%s, ax: %.2f, ay:%.2f, az:%.2f, gx:%.2f, gy:%.2f, gz:%.2f, us_back:%.2f, us_left:%.2f, us_front:%.2f}\n",
		tel.Timestamp,
		eventUuid,
		tel.AccelX,
		tel.AccelY,
		tel.AccelZ,
		tel.GyroX,
		tel.GyroY,
		tel.GyroZ,
		tel.UsBack,
		tel.UsLeft,
		tel.UsFront,
	)
	p := influxdb2.NewPointWithMeasurement(measurement).
		//AddTag("unit", "temperature").
		AddField("Event ID", eventUuid).
		AddField("US_Left", tel.UsLeft).
		AddField("US_Front", tel.UsFront).
		AddField("US_Back", tel.UsBack).
		SetTime(time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()
	// error handling
	//err := writeAPI.Errors()
	//if err != nil {
	//	panic(err)
	//}

	log.Println("Uploaded event")
}

func getInfluxData(backup *telemetry.ReadBackup, client influxdb2.Client) {
	// TODO return a collection of data from here
	//s := string(backup.ReadFrom)
	// Get query client
	queryAPI := client.QueryAPI("380")
	queries := make(map[string]string)

	queries["US_Back"] = `from(bucket:"380")|> range(start: -1m) |> filter(fn: (r) => r._measurement == "Influx_Test_Event") |> filter(fn: (r) => r._field == "US_Back") |> last()`
	//query_us_front := `from(bucket:"380")|> range(start: -%sm) |> filter(fn: (r) => r._measurement == "Influx_Test_Event") |> filter(fn: (r) => r._field == "US_Front")`
	results := make(map[string]string)
	// get QueryTableResult
	for queryType, query := range queries {
		result, err := queryAPI.Query(context.Background(), query)
		if err != nil {
			panic(err)
		}

		// Iterate over query response
		for result.Next() {
			// Notice when group key has changed
			if result.TableChanged() {
				log.Printf("table: %s\n", result.TableMetadata().String())
			}
			// Access data
			log.Printf("value: %v\n", result.Record().Value())
		}
	}

	// check for an error
	//if result.Err() != nil {
	//	fmt.Printf("query parsing error: %\n", result.Err().Error())
	//}

}
