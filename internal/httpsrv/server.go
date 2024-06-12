package httpsrv

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.hollow.sh/toolbox/events"
	"go.hollow.sh/toolbox/ginauth"
	"go.hollow.sh/toolbox/ginjwt"
	"go.infratographer.com/x/versionx"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gocloud.dev/secrets"

	"github.com/metal-toolbox/fleetdb/internal/metrics"
	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
)

// Server implements the HTTP Server
type Server struct {
	Logger        *zap.Logger
	Listen        string
	Debug         bool
	DB            *sqlx.DB
	OIDCEnabled   bool
	AuthConfigs   []ginjwt.AuthConfig
	SecretsKeeper *secrets.Keeper
	EventStream   events.Stream
}

var (
	readTimeout        = 10 * time.Second
	writeTimeout       = 20 * time.Second
	corsMaxAge         = 12 * time.Hour
	noOpAuthMiddleware = &ginauth.MultiTokenMiddleware{}
)

func getGinPrometheusAdapter(h http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *Server) setup() *gin.Engine {
	var err error
	authMW := noOpAuthMiddleware

	if s.OIDCEnabled {
		authMW, err = ginjwt.NewMultiTokenMiddlewareFromConfigs(s.AuthConfigs...)
		if err != nil {
			s.Logger.With(
				zap.Error(err),
			).Fatal("failed to initialize auth middleware")
		}
	}

	// Setup default gin router
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowAllOrigins:  true,
		AllowCredentials: true,
		MaxAge:           corsMaxAge,
	}))

	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next() // call the rest of the handler chain
		metrics.APICallEpilog(start, c.FullPath(), c.Writer.Status())
	})

	v1Rtr := fleetdbapi.Router{
		DB:            s.DB,
		AuthMW:        authMW,
		SecretsKeeper: s.SecretsKeeper,
		Logger:        s.Logger,
		EventStream:   s.EventStream,
	}

	r.Use(ginzap.GinzapWithConfig(s.Logger.With(zap.String("component", "httpsrv")), &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			return []zapcore.Field{
				zap.String("path", c.Request.URL.Path),
				zap.String("query", c.Request.URL.RawQuery),
				zap.String("jwt_subject", ginjwt.GetSubject(c)),
				zap.String("jwt_user", ginjwt.GetUser(c)),
			}
		}),
		SkipPaths: []string{
			"/healthz",
			"/healthz/liveness",
			"/healthz/readiness",
			"/metrics",
		},
	}))

	r.Use(ginzap.RecoveryWithZap(s.Logger.With(zap.String("component", "httpsrv")), true))

	tp := otel.GetTracerProvider()
	if tp != nil {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}

		r.Use(otelgin.Middleware(hostname, otelgin.WithTracerProvider(tp)))
	}

	// Version endpoint returns build information
	r.GET("/version", s.version)

	// Health endpoints
	r.GET("/healthz", s.livenessCheck)
	r.GET("/healthz/liveness", s.livenessCheck)
	r.GET("/healthz/readiness", s.readinessCheck)

	r.GET("/metrics", getGinPrometheusAdapter(promhttp.Handler()))

	v1 := r.Group("/api/v1")
	{
		v1Rtr.Routes(v1)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "invalid request - route not found"})
	})

	return r
}

// NewServer returns a configured server
func (s *Server) NewServer() *http.Server {
	if !s.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	return &http.Server{
		Handler:      s.setup(),
		Addr:         s.Listen,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

// Run will start the server listening on the specified address
func (s *Server) Run() error {
	if !s.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	return s.setup().Run(s.Listen)
}

// livenessCheck ensures that the server is up and responding
func (s *Server) livenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}

// readinessCheck ensures that the server is up and that we are able to process
// requests. Currently our only dependency is the DB so we just ensure that is
// responding.
func (s *Server) readinessCheck(c *gin.Context) {
	if err := s.DB.PingContext(c.Request.Context()); err != nil {
		s.Logger.Sugar().Errorf("readiness check db ping failed", "err", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "DOWN",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}

// version returns the fleetdb build information.
func (s *Server) version(c *gin.Context) {
	c.JSON(http.StatusOK, versionx.BuildDetails().String())
}
