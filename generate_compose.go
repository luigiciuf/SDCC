package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Config struct rappresenta la struttura del file di configurazione
type Config struct {
	Nodes int `json:"nodes"`
}

func getNumNodesFromConfig() int {
	// Leggi il file di configurazione
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Errore nella lettura del file di configurazione: %v", err)
	}

	// Parse del file JSON
	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Errore nel parsing del file di configurazione JSON: %v", err)
	}

	return config.Nodes
}

func generateDockerCompose(numNodes int) {
	// Apertura del file docker-compose.yml in modalità scrittura
	file, err := os.Create("docker-compose.yml")
	if err != nil {
		log.Fatalf("Errore nell'apertura del file docker-compose.yml: %v", err)
	}
	defer file.Close()

	// Scrittura del contenuto nel file
	file.WriteString("version: \"3\"\n\n")
	file.WriteString("services:\n")
	file.WriteString("  registry:\n")
	file.WriteString("    build:\n")
	file.WriteString("      context: ./registry\n") // Riferimento alla directory registry
	file.WriteString("    ports:\n")
	file.WriteString("      - \"1234:1234\"\n")
	file.WriteString("    networks:\n")
	file.WriteString("      - my_network\n\n")

	// Aggiungere i servizi per i nodi
	for i := 1; i <= numNodes; i++ {
		file.WriteString(fmt.Sprintf("  node%d:\n", i))
		file.WriteString("    build:\n")
		file.WriteString("      context: ./node\n") // Riferimento alla directory node
		file.WriteString(fmt.Sprintf("    ports:\n      - \"%d:%d\"\n", 8000+i, 8000+i))
		file.WriteString(fmt.Sprintf("    environment:\n      - NODE_ID=node%d\n", i))
		file.WriteString("    networks:\n")
		file.WriteString("      - my_network\n")
		file.WriteString("    depends_on:\n")
		file.WriteString("      - registry\n\n")
	}

	file.WriteString("networks:\n")
	file.WriteString("  my_network:\n")
	file.WriteString("    driver: bridge\n")
}

func main() {
	numNodes := getNumNodesFromConfig()
	generateDockerCompose(numNodes)
}
