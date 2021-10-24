package middleware

import (
	"encoding/json"
	"net/http"
	"server/controller"
)

func writeCommonHeaders(w http.ResponseWriter) {
	acceptedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token"
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", acceptedHeaders)
}

//Options eats options requests
func Options(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, PUT, POST")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("")
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(controller.HealthCheck())
}
