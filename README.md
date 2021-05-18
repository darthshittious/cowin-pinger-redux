## Cowin Alerts

Cowin Alerts

## Build

```
go get -d ./...
go build main/main.go
```

```
make <build_type>
```
e.g.

```
make windows_64
```

## Instruction to Execute

(Replace `main-windows` with your build below, if necessary)

* If you want to search by pincode ``./bin/main-windows.exe 1 03-05-2021 45 560078 true 5m``

* If you want to search by district id ``./bin/main-windows.exe 2 03-05-2021 45 265 false``

* first Argument - 1: search by pincode 2: search by district id
* second argument - date dd-mm-yyyy
* third argument - minimum age limit - 18/45
* fourth argument - pincode/district id
* fifth argument - whether to ping constantly: `true/false`
* last argument - if pinging constantly, interval between pings - optional

* For State ID - https://cdn-api.co-vin.in/api/v2/admin/location/states
* For District ID - https://cdn-api.co-vin.in/api/v2/admin/location/districts/`stateid`

Alternatively, district names and IDs are mapped out in [districts.csv](districts.csv)
## Output
Console output will look like this:
```
Center Name Manipal Clinic
No Slots Available for given date | pincode | district !
No Slots Available for given date | pincode | district !
--------------------
Center Name Jarganahalli Corporator Office
Date 03-05-2021
Slot Information [09:00AM-11:00AM 11:00AM-01:00PM 01:00PM-03:00PM 03:00PM-06:00PM]
Count Available: 6 Minimum Age: 45
--------------------
```

You will also get OS level alerts

## Issues/PRs

Go crazy. I'm sure a lot of you know this stuff better than me plus this reduced thing was born out of 20
minutes of furious typing so there's a high chance it might suck. But in my defence, it works, and is pretty
noob-friendly.
