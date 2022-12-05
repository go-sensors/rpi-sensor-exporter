# go-sensors/rpi-sensor-exporter

A program to read data from sensors using the [go-sensors] libraries and reporting with Prometheus.

[go-sensors]: https://github.com/go-sensors

## Quickstart

[![Deploy with Balena](https://www.balena.io/deploy.svg)](https://dashboard.balena-cloud.com/deploy?repoUrl=https://github.com/go-sensors/rpi-sensor-exporter)

By default, no sensors are enabled. When deploying with Balena, set the `*_ENABLED` environment variables for your sensors to start reading from them.

| Environment variable                    | Description                                                                                | Default value   | Valid values                                                                                                                                                                      |
| --------------------------------------- | ------------------------------------------------------------------------------------------ | --------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `EXPORTER_METRICS_SERVER_ADDR`          | Address (host:port) to which to bind for hosting the Prometheus metrics server             | `:9000`         | Any valid host and port combination; host may be empty to bind on all available addresses                                                                                         |
| `EXPORTER_LOG_LEVEL`                    | Determines the amount of detail included in the log output                                 | `warn`          | `fatal`, `error`, `warn`, `info`, `debug`                                                                                                                                         |
| `EXPORTER_ASAIRAHT10_ENABLED`           | Enable the Asair AHT10/AHT20 temperature and relative humidity sensors                     | `false`         | `true`, `false`                                                                                                                                                                   |
| `EXPORTER_ASAIRAHT10_I2C_BUS`           | Number of I2C bus on which to find the sensor                                              | `1`             | Any valid I2C bus available to the device. [On Raspberry Pi, I2C buses `0` and `1`][pinout-i2c] are typically configured on the GPIO header                                       |
| `EXPORTER_ASAIRAHT10_I2C_ADDR`          | 7-bit I2C address of the sensor                                                            | `0x38` (dec 56) | A valid 7-bit I2C address. May be specified in decimal or in hexadecimal when prefixed with `0x`                                                                                  |
| `EXPORTER_CUBICPM1003_ENABLED`          | Enable the Cubic PM1003 particulate matter sensor for measuring air quality                | `false`         | `true`, `false`                                                                                                                                                                   |
| `EXPORTER_CUBICPM1003_DEVICE_PATH`      | Path or name of block device through which to communicate with the sensor's UART interface | `/dev/ttyAMA0`  | Any valid path to a block device where the sensor is connected. Depending on the Raspberry Pi model, [there may be two or four UARTs][pinout-uart] configured on the GPIO header. |
| `EXPORTER_PLANTOWERPMS5003_ENABLED`     | Enable the Plantower PMS5003 particulate matter sensor for measuring air quality           | `false`         | `true`, `false`                                                                                                                                                                   |
| `EXPORTER_PLANTOWERPMS5003_DEVICE_PATH` | Path or name of block device through which to communicate with the sensor's UART interface | `/dev/ttyAMA0`  | Any valid path to a block device where the sensor is connected. Depending on the Raspberry Pi model, [there may be two or four UARTs][pinout-uart] configured on the GPIO header. |
| `EXPORTER_SENSIRONSCD30_ENABLED`        | Enable the Sensiron SCD30 gas sensor                                                       | `false`         | `true`, `false`                                                                                                                                                                   |
| `EXPORTER_SENSIRONSCD30_I2C_BUS`        | Number of I2C bus on which to find the sensor                                              | `1`             | Any valid I2C bus available to the device. [On Raspberry Pi, I2C buses `0` and `1`][pinout-i2c] are typically configured on the GPIO header                                       |
| `EXPORTER_SENSIRONSCD30_I2C_ADDR`       | 7-bit I2C address of the sensor                                                            | `0x61` (dec 97) | A valid 7-bit I2C address. May be specified in decimal or in hexadecimal when prefixed with `0x`                                                                                  |
| `EXPORTER_SENSIRONSGP30_ENABLED`        | Enable the Sensiron SGP30 gas sensor                                                       | `false`         | `true`, `false`                                                                                                                                                                   |
| `EXPORTER_SENSIRONSGP30_I2C_BUS`        | Number of I2C bus on which to find the sensor                                              | `1`             | Any valid I2C bus available to the device. [On Raspberry Pi, I2C buses `0` and `1`][pinout-i2c] are typically configured on the GPIO header                                       |
| `EXPORTER_SENSIRONSGP30_I2C_ADDR`       | 7-bit I2C address of the sensor                                                            | `0x58` (dec 88) | A valid 7-bit I2C address. May be specified in decimal or in hexadecimal when prefixed with `0x`                                                                                  |

[pinout-i2c]: https://pinout.xyz/pinout/i2c
[pinout-uart]: https://pinout.xyz/pinout/uart

## Building

TBD

## Code of Conduct

We are committed to fostering an open and welcoming environment. Please read our [code of conduct](CODE_OF_CONDUCT.md) before participating in or contributing to this project.

## Contributing

We welcome contributions and collaboration on this project. Please read our [contributor's guide](CONTRIBUTING.md) to understand how best to work with us.

## License and Authors

[![Daniel James logo](https://secure.gravatar.com/avatar/eaeac922b9f3cc9fd18cb9629b9e79f6.png?size=16) Daniel James](https://github.com/thzinc)

[![license](https://img.shields.io/github/license/go-sensors/rpi-sensor-exporter.svg)](https://github.com/go-sensors/rpi-sensor-exporter/blob/master/LICENSE)
[![GitHub contributors](https://img.shields.io/github/contributors/go-sensors/rpi-sensor-exporter.svg)](https://github.com/go-sensors/rpi-sensor-exporter/graphs/contributors)

This software is made available by Daniel James under the MIT license.
