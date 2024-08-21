package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	NumberOfNodes int `json:"number_of_nodes"`
}

func main() {
	// Leggi il file di configurazione
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("Errore nell'apertura del file config.json:", err)
		return
	}
	defer file.Close()

	// Decodifica il file JSON
	configData, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Errore nella lettura del file config.json:", err)
		return
	}

	var config Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		fmt.Println("Errore nella decodifica del JSON:", err)
		return
	}

	// Genera il contenuto di docker-compose.yml
	dockerComposeContent := fmt.Sprintf(`
version: "3"
services:
  registry:
    build:
      context: ${PROJECT_PATH}
      dockerfile: ${DOCKERFILE_PATH}/registry/Dockerfile
    ports:
      - "2020:2020"
    networks:
      - my_network
    cap_add:
      - NET_ADMIN
`)

	for i := 1; i <= config.NumberOfNodes; i++ {
		dockerComposeContent += fmt.Sprintf(`
  node%d:
    build:
      context: ${PROJECT_PATH}
      dockerfile: ${DOCKERFILE_PATH}/node/Dockerfile
    privileged: true
    ports:
      - "800%d:800%d"
    environment:
      - NODE_ID=node%d
    networks:
      - my_network
    depends_on:
      - registry
    cap_add:
      - NET_ADMIN
`, i, i, i, i)
	}

	dockerComposeContent += `
networks:
  my_network:
    driver: bridge
`

	// Scrivi il contenuto generato in docker-compose.yml
	err = ioutil.WriteFile("docker-compose.yml", []byte(dockerComposeContent), 0644)
	if err != nil {
		fmt.Println("Errore nella scrittura del file docker-compose.yml:", err)
		return
	}

	fmt.Println("docker-compose.yml generato con successo.")
}
