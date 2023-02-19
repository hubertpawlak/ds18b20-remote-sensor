# What is ds18b20-remote-sensor?
This is a simple program to periodically read temperatures from 1-Wire sensors and send them to a specified endpoint. Easy to run as a systemd service. Tested on Raspberry Pi Zero.

# Request content
This program send an HTTP POST request with JSON containing  multiple readings. Here is an example:
```json
[
  {
    "hw_id": "sensorId1",
    "resolution": 12,
    "temperature": 100
  },{
    "hw_id": "sensorId2",
    "resolution": 10,
    "temperature": 23
  },
]
```
- **hw_id** - 1-Wire sensor id
- **resolution** - sensor resolution in bits
- **temperature** - temperature in Celsius

# rsConfig.yaml - example
```yaml
endpoint: https://localhost/api/storeTemperature
token: AUTH-TOKEN
interval: 8000
```
- **endpoint** - POST request API URL
- **token** - `Bearer` token sent in `Authorization` header
- **interval** - time between reading all temperatures in milliseconds

# How to build?
You need [Go](https://go.dev/) to compile this project. [Install Go](https://go.dev/doc/install) before proceeding.
## Step 1 - Get all dependencies
```bash
go get
```
## Step 2 - Compile
```bash
GOOS=linux GOARCH=arm go build
```
Adjust `GOOS` and `GOARCH` if you are cross-compiling (like me) or simply run
```bash
go build
```
on your machine with sensors
## Step 3 - Run
```bash
./ds18b20-remote-sensor --help
```
I suggest running it as a systemd service with `rsConfig.yaml` to prevent leaking secrets.