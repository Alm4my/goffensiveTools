package main

import (
	"fmt"
	"net"
)

func main() {

	for i := 0; i <= 1024; i++ {
		adresse := fmt.Sprintf("inphb.ci:%d", i)
		conn, err := net.Dial("tcp", adresse)
		if err == nil {
			fmt.Println("Le port", i, "de l'adresse", adresse, "est ouvert.")
			err := conn.Close()
			if err != nil {
				fmt.Println("Nous avons eu une erreur en fermant la connection à", adresse)
				continue
			}
		}
		// nous pouvons avoir un else ici si on veut les ports fermes ou filtrés
	}
}
