// Copyright 2023 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

const configFile = ".config/zero/zero.yaml"

type Config struct {
	ServiceName     string `yaml:"serviceName"`
	ServiceConfig   string `yaml:"serviceConfig"`
	ApiKey          string `yaml:"apiKey"`
	Summary         string `yaml:"summary"`
	Title           string `yaml:"title"`
	ProducerProject string `yaml:"producerProject"`
	ConsumerProject string `yaml:"consumerProject"`
}

func GetConfig() (*Config, error) {
	var config Config
	dirname, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path.Join(dirname, configFile))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	return &config, err
}
