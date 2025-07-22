// Package requester contains the core logic for creating and sending
// HTTPS requests with a spoofed TLS (JA3) fingerprint.
package requester

import (
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"

	utls "github.com/refraction-networking/utls"
)

// RequestParams holds all the parameters for building and sending a request.
type RequestParams struct {
	URL         string
	Method      string
	JA3Profile  string // e.g., "Chrome", "Firefox", "iOS", "Safari", "Random"
	Headers     map[string]string
	RequestBody string
}

// SendRequest creates an HTTP client with a specified TLS fingerprint,
// builds the request, and returns the HTTP response.
func SendRequest(params RequestParams) (*http.Response, error) {
	// --- 1. Select the ClientHelloID based on the desired profile ---
	var clientHello utls.ClientHelloID
	switch strings.ToLower(params.JA3Profile) {
	case "chrome":
		clientHello = utls.HelloChrome_108
	case "firefox":
		clientHello = utls.HelloFirefox_108
	case "ios":
		clientHello = utls.HelloIOS_16
	case "safari":
		clientHello = utls.HelloSafari_16_0
	case "random":
		clientHello = utls.HelloRandomized
	default:
		// Default to Chrome for safety
		clientHello = utls.HelloChrome_108
	}

	// --- 2. Create a custom dialer for the HTTP transport ---
	// This is the core of the spoofing. We replace the standard TLS dialer
	// with one that uses the utls library.
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	// The custom DialTLS function
	dialTLS := func(network, addr string) (net.Conn, error) {
		// Establish a raw TCP connection
		rawConn, err := dialer.Dial(network, addr)
		if err != nil {
			return nil, err
		}

		// Extract the host for SNI (Server Name Indication)
		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			host = addr // Assume no port if split fails
		}

		// Configure the uTLS connection
		config := &utls.Config{
			ServerName:         host,
			InsecureSkipVerify: true, // Often necessary for scraping non-standard sites
		}

		// Create the uTLS client connection
		uconn := utls.UClient(rawConn, config, clientHello)

		// Perform the handshake to establish the TLS session
		if err := uconn.Handshake(); err != nil {
			uconn.Close()
			return nil, err
		}
		return uconn, nil
	}

	// --- 3. Create the HTTP client with the custom transport ---
	client := &http.Client{
		Transport: &http.Transport{
			DialTLS:         dialTLS,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Redundant but safe
		},
		Timeout: 60 * time.Second,
	}

	// --- 4. Build the HTTP request ---
	var bodyReader *strings.Reader
	if params.RequestBody != "" {
		bodyReader = strings.NewReader(params.RequestBody)
	} else {
		bodyReader = strings.NewReader("")
	}

	req, err := http.NewRequest(params.Method, params.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set a default User-Agent if none is provided, matching the JA3 profile
	if _, ok := params.Headers["User-Agent"]; !ok {
		switch strings.ToLower(params.JA3Profile) {
		case "chrome":
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
		case "firefox":
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0")
		case "safari":
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15")
		default:
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
		}
	}

	// Add all other custom headers
	for key, value := range params.Headers {
		req.Header.Set(key, value)
	}

	// --- 5. Execute the request ---
	return client.Do(req)
}
