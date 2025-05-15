package main

import (
	"fmt"
	"log"
	"net/http"

	ghttp "github.com/ahmadrosid/heline/http"
	"github.com/ahmadrosid/heline/core/module/solr"
)

func main() {
	// Set up Solr schema if needed
	fmt.Println("🔄 Checking and setting up Solr schema...")
	if err := solr.SetupSchema(); err != nil {
		log.Printf("⚠️ Warning: Failed to set up Solr schema: %v\n", err)
		// Continue anyway, as the schema might already be set up or will be set up later
	}

	port := "8000"

	fmt.Printf("🚀 Starting server on http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, ghttp.Handler())
	if err != nil {
		println("❌ Server already started!")
		println(err.Error())
	}
}
