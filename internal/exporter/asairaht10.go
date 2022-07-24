package exporter

import (
	"github.com/go-sensors/asairaht10"
	"github.com/go-sensors/core/humidity"
	"github.com/go-sensors/core/i2c"
	"github.com/go-sensors/core/temperature"
	"github.com/go-sensors/prometheus"
	"github.com/go-sensors/rpi-sensor-exporter/internal/log"
	"github.com/go-sensors/rpii2c"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syncromatics/go-kit/v2/cmd"
)

func RegisterAsairAHT10(cmd *cobra.Command) {
	cmd.Flags().Bool("asairaht10-enabled", false, "Enable the Asair AHT10/AHT20 temperature and relative humidity sensors")
	cmd.Flags().Int("asairaht10-i2c-bus", 1, "Number of I2C bus on which to find the sensor")
	cmd.Flags().Uint8("asairaht10-i2c-addr", asairaht10.GetDefaultI2CPortConfig().Address, "7-bit I2C address of the sensor")
}

func TryStartAsairAHT10(group *cmd.ProcessGroup) error {
	if !viper.GetBool("asairaht10-enabled") {
		log.Info("Skipping Asair AHT10")
		return nil
	}

	log.Info("Starting Asair AHT10")

	i2cPortFactory, err := rpii2c.NewI2CPort(
		viper.GetInt("asairaht10-i2c-bus"),
		&i2c.I2CPortConfig{
			Address: byte(viper.GetUint("asairaht10-i2c-addr")),
		})
	if err != nil {
		return errors.Wrap(err, "failed to initialize I2C port factory")
	}

	sensor := asairaht10.NewSensor(i2cPortFactory,
		asairaht10.WithRecoverableErrorHandler(shouldTerminate))

	metricHandler := prometheus.NewMetricHandlerWithLabels(&prometheus.Labels{
		Source: "asairaht10",
	})
	temperatureConsumer := temperature.NewConsumer(sensor, metricHandler)
	humidityConsumer := humidity.NewConsumer(sensor, metricHandler)

	group.Start(temperatureConsumer.Run)
	group.Start(humidityConsumer.Run)
	group.Start(sensor.Run)

	return nil
}
