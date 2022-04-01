package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	var (
		LocalReplicaNode string
		SequencerNode    string
		TimeoutSecond    int64
	)

	flag.StringVar(&LocalReplicaNode, "replica", "http://l2geth:8545", "local replica rpc endpoint")
	flag.StringVar(&SequencerNode, "sequencer", "https://andromeda.metis.io/?owner=1088", "sequencer rpc endpoint")
	flag.Int64Var(&TimeoutSecond, "timeout", 5, "request timeout")
	flag.Parse()

	sequencer, err := rpc.Dial(SequencerNode)
	if err != nil {
		panic(err)
	}
	defer sequencer.Close()

	replica, err := rpc.Dial(LocalReplicaNode)
	if err != nil {
		panic(err)
	}
	defer replica.Close()

	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if req.Method == http.MethodGet {
			return
		}
		newctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(TimeoutSecond))
		defer cancel()

		w.Header().Set("Content-Type", "application/json")
		var reqdata JsonRequest
		_ = json.NewDecoder(req.Body).Decode(&reqdata)

		var result interface{}
		var err error
		if reqdata.Method == "eth_sendRawTransaction" {
			err = sequencer.CallContext(newctx, &result, reqdata.Method, reqdata.Params...)
		} else {
			err = replica.CallContext(newctx, &result, reqdata.Method, reqdata.Params...)
		}
		_ = json.NewEncoder(w).Encode(&JsonResponse{
			ID: reqdata.ID, Version: reqdata.Version,
			Result: result, Error: UnwrapErrorMessage(err),
		})
	})

	server := http.Server{Addr: ":8545", Handler: handler}
	go func() {
		stopSig := make(chan os.Signal, 1)
		signal.Notify(stopSig, syscall.SIGINT, syscall.SIGTERM)
		<-stopSig
		log.Println("graceful stopping")
		_ = server.Shutdown(context.TODO())
	}()

	log.Println("replica proxy server is starting")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalln("server starts failed", err)
	}
}
