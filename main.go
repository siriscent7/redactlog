package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/siriscent7/redactlog/redactor"
	"github.com/siriscent7/redactlog/web"
)

func main() {
	serve := flag.Bool("serve", false, "run as a web server instead of CLI")
	modeFlag := flag.String("mode", "mask", "redaction mode: mask | hash | drop")
	flag.Parse()

	if *serve {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("RedactLog web UI on :%s", port)
		log.Fatal(http.ListenAndServe(":"+port, web.Handler()))
		return
	}

	// CLI mode: stream stdin -> redacted stdout
	mode := redactor.Mask
	switch *modeFlag {
	case "hash":
		mode = redactor.Hash
	case "drop":
		mode = redactor.Drop
	}
	r := redactor.New(mode)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	total := 0
	for scanner.Scan() {
		clean, n := r.Redact(scanner.Text())
		total += n
		fmt.Println(clean)
	}
	fmt.Fprintf(os.Stderr, "redacted %d PII items\n", total)
}
