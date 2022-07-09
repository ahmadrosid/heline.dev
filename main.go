package main

import (
	"fmt"
	"net/http"

	ghttp "github.com/ahmadrosid/heline/http"
	"github.com/nullitics/nullitics"
)

func main() {
	option := nullitics.BlacklistPrefix([]string{"/_next/"}...)
	context := nullitics.New(option)
	handler := ghttp.Handler(context.Report(nil))

	port := "80"

	fmt.Printf("🚀 Starting server on http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, context.Collect(handler))
	if err != nil {
		println("❌ Server already started!")
		println(err.Error())
	}
}
