package routes

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
		ExposedHeaders: []string{"Link"},
		// AllowCredentials: true,
		MaxAge: 1200,
	}))

	r.Get("/img/*", proxyHandlerImg)
	r.Get("/*", proxyHandler)
	return r
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	apiURL := "https://api.mangadex.org" + r.RequestURI
	log.Println("API URL:", apiURL)

	req, err := http.NewRequest(r.Method, apiURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to perform request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	var reader io.Reader = resp.Body
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			http.Error(w, "Failed to decompress response", http.StatusInternalServerError)
			return
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	io.Copy(w, reader)
}

func proxyHandlerImg(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path[len("/img/"):]

	log.Println("FIRS>>>>", urlPath)

	var apiURL string
	if strings.HasPrefix(urlPath, "http://") || strings.HasPrefix(urlPath, "https://") {
		apiURL = urlPath
	} else {
		apiURL = "https://" + urlPath
	}

	log.Println("Image API URL:", apiURL)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to perform request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	io.Copy(w, resp.Body)
}
