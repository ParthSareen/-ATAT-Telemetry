package main

import (
	telemetry "GoTelemetry/server/pb"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	uuid "github.com/satori/go.uuid"
	"log"
	"net"
	"os"
	"time"
)

var (
	addr, network          string
	dbAddr, dbUname, dbPwd string
	//influx         influxdb2.NewClient(url, token)
)

func main() {
	network = "tcp"
	flag.StringVar(&addr, "e", ":10101", "service endpoint")
	flag.StringVar(&dbAddr, "r", "http://localhost:8086", "influxDB endpoint")
	flag.StringVar(&dbUname, "u", "admin", "influxDB username")
	flag.StringVar(&dbPwd, "p", "admin", "influxDB password")
	flag.Parse()

	client := influxdb2.NewClient("http://localhost:8086", "kq-DUuEgrNFdBcTbb_AaAzJ5oBSWU5cOCYcLAnHU_oZV3T-4fTIzpk6salOSvQrKnz_de0rUXZcvfs3MGtmOlw==")

	ln, err := net.Listen(network, addr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer ln.Close()

	log.Printf("Telemetry Service Initialized: (%s) %s\n", network, addr)
	log.Printf("Uploading event")
	// get non-blocking write client
	writeAPI := client.WriteAPI("380", "380")
	// write line protocol
	uuid := uuid.NewV4().String()
	log.Println(uuid)
	tel := telemetry.TelemetryEvent{
		Timestamp: 0,
		UsFront:   0,
		UsLeft:    0,
		UsBack:    0,
		AccelX:    0,
		AccelY:    0,
		AccelZ:    0,
		GyroX:     0,
		GyroY:     0,
		GyroZ:     0,
	}
	// get non-blocking write client
	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": 24.5, "max": 45},
		time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// create point using fluent style
	p = influxdb2.NewPointWithMeasurement("stat").
		AddTag("unit", "temperature").
		AddField("avg", 23.2).
		AddField("max", 45).
		SetTime(time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()

	//myMeasurement,tag1=value1,tag2=value2 fieldKey="fieldValue" 1556813561098000000
	//writeAPI.WriteRecord(fmt.Sprintf("Event,EventID=\"test\", ax=%.2f,ay=%.2f,az=%.2f,gx=%.2f,gy=%.2f,gz=%.2f,us_back=%.2f,us_left=%.2f,us_front=%.2f %d",
	//	//uuid,
	//	tel.AccelX,
	//	tel.AccelY,
	//	tel.AccelZ,
	//	tel.GyroX,
	//	tel.GyroY,
	//	tel.GyroZ,
	//	tel.UsBack,
	//	tel.UsLeft,
	//	tel.UsFront,
	//	tel.Timestamp,
	//))
	// Close at end
	client.Close()
	log.Printf("Uploading event done")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			conn.Close()
			continue
		}
		log.Println("Connected to ", conn.RemoteAddr())
		go handleConnection(conn, client)
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
	uuid := uuid.NewV4()

	log.Println("Telemetry data:")
	fmt.Printf("{Timestamp:%d, EventID:%d, ax: %.2f, ay:%.2f, az:%.2f, gx:%.2f, gy:%.2f, gz:%.2f, us_back:%.2f, us_left:%.2f, us_front:%.2f}\n",
		tel.Timestamp,
		uuid,
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

	// get non-blocking write client
	writeAPI := client.WriteAPI("380", "380")
	// write line protocol
	writeAPI.WriteRecord(fmt.Sprintf("Timestamp=%d, EventID=%d, ax= %.2f, ay=%.2f, az=%.2f, gx=%.2f, gy=%.2f, gz=%.2f, us_back=%.2f, us_left=%.2f, us_front=%.2f",
		tel.Timestamp,
		uuid,
		tel.AccelX,
		tel.AccelY,
		tel.AccelZ,
		tel.GyroX,
		tel.GyroY,
		tel.GyroZ,
		tel.UsBack,
		tel.UsLeft,
		tel.UsFront,
	))
	// Flush writes
	writeAPI.Flush()

}

func handleTelemetryData() {

}
