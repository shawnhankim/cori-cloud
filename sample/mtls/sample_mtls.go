package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ca   = "/tmp/myCA.pem"
	cert = "/tmp/myCA.cert"
	key  = "/tmp/myCA.key"
)

func main() {
	go HttpServer()
	HttpClient()
}

func HttpServer() {
	http.HandleFunc("/hello", HelloServer)
	if err := http.ListenAndServeTLS(":9191", cert, key, nil); err != nil {
		fmt.Println("ERR:", err)
	}
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Hello World.\n"))
}

func HttpClient() {

	clientCACert, err := ioutil.ReadFile(ca)
	if err != nil {
		fmt.Println("ERR", err)
	}
	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	clientKeyPair, err := tls.LoadX509KeyPair(cert, key)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientKeyPair},
		RootCAs:      clientCertPool,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", "https://localhost:9191/hello", nil)
	if err != nil {
		fmt.Println("ERR:", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERR:", err)
	}
	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(htmlData))
}
