package main

import (
	"C"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const Kext = "/System/Library/Extensions/HoRNDIS.kext"

var c = make(chan int, 10)

//export DeviceAddedCB
func DeviceAddedCB(vendorID, productID int32, name *C.char) {
	fmt.Printf("Name: %s, Vendor: %d, Product %d\n", C.GoString(name), vendorID, productID)
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
				fmt.Println("Unloading Kext")
				exec.Command("/sbin/kextunload", Kext).CombinedOutput()
				time.Sleep(time.Second)
				fmt.Println("Loading Kext")
				exec.Command("/sbin/kextload", Kext).CombinedOutput()
			}()
		}
	}
}

func main() {
	vendorID := 0x2717
	productID := 0x0388

	if os.Geteuid() != 0 {
		panic("Root needed for kextload/unload")
	}

	go DeviceAddedDebounced()

	fmt.Println("HoRNDIS Reloader by Gittu")

	InitRunLoop()

	AddDeviceMatch(vendorID, productID)
	// AddDeviceMatch(0, 0)

	RunLoopRun()
}
