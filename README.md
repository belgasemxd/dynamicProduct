# Templ + HTMX + Gin + MongoDBProjekt

Dieses Projekt ist ein kleiner Webserver, gebaut mit:

- **[Templ](https://templ.guide/)** – type-safe HTML Templates für Go.  
- **[HTMX](https://htmx.org/)** – ermöglicht dynamische Webseiten ohne viel JavaScript; z. B. für asynchrone Updates von Teilen der Seite.
- **[Gin](https://github.com/gin-gonic/gin)** – schnelles Webframework für Go.  
- **[MongoDB](https://www.mongodb.com/)** – Datenbank für die Speicherung von Daten.  


Der Server läuft lokal und ist nach dem Starten unter **http://localhost:8080** erreichbar.  

---

## Funktionen
Der Großteil des Projekts ist selbsterklärend. Ein besonders interessantes Feature ist jedoch der **Migreren-Button**.  
Im Projektverzeichnis unter `resources/config` befindet sich die Datei `model.yaml`, in der die Tabelle definiert wird.  
Nimmt man Änderungen an dieser Definition vor und klickt anschließend auf **Migrieren**, werden die Anpassungen sofort übernommen – ohne dass der Server neu gestartet werden muss.

---

## ▶️ Projekt starten
   ```bash
    go mod tidy
    cd cmd
    cd run
    go run main.go
