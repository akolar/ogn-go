package main

import (
	"context"
	"fmt"

	"github.com/akolar/ogn-go"
)

func main() {
	serv := ogn.NewServer("aprs.glidernet.org", 10152)
	settings := ogn.NewSettings("a", "", "")

	aprs := ogn.Connect(serv, settings, true)
	ctx := context.Background()
	ctx, _ = context.WithCancel(ctx)

	go func() {
		err := aprs.Receive(ctx)
		if err != nil {
			fmt.Printf("An error occured: %s\n", err)
			return
		}
	}()

	for msg := aprs.Read(); ; msg = aprs.Read() {
		fmt.Println(msg)
	}
}
