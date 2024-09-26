package cmd

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metal-toolbox/rivets/ginjwt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/volatiletech/sqlboiler/boil"
	"go.infratographer.com/x/crdbx"
	"go.infratographer.com/x/otelx"
	"go.infratographer.com/x/viperx"
	"go.uber.org/zap"
	"gocloud.dev/secrets"

	// import gocdk secret drivers
	_ "gocloud.dev/secrets/localsecrets"

	"github.com/metal-toolbox/fleetdb/internal/config"
	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	"github.com/metal-toolbox/fleetdb/internal/httpsrv"
)

var (
	apiDefaultListen = "0.0.0.0:8000"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "starts the hollow server",
	Run: func(cmd *cobra.Command, _ []string) {
		serve(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("listen", apiDefaultListen, "address to listen on")
	viperx.MustBindFlag(viper.GetViper(), "listen", serveCmd.Flags().Lookup("listen"))

	otelx.MustViperFlags(viper.GetViper(), serveCmd.Flags())
	crdbx.MustViperFlags(viper.GetViper(), serveCmd.Flags())

	// OIDC Flags
	serveCmd.Flags().Bool("oidc", true, "use oidc auth")
	ginjwt.BindFlagFromViperInst(viper.GetViper(), "oidc.enabled", serveCmd.Flags().Lookup("oidc"))

	// DB Flags
	serveCmd.Flags().String("db-encryption-driver", "", "encryption driver uri; 32 byte base64 encoded string, (example: base64key://your-encoded-secret-key)")
	viperx.MustBindFlag(viper.GetViper(), "db.encryption_driver", serveCmd.Flags().Lookup("db-encryption-driver"))
}

func serve(ctx context.Context) {
	err := otelx.InitTracer(config.AppConfig.Tracing, appName, logger)
	if err != nil {
		logger.Fatalw("unable to initialize tracing system", "error", err)
	}

	db := initDB()

	dbtools.RegisterHooks()

	if err := dbtools.SetupComponentTypes(ctx, db); err != nil {
		logger.With(
			zap.Error(err),
		).Fatal("set up component types")
	}

	keeper, err := secrets.OpenKeeper(ctx, viper.GetString("db.encryption_driver"))
	if err != nil {
		logger.Fatalw("failed to open secrets keeper", "error", err)
	}
	defer keeper.Close()

	logger.Infow("starting server",
		"address", viper.GetString("listen"),
	)

	var oidcEnabled bool
	if viper.GetViper().GetBool("oidc.enabled") {
		logger.Infow("OIDC enabled")

		if len(config.AppConfig.APIServerJWTAuth) == 0 {
			logger.Fatal("OIDC enabled without configuration")
		}
		oidcEnabled = true
	} else {
		logger.Infow("OIDC disabled")
	}

	hs := &httpsrv.Server{
		Logger:        logger.Desugar(),
		Listen:        viper.GetString("listen"),
		Debug:         config.AppConfig.Logging.Debug,
		DB:            db,
		OIDCEnabled:   oidcEnabled,
		SecretsKeeper: keeper,
		AuthConfigs:   config.AppConfig.APIServerJWTAuth,
	}

	if err := hs.Run(); err != nil {
		logger.Fatalw("failed starting server", "error", err)
	}
}

func initDB() *sqlx.DB {
	dbDriverName := "postgres"

	sqldb, err := crdbx.NewDB(config.AppConfig.CRDB, config.AppConfig.Tracing.Enabled)
	if err != nil {
		logger.Fatalw("failed to initialize database connection", "error", err)
	}

	boil.SetDB(sqldb)

	db := sqlx.NewDb(sqldb, dbDriverName)

	return db
}
