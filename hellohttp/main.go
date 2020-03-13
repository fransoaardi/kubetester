package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/hellohttp", func(w http.ResponseWriter, r *http.Request) {
		version := "v1-hellohttp"
		hostname, _ := os.Hostname()

		name := r.URL.Query().Get("name")
		type Introduction struct{
			Name string `json:"name"`
			Version string `json:"version"`
			Hostname string `json:"hostname"`
		}
		output := Introduction{
			Name:     name,
			Version:  version,
			Hostname: hostname,
		}

		out, _ := json.Marshal(output)
		w.Write(out)
	})

	server := http.Server{
		Addr: "0.0.0.0:8100",
		Handler: mux,
	}

	server.ListenAndServe()
}
