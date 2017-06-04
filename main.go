package main

import (
	"C"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)
import "strconv"

const Kext = "/System/Library/Extensions/HoRNDIS.kext"

var c = make(chan int, 10)

//export DeviceAddedCB
func DeviceAddedCB(vendorID, productID int32, name *C.char) {
	log.Printf("Inserted Device - Name: %s, Vendor: 0x%04x, Product: 0x%04x\n", C.GoString(name), vendorID, productID)
	c <- 1
}

func DeviceAddedDebounced() {
	var ct <-chan time.Time
	for {
		select {
		case <-c:
			ct = time.After(time.Second)
		case <-ct:
			go func() {
				time.Sleep(time.Second)
				err := exec.Command("/sbin/kextunload", Kext).Run()
				if err != nil {
					log.Println("Unloading Kext Failed")
				} else {
					log.Println("Unloading Kext Success")
				}
				time.Sleep(time.Second)
				err = exec.Command("/sbin/kextload", Kext).Run()
				if err != nil {
					log.Println("Loading Kext Failed")
				} else {
					log.Println("Loading Kext Success")
				}
			}()
		}
	}
}

func main() {

	if os.Geteuid() != 0 {
		log.Fatal("Root needed for kextload/unload")
	}

	go DeviceAddedDebounced()

	InitRunLoop()

	if len(os.Args) < 2 {
		log.Fatal("At least one vendor:product needed as argument")
	}
	args := os.Args[1:]
	for _, v := range args {
		a := strings.Split(v, ":")
		if len(a) != 2 {
			log.Fatal("Invalid vendor:product arg - %s", v)
		}
		vendorID, err := strconv.ParseInt(a[0], 0, 32)
		if err != nil {
			log.Fatal("Invalid vendor arg - %s", a[0])
		}
		productID, err := strconv.ParseInt(a[1], 0, 32)
		if err != nil {
			log.Fatal("Invalid product arg - %s", a[1])
		}

		log.Printf("Listening for device 0x%04x:0x%04x", vendorID, productID)
		AddDeviceMatch(int(vendorID), int(productID))
	}

	log.Println("HoRNDIS Reloader by Gittu Started")

	RunLoopRun()
}
