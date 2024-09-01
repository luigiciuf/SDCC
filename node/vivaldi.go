package main

import "math"

const cc = 0.25
const ce = 0.5
const InitialError = 10
const scaleFactor = 100.0 // Fattore di scala da applicare alle coordinate

type Context struct {
	Vec   *HVector
	Error float64
}

func NewContext() *Context {
	return NewContextFromValues(
		NewHVector(0, 0, 0),
		InitialError,
	)
}

func NewContextFromValues(vec *HVector, error float64) *Context {
	return &Context{
		Vec:   vec,
		Error: error,
	}
}

// Funzione che aggiorna il contesto del nodo in base al tempo di round-trip (RTT) e al contesto del nodo pari
func (ctx *Context) Update(rtt float64, peer *Context) *HVector {
	w := ctx.Error / (ctx.Error + peer.Error) // Calcola il peso w basato sugli errori relativi dei nodi
	ab := ctx.Vec.Sub(peer.Vec)               // Calcola la differenza vettoriale tra le posizioni dei nodi
	re := rtt - ab.Magnitude()                // Calcola la differenza tra RTT osservato e la distanza stimata
	es := math.Abs(re) / rtt

	if es != es || w != w || re != re {
		return ctx.Vec
	}
	// Aggiorna l'errore del nodo in base all'errore scalare e al peso calcolato
	ctx.Error = es*ce*w + ctx.Error*(1-ce*w) // e_i = e_s*c_e*w + e_i*(1 - c_e*w)
	// Calcola il fattore di correzione delle coordinate
	d := cc * w
	// Aggiorna il vettore del nodo in base alla correzione calcolata e alla differenza vettoriale
	ctx.Vec = ctx.Vec.Add(ab.Unit().Scale(d * re * 100.0))
	return ctx.Vec
}
