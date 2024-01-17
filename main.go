package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gosuri/uilive"
)

// TODO add keyboards controls (space pause/resume, enter new timer)
// TODO compatible with figlet

type CliOptions struct {
	downFlag     int
	upFlag       int
	countTarget  time.Time
	tickInterval int
}

func Usage() {
	fmt.Printf(`Usage: go-timer [s]

Count up or down to a optional specific time or duration and exit.

Options:
`)
	flag.PrintDefaults()
	fmt.Println(`
Examples:

  Count up to 10 seconds
    go-timer 10

  Count down 50 seconds
    go-timer -down 50`)
}

func getParsedFlags(args []string) *CliOptions {
	defaultTickInterval := 25
	options := &CliOptions{downFlag: 0, upFlag: 0, tickInterval: defaultTickInterval}

	// down an up with different targets in seconds or a date & time string
	flag.IntVar(&options.downFlag, "down", 0, "count down the given number of ms and exit")
	flag.IntVar(&options.upFlag, "up", 0, "count up till this number and exit")
	flag.Usage = Usage
	flag.Parse()

	// no up/down flags used but a number passed as argument
	if options.downFlag == 0 && options.upFlag == 0 && len(args) > 1 && args[1] != "" {
		parsedValue, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil {
			panic(err)
		}
		if parsedValue < 0 {
			options.downFlag = int(parsedValue)
		} else {
			options.upFlag = int(parsedValue)
		}
	}
	if options.downFlag > 0 {
		options.upFlag = 0
	}
	if options.upFlag > 0 || options.downFlag > 0 {
		options.countTarget = time.Now().Add(time.Duration(options.upFlag+options.downFlag) * time.Second)
	}

	return options
}

// https://github.com/gosuri/uilive/blob/master/example/main.go
func main() {
	startTime := time.Now()

	cliOptions := getParsedFlags(os.Args)

	writer := uilive.New()
	writer.Start()

	tick := time.Tick(time.Duration(cliOptions.tickInterval) * time.Millisecond)
	for range tick {
		beat(writer, cliOptions.upFlag > 0, startTime, cliOptions.countTarget)
	}

	writer.Stop()
}

func beat(writer *uilive.Writer, countUp bool, startTime, countTarget time.Time) {
	var duration time.Duration

	switch countTarget.IsZero() {
	case true:
		duration = time.Since(startTime)
	case countUp:
		duration = time.Since(startTime)
	default:
		duration = time.Until(countTarget)
	}

	fmt.Fprintf(writer, "%.3fs\r\n", duration.Seconds())

	// exit if countdown target is reached
	if !countTarget.IsZero() && time.Until(countTarget).Milliseconds() < 0 {
		os.Exit(0)
	}
}
