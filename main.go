package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var countError int

type Data struct {
	LoadAverage  int
	RAM          int
	RAMUsage     int
	ROM          int
	ROMUsage     int
	Network      int
	NetworkUsage int
}

func main() {

	for {
		resp, err := http.Get("http://srv.msk01.gigacorp.local/_stats")
		if err != nil {
			errorResponce()
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			errorResponce()
			continue
		}
		if resp.StatusCode != 200 {
			errorResponce()
			continue
		}

		values := strings.Split(string(body), ",")
		if len(values) != 7 {
			errorResponce()
			continue
		}
		data, err := NewData(values)
		if err != nil {
			errorResponce()
			continue
		}
		data.check()
		countError = 0
	}

}

func NewData(values []string) (*Data, error) {
	LoadAverage, err := strconv.Atoi(values[0])
	if err != nil {
		return nil, err
	}
	RAM, err := strconv.Atoi(values[1])
	if err != nil {
		return nil, err
	}
	RAMUsage, err := strconv.Atoi(values[2])
	if err != nil {
		return nil, err
	}
	ROM, err := strconv.Atoi(values[3])
	if err != nil {
		return nil, err
	}
	ROMUsage, err := strconv.Atoi(values[4])
	if err != nil {
		return nil, err
	}
	Network, err := strconv.Atoi(values[5])
	if err != nil {
		return nil, err
	}
	NetworkUsage, err := strconv.Atoi(values[6])
	if err != nil {
		return nil, err
	}

	return &Data{
		LoadAverage:  LoadAverage,
		RAM:          RAM,
		RAMUsage:     RAMUsage,
		ROM:          ROM,
		ROMUsage:     ROMUsage,
		Network:      Network,
		NetworkUsage: NetworkUsage,
	}, nil

}

func errorResponce() {
	if countError >= 3 {
		fmt.Println("Unable to fetch server statistic")
	}
}

func (d Data) check() {
	if d.LoadAverage > 30 {
		fmt.Printf("Load Average is too high: %d\n", d.LoadAverage)
	}
	percentRamLoad := d.RAMUsage * 100 / d.RAM
	if percentRamLoad > 80 {
		msg := fmt.Sprintf("Memory usage too high: %d", percentRamLoad) + "%"
		fmt.Println(msg)
	}
	percentRomLoad := d.ROMUsage * 100 / d.ROM
	if percentRomLoad > 90 {
		fmt.Printf("Free disk space is too low: %d Mb left\n", (d.ROM-d.ROMUsage)/(1024*1024))
	}
	percentNetwork := d.NetworkUsage * 100 / d.Network
	if percentNetwork > 90 {
		fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n",
			(d.Network-d.NetworkUsage)/(1000*1000))
	}
}
