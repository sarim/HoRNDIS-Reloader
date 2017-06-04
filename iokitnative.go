package main

/*
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit

#include <CoreFoundation/CoreFoundation.h>
#include <IOKit/IOKitLib.h>
#include <IOKit/IOMessage.h>
#include <IOKit/IOCFPlugIn.h>
#include <IOKit/usb/IOUSBLib.h>
#include <os/log.h>

// This vars are here in C because Go can't call c MACRO.
// Also CFSTR need compile time string.
CFStringRef CFSTR_USBVendorID = CFSTR(kUSBVendorID);
CFStringRef CFSTR_USBProductID = CFSTR(kUSBProductID);

IONotificationPortRef gNotifyPort;
io_iterator_t gAddedIter;
CFRunLoopRef gRunLoop;
CFRunLoopSourceRef runLoopSource;

extern void DeviceAddedCB(int32_t p0, int32_t p1, char* p2);

void macLog(char *txt) {
    os_log(OS_LOG_DEFAULT, "%{public}s", txt);
}

void initRunLoop()
{
    gNotifyPort = IONotificationPortCreate(kIOMasterPortDefault);
    runLoopSource = IONotificationPortGetRunLoopSource(gNotifyPort);

    gRunLoop = CFRunLoopGetCurrent();
    CFRunLoopAddSource(gRunLoop, runLoopSource, kCFRunLoopDefaultMode);
}

void deviceAdded(void* refCon, io_iterator_t iterator)
{
    kern_return_t kr;
    io_service_t usbDevice;
    while ((usbDevice = IOIteratorNext(iterator))) {
        io_name_t deviceName;
        int32_t vendor, product;

        kr = IORegistryEntryGetName(usbDevice, deviceName);
        if (KERN_SUCCESS != kr) {
            deviceName[0] = '\0';
        }

        CFTypeRef vendorID = IORegistryEntryCreateCFProperty(usbDevice, CFSTR_USBVendorID, kCFAllocatorDefault, 0);
        CFTypeRef productID = IORegistryEntryCreateCFProperty(usbDevice, CFSTR_USBProductID, kCFAllocatorDefault, 0);

        CFNumberGetValue(vendorID, kCFNumberSInt32Type, &vendor);
        CFNumberGetValue(productID, kCFNumberSInt32Type, &product);

        CFRelease(vendorID);
        CFRelease(productID);

        kr = IOObjectRelease(usbDevice);

        DeviceAddedCB(vendor, product, deviceName);
    }
}

bool addDeviceMatch(CFMutableDictionaryRef matchingDict)
{
    kern_return_t kr = IOServiceAddMatchingNotification(gNotifyPort,
                                                        kIOFirstMatchNotification,
                                                        matchingDict,
                                                        deviceAdded,
                                                        NULL,
                                                        &gAddedIter);
    return kr == KERN_SUCCESS;
}

*/
import "C"
import (
	"log"
	"unsafe"
)

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	C.macLog(C.CString(string(bytes)))
	return len(bytes), nil
}

func init() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
}

func InitRunLoop() {
	C.initRunLoop()
}

func RunLoopRun() {
	C.deviceAdded(nil, C.gAddedIter)
	C.CFRunLoopRun()
}

func AddDeviceMatch(vendorID, productID int) {

	matchingDict := C.IOServiceMatching(C.CString(C.kIOUSBDeviceClassName))

	numberRef := C.CFNumberCreate(C.kCFAllocatorDefault, C.kCFNumberSInt32Type,
		unsafe.Pointer(&vendorID))

	if vendorID != 0 {
		C.CFDictionarySetValue(matchingDict,
			unsafe.Pointer(C.CFSTR_USBVendorID),
			unsafe.Pointer(numberRef))
	}

	C.CFRelease(numberRef)

	numberRef = C.CFNumberCreate(C.kCFAllocatorDefault, C.kCFNumberSInt32Type,
		unsafe.Pointer(&productID))

	if productID != 0 {
		C.CFDictionarySetValue(matchingDict,
			unsafe.Pointer(C.CFSTR_USBProductID),
			unsafe.Pointer(numberRef))
	}

	C.CFRelease(numberRef)

	numberRef = nil

	r := C.addDeviceMatch(matchingDict)

	if r == false {
		log.Panic("IOServiceAddMatchingNotification Failed")
	}
}
