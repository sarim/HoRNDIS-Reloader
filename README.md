# HoRNDIS-Reloader
Its a little daemon that listen to USB subsystem and reload HoRNDIS kext when a usb tether device is inserted.

### WHY
Since OS X 10.11 El Capitan HoRNDIS doesn't work as is. You have to unload and reload the kext after inserting your phone. Same issue remains in macOS 10.12 Sierra.

It is painful to manually unload/load kext every time you plugin your phone for tethering. So this little utility does it for you.

### HOW

Edit the included plist file and add your phone's vendor and product id to it. You can add multiple device ids. Also change the Path to where your `HoRNDIS-Reloader` binary is located.

```xml
<array>
    <string>/Users/FooBar/HoRNDIS-Reloader/HoRNDIS-Reloader</string>
    <string>0x2717:0x0100</string>
    <string>0x2717:0x0200</string>
    <string>0x2717:0x0300</string>
</array>
```
You can find your phone's device ids in `About This Mac` -> `System Information` -> `USB`.

Now copy it to /Library/LaunchDaemons enable the daemon.

```sh
sudo cp com.sarimkhan.horndis.reloader.plist /Library/LaunchDaemons/
sudo launchctl load /Library/LaunchDaemons/com.sarimkhan.horndis.reloader.plist
```

It will run after every boot. You can check its log in Console app.


### Screenshot

![HoRNDIS-Reloader ScreenShot](https://raw.githubusercontent.com/sarim/HoRNDIS-Reloader/screenshot/HoRNDIS-Reloader%20ScreenShot.png)


### Source

Checkout the repo and run
```
go build
```

You need to have xcode installed as it uses CoreFoundation and IOKit libraries via cgo.


### Licence
Licensed under Mozilla Public License 1.1 ("MPL"), an open source/free software license.