package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Struttura per mantenere le informazioni sui nodi
type NodeInfo struct {
	ID          string
	Address     string
	IP          string
	LastSeen    time.Time
	Coordinates *HVector
}

// Struttura del registro
type Registry struct {
	mu    sync.Mutex
	nodes map[string]*NodeInfo
}

// Funzione per creare un nuovo registro
func NewRegistry() *Registry {
	return &Registry{
		nodes: make(map[string]*NodeInfo),
	}
}

// Funzione per registrare un nuovo nodo
func (r *Registry) registerNode(w http.ResponseWriter, req *http.Request) {
	nodeID := req.URL.Query().Get("id")
	nodeAddress := req.URL.Query().Get("address")
	nodeIP := req.URL.Query().Get("ip")
	x := req.URL.Query().Get("x")
	y := req.URL.Query().Get("y")
	h := req.URL.Query().Get("h")

	if nodeID == "" || nodeAddress == "" || nodeIP == "" || x == "" || y == "" || h == "" {
		http.Error(w, "Missing node ID, address, IP, or coordinates", http.StatusBadRequest)
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	coordinates := NewHVectorFromString(x, y, h)

	r.nodes[nodeID] = &NodeInfo{
		ID:          nodeID,
		Address:     nodeAddress,
		IP:          nodeIP,
		LastSeen:    time.Now(),
		Coordinates: coordinates,
	}

	fmt.Printf("Node registered: ID=%s, Address=%s, IP=%s, Coordinates=%+v\n", nodeID, nodeAddress, nodeIP, coordinates)
	fmt.Fprintln(w, "Node registered")
}

// Funzione per gestire i ping dei nodi
func (r *Registry) handlePing(w http.ResponseWriter, req *http.Request) {
	nodeID := req.URL.Query().Get("id")
	x := req.URL.Query().Get("x")
	y := req.URL.Query().Get("y")
	h := req.URL.Query().Get("h")

	if nodeID == "" || x == "" || y == "" || h == "" {
		http.Error(w, "Missing node ID or coordinates", http.StatusBadRequest)
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	node, exists := r.nodes[nodeID]
	if !exists {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	node.LastSeen = time.Now()
	node.Coordinates = NewHVectorFromString(x, y, h)
	fmt.Printf("Node pinged: ID=%s, Coordinates=%+v\n", nodeID, node.Coordinates)
	fmt.Fprintln(w, "Ping received")
}

// Funzione per ottenere la lista dei nodi
func (r *Registry) listNodes(w http.ResponseWriter, req *http.Request) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, node := range r.nodes {
		fmt.Fprintf(w, "%s,%s,%s,%f,%f,%f\n", node.ID, node.Address, node.IP, node.Coordinates.X, node.Coordinates.Y, node.Coordinates.H)
	}
}

// Funzione per rimuovere i nodi inattivi
func (r *Registry) removeInactiveNodes() {
	for {
		time.Sleep(30 * time.Second) // Controlla ogni 30 secondi
		r.mu.Lock()
		for id, node := range r.nodes {
			if time.Since(node.LastSeen) > 1*time.Minute { // Se un nodo non pinga da più di 1 minuto
				delete(r.nodes, id)
				fmt.Printf("Node removed due to inactivity: ID=%s\n", id)
			}
		}
		r.mu.Unlock()
	}
}

func main() {
	registry := NewRegistry()

	// Handlers per le varie richieste
	http.HandleFunc("/register", registry.registerNode)
	http.HandleFunc("/ping", registry.handlePing)
	http.HandleFunc("/nodes", registry.listNodes)

	// Avvia la goroutine per rimuovere i nodi inattivi
	go registry.removeInactiveNodes()

	// Avvia il server HTTP
	fmt.Println("Registry server started on port 2020")
	http.ListenAndServe(":2020", nil)
}

// NewHVectorFromString è una funzione di supporto per convertire le coordinate dalla stringa ai float64
func NewHVectorFromString(x, y, h string) *HVector {
	xVal := parseStringToFloat64(x)
	yVal := parseStringToFloat64(y)
	hVal := parseStringToFloat64(h)
	return NewHVector(xVal, yVal, hVal)
}

// parseStringToFloat64 è una funzione di supporto per convertire una stringa in float64
func parseStringToFloat64(val string) float64 {
	res, err := strconv.ParseFloat(val, 64)
	if err != nil {
		fmt.Printf("Error converting string to float64: %v\n", err)
		return 0
	}
	return res
}
