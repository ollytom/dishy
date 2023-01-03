package main

import (
	"flag"
	"log"

	"olowe.co/dishy"
)

const usage = "usage: dishy [-a address] command"

var aFlag = flag.String("a", dishy.DefaultDishyAddr, "dishy device IP address")

func main() {
	log.SetFlags(0)
	log.SetPrefix("dishy:")

	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal(usage)
	}
	cmd := flag.Args()[0]
	addr := *aFlag
	client, err := dishy.Dial(addr)
	if err != nil {
		log.Fatalf("dial %s: %v", addr, err)
	}

	switch cmd {
	default:
		log.Fatalf("unknown command %s", cmd)
	case "reboot":
		err = client.Reboot()
	case "stow":
		err = client.Stow()
	case "unstow":
		err = client.Unstow()
	}
	if err != nil {
		log.Fatalf("%s: %v", cmd, err)
	}
}
