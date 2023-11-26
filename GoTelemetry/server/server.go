package main

import (
	handler "GoTelemetry/server/influxHandlers"
	telemetry "GoTelemetry/server/pb"
	"flag"
	"time"

	//"fmt"
	"github.com/golang/protobuf/proto"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log"
	"net"
	"os"
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
	//authToken = os.Getenv("INFLUX_TOKEN")
	// TODO: Remove
	authToken = "kq-DUuEgrNFdBcTbb_AaAzJ5oBSWU5cOCYcLAnHU_oZV3T-4fTIzpk6salOSvQrKnz_de0rUXZcvfs3MGtmOlw=="
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
		go handleConnection(conn, client)
		// TODO probably should close this at some point
		//client.Close()
	}
}

func writeToFile(tel *telemetry.TelemetryEvent) {
	file, err := os.OpenFile("gyro2.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := file.Write([]byte("\n\n" + time.Now().String() + " " + tel.String())); err != nil {
		log.Fatal(err)
	}
	if err := file.Close(); err != nil {
		log.Fatal(err)
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
	//TODO: add timeout
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

	// TODO: need to make fault tolerant
	log.Println(tel.String())
	switch tel.TelCmd {
	case telemetry.TelemetryEvent_CMD_READ_DATA:
		handler.HandleReadBackup(conn, &tel)
	case telemetry.TelemetryEvent_CMD_ULTRASONIC:
		go handler.UltrasonicData(&tel, writeAPI)
		go writeToFile(&tel)
	case telemetry.TelemetryEvent_CMD_ACCELERATION:
		handler.ImuDataAccel(&tel, writeAPI)
	case telemetry.TelemetryEvent_CMD_GYRO:
		handler.ImuDataGyro(&tel, writeAPI)
	case telemetry.TelemetryEvent_CMD_ENCODER:
		handler.MotorDataEncoder(&tel, writeAPI)
	case telemetry.TelemetryEvent_CMD_MOTOR_SPEED:
		handler.MotorDataSpeed(&tel, writeAPI)
	case telemetry.TelemetryEvent_CMD_ORIENTATION:
		handler.RobotOrientation(&tel, writeAPI)
	case telemetry.TelemetryEvent_CMD_LOCATION:
		handler.LocationData(&tel, writeAPI)
	case telemetry.TelemetryEvent_CMD_SHUTDOWN:
		handler.ShutdownData(&tel, writeAPI)
	default:
		log.Println("Error, no command matched")
	}
	// Do not run read in goroutine as it will not have enough time to read
	//handleReadBackup(conn, tel)
	return

}
