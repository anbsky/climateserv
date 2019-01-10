package climateserv

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tarm/serial"
	"github.com/urfave/cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type AirQualityData struct {
	Timestamp   time.Time `json:"timestamp"`
	PM25        float64   `json:"pm25"`
	PM10        float64   `json:"pm10"`
	AQIPM25     int       `json:"aqi_pm25"`
	AQICategory string    `json:"aqi_category"`
}

var samples = make([]AirQualityData, 10)

func main() {
	var device_path string
	var port string
	var interval string

	app := cli.NewApp()
	app.Name = "ClimateServ"
	app.Version = fmt.Sprintf("%v, commit %v, built at %v", version, commit, date)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "device, d",
			Usage:       "SDS device path (i.e. /dev/ttyUSB0 on linux, /dev/tty.usbserial-1430 on mac)",
			Destination: &device_path,
		},
		cli.StringFlag{
			Name:        "port, p",
			Usage:       "HTTP port to listen on",
			Destination: &port,
			Value:       "2510",
		},
		cli.StringFlag{
			Name:        "interval, i",
			Usage:       "Interval between device polling (in seconds)",
			Destination: &interval,
		},
	}
	app.Action = func(c *cli.Context) error {
		interval_numeric, err := strconv.Atoi(interval)
		if err != nil {
			log.Fatal("interval argument has to be a number")
		}
		go serve(port)
		start_reading(device_path, interval_numeric)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func open_device(device_path string) *serial.Port {
	serial_config := &serial.Config{Name: device_path, Baud: 9600}
	serial_port, err := serial.OpenPort(serial_config)
	if err != nil {
		log.Fatal(err)
	}
	return serial_port
}

func start_reading(device_path string, interval int) {
	for {
		entry, err := read_device(device_path)
		if err != nil {
			log.Printf("%s (PM2.5: %.2f, PM10: %.2f)\n", err, entry.PM25, entry.PM10)
			continue
		}
		log.Printf("PM2.5: %.2f, PM10: %.2f\n", entry.PM25, entry.PM10)
		samples = append(samples, entry)
		samples = samples[1:]
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func read_device(device_path string) (entry AirQualityData, err error) {
	device := open_device(device_path)
	buffer := make([]byte, 10)
	_, err = device.Read(buffer)
	device.Close()

	if err != nil {
		return AirQualityData{PM25: 0.0, PM10: 0.0, Timestamp: time.Now().UTC()}, err
	}
	PM25 := parse_and_convert(buffer[2], buffer[3])
	PM10 := parse_and_convert(buffer[4], buffer[5])
	AQIPM25, category := CurrentPM25toAQI(PM25)
	entry = AirQualityData{
		PM25:        PM25,
		PM10:        PM10,
		AQIPM25:     AQIPM25,
		AQICategory: category,
		Timestamp:   time.Now().UTC()}
	if entry.PM25 > 400 || entry.PM10 > 400 {
		err = errors.New("crazy value read from the port")
	}
	return entry, err
}

func parse_and_convert(raw_value_low byte, raw_value_high byte) float64 {
	low, _ := strconv.ParseInt(fmt.Sprintf("%d", raw_value_low), 10, 32)
	high, _ := strconv.ParseInt(fmt.Sprintf("%d", raw_value_high), 10, 32)
	return ((float64(high) * 256) + float64(low)) / 10
}

func serve(port string) {
	http.HandleFunc("/api/v1/air_quality/samples", handle_samples_view)
	http.HandleFunc("/api/v1/air_quality/current", handle_current_view)
	http.ListenAndServe(":"+port, nil)
}

func reply_with_json(payload []byte, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/javascript")
	fmt.Fprint(writer, string(payload))
}

func handle_samples_view(writer http.ResponseWriter, request *http.Request) {
	serialized_samples, _ := json.Marshal(samples)
	reply_with_json(serialized_samples, writer)
}

func handle_current_view(writer http.ResponseWriter, request *http.Request) {
	serialized_current, _ := json.Marshal(samples[len(samples)-1])
	reply_with_json(serialized_current, writer)
}
