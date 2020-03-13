package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "github.com/fransoaardi/helloserve/proto"
)

type Output struct{
	Name string `json:"name"`
	Version string `json:"version"`
	Hostname string `json:"hostname"`
	Response interface{} `json:"response"`
}

func main(){
	mux := http.NewServeMux()
	version := "v1-helloserve"
	hostname, _ := os.Hostname()

	mux.HandleFunc("/hellogrpc", func(w http.ResponseWriter, r *http.Request){
		name := r.URL.Query().Get("name")

		conn, err := grpc.Dial("grpc-svc.default.svc.cluster.local",
			grpc.WithInsecure(), grpc.WithBlock(), grpc.WithBalancerName("round_robin"))

		//conn, err := grpc.Dial("localhost:8000",
		//	grpc.WithInsecure(), grpc.WithBlock(), grpc.WithBalancerName("round_robin"))

		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}

		defer conn.Close()
		c := pb.NewHelloClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		greeting := pb.Greeting{}
		greeting.Name = name

		introduction, err := c.SayHello(ctx, &greeting)
		if err != nil {
			log.Fatal(err)
		}

		out := Output{
			Name: name,
			Version: version,
			Hostname:hostname,
			Response: introduction,
		}

		write, _ := json.Marshal(out)
		w.Write(write)
	})

	mux.HandleFunc("/hellohttp", func(w http.ResponseWriter, r *http.Request) {


		name := r.URL.Query().Get("name")

		cli := http.Client{}
		req, err := http.NewRequest(http.MethodGet, "http://http-svc.default.svc.cluster.local/hellohttp", nil)
		if err != nil {
			fmt.Println(err)
		}
		resp, err := cli.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		read, _ := ioutil.ReadAll(resp.Body)
		var respBody interface{}
		json.Unmarshal(read, &respBody)

		out := Output{
			Name: name,
			Version: version,
			Hostname:hostname,
			Response: respBody,
		}

		write, _ := json.Marshal(out)
		w.Write(write)
	})

	server := http.Server{
		Addr: "0.0.0.0:8200",
		Handler: mux,
	}

	server.ListenAndServe()
}
