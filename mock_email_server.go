package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type EmailPayload struct {
	MailTo  string      `json:"mailTo"`
	Subject string      `json:"subject"`
	Body    interface{} `json:"body"`
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	// Print raw body for debugging
	fmt.Printf("sendEmailHandler")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fmt.Printf("\n📩 Raw JSON Received:\n%s\n", string(body))

	// Reset the body for decoding
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	// Decode into EmailPayload struct
	var payload EmailPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		fmt.Printf("❌ JSON decode error: %v\n", err)
		return
	}

	// Simulate email processing
	log.Printf("✅ Mock Email Sent!\n📨 To: %s\n📌 Subject: %s\n📄 Body: %+v\n",
		payload.MailTo, payload.Subject, payload.Body)

	// Send OK response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email sent successfully (mock)"))
}

func main() {
	fmt.Println("📧 Mock Email Service running at: http://127.0.0.1:8010/send-email")
	http.HandleFunc("/send-email", sendEmailHandler)

	fmt.Println("erwtetgtrhy")
	if err := http.ListenAndServe(":8010", nil); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
