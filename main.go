package main

import (
	"log"
	"main/client"
	_ "main/modules"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("🚀 PbxGo starting...")
	err := client.InitClient()
	if err != nil {
		log.Fatalf("❌ Failed to initialize client: %v", err)
	}
	log.Println("✅ Client initialized")
	log.Println("✅ Handlers registered")

	client.Client.SetCommandPrefixes(".")
	client.RegisterHandlers()
	client.Client.Idle()
}
