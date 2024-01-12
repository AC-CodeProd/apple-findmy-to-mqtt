# Apple FindMy to MQTT
A go script that reads local FindMy cache files to publish device locations and metadata (including those of AirTags, AirPods, Apple Watches, iPhones) to MQTT.

- [Apple FindMy to MQTT](#apple-findmy-to-mqtt)
  - [Disclaimer](#disclaimer)
  - [Description](#description)
  - [Prerequisites](#prerequisites)
  - [Configuration](#configuration)
  - [Run](#run)
    - [Development](#development)
    - [Production](#production)
  - [Credits](#credits)
---
## Disclaimer
This script is provided as-is, without any warranty. Use at your own risk. This code is not tested and should only be used for experimental purposes. Loading the FindMy cache files is not intended by Apple and might cause problems. The project is not affiliated to Apple Inc., MQTT.

## Description
This go script reads FindMy cache files and publishes location data over MQTT. The script must be run on macOS with a FindMy installation running. It must be run in a terminal with full disk access in order to read the cache files.

## Prerequisites
- Golang 1.20 or higher is required to build and run the project. You can find the installer on
  the official Golang [download](https://go.dev/doc/install) page.

## Configuration

This application uses a JSON configuration file to set up various parameters. Here's a brief overview of the configuration fields:

| Key | Description | Default Value |
| --- | ----------- | ------------- |
| `environment` | Sets the environment mode for the application. | `development` |
| `log_output` | Specifies the file where the application's logs will be written. | `logs/development.log` |
| `log_level` | Sets the level of logs that will be written. | `debug` |
| `tz` | Sets the timezone for the application. | `Europe/Paris` |
| `scan_timer` | Sets the interval (in seconds) at which the script scans the Apple FindMy cache. | `5` |

You should adjust these settings according to your needs and environment. Please ensure to replace all the placeholders with your actual data.
## Run

### Development
Run the following
```sh
$ git clone git@github.com:AC-CodeProd/apple-findmy-to-mqtt.git
$ cd apple-findmy-to-mqtt
```
With docker-compose
```sh
$ docker-compose -f docker-compose.yml up -d
```
OR Run the Go program in live-reload mode using the 'air'
```sh
$ make run-live
```
### Production
Download binaries for Linux, macOS, and Windows are available as [Github Releases](https://github.com/AC-CodeProd/apple-findmy-to-mqtt/releases/latest).
Using binarie:
```sh
$ apple-findmy-to-mqtt-v{{VERSION}}-{{ARCHITECTURES}} scan
```
OR Build binarie
```sh
$ git clone git@github.com:AC-CodeProd/apple-findmy-to-mqtt.git
$ cd apple-findmy-to-mqtt
$ make build
$ cd build
$ apple-findmy-to-mqtt-v{{VERSION}}-{{ARCHITECTURES}} scan
```
## Credits
This work was inspired by <a href="https://github.com/muehlt/home-assistant-findmy" target="_blank">muehlt</a>