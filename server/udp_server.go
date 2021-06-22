package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/thearyanahmed/logflow/rpc/client"
	"net"
	"os"
)


var rpcClient *client.RpcClient

func main() {

	var err error

	rpcClient, err = client.NewRpcClient()

	if err != nil {
		fmt.Printf("could not connect to rpc, %v\n",err.Error())
		os.Exit(1)
	}

	if err := serve(context.Background(), ":6060"); err != nil {
		os.Exit(1)
	}
}

// serve is capable of answering to a single client at a time
func serve(ctx context.Context, address string) error {

	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		log.Errorf("failed to UDP listen on '%s' with '%v'", address, err)
		return err
	}
	defer func() {
		if err := pc.Close(); err != nil {
			log.Errorf("failed to close packet connection with '%v'", err)
		}
	}()

	errChan := make(chan error, 1)
	// maxBufferSize specifies the size of the buffers that
	// are used to temporarily hold data from the UDP packets
	// that we receive.
	buffer := make([]byte, 2048)
	go func() {
		for {
			_, _, err := pc.ReadFrom(buffer)
			if err != nil {
				errChan <- err
				return
			}
			fmt.Printf("req : %v\n", string(buffer))

			rpcClient.Add()

			err = rpcClient.Send(string(buffer))

			if err != nil {
				return
			}
		}
	}()

	rpcClient.Wait()

	var ret error
	select {
	case <-ctx.Done():
		ret = ctx.Err()
		log.Infof("cancelled with '%v'", err)

		terminate(errChan)
	case ret = <-errChan:
		terminate(errChan) // bleh

	}

	return ret
}

func terminate(errChan chan<- error) {
	reply , err := rpcClient.Terminate()

	if err != nil {
		log.Fatalf("error reply %v\n",err.Error())
		errChan <- err
	} else {
		fmt.Printf("Reply %v\n : %v\n",reply.GetStreamedCount(),reply.GetMessage())
		fmt.Printf("Session terminated")
	}
}