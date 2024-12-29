package routes

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func Routes() http.Handler {
	r := chi.NewRouter()

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	corsSite := os.Getenv("CORS_ALLOWED_ORIGINS")
	allowedOrigins := strings.Split(corsSite, ",")

	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           1200,
	}))

	r.Get("/img/*", proxyHandlerImg)
	r.Get("/*", proxyHandler)

	return r
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	apiURL := "https://api.mangadex.org" + r.RequestURI
	req, err := http.NewRequest(r.Method, apiURL, r.Body)
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

	for key, values := range resp.Header {
		if key == "Access-Control-Allow-Origin" || key == "Access-Control-Allow-Credentials" {
			continue
		}
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Cache-Control", "public, max-age=900, stale-while-revalidate=60")
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
}

func proxyHandlerImg(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path[len("/img/"):]

	var apiURL string
	if strings.HasPrefix(urlPath, "http://") || strings.HasPrefix(urlPath, "https://") {
		apiURL = urlPath
	} else {
		apiURL = "https://" + urlPath
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("User-Agent", "YourCustomUserAgent/Chimas")
	req.Header.Set("Accept", "image/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to perform request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		if key == "Access-Control-Allow-Origin" || key == "Access-Control-Allow-Credentials" {
			continue
		}
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, "Failed to copy response body", http.StatusInternalServerError)
	}
}
