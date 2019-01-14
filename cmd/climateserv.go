package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/sayplastic/climateserv"
	"github.com/urfave/cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

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
		go climateserv.Serve(port)
		climateserv.StartReading(device_path, interval_numeric)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
