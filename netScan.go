package main

import (
	"fmt"
	"net"
)

func main() {
	_, err := net.Dial("tcp", "google.com:443")
	if err == nil {
		fmt.Println("Connexion Établie!")
	} else {
		fmt.Println("Nous n'avons pas pu nous connecter au serveur!")
	}
}
