// main.go is the entry point for the StealthyScraper CLI.
// It parses command-line arguments, constructs a request,
// and uses the requester package to send it.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Djmon007/stealthysraper/requester"
)

// headersFlag is a custom flag type to handle multiple header arguments.
type headersFlag map[string]string

func (h *headersFlag) String() string {
	return "Custom headers"
}

func (h *headersFlag) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("header must be in 'Key: Value' format")
	}
	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])
	(*h)[key] = val
	return nil
}

func main() {
	// Set up logging to print without timestamp prefixes for cleaner output
	log.SetFlags(0)

	// --- Define CLI Flags ---
	method := flag.String("method", "GET", "HTTP method (GET, POST, PUT, DELETE)")
	ja3Profile := flag.String("ja3", "Chrome", "TLS fingerprint to use. Options: Chrome, Firefox, iOS, Safari, Random")
	data := flag.String("data", "", "Request body for POST, PUT, etc.")

	// Custom flag for headers
	headers := make(headersFlag)
	flag.Var(&headers, "H", "Set custom headers (can be used multiple times, e.g., -H \"Key1: Value1\" -H \"Key2: Value2\")")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [flags] <URL>\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// --- Validate Input ---
	if flag.NArg() != 1 {
		log.Println("Error: You must specify exactly one URL.")
		flag.Usage()
		os.Exit(1)
	}
	url := flag.Arg(0)

	// --- Prepare and Send Request ---
	params := requester.RequestParams{
		URL:         url,
		Method:      strings.ToUpper(*method),
		JA3Profile:  *ja3Profile,
		Headers:     headers,
		RequestBody: *data,
	}

	resp, err := requester.SendRequest(params)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// --- Print Response ---
	// Print status and headers to stderr to separate them from the body
	fmt.Fprintf(os.Stderr, "HTTP/1.1 %s\n", resp.Status)
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Fprintf(os.Stderr, "%s: %s\n", key, value)
		}
	}
	fmt.Fprintln(os.Stderr) // Separator line

	// Print the body to stdout
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	fmt.Println(string(body))
}
