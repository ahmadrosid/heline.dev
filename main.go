package main

import (
	"fmt"
	"log"
	"net/http"

	ghttp "github.com/ahmadrosid/heline/http"
	"github.com/ahmadrosid/heline/core/module/solr"
	"github.com/nullitics/nullitics"
)

func main() {
	// Set up Solr schema if needed
	fmt.Println("🔄 Checking and setting up Solr schema...")
	if err := solr.SetupSchema(); err != nil {
		log.Printf("⚠️ Warning: Failed to set up Solr schema: %v\n", err)
		// Continue anyway, as the schema might already be set up or will be set up later
	}

	option := nullitics.BlacklistPrefix([]string{"/_next/", "favicon.png"}...)
	context := nullitics.New(option)
	handler := ghttp.Handler(context.Report(nil))

	port := "8000"

	fmt.Printf("🚀 Starting server on http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, context.Collect(handler))
	if err != nil {
		println("❌ Server already started!")
		println(err.Error())
	}
}
