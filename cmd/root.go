package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.infratographer.com/x/goosex"
	"go.infratographer.com/x/loggingx"
	"go.infratographer.com/x/versionx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/metal-toolbox/fleetdb/db"
	"github.com/metal-toolbox/fleetdb/internal/config"
)

var (
	appName = "fleetdb"
	cfgFile string
	logger  *zap.SugaredLogger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fleetdb",
	Short: "Server Service for Hollow ecosystem",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	// Logging flags
	loggingx.MustViperFlags(viper.GetViper(), rootCmd.PersistentFlags())

	// Register version command
	versionx.RegisterCobraCommand(rootCmd, func() { versionx.PrintVersion(logger) })

	// Setup migrate command
	goosex.RegisterCobraCommand(rootCmd, func() {
		goosex.SetBaseFS(db.Migrations)
		goosex.SetDBURI(config.AppConfig.CRDB.URI)
		goosex.SetLogger(logger)
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("fleetdb")
	viper.AutomaticEnv() // read in environment variables that match

	logger = initLogger()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			logger.With(
				zap.Error(err),
				zap.String("file", viper.ConfigFileUsed()),
			).Fatal("failed reading configuration")
		}
		logger.Infow("using config file", "file", viper.ConfigFileUsed())
	}

	setupAppConfig()
}

func initLogger() *zap.SugaredLogger {
	logCfg := zap.NewProductionConfig()

	if viper.GetViper().GetBool("logging.pretty") {
		logCfg.Development = true
	}

	logCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if viper.GetViper().GetBool("logging.debug") {
		logCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	log := zap.Must(logCfg.Build())
	_ = zap.ReplaceGlobals(log) // make the logger accessible globally by zap.L() (sugared with zap.S())
	return log.Sugar()
}

// setupAppConfig loads our config.AppConfig struct with the values bound by
// viper. Then, anywhere we need these values, we can just return to AppConfig
// instead of performing viper.GetString(...), viper.GetBool(...), etc.
func setupAppConfig() {
	err := viper.Unmarshal(&config.AppConfig)
	if err != nil {
		fmt.Printf("unable to decode app config: %s", err)
		os.Exit(1)
	}
}
