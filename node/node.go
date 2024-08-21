package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Struttura del nodo con le coordinate
type Node struct {
	ID      string
	Address string
	Context *Context
}

// Coordinate del nodo in uno spazio euclideo
type Coordinate struct {
	X, Y float64
}

// Funzione per ottenere l'indirizzo IP locale del nodo
func getLocalIPAddress() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80") // Connessione temporanea per ottenere l'IP
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil // Restituisce l'indirizzo IP
}

// Funzione per generare un ritardo casuale usando una distribuzione normale
func simulateNetworkDelay(mean, stddev float64) {
	delay := rand.NormFloat64()*stddev + mean
	if delay < 0 {
		delay = 0 // Evita ritardi negativi
	}
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

// Funzione per inviare un ping al registro
func sendPingToRegistry(registryAddress, nodeID string, coordinates *HVector) error {
	simulateNetworkDelay(1000, 20) // Simula un ritardo medio di 100ms con deviazione standard di 20ms
	url := fmt.Sprintf("http://%s/ping?id=%s&x=%f&y=%f&h=%f", registryAddress, nodeID, coordinates.X, coordinates.Y, coordinates.H)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error sending ping to registry: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to ping registry: %s", resp.Status)
	}

	return nil
}

// Funzione per ottenere la lista degli indirizzi dei nodi registrati, escludendo il nodo corrente
func getNodeAddresses(registryAddress string, currentNodeID string) ([]string, error) {
	url := fmt.Sprintf("http://%s/nodes", registryAddress) // URL per ottenere i nodi
	resp, err := http.Get(url)                             // Esegui la richiesta HTTP
	if err != nil {
		return nil, fmt.Errorf("Error retrieving node addresses: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to retrieve node addresses: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body) // Leggi il contenuto della risposta
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}

	nodeLines := strings.Split(string(body), "\n") // Estrai le righe
	nodeAddresses := []string{}

	// Estrarre indirizzi validi e assicurarsi di escludere il nodo corrente
	for _, line := range nodeLines {
		parts := strings.Split(line, ",") // Dividi ogni riga in parti
		if len(parts) >= 2 {
			nodeID := strings.TrimSpace(parts[0])  // ID del nodo
			address := strings.TrimSpace(parts[1]) // Indirizzo del nodo

			if nodeID == currentNodeID { // Escludi il nodo corrente
				continue
			}

			// Aggiungi il prefisso "http://" se manca
			if !strings.HasPrefix(address, "http://") {
				address = "http://" + address
			}

			// Assicurati che l'indirizzo non contenga spazi
			if strings.Contains(address, " ") {
				fmt.Printf("Invalid address: %s\n", address) // Indirizzo non valido
				continue
			}

			nodeAddresses = append(nodeAddresses, address) // Aggiungi indirizzi validi
		}
	}

	return nodeAddresses, nil
}

// Funzione per contattare altri nodi e fare gossiping
func (n *Node) contactOtherNodes(registryAddress string) {
	simulateNetworkDelay(1000, 20)                                // Simula un ritardo medio di 100ms con deviazione standard di 20ms
	nodeAddresses, err := getNodeAddresses(registryAddress, n.ID) // Ottieni la lista dei nodi
	if err != nil {
		fmt.Printf("Error getting node addresses: %v\n", err)
		return
	}

	// Scegli casualmente un nodo dalla lista
	selectedNode := chooseRandomNode(nodeAddresses)
	if selectedNode == "" {
		fmt.Println("No other active nodes found to gossip with.")
		return
	}
	// Esegui lo scambio di informazioni con il nodo selezionato
	rtt, _, _ := pingNode(selectedNode)
	// Aggiorna le coordinate del nodo con l'algoritmo di Vivaldi
	n.updateCoordinate(rtt, selectedNode)
}

func pingNode(nodeAddress string) (time.Duration, *Context, error) {
	// Simula un ritardo di rete
	simulateNetworkDelay(1000, 20)

	// Invia il ping al nodo target e ottieni le sue coordinate
	start := time.Now()
	resp, err := http.Get(nodeAddress + "/ping")
	if err != nil {
		return 0, nil, fmt.Errorf("error pinging node: %v", err)
	}
	defer resp.Body.Close()

	// Estrai le coordinate del nodo target dalla risposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Supponiamo che il corpo della risposta contenga le coordinate in formato appropriato
	// Converti il corpo della risposta in coordinate e crea un Context
	targetContext := parseResponseToContext(string(body))

	// Restituisci il tempo RTT e il contesto del nodo target
	return time.Since(start), targetContext, nil
}
func parseResponseToContext(responseBody string) *Context {
	parts := strings.Split(responseBody, ",")
	if len(parts) < 4 {
		return nil // Invalid response
	}

	x, errX := strconv.ParseFloat(parts[1], 64)
	y, errY := strconv.ParseFloat(parts[2], 64)
	h, errH := strconv.ParseFloat(parts[3], 64)

	if errX != nil || errY != nil || errH != nil {
		return nil // Handle parsing errors
	}

	return NewContextFromValues(NewHVector(x, y, h), InitialError)
}

// Funzione per aggiornare le coordinate del nodo usando l'algoritmo Vivaldi
func (n *Node) updateCoordinate(rtt time.Duration, targetNode string) {
	// Ottenere il contesto del nodo target (ipotizzando che sia possibile)
	targetContext := &Context{
		Vec:   NewHVector(rand.Float64(), rand.Float64(), rand.Float64()), // Placeholder per il vettore del nodo target
		Error: InitialError,                                               // Placeholder per l'errore del nodo target
	}
	if targetContext == nil {
		fmt.Println("Target context is nil, cannot update coordinates.")
		return
	}
	n.Context.Update(float64(rtt.Milliseconds()), targetContext)
}

// Funzione per scegliere casualmente un nodo dalla lista escludendo il nodo corrente
func chooseRandomNode(nodeAddresses []string) string {
	rand.Seed(time.Now().UnixNano()) // Inizializza il generatore di numeri casuali

	// Ottieni l'indirizzo IP locale del nodo corrente
	localIP, err := getLocalIPAddress()
	if err != nil {
		fmt.Printf("Error getting local IP address: %v\n", err)
		return ""
	}

	// Creiamo una slice di indici degli indirizzi dei nodi attivi escludendo il nodo corrente
	var activeNodeIndices []int
	for i, address := range nodeAddresses {
		if address != "" && !strings.Contains(address, localIP) {
			activeNodeIndices = append(activeNodeIndices, i)
		}
	}

	if len(activeNodeIndices) == 0 {
		return "" // Non ci sono altri nodi attivi da scegliere
	}

	// Scegliamo casualmente un indice dalla slice degli indici attivi
	randomIndex := rand.Intn(len(activeNodeIndices))
	chosenIndex := activeNodeIndices[randomIndex]

	// Restituiamo l'indirizzo del nodo corrispondente all'indice scelto casualmente
	return nodeAddresses[chosenIndex]
}

// Funzione per registrare il nodo al registro
func registerNode(registryAddress, nodeID, nodeAddress string, coordinates *HVector) error {
	ip, err := getLocalIPAddress()
	if err != nil {
		return fmt.Errorf("Error getting local IP address: %v", err)
	}

	url := fmt.Sprintf("http://%s/register?id=%s&address=%s&ip=%s&x=%f&y=%f&h=%f", registryAddress, nodeID, nodeAddress, ip, coordinates.X, coordinates.Y, coordinates.H)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error registering node: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register node: %s", resp.Status)
	}

	return nil
}

