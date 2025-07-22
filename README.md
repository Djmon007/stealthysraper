# StealthyScraper: Custom HTTPS & TLS Fingerprint Spoofing Tool

**StealthyScraper** is a command-line tool and Go library for making highly customizable HTTPS requests with a key feature for modern data scraping: **TLS fingerprint spoofing**.

Modern web application firewalls (WAFs) and bot detection systems don't just look at IP addresses and headers; they analyze the client's TLS handshake to create a "fingerprint" (often a [JA3 hash](https://github.com/salesforce/ja3)). Standard HTTP libraries in languages like Python or Go have a unique, easily identifiable fingerprint. StealthyScraper allows you to mimic the TLS fingerprint of common web browsers, drastically reducing the chances of being blocked.

## Key Features

* **Customizable HTTPS Requests:** Control the method, headers, and body of your requests.

* **TLS Fingerprint (JA3) Spoofing:** Masquerade as popular browsers (Chrome, Firefox, Safari, etc.) to bypass advanced bot detectors.

* **User-Friendly CLI:** Easy-to-use command-line interface for quick tests and scraping tasks.

* **Usable as a Go Library:** Integrate the core request logic directly into your own Go applications.

* **High Performance:** Built with Go for speed and concurrency.

## Why TLS Fingerprinting Matters

When your browser connects to a secure website, it sends a `ClientHello` message. This message contains details about the TLS version, supported cipher suites, extensions, and elliptic curves. The specific combination and order of these parameters create a unique signature. Services can block requests from any signature that doesn't match that of a real browser. This tool allows you to generate a `ClientHello` that matches a real browser, making your requests look legitimate at the transport layer.

## Installation

Ensure you have Go installed (version 1.18 or later).

```
# Clone the repository
git clone [https://github.com/your-/stealthysraper.git](https://github.com/your-Djmon007/stealthysraper.git)
cd stealthysraper

# Build the binary
go build

# Or install it directly to your GOPATH
go install .

```

## CLI Usage

The tool is simple to use. The basic syntax is `stealthysraper [flags] <URL>`.

```
./stealthysraper -h
Usage of ./stealthysraper:
  -data string
        Request body for POST, PUT, etc.
  -H value
        Set custom headers (can be used multiple times, e.g., -H "Key1: Value1" -H "Key2: Value2")
  -ja3 string
        TLS fingerprint to use. Options: Chrome, Firefox, iOS, Safari, Random (default "Chrome")
  -method string
        HTTP method (GET, POST, PUT, DELETE) (default "GET")

```

### Examples

**1. Simple GET Request (emulating Chrome)**

This will make a GET request to a TLS inspection endpoint. You can check the reported JA3 hash to confirm it's working.

```
./stealthysraper [https://tls.browserleaks.com/json](https://tls.browserleaks.com/json)

```

**2. POST Request with Custom Data and Headers (emulating Firefox)**

```
./stealthysraper -method POST \
                 -ja3 Firefox \
                 -H "Content-Type: application/json" \
                 -H "Authorization: Bearer mysecrettoken" \
                 -data '{"key": "value", "login": "user1"}' \
                 [https://api.example.com/login](https://api.example.com/login)

```

**3. GET Request emulating an iPhone (Safari)**

```
./stealthysraper -ja3 iOS [https://httpbin.org/headers](https://httpbin.org/headers)

```

## Library Usage

You can also import the `requester` package into your own Go projects for more complex scraping logic.

```
package main

import (
	"fmt"
	"io"
	"log"

	"[github.com/your-/stealthysraper/requester](https://github.com/your-/stealthysraper/requester)"
)

func main() {
	// Define the request parameters
	params := requester.RequestParams{
		URL:         "[https://tls.browserleaks.com/json](https://tls.browserleaks.com/json)",
		Method:      "GET",
		JA3Profile:  "Firefox", // Can be "Chrome", "Firefox", "iOS", "Safari", "Random"
		Headers:     map[string]string{"Accept-Language": "en-US,en;q=0.9"},
	}

	// Send the request
	resp, err := requester.SendRequest(params)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read and print the response
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	fmt.Println("Response Body:")
	fmt.Println(string(body))
}

```

## Disclaimer

This tool is intended for educational purposes and for legitimate data scraping activities. Users are responsible for complying with the terms of service of any website they target. Please scrape responsibly.

## License

This project is licensed under the MIT License - see the [LICENSE](https://www.google.com/search?q=LICENSE) file for details.
