package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"os"
	"time"

	"github.com/fatih/color"
)

const (
	infoHeader  = "\n----------------[INFO]---------------------\n"
	testHeader  = "\n----------------[TEST]---------------------\n"
	extraHeader = "\n----------------[EXTRA]---------------------\n"
)

func Shoot(url string, showHeaders bool) {
	var dns, tlsHandshake, connect, start time.Time
	var redirects int

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	fmt.Printf(
		"%sURL: \t%s\nPROTO: \t%s\nMETHOD: %s\n%s\n",
		infoHeader, req.URL, req.Proto, req.Method, testHeader,
	)

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) {
			redirects++
			if redirects > 1 {
				fmt.Printf("\u2757 Redirect to %s, taking another round\n", color.YellowString(dsi.Host))
			}
			dns = time.Now()
		},
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Done: \t%v\n", time.Since(dns))
		},
		ConnectStart: func(network, addr string) {
			connect = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			if err != nil {
				log.Printf("Failed to connect: %v", err)
			} else {
				fmt.Printf(color.GreenString("Connected ")+"to: \t%s\n", addr)
				fmt.Printf("Connect time: \t%v\n", time.Since(connect))
			}
		},
		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			if err != nil {
				log.Printf("TLS handshake failed: %v", err)
			} else {
				fmt.Printf("TLS Handshake: \t%v\n", time.Since(tlsHandshake))
			}
		},
		GotFirstResponseByte: func() {
			color.Cyan("TTFB: \t\t%v\n", time.Since(start))
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	client := &http.Client{}
	start = time.Now()
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("\nStatus code: \t")
	switch resp.StatusCode {
	case 200:
		color.Green("%v", resp.StatusCode)
	case 404:
		color.Yellow("%v", resp.StatusCode)
	default:
		color.Red("%v", resp.StatusCode)
	}

	color.Cyan("Total time: \t%v\n", time.Since(start))

	if showHeaders {
		fmt.Print(extraHeader)
		fmt.Println("Response Headers:")
		for key, value := range resp.Header {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
}

func printUsageAndExit() {
	fmt.Println("Usage: raiden [--headers] <url>")
	fmt.Println("\nRaiden is a simple HTTP client that traces the request lifecycle.")
	fmt.Println("\nArguments:")
	fmt.Println("  <url>        The URL to send the GET request to.")
	fmt.Println("  --headers    Show response headers.")
	fmt.Println("\nExample:")
	fmt.Println("  raiden --headers https://www.example.com")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printUsageAndExit()
	}

	showHeaders := false
	var url string

	if len(os.Args) == 3 && os.Args[1] == "--headers" {
		showHeaders = true
		url = os.Args[2]
	} else if len(os.Args) == 2 {
		url = os.Args[1]
	} else {
		printUsageAndExit()
	}

	Shoot(url, showHeaders)
}
