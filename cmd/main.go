package main

import (
	"context"
	"os"
	"strings"

	"github.com/go-sensors/core/humidity"
	"github.com/go-sensors/prometheus"
	"github.com/go-sensors/rpi-sensor-exporter/internal/exporter"
	"github.com/go-sensors/rpi-sensor-exporter/internal/log"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syncromatics/go-kit/v2/cmd"
)

var (
	rootCmd = &cobra.Command{
		Use:   "rpi-sensor-exporter",
		Short: "Reads data from sensors and reports using Prometheus",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			isTerminal := isatty.IsTerminal(os.Stdout.Fd())
			log.InitializeLogger(isTerminal, viper.GetString("log-level"))
			cmd.SilenceErrors = !isTerminal
			cmd.SilenceUsage = !isTerminal
		},
		RunE: func(*cobra.Command, []string) error {
			group := cmd.NewProcessGroup(context.Background())

			metricsServer := prometheus.NewMetricsServer(viper.GetString("metrics-server-addr"))
			group.Start(metricsServer.Run)

			var err error
			humidityHandlers := []humidity.Handler{}
			var humidityHandler humidity.Handler
			err, humidityHandler = exporter.TryStartSensironSGP30(group)
			if err != nil {
				return err
			}

			if humidityHandler != nil {
				humidityHandlers = append(humidityHandlers, humidityHandler)
			}

			err = exporter.TryStartSensironSCD30(group)
			if err != nil {
				return err
			}

			err = exporter.TryStartAsairAHT10(group, humidityHandlers)
			if err != nil {
				return err
			}

			err = exporter.TryStartCubicPM1003(group)
			if err != nil {
				return err
			}

			err = exporter.TryStartPlantowerPMS5003(group)
			if err != nil {
				return err
			}

			return group.Wait()
		},
	}
)

func init() {
	rootCmd.Flags().String("log-level", "warn", "Determines the amount of detail included in the log output; valid options are: fatal, error, warn, info, debug")
	rootCmd.Flags().String("metrics-server-addr", ":9000", "Address (host:port) to which to bind for hosting the Prometheus metrics server")

	exporter.RegisterAsairAHT10(rootCmd)
	exporter.RegisterCubicPM1003(rootCmd)
	exporter.RegisterPlantowerPMS5003(rootCmd)
	exporter.RegisterSensironSCD30(rootCmd)
	exporter.RegisterSensironSGP30(rootCmd)

	viper.SetEnvPrefix("EXPORTER")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	viper.BindPFlags(rootCmd.Flags())
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("failed to terminate cleanly",
			"err", err)
	}
}
