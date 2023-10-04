package main

import (
	"context"
	"fmt"
	"github.com/pavankpdev/goaa"
)

func main() {
	var client = goaa.CreateClient("https://eth-sepolia.g.alchemy.com/v2/3QDeRBQwnvG0DyXXVyFds7I7w0f8dUfe")
	fmt.Println(client.ChainID(context.Background()))
}
