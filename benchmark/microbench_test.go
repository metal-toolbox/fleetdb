package benchmark_test

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/metal-toolbox/fleetdb/internal/config"
	"github.com/metal-toolbox/fleetdb/internal/dbtools"
	"github.com/metal-toolbox/fleetdb/internal/httpsrv"
	fleetdbapi "github.com/metal-toolbox/fleetdb/pkg/api/v1"
	"github.com/vattle/sqlboiler/boil"
	"go.infratographer.com/x/crdbx"
	"go.uber.org/zap"
)

// go test -benchmem -run=^$ -tags testtools -bench . github.com/metal-toolbox/fleetdb/benchmark -benchtime=30s
/*
BenchmarkGetComponentsJson
p50 latency: 3.597459ms
p90 latency: 6.332708ms
p99 latency: 20.509166ms

BenchmarkGetComponentsProto
p50 latency: 3.574375ms
p90 latency: 5.545833ms
p99 latency: 17.993708ms
*/

func init() {
	return

	go func() {
		sqldb, err := crdbx.NewDB(crdbx.Config{
			URI: "postgresql://root@fleetdb-crdb:26257/defaultdb?sslmode=disable",
		}, false)
		if err != nil {
			fmt.Errorf("failed to initialize database connection %v", err)
		}

		boil.SetDB(sqldb)

		db := sqlx.NewDb(sqldb, "postgres")

		logger, _ := zap.NewDevelopment()
		hs := &httpsrv.Server{
			Logger:        logger,
			Listen:        "localhost:12345",
			Debug:         config.AppConfig.Logging.Debug,
			DB:            db,
			OIDCEnabled:   false,
			SecretsKeeper: nil,
			AuthConfigs:   config.AppConfig.APIServerJWTAuth,
		}

		if err := hs.Run(); err != nil {
			fmt.Errorf("failed starting server %v", err)
		}
	}()

	client, err := fleetdbapi.NewClient("http://localhost:8000", nil)
	if err != nil {
		// b.Fatalf("failed to create client %v\n", err)
	}
	componentTypeSlice, _, err := client.ListServerComponentTypes(context.Background(), nil)
	if err != nil {
		// b.Fail()
		fmt.Errorf("%v", err)
	}

	csFixtureCreate := fleetdbapi.ServerComponentSlice{
		{
			ServerUUID:        uuid.MustParse("224c776c-239a-4853-b4b8-29fe7942f8cf"),
			Name:              "My Lucky Fin5",
			Vendor:            "barracuda5",
			Model:             "a lucky fin5",
			Serial:            "right5",
			ComponentTypeID:   componentTypeSlice.ByName("GPU").ID,
			ComponentTypeName: componentTypeSlice.ByName("GPU").Name,
			ComponentTypeSlug: componentTypeSlice.ByName("GPU").Slug,
			VersionedAttributes: []fleetdbapi.VersionedAttributes{
				{
					Namespace: dbtools.FixtureNamespaceVersioned,
					Data:      json.RawMessage(`{"version":"1.0"}`),
				},
				{
					Namespace: dbtools.FixtureNamespaceVersioned,
					Data:      json.RawMessage(`{"version":"2.0"}`),
				},
			},
		},
	}
	res, err := client.CreateComponents(context.TODO(), uuid.MustParse("224c776c-239a-4853-b4b8-29fe7942f8cf"), csFixtureCreate)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("============== res = %v\n", res)
}

func printHistogram(times []time.Duration) {
	// Set bin sizes (e.g., 10 ms intervals)
	n := len(times)
	sort.Slice(times, func(i, j int) bool {
		return times[i] < times[j]
	})
	fmt.Printf("\np50 latency: %v\np90 latency: %v\np99 latency: %v\n", times[n/2], times[n/100*90], times[n/100*99])
}

func BenchmarkGetComponentsJson(b *testing.B) {
	client, err := fleetdbapi.NewClient("http://localhost:8000", nil)
	if err != nil {
		b.Fatalf("failed to create client %v\n", err)
	}

	// resp, _, err := client.GetComponents(context.TODO(), uuid.MustParse("224c776c-239a-4853-b4b8-29fe7942f8cf"), nil)
	// fmt.Printf("resp= %v\b", resp)

	var times []time.Duration
	for i := 0; i < b.N; i++ {
		start := time.Now()
		resp, _, err := client.GetComponents(context.TODO(), uuid.MustParse("224c776c-239a-4853-b4b8-29fe7942f8cf"), nil)
		if err != nil {
			b.Fatal()
		}
		_ = resp
		duration := time.Since(start)
		times = append(times, duration)
	}
	printHistogram(times)
}

func BenchmarkGetComponentsProto(b *testing.B) {
	client, err := fleetdbapi.NewClient("http://localhost:8000", nil)
	if err != nil {
		b.Fatalf("failed to create client %v\n", err)
	}

	resp, _, err := client.GetComponentsProto(context.TODO(), uuid.MustParse("224c776c-239a-4853-b4b8-29fe7942f8cf"), nil)
	fmt.Printf("resp= %v, err = %v\n", resp, err)

	var times []time.Duration
	for i := 0; i < b.N; i++ {
		start := time.Now()
		resp, _, err := client.GetComponentsProto(context.TODO(), uuid.MustParse("224c776c-239a-4853-b4b8-29fe7942f8cf"), nil)
		if err != nil {
			b.Fatal()
		}
		_ = resp
		duration := time.Since(start)
		times = append(times, duration)
	}
	printHistogram(times)
}
