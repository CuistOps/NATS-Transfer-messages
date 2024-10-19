package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nats-io/nats.go"
)

const (
	subject = "app.message"
)

var (
	nc *nats.Conn
)

func sendOrderToConsumers(message string) error {

	//log.Printf("☕ New order received for %s", subject)

	err := nc.Publish(subject, []byte(message))
	if err != nil {
		log.Printf("Error sending order: %v", err)
		return err
	}

	return nil
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	log.Println("🚀 New message received")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	message := r.FormValue("message")
	sendOrderToConsumers(message)

}

func handleHome(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	html := `
        <html>
            <body>
                <h1>Send a Message</h1>
                <form action="/send" method="POST">
                    <label for="message">Message:</label>
                    <input type="text" id="message" name="message">
                    <input type="submit" value="Send">
                </form>
            </body>
        </html>
    `
	w.Write([]byte(html))
}

func main() {

	err := error(nil)

	nc, err = nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatal("connect to nats: ", err)
	}

	log.Printf("🚀 Starting sender service with %s as NATS URI\n", os.Getenv("NATS_URL"))

	defer nc.Drain()

	http.HandleFunc("/send", handleSend)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/index", http.StatusMovedPermanently)
	})

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		handleHome(w)
	})

	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}

}
