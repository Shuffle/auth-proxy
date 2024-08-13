package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Should validate if the requestuest is just sent as a proxy with the right header
func validateProxyRequest(request *http.Request) error {
	if request.Header.Get("X-Proxy-Request") != "true" {
		return fmt.Errorf("Invalid requestuest")
	}

	// FIXME: Do we do the Proxy injection on the
	// Check connection to Shuffle backend
	// Environments:
	// - Url: SHUFFLE_URL
	// - Apikey: SHUFFLE_API_KEY

	return nil
}

func handleProxy(resp http.ResponseWriter, request *http.Request) {
	log.Printf("PROXY REQUEST URL: %s", request.URL.String())
	/*
		err := validateProxyRequest(request)
		if err != nil {
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write([]byte(fmt.Sprintf(`{"success": false, "message": "Invalid proxy request: %s"}`, err)))
			return
		}
	*/

	isAuthInjected := false


	proxyHost := request.Header.Get("X-Proxy-Host")
	log.Printf("[DEBUG] Proxy request to: %s", request.URL.String())

	httpClient := &http.Client{
		Timeout: 15 * time.Second,
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("[ERROR] Issue in SSR body proxy: %s", err)
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	request.Body = ioutil.NopCloser(bytes.NewReader(body))
	url := fmt.Sprintf("%s%s", proxyHost, request.RequestURI)
	proxyReq, err := http.NewRequest(request.Method, url, bytes.NewReader(body))

	// We may want to filter some headers, otherwise we could just use a shallow copy
	proxyReq.Header = make(http.Header)
	for h, val := range request.Header {
		proxyReq.Header[h] = val
	}

	newresp, err := httpClient.Do(proxyReq)
	if err != nil {
		log.Printf("[ERROR] Issue in SSR newresp for %s - should retry: %s", url, err)
		http.Error(resp, err.Error(), http.StatusBadGateway)
		return
	}

	defer newresp.Body.Close()

	urlbody, err := ioutil.ReadAll(newresp.Body)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusBadGateway)
		return
	}

	//log.Printf("RESP: %s", urlbody)
	for key, value := range newresp.Header {
		//log.Printf("%s %s", key, value)
		for _, item := range value {
			resp.Header().Set(key, item)
		}
	}



	resp.Header().Set("X-Proxy-Response", "true")
	resp.Header().Set("X-Proxy-Host", proxyHost)
	resp.Header().Set("X-Proxy-Auth-Injection", fmt.Sprintf("%t", isAuthInjected))

	resp.WriteHeader(newresp.StatusCode)
	resp.Write(urlbody)
}

func init() {
	http.HandleFunc("/", handleProxy)
}

func main() {

	port := "5004"

	log.Printf("[DEBUG] Started proxy on port 5004")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Printf("[WARNING] Error in requestuest handler: %s", err)
	}
}
