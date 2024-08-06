# Usa un'immagine base di Go con Alpine Linux
FROM golang:alpine

# Imposta la directory di lavoro all'interno del container
WORKDIR /app

#prova

# Copia i file delle dipendenze Go nel container
COPY ../../go.mod /app/go.mod
COPY ../../go.sum /app/go.sum

# Copia i file sorgente Go nel container
COPY ../../node/hvector.go /app/hvector.go
#COPY ../../node/node.go /app/node.go
COPY ../../node/registry.go /app/registry.go
COPY ../../node/utils.go /app/utils.go
COPY ../../node/vivaldi.go /app/vivaldi.go

# Scarica i moduli Go specificati in go.mod
RUN go mod download

# Compila l'applicazione Go
RUN go build -o registry .

# Definisci il punto di ingresso per il container
ENTRYPOINT [ "./registry" ]
