package main

import (
	"fmt"
	"net"
	"sync"
)

func main() {

	var waitGroup sync.WaitGroup

	for i := 0; i <= 65535; i++ {
		waitGroup.Add(1)
		go func(j int) {
			defer waitGroup.Done()
			adresse := fmt.Sprintf("127.0.0.1:%d", j)
			conn, err := net.Dial("tcp", adresse)
			if err == nil {
				fmt.Println("Le port", j, "de l'adresse", adresse, "est ouvert.")
				err := conn.Close()
				if err != nil {
					fmt.Println("Nous avons eu une erreur en fermant la connection à", adresse)
				}
			}
			// nous pouvons avoir un else ici si on veut les ports fermes ou filtrés
		}(i)
	}
	waitGroup.Wait()
}
