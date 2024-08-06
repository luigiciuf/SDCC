package main

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type DockerCompose struct {
	Services map[string]Service `yaml:"services"`
}

type Service struct {
	Ports       []string `yaml:"ports"`
	Environment []string `yaml:"environment"`
}

// ParseDockerCompose legge e analizza il file docker-compose.yml restituendo una mappa
// contenente il nome del nodo e la porta associata.
func ParseDockerCompose() (map[string]string, error) {
	filePath := "/app/config/docker-compose.yml"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parsa il file YAML
	var compose DockerCompose
	err = yaml.Unmarshal(data, &compose)
	if err != nil {
		return nil, err
	}

	// Inizializza una mappa per memorizzare il nome del nodo e la porta associata
	services := make(map[string]string)

	// Scansiona i servizi e memorizza il nome del nodo e la porta associata
	for name, service := range compose.Services {
		if len(service.Ports) > 0 {
			port := strings.Split(service.Ports[0], ":")[0]
			services[name] = port
		} else {
			// Se la porta non è definita, impostala su "non definita"
			services[name] = "non definita"
		}
	}

	return services, nil
}
