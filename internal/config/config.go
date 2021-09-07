package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type EndpointConfiguration struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (s *EndpointConfiguration) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type DatabaseConfiguration struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type KafkaConfiguration struct {
	Topic   string   `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
}

type PrometheusConfiguration struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
}

type BatchConfiguration struct {
	Size uint `yaml:"size"`
}

type ApplicationConfiguration struct {
	Name string `yaml:"name"`
}

type Configuration struct {
	Application *ApplicationConfiguration `yaml:"application"`
	Grpc        *EndpointConfiguration    `yaml:"grpc"`
	Gateway     *EndpointConfiguration    `yaml:"gateway"`
	Db          *DatabaseConfiguration    `yaml:"db"`
	Batch       *BatchConfiguration       `yaml:"batch"`
	Jaeger      *EndpointConfiguration    `yaml:"jaeger"`
	Kafka       *KafkaConfiguration       `yaml:"kafka"`
	Prometheus  *PrometheusConfiguration  `yaml:"prometheus"`
}

func LoadConfiguration(path string) (*Configuration, error) {
	var file *os.File
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if defErr := file.Close(); defErr != nil {
			err = defErr
		}
	}()

	confContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	confContent = []byte(os.ExpandEnv(string(confContent)))
	conf := &Configuration{}
	if err = yaml.Unmarshal(confContent, conf); err != nil {
		return nil, err
	}
	return conf, nil
}
