# What is ds18b20-remote-sensor?
This is a simple program to periodically read temperatures from 1-Wire sensors and send them to a specified endpoint. Easy to run as a systemd service.

# Request content
This program send an HTTP POST request with JSON containing  multiple readings. Here is an example:
```json
[
  {
    "hwId": "sensorId1",
    "resolution": 12,
    "temperature": 100
  },{
    "hwId": "sensorId2",
    "resolution": 10,
    "temperature": 23
  },
]
```
- **hwId** - 1-Wire sensor id
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

# How to use?
## Step 1 - Build
```bash
GOOS=linux GOARCH=arm go build
```
Adjust `GOOS` and `GOARCH` if you are cross-compiling or simply run
```bash
go build
```
on your machine with sensors
## Step 2 - Use the compiled binary
```bash
./ds18b20-remote-sensor --help
```
I suggest running it as a systemd service with `rsConfig.yaml` to prevent leaking secrets.