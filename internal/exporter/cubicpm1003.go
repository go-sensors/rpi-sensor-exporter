package exporter

import (
	"github.com/go-sensors/core/pm"
	"github.com/go-sensors/cubicpm1003"
	"github.com/go-sensors/prometheus"
	"github.com/go-sensors/rpi-sensor-exporter/internal/log"
	"github.com/go-sensors/serial"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syncromatics/go-kit/v2/cmd"
)

func RegisterCubicPM1003(cmd *cobra.Command) {
	cmd.Flags().Bool("cubicpm1003-enabled", false, "Enable the Cubic PM1003 particulate matter sensor for measuring air quality")
	cmd.Flags().String("cubicpm1003-device-path", "/dev/serial0", "Path or name of block device through which to communicate with the sensor's UART interface")
}

func TryStartCubicPM1003(group *cmd.ProcessGroup) error {
	if !viper.GetBool("cubicpm1003-enabled") {
		log.Info("Skipping Cubic PM1003")
		return nil
	}

	log.Info("Starting Cubic PM1003")

	serialPortFactory, err := serial.NewSerialPort(
		viper.GetString("cubicpm1003-device-path"),
		cubicpm1003.GetDefaultSerialPortConfig())
	if err != nil {
		return errors.Wrap(err, "failed to initialize serial port factory")
	}

	sensor := cubicpm1003.NewSensor(serialPortFactory,
		cubicpm1003.WithRecoverableErrorHandler(shouldTerminate))

	metricHandler := prometheus.NewMetricHandlerWithLabels(&prometheus.Labels{
		Source: "cubicpm1003",
	})
	consumer := pm.NewConsumer(sensor, metricHandler)

	group.Start(consumer.Run)
	group.Start(sensor.Run)

	return nil
}
