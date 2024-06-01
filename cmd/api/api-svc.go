package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var authData = Auth{
	Login:    "Admin",
	Password: "admin123",
}

type Auth struct {
	Login    string
	Password string
}

type ResponseClient struct {
	Status  int
	Service string
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1", health)
	mux.HandleFunc("/api/v1/protected", basicAuth(protected))

	s := &http.Server{
		Addr:         ":9999",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 15,
	}

	// обработка того, что сервер успешно запущен
	log.Printf("👌 Starting server on  %s👌 ", s.Addr)

	err := s.ListenAndServe()

	if err != nil {
		log.Fatalf("Err %v\n", err.Error())
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Content-Type", "application/json")
	(w).WriteHeader(http.StatusOK)
	// базовая дессерелиазация массив байт
	dataBytes, _ := json.Marshal(ResponseClient{Status: http.StatusOK, Service: "health"})
	(w).Write(dataBytes)

}

// оборачиваем обработку для функций handle func

func basicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Content-Type", "application/json")
		decoder := json.NewDecoder(r.Body)
		var params Auth
		err := decoder.Decode(&params)

		if err != nil {
			fmt.Println("Error")
			(w).WriteHeader(http.StatusUnauthorized)
			return
		}

		// login, password, ok := r.BasicAuth()

		log.Printf("login %s\n", params.Login)
		log.Printf("password %s", params.Password)

		// if !ok {
		// 	(w).WriteHeader(http.StatusUnauthorized)
		// 	return
		// }

		if authData.Login != params.Login || authData.Password != params.Password {
			log.Println("true")
			(w).WriteHeader(http.StatusUnauthorized)
			return
		}
		(w).WriteHeader(http.StatusOK)
		handler.ServeHTTP(w, r)

	}
}

func protected(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Content-Type", "application/json")
	(w).WriteHeader(http.StatusOK)
	dataBytes, _ := json.Marshal(ResponseClient{Status: http.StatusOK, Service: "protected"})
	log.Println("protected")
	(w).Write(dataBytes)
}