// Esecuzione del nodo
func main() {
	registryAddress := "registry:2020"     // Porta 2020 del registro
	current_nodeID := os.Getenv("NODE_ID") // Ottieni l'ID del nodo dall'ambiente

	// Ottieni l'indirizzo IP locale
	localIP, err := getLocalIPAddress()
	if err != nil {
		fmt.Printf("Error getting local IP address: %v\n", err)
		return
	}

	nodePort := fmt.Sprintf("800%d", current_nodeID[4]-'0') // Determina la porta del nodo
	nodeAddress := fmt.Sprintf("%s:%s", localIP, nodePort)  // Crea l'indirizzo del nodo
	// Crea il nodo con coordinate iniziali casuali
	node := &Node{
		ID:      current_nodeID,
		Address: nodeAddress,
		Context: NewContext(),
	}
	// Registra il nodo al registro
	err = registerNode(registryAddress, current_nodeID, nodeAddress, node.Context.Vec)
	if err != nil {
		fmt.Printf("Error registering node: %v\n", err)
		return
	}

	// Invio periodico di ping al registro
	go func() {
		for {
			err := sendPingToRegistry(registryAddress, current_nodeID, node.Context.Vec) // Invia un ping
			if err != nil {
				fmt.Printf("Error sending ping to registry: %v\n", err)
			}
			time.Sleep(15 * time.Second) // Intervallo tra i ping
			node.contactOtherNodes(registryAddress)
		}

	}()

	fmt.Printf("Node %s registered successfully at address %s\n", current_nodeID, nodeAddress)

	// Mantieni il nodo in esecuzione
	http.ListenAndServe(fmt.Sprintf(":%s", nodePort), nil) // Ascolta sulla porta specificata
}
