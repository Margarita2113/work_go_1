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
		fmt.Println(fmt.Sprintf("Load Average is too high: %d", d.LoadAverage))
	}
	percentRamLoad := (float64(d.RAMUsage) / float64(d.RAM)) * 100
	if percentRamLoad > 80 {
		fmt.Println(fmt.Sprintf("Memory usage too high: %.f", percentRamLoad) + "%")
	}
	percentRomLoad := (float64(d.ROMUsage) / float64(d.ROM)) * 100
	if percentRomLoad > 90 {
		fmt.Println(fmt.Sprintf("Free disk space is too low: %d Mb left", (d.ROM-d.ROMUsage)/(1024*1024)))
	}
	percentNetwork := (float64(d.NetworkUsage) / float64(d.Network)) * 100
	if percentNetwork > 90 {
		fmt.Println(fmt.Sprintf("Network bandwidth usage high: %d Mbit/s available",
			(d.Network-d.NetworkUsage)/(1024*1024)))
	}
}
