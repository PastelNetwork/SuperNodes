package main

import (
	"fmt"
	"github.com/a-ok123/go-psl/internal/psl_tools"
)

func main() {

	var psl psl_tools.PslNode
	psl.Connect()

	val, err := psl.GetMNRegFee()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", val)

	/*	client := ipfs_client.IPFSClient{}

		common.Run("config.yml", func(ctx context.Context, config *common.Config){

			// connect to ipfs node and listen
			client.Connect2Node()
			client.Subscribe2Topic("mn1")
			client.ListenToTopic(ctx, "mn1", func(message string){println(message)})

			// Connect to cNode - ":7000"

			// Start Rest Listener

		})
	*/
}
