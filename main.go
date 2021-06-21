package main

import (
	"flag"
	"fmt"
	"github.com/thearyanahmed/logflow/collectors/file"
	"github.com/thearyanahmed/logflow/rpc/client"
	"github.com/thearyanahmed/logflow/rpc/server"
	"github.com/thearyanahmed/logflow/utils/env"
)

var (
	action string
)

func main()  {
	env.LoadEnv()

	flag.StringVar(&action, "action", "help", "start rpc server. available option: server, client")

	flag.Parse()

	switch action {
	case "serve":
		server.Run()
	case "client":

		// create a client instance
		// when creating an instance, it will create a grpc dial
		// which will hold the connect
		// a lock and possibly an wg

		// client should have a function


		client.Run()
	case "collector:file":
		file.Run(env.Get("TEST_DATA_FILE"))
	case "help":
		printHelp()
	default:
		fmt.Printf("You did not select any proper option.quiting\n")
	}
}

func printHelp() {
	fmt.Printf(`
	
About : Logflow reads nginx log data and sends to kafka using grpc. 
	
Dependencies: Logflow requires kafka to be install and running on your machine. Make sure .env file is setup correctly.

Available commands    
------------------------------------------------------------------------------------------------
serve				: Will star the gRPC server
client			 	: Will start a gRPC client that reads from nginx log (specified in the .env) and streams to gRPC server. Everything is subject to change.

help				: Displayes help menu.	

How to use : For now, you can use it like
"go run main.go --action serve"
`)
}