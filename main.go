package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ConnectionInfo represents an answer of the whoami command.
type YoulessRealtime struct {
	TotalPowerConsumption   string `json:"cnt"`
	CurrentPower            int    `json:"pwr"`
	TotalS0PowerConsumption string `json:"cs0"`
	CurrentS0Power          int    `json:"ps0"`
}

func main() {
	ip := flag.String("ip", `192.168.178.20`, "Youless ip")
	refreshInSeconds := flag.Duration("refreshInSeconds", 1, "How often to update in seconds")
	name := flag.String("name", "meter1", "Name of your meter")
	s0name := flag.String("s0name", "meter1s0", "Name of your meter")
	flag.Parse()

	var (
		totalPowerConsumption = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "youless_total_power_consumption",
			ConstLabels: prometheus.Labels{"name": *name},
		})
		currentPower = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "youless_current_power",
			ConstLabels: prometheus.Labels{"name": *name},
		})
		totalS0PowerConsumption = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "youless_total_s0_power_consumption",
			ConstLabels: prometheus.Labels{"name": *s0name},
		})
		currentS0Power = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "youless_current_s0_power",
			ConstLabels: prometheus.Labels{"name": *s0name},
		})
	)

	prometheus.MustRegister(totalPowerConsumption)
	prometheus.MustRegister(currentPower)
	prometheus.MustRegister(totalS0PowerConsumption)
	prometheus.MustRegister(currentS0Power)

	go func() {
		for {
			req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/a?f=j", *ip), nil)
			if err != nil {
				break
			}

			client := &http.Client{}
			response, err := client.Do(req)

			defer response.Body.Close()
			if err != nil {
				break
			}

			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				break
			}
			var data YoulessRealtime
			json.Unmarshal([]byte(contents), &data)

			totalPowerConsumptionParsedValued, _ := strconv.ParseFloat(strings.TrimSpace(strings.Replace(data.TotalPowerConsumption, ",", ".", 1)), 64)
			totalPowerConsumption.Set(totalPowerConsumptionParsedValued)

			currentPower.Set(float64(data.CurrentPower))

			totalS0PowerConsumptionParsedValued, err := strconv.ParseFloat(strings.TrimSpace(strings.Replace(data.TotalS0PowerConsumption, ",", ".", 1)), 64)
			totalS0PowerConsumption.Set(totalS0PowerConsumptionParsedValued)

			currentS0Power.Set(float64(data.CurrentS0Power))

			time.Sleep(*refreshInSeconds * time.Second)
		}
	}()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))

}
