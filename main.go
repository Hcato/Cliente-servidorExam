package main

import (
	servidorp "examen1/servidor_p"
	servidorr "examen1/servidor_r"
)

func main() {
	go servidorp.Run()
	go servidorr.Run()
	select {}
}
