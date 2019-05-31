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
type Configuration struct {
	Ip               *string
	RefreshInSeconds *int
	Name, NameS0     *string
}

func main() {
	config := Configuration{}
	config.Ip = flag.String("ip", "127.0.0.1", "Youless device ip")
	config.RefreshInSeconds = flag.Int("refreshInSeconds", 1, "How often to update in seconds")
	config.Name = flag.String("name", "meter1", "Name of your meter")
	config.NameS0 = flag.String("s0name", "meter1s0", "Name of your s0 meter")

	flag.Parse()

	var (
		totalPowerConsumption = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "youless_total_value",
			ConstLabels: prometheus.Labels{"name": *config.Name},
		})
		currentPower = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "youless_current_value",
			ConstLabels: prometheus.Labels{"name": *config.Name},
		})
		totalS0PowerConsumption = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "youless_s0_total_value",
			ConstLabels: prometheus.Labels{"name": *config.NameS0},
		})
		currentS0Power = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "youless_s0_current_value",
			ConstLabels: prometheus.Labels{"name": *config.NameS0},
		})
	)

	prometheus.MustRegister(totalPowerConsumption)
	prometheus.MustRegister(currentPower)
	prometheus.MustRegister(totalS0PowerConsumption)
	prometheus.MustRegister(currentS0Power)

	go func() {
		for {
			for {
				req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/a?f=j", *config.Ip), nil)
				if err != nil {
					log.Println("Request building failure: ", err)
					break
				}

				client := &http.Client{}
				response, err := client.Do(req)

				if err != nil {
					log.Println("Connection failure: ", err)
					break
				}

				contents, err := ioutil.ReadAll(response.Body)
				response.Body.Close()
				if err != nil {
					log.Println("Reading data failure: ", err)
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

				time.Sleep(time.Duration(*config.RefreshInSeconds) * time.Second)
			}
			time.Sleep(time.Duration(*config.RefreshInSeconds) * time.Second)
		}
	}()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))

}
