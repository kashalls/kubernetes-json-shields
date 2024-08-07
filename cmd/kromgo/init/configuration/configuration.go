package configuration

import (
	"flag"
	"os"
	"time"

	"github.com/kashalls/kromgo/cmd/kromgo/init/log"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Config struct for configuration environmental variables
type Config struct {
	ServerHost         string        `env:"SERVER_HOST" envDefault:"localhost"`
	ServerPort         int           `env:"SERVER_PORT" envDefault:"8888"`
	ServerReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT"`
	Prometheus         string        `yaml:"prometheus,omitempty" json:"prometheus,omitempty"`
	Metrics            []Metric      `yaml:"metrics" json:"metrics"`
}

type Metric struct {
	Name   string        `yaml:"name" json:"name"`
	Query  string        `yaml:"query" json:"query"`
	Label  string        `yaml:"label,omitempty" json:"label,omitempty"`
	Prefix string        `yaml:"prefix,omitempty" json:"prefix,omitempty"`
	Suffix string        `yaml:"suffix,omitempty" json:"suffix,omitempty"`
	Colors []MetricColor `yaml:"colors,omitempty" json:"colors,omitempty"`
}

type MetricColor struct {
	Min           float64 `yaml:"min" json:"min"`
	Max           float64 `yaml:"max" json:"max"`
	Color         string  `yaml:"color,omitempty" json:"color,omitempty"`
	ValueOverride string  `yaml:"valueOverride,omitempty" json:"valueOverride,omitempty"`
}

var configPath = "/kromgo/config.yaml" // Default config file path
var ProcessedMetrics map[string]Metric

// Init sets up configuration by reading set environmental variables
func Init() Config {

	// Check if a custom config file path is provided via command line argument
	configPathFlag := flag.String("config", "", "Path to the YAML config file")
	flag.Parse()
	if *configPathFlag != "" {
		configPath = *configPathFlag
	}

	// Read file from path.
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Error("error reading config file", zap.Error(err))
		os.Exit(1)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Error("error unmarshalling config yaml", zap.Error(err))
		os.Exit(1)
	}

	ProcessedMetrics = preprocess(config.Metrics)
	return config
}

func preprocess(metrics []Metric) map[string]Metric {
	reverseMap := make(map[string]Metric)
	for _, obj := range metrics {
		reverseMap[obj.Name] = obj
	}
	return reverseMap
}
