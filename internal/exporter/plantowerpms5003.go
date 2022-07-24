package exporter

import (
	"github.com/go-sensors/core/pm"
	"github.com/go-sensors/plantowerpms5003"
	"github.com/go-sensors/prometheus"
	"github.com/go-sensors/rpi-sensor-exporter/internal/log"
	"github.com/go-sensors/serial"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syncromatics/go-kit/v2/cmd"
)

func RegisterPlantowerPMS5003(cmd *cobra.Command) {
	cmd.Flags().Bool("plantowerpms5003-enabled", false, "Enable the Plantower PMS5003 particulate matter sensor for measuring air quality")
	cmd.Flags().String("plantowerpms5003-device-path", "/dev/serial0", "Path or name of block device through which to communicate with the sensor's UART interface")
}

func TryStartPlantowerPMS5003(group *cmd.ProcessGroup) error {
	if !viper.GetBool("plantowerpms5003-enabled") {
		log.Info("Skipping Plantower PMS5003")
		return nil
	}

	log.Info("Starting Plantower PMS5003")

	serialPortFactory, err := serial.NewSerialPort(
		viper.GetString("plantowerpms5003-device-path"),
		plantowerpms5003.GetDefaultSerialPortConfig())
	if err != nil {
		return errors.Wrap(err, "failed to initialize serial port factory")
	}

	sensor := plantowerpms5003.NewSensor(serialPortFactory,
		plantowerpms5003.WithRecoverableErrorHandler(shouldTerminate))

	metricHandler := prometheus.NewMetricHandlerWithLabels(&prometheus.Labels{
		Source: "plantowerpms5003",
	})
	consumer := pm.NewConsumer(sensor, metricHandler)

	group.Start(consumer.Run)
	group.Start(sensor.Run)

	return nil
}
