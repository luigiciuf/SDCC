# SDCC Gossiping-Based System
Il progetto finale per il corso di Sistemi Distribuiti e Cloud Computing (SDCC) prevede la realizzazione 
di un sistema distribuito basato su algoritmi di gossiping. Questo sistema è 
progettato per stimare le distanze tra nodi in una rete distribuita e 
rilevare eventuali guasti attraverso l'algoritmo di Vivaldi. Utilizzando 
container Docker, il sistema è stato distribuito su istanze Amazon EC2, 
garantendo scalabilità e resilienza. L'emulazione dei ritardi di rete è stata 
implementata per simulare condizioni realistiche di latenza, permettendo di 
testare l'efficacia del sistema in scenari complessi.

## Esecuzione del Progetto

1. **Configurazione Iniziale**:
   Utilizza il file `config.json` per definire il numero di nodi da inizializzare all'interno della rete.

2. **Generazione del Docker Compose**:
   Esegui il seguente comando per generare il file `docker-compose.yml` basato sul numero di nodi configurati:

   ```bash
   go run generate_compose.go
3. **Build dell'Immagine Docker**: Specifica le directory nel file path.env e costruisci l'immagine Docker con il comando:
    ```bash
   docker-compose --env-file path.env build
4. **Avvio dei Container**: Avvia i container Docker con il comando:
   ```bash
   docker-compose --env-file path.env up -d
Dopo questi comandi il sistema sarà pronto per essere utillizzato.

## EC2
Per eseguire il codice si EC2 dobbiamo:
1. **Creare una nuova istanza EC2** direttamente dal sito di AWS.
2. **Avviare il seguente file**:
   ```bash
   connection_to_aws.bat
3. **Preparare l'ambiente** per l'esecuzione del codice eseguendo il seguente file su una scheda del terminale:
      ```bash
   configure_aws.sh
4. **Avviare la build** con il seguente comando:
      ```bash
   sudo docker-compose --env-file path.env build
5. **Avviare i container** :
      ```bash
   sudo docker-compose --env-file path.env up -d

Infine per scollegarsi dall'istanza digitare 
```bash
exit