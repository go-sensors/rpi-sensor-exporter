package exporter

import (
	"github.com/go-sensors/core/gas"
	"github.com/go-sensors/core/humidity"
	"github.com/go-sensors/core/i2c"
	"github.com/go-sensors/core/temperature"
	"github.com/go-sensors/prometheus"
	"github.com/go-sensors/rpi-sensor-exporter/internal/log"
	"github.com/go-sensors/rpii2c"
	"github.com/go-sensors/sensironscd30"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syncromatics/go-kit/v2/cmd"
)

func RegisterSensironSCD30(cmd *cobra.Command) {
	cmd.Flags().Bool("sensironscd30-enabled", false, "Enable the Sensiron SCD30 gas sensor")
	cmd.Flags().Int("sensironscd30-i2c-bus", 1, "Number of I2C bus on which to find the sensor")
	cmd.Flags().Uint8("sensironscd30-i2c-addr", sensironscd30.GetDefaultI2CPortConfig().Address, "7-bit I2C address of the sensor")
}

func TryStartSensironSCD30(group *cmd.ProcessGroup) error {
	if !viper.GetBool("sensironscd30-enabled") {
		log.Info("Skipping Sensiron SCD30")
		return nil
	}

	log.Info("Starting Sensiron SCD30")

	i2cPortFactory, err := rpii2c.NewI2CPort(
		viper.GetInt("sensironscd30-i2c-bus"),
		&i2c.I2CPortConfig{
			Address: byte(viper.GetUint("sensironscd30-i2c-addr")),
		})
	if err != nil {
		return errors.Wrap(err, "failed to initialize I2C port factory")
	}

	sensor := sensironscd30.NewSensor(i2cPortFactory,
		sensironscd30.WithRecoverableErrorHandler(shouldTerminate))

	metricHandler := prometheus.NewMetricHandlerWithLabels(&prometheus.Labels{
		Source: "sensironscd30",
	})
	gasConsumer := gas.NewConsumer(sensor, metricHandler)
	temperatureConsumer := temperature.NewConsumer(sensor, metricHandler)
	humidityConsumer := humidity.NewConsumer(sensor, metricHandler)

	group.Start(gasConsumer.Run)
	group.Start(temperatureConsumer.Run)
	group.Start(humidityConsumer.Run)
	group.Start(sensor.Run)

	return nil
}
