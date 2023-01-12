package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"olowe.co/dishy"
)

const usage = "usage: dishy [-a address] command"

var aFlag = flag.String("a", dishy.DefaultDishyAddr, "dishy device IP address")

func printStatus(client *dishy.Client) error {
	stat, err := client.Status()
	if err != nil {
		return fmt.Errorf("read status: %w", err)
	}
	fmt.Fprintln(os.Stderr, stat)
	fmt.Println("alerts:", stat.Alerts)
	fmt.Println("id:", stat.DeviceInfo.Id)
	fmt.Println("ready:", stat.ReadyStates)
	fmt.Println("outage:", stat.Outage)
	fmt.Println("gps:", stat.GpsStats)
	fmt.Println("stowed:", stat.StowRequested)
	fmt.Println("updates:", stat.SoftwareUpdateState)
	fmt.Println("lowsnr:", stat.IsSnrPersistentlyLow)
	return nil
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("dishy: ")

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
	case "stat":
		err = printStatus(client)
	case "metrics":
		stat, err := client.Status()
		if err != nil {
			log.Fatalln("read status:", err)
		}
		err = dishy.WriteOpenMetrics(os.Stdout, stat)
	}
	if err != nil {
		log.Fatalf("%s: %v", cmd, err)
	}
}
