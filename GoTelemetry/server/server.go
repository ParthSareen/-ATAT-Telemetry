package main

import (
	telemetry "GoTelemetry/server/pb"
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
	addr, network          string
	dbAddr, dbUname, dbPwd string
)

func main() {
	network = "tcp"
	flag.StringVar(&addr, "e", ":10101", "service endpoint")
	flag.StringVar(&dbAddr, "r", "http://localhost:8086", "influxDB endpoint")
	flag.StringVar(&dbUname, "u", "admin", "influxDB username")
	flag.StringVar(&dbPwd, "p", "admin", "influxDB password")
	flag.Parse()

	ln, err := net.Listen(network, addr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer ln.Close()

	log.Printf("Telemetry Service Initialized: (%s) %s\n", network, addr)

	for {
		// TODO: Look into cleaning client setup code
		client := influxdb2.NewClient("http://localhost:8086", "")
		writeAPI := client.WriteAPI("380", "380")

		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			conn.Close()
			continue
		}
		log.Println("Connected to ", conn.RemoteAddr())
		go handleConnection(conn, writeAPI)
		//client.Close()
	}
}

func handleConnection(conn net.Conn, writeAPI api.WriteAPI) {
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
	go handleTelemetryData(&tel, writeAPI)

}

func handleTelemetryData(tel *telemetry.TelemetryEvent, writeAPI api.WriteAPI) {
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
	p := influxdb2.NewPointWithMeasurement("Test_Event").
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
	log.Println("Uploaded event")
}
