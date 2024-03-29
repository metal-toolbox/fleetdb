package cmd

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.infratographer.com/x/goosex"
	"go.infratographer.com/x/loggingx"
	"go.infratographer.com/x/versionx"
	"go.uber.org/zap"

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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hollow.yaml)")

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
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".hollow" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".hollow")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("fleetdb")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error: using config file, failed to read in!: %s\n", viper.ConfigFileUsed())
		os.Exit(1)
	} else {
		fmt.Printf("using config file: %s\n", viper.ConfigFileUsed())
	}

	setupAppConfig()

	// setupLogging()
	logger = loggingx.InitLogger(appName, config.AppConfig.Logging)

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
