package main
import (
	"log";"net/http";"headless-db-cms/internal"
)
func main(){ mux:=http.NewServeMux(); mux.HandleFunc("/", func(w http.ResponseWriter,r *http.Request){ w.Write([]byte("ok"))}); log.Fatal(http.ListenAndServe(":8080", internal.AuthMiddleware(mux))) }