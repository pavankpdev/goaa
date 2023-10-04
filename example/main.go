package main

import (
	"context"
	"fmt"
	"github.com/pavankpdev/goaa"
)

func main() {
	var client = goaa.CreateClient("https://eth-sepolia.g.alchemy.com/v2/3QDeRBQwnvG0DyXXVyFds7I7w0f8dUfe")
	fmt.Println(client.ChainID(context.Background()))
	fmt.Println(goaa.PrivateKeyToAddress("af5ead4413ff4b78bc94191a2926ae9ccbec86ce099d65aaf469e9eb1a0fa87f"))
}
