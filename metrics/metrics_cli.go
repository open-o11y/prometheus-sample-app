package metrics

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Host string `yaml:"Host"`
	Port string `yaml:"Port"`

	Type      string `yaml:"Type"`
	Count     int    `yaml:"Count"`
	Frequency int    `yaml:"Frequency"`
	Random    bool   `yaml:"Random"`
}

type CommandLine struct{}

func (c *Config) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "healthy")
}

// initConnection handles the metric creation and also updates the metrics via go routines
// The delegation logic is handled here
func (cli *CommandLine) initConnection(str string, count int, freq int, isRandom bool, conf Config) {

	rand.Seed(time.Now().Unix())
	mc := newMetricCollector()
	mc.interval = freq
	addressPtr := flag.String("listen_address", conf.Host+":"+conf.Port, "server listening address")
	address := *addressPtr
	switch str {
	case "counter":
		cli.createCounter(count, mc)
	case "gauge":
		cli.createGauge(count, mc)
	case "histogram":
		cli.createHistogram(count, mc)
	case "summary":
		cli.createSummary(count, mc)
	case "all":
		cli.createAll(count, mc, isRandom)
	default:
		log.Println("Invalid type")
	}
	log.Println("Serving on address: " + address)
	if isRandom {
		log.Println("Producing randomized metrics per type")

	} else {
		log.Println("Producing " + fmt.Sprintf("%d", count) + " metric(s) per type")
	}
	log.Println("Updating at a frequency of " + fmt.Sprintf("%d", mc.interval) + " seconds")
	http.HandleFunc("/", healthCheckHandler)
	http.Handle("/metrics", promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}))

	log.Fatal(http.ListenAndServe(address, nil))
}

func (cli *CommandLine) createCounter(count int, mc metricCollector) {
	mc.registerCounter(count)
	go mc.updateCounter()
}

func (cli *CommandLine) createGauge(count int, mc metricCollector) {
	mc.registerGauge(count)
	go mc.updateGauge()
}

func (cli *CommandLine) createHistogram(count int, mc metricCollector) {
	mc.registerHistogram(count)
	go mc.updateHistogram()

}

func (cli *CommandLine) createSummary(count int, mc metricCollector) {
	mc.registerSummary(count)
	go mc.updateSummary()
}

// createAll generates all 4 metric types
// If isRandom is sent as true, createAll will generate randomized metrics. Otherwise createALl will steadily create the 4 types of metrics with a fixed count (provided by the user
func (cli *CommandLine) createAll(count int, mc metricCollector, isRandom bool) {

	if isRandom {
		idx := rand.Intn(4)
		lower := 1
		upper := 4
		amount := rand.Intn(upper-lower) + lower
		metrics := []string{"counter", "gauge", "histogram", "summary"}
		rands := []int{rand.Intn(200), rand.Intn(200), rand.Intn(200), rand.Intn(200)}
		max := rands[0]
		for i := 0; i < len(rands); i++ {
			if rands[i] > max {
				max = rands[i]
			}
		}
		for i := 0; i <= amount; i++ {
			if idx >= len(metrics) {
				idx = 0
			}
			str := metrics[idx]
			idx++
			switch str {
			case "counter":
				cli.createCounter(rands[0], mc)
			case "gauge":
				cli.createGauge(rands[1], mc)
			case "histogram":
				cli.createHistogram(rands[2], mc)
			case "summary":
				cli.createSummary(rands[3], mc)
			}
		}

	} else {
		mc.registerCounter(count)
		mc.registerGauge(count)
		mc.registerHistogram(count)
		mc.registerSummary(count)
		go mc.updateMetrics(count)

	}

}

// Run reads the config file and uses the data as default arguments.
// These arguments can be overriden by CLI input (see README)
func (cli *CommandLine) Run() {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	var conf Config
	if err := conf.Parse(data); err != nil {
		log.Fatal(err)
	}
	generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)

	metricType := generateCmd.String("metric_type", conf.Type, "Type of metric (counter, gauge, histogram, summary)")
	metricCount := generateCmd.Int("metric_count", conf.Count, "Amount of metrics to create")
	metricFreq := generateCmd.Int("metric_frequency", conf.Frequency, "Refresh interval in seconds")
	rand := generateCmd.Bool("is_random", conf.Random, "Metrics specification")

	if len(os.Args) > 2 {
		err := generateCmd.Parse(os.Args[3:])
		if err != nil {
			log.Panic(err)
		}
	}

	cli.initConnection(*metricType, *metricCount, *metricFreq, *rand, conf)

}
