package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"time"
)

var portfolio *Portfolio

type Check struct {
	Name             string  `json:"name"`
	URL              string  `json:"url"`
	Frequency        string  `json:"frequency"`
	IsUp             bool    `json:"is_up"`
	UpTime           float64 `json:"uptime"`
	AvgResponseTime  int64   `json:"avg_response_time"`
	ApdexT           string  `json:"apdex_threshold"`
	Apdex            float64 `json:"apdex"`
	SatisfiedCount   float64 `json:"satisfied_count"`
	ToleratingCount  float64 `json:"tolerating_count"`
	FrustratingCount float64 `json:"frustrating_count"`
	ErrorCount       float64 `json:"error_count"`
	RequestCount     float64 `json:"request_count"`
}

type Portfolio struct {
	checks []Check
}

func NewPortfolio(configFileName string) *Portfolio {
	var config []Check
	configFile, err := os.Open(configFileName)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	p := &Portfolio{config}
	return p
}

func (p *Portfolio) Start() {
	fmt.Printf("Checks: %d = %v\n", len(p.checks), p.checks)
	for i, check := range p.checks {
		log.Printf("Initiating check %d: '%s' [%s] (%s)\n", i+1, check.Name, check.URL, check.Frequency)
		go p.checks[i].monitor()
	}
}

func (c *Check) monitor() {
	for {
		c.ping()

		nextCheck, _ := time.ParseDuration(c.Frequency)
		time.Sleep(nextCheck)
	}
}

func (c *Check) ping() {
	timeout := time.Duration(60 * time.Second)

	transport := http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		},
	}
	client := http.Client{
		Transport: &transport,
	}

	start := time.Now()
	resp, err := client.Head(c.URL)
	elapsed := time.Now().UnixNano() - start.UnixNano()

	c.RequestCount++

	if err != nil {
		if c.IsUp {
			log.Printf("DOWN '%s' [%s] with error: %v\n", c.Name, c.URL, err.Error())
		}
		c.ErrorCount++
		c.updateMetrics()
		c.IsUp = false
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if c.IsUp {
			log.Printf("DOWN '%s' [%s] with code: %d\n", c.Name, c.URL, resp.StatusCode)
		}
		c.ErrorCount++
		c.updateMetrics()
		c.IsUp = false
		return
	}

	if !c.IsUp {
		log.Printf("UP '%s' [%s]\n", c.Name, c.URL)
	}

	threshold, err := time.ParseDuration(c.ApdexT)

	if err != nil {
		threshold = time.Duration(2 * time.Second)
		c.ApdexT = threshold.String()
		log.Printf("Setting Apdex threshold for '%s' to %s\n", c.Name, c.ApdexT)
	}

	c.AvgResponseTime = (c.AvgResponseTime*(int64(c.RequestCount)-1) + elapsed) / int64(c.RequestCount)

	if elapsed < threshold.Nanoseconds() {
		c.SatisfiedCount++
	} else if elapsed < (threshold.Nanoseconds() * 4) {
		c.ToleratingCount++
	} else {
		c.FrustratingCount++
	}

	c.updateMetrics()
	c.IsUp = true
}

func (c *Check) updateMetrics() {
	c.Apdex = (c.SatisfiedCount + (c.ToleratingCount / 2.0)) / c.RequestCount
	c.UpTime = math.Floor(((c.RequestCount - c.ErrorCount) / c.RequestCount) * 100)
}
