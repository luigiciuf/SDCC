# Usa un'immagine di base leggera per Go
FROM golang:alpine
# Imposta il working directory all'interno del container per il nodo
WORKDIR /app


# Copia i file del nodo nella directory di lavoro del container
COPY ../../node/node.go /app/node.go
COPY ../../node/utils.go /app/utils.go
COPY ../../node/vivaldi.go /app/vivaldi.go
COPY ../../node/hvector.go /app/hvector.go
#COPY ../../node/registry.go /app/registry.go
COPY ../../go.mod /app/go.mod
COPY ../../go.sum /app/go.sum

# Compila il codice Go per il nodo
RUN go build -o node .
RUN go build -o utils .
RUN go build -o vivaldi .
RUN go build -o hvector .
#RUN go build -o registry .
ENTRYPOINT [ "./node" ]