package main

import (
	"fmt"
	"os"

	peasant "github.com/cafofodopirata/peasant/internal"
	gopeasant "github.com/candango/gopeasant"
)

func main() {
	ht := gopeasant.NewHttpTransport("http://localhost:8080", "Nonce")
	transport := peasant.NewCafofoTransport(ht)
	dir, err := transport.Directory()
	if err != nil {
		fmt.Printf("error getting the directory: %v", err)
		os.Exit(1)
	}
	fmt.Println(dir)

	nonce, err := transport.NewNonce()
	if err != nil {
		fmt.Printf("error getting nonce: %v", err)
		os.Exit(1)
	}
	fmt.Println(nonce)

}
