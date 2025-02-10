package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
)

var legacyServer string
var port string

var certFile string
var keyFile string

func ProxyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetURL, err := url.Parse(legacyServer)
		if err != nil {
			http.Error(w, "Error while parsing the target URL", http.StatusInternalServerError)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error while reading the request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		if r.Method == http.MethodPost {
			r.Method = http.MethodGet
		}

		proxyReq, err := http.NewRequest(r.Method, targetURL.ResolveReference(r.URL).String(), bytes.NewBuffer(body))
		if err != nil {
			http.Error(w, "Error while creating the request", http.StatusInternalServerError)
			return
		}

		for key, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}

		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, "Error while sending the request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})
}

func main() {
	flag.StringVar(&port, "port", "8080", "The port to run the proxy server on (default: 8080)")
	flag.StringVar(&legacyServer, "legacy-server", "", "The legacy server to proxy requests to")
	flag.StringVar(&certFile, "cert-file", "", "The path to the certificate file")
	flag.StringVar(&keyFile, "key-file", "", "The path to the key file")
	flag.Parse()

	if legacyServer == "" {
		log.Fatal("Usage: post-to-get-middleware --legacy-server=http://localhost:8081")
	}

	isHttps := certFile != "" && keyFile != ""

	r := http.NewServeMux()

	r.Handle("/", ProxyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Post to GET!"))
	})))

	log.Printf("Proxy server is running on :%s\n", port)

	if isHttps {
		log.Fatal(http.ListenAndServeTLS(":"+port, certFile, keyFile, r))
		return
	}

	log.Fatal(http.ListenAndServe(":"+port, r))
}
