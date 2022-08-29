package exporter

import (
	"github.com/go-sensors/core/gas"
	"github.com/go-sensors/core/humidity"
	"github.com/go-sensors/core/i2c"
	"github.com/go-sensors/prometheus"
	"github.com/go-sensors/rpi-sensor-exporter/internal/log"
	"github.com/go-sensors/rpii2c"
	"github.com/go-sensors/sensironsgp30"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syncromatics/go-kit/v2/cmd"
)

func RegisterSensironSGP30(cmd *cobra.Command) {
	cmd.Flags().Bool("sensironsgp30-enabled", false, "Enable the Sensiron SGP30 gas sensor")
	cmd.Flags().Int("sensironsgp30-i2c-bus", 1, "Number of I2C bus on which to find the sensor")
	cmd.Flags().Uint8("sensironsgp30-i2c-addr", sensironsgp30.GetDefaultI2CPortConfig().Address, "7-bit I2C address of the sensor")
}

func TryStartSensironSGP30(group *cmd.ProcessGroup) (error, humidity.Handler) {
	if !viper.GetBool("sensironsgp30-enabled") {
		log.Info("Skipping Sensiron SGP30")
		return nil, nil
	}

	log.Info("Starting Sensiron SGP30")

	i2cPortFactory, err := rpii2c.NewI2CPort(
		viper.GetInt("sensironsgp30-i2c-bus"),
		&i2c.I2CPortConfig{
			Address: byte(viper.GetUint("sensironsgp30-i2c-addr")),
		})
	if err != nil {
		return errors.Wrap(err, "failed to initialize I2C port factory"), nil
	}

	sensor := sensironsgp30.NewSensor(i2cPortFactory,
		sensironsgp30.WithRecoverableErrorHandler(shouldTerminate))

	metricHandler := prometheus.NewMetricHandlerWithLabels(&prometheus.Labels{
		Source: "sensironsgp30",
	})
	consumer := gas.NewConsumer(sensor, metricHandler)

	group.Start(consumer.Run)
	group.Start(sensor.Run)

	return nil, sensor
}
