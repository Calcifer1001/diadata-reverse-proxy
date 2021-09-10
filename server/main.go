package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func main() {
	demoURL, err := url.Parse("https://localhost")
	if err != nil {
		log.Fatal(err)
	}

	demoGoogleURL, errGoogle := url.Parse("https://www.google.com")
	if errGoogle != nil {
		log.Fatal(errGoogle)
	}
	var selected = ""
	proxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var resp *http.Response
		respDefault, err := getResponse(demoURL)
		respGoogle, errGoogle := getResponse(demoGoogleURL)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(rw, err)
			return
		}
		if errGoogle != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(rw, errGoogle)
			return
		}

		if selected == "" {
			var rand = rand.Float32()
			if rand < 0.5 {
				fmt.Println("Choosing default")
				resp = respDefault
				selected = "default"
			} else {
				fmt.Println("Choosing google")
				resp = respGoogle
				selected = "google"
			}
		} else {
			if selected == "google" {
				resp = respGoogle
			} else {
				resp = respDefault
			}
		}

		for key, values := range resp.Header {
			for _, value := range values {
				rw.Header().Set(key, value)
			}
		}

		done := make(chan bool)
		go func() {
			for {
				select {
				case <-time.Tick(time.Millisecond * 10):
					rw.(http.Flusher).Flush()
				case <-done:
					selected = ""
					return
				}
			}
		}()

		var trailerKeys []string
		for key := range resp.Trailer {
			trailerKeys = append(trailerKeys, key)
		}

		rw.Header().Set("Trailer", strings.Join(trailerKeys, ","))

		rw.WriteHeader(resp.StatusCode)
		io.Copy(rw, resp.Body)
		close(done)

		for key, values := range resp.Trailer {
			for _, value := range values {
				rw.Header().Set(key, value)
			}
		}

	})
	http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", proxy)
}

func setReq(req *http.Request, url *url.URL) {
	req.Host = url.Host
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.RequestURI = ""
}

func getResponse(url *url.URL) (*http.Response, error) {
	fmt.Println(url.RequestURI())
	var req, err = http.NewRequest("GET", url.RequestURI(), nil)
	setReq(req, url)
	s, _, _ := net.SplitHostPort(req.RemoteAddr)
	req.Header.Set("X-Forwarded-For", s)
	resp, err := http.DefaultClient.Do(req)
	return resp, err
}
