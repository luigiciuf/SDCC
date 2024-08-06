package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
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

// Funzione per inviare un ping al registro
func sendPingToRegistry(registryAddress, nodeID string, coordinates *HVector) error {
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
	rtt, err := pingNode(selectedNode)
	if err != nil {
		fmt.Printf("Error pinging node %s: %v\n", selectedNode, err)
		return
	}
	// Aggiorna le coordinate del nodo con l'algoritmo di Vivaldi
	n.updateCoordinate(rtt, selectedNode)
}

// Funzione per inviare un ping al nodo selezionato e ottenere il tempo di ping
func pingNode(nodeAddress string) (time.Duration, error) {
	start := time.Now()
	_, err := http.Get(nodeAddress + "/ping") // Invia il ping al nodo
	if err != nil {
		return 0, fmt.Errorf("error pinging node: %v", err)
	}
	return time.Since(start), nil // Restituisci il tempo trascorso dal ping
}

// Funzione per aggiornare le coordinate del nodo usando l'algoritmo Vivaldi
func (n *Node) updateCoordinate(rtt time.Duration, targetNode string) {
	// Ottenere il contesto del nodo target (ipotizzando che sia possibile)
	targetContext := &Context{
		Vec:   NewHVector(rand.Float64(), rand.Float64(), rand.Float64()), // Placeholder per il vettore del nodo target
		Error: InitialError,                                               // Placeholder per l'errore del nodo target
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

// Funzione per scambiare informazioni con un altro nodo
func exchangeInfo(sourceNodeID, targetNodeAddress string) {
	// Esempio di implementazione dello scambio di informazioni tra nodi
	fmt.Printf("Node %s exchanging information with node at address: %s\n", sourceNodeID, targetNodeAddress)
	// Esegui qui la logica per lo scambio di informazioni con il nodo selezionato
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
