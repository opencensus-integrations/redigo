package observability

import (
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// Pool metrics:
// 1. Connections taken
// 2. Connections closed
// 3. Connections usetime -- how long is a connection used until it is closed, discarded or returned
// 4. Connections reused
// 4. Connections stale
// 5. Dial errors

const dimensionless = "1"
const milliseconds = "ms"

var (
	// The measures to record metrics
	MBytesRead    = stats.Int64("redis/bytes_read", "The number of bytes read from the server", stats.UnitBytes)
	MBytesWritten = stats.Int64("redis/bytes_written", "The number of bytes written out to the server", stats.UnitBytes)
	MErrors       = stats.Int64("redis/errors", "The number of errors encountered", dimensionless)
	MWrites       = stats.Int64("redis/writes", "The number of write invocations", dimensionless)
	MReads        = stats.Int64("redis/reads", "The number of read invocations", dimensionless)

	MRoundtripLatencyMs = stats.Float64("redis/roundtrip_latency", "The latency in milliseconds of a method/operation", milliseconds)

	MConnectionsClosed = stats.Int64("redis/connections_closed", "The number of connections that have been closed", dimensionless)
	MConnectionsOpen   = stats.Int64("redis/connections_new", "The number of open connections", dimensionless)
)

var (
	// The tag keys to record alongside measures
	KeyCommandName, _ = tag.NewKey("cmd")
	KeyKind, _        = tag.NewKey("kind")
	KeyDetail, _      = tag.NewKey("detail")
	KeyState, _       = tag.NewKey("state")
)

var defaultMillisecondsDistribution = view.Distribution(
	// [0ms, 0.001ms, 0.005ms, 0.01ms, 0.05ms, 0.1ms, 0.5ms, 1ms, 1.5ms, 2ms, 2.5ms, 5ms, 10ms, 25ms, 50ms, 100ms, 200ms, 400ms, 600ms, 800ms, 1s, 1.5s, 2.5s, 5s, 10s, 20s, 40s, 100s, 200s, 500s]
	0, 0.000001, 0.000005, 0.00001, 0.00005, 0.0001, 0.0005, 0.001, 0.0015, 0.002, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.2, 0.4, 0.6, 0.8, 1, 1.5, 2.5, 5, 10, 20, 40, 100, 200, 500,
)

var defaultBytesDistribution = view.Distribution(
	// [0, 1KB, 2KB, 4KB, 16KB, 64KB, 256KB,   1MB,     4MB,     16MB,     64MB,     256MB,     1GB,        4GB]
	0, 1024, 2048, 4096, 16384, 65536, 262144, 1048576, 4194304, 16777216, 67108864, 268435456, 1073741824, 4294967296,
)

var Views = []*view.View{
	{
		Name:        "redis/client/bytes_written",
		Description: "The distribution of bytes written out to the server",
		Aggregation: defaultBytesDistribution,
		Measure:     MBytesWritten,
	},
	{
		Name:        "redis/client/bytes_read",
		Description: "The distribution of bytes read from the server",
		Aggregation: defaultBytesDistribution,
		Measure:     MBytesRead,
	},
	{
		Name:        "redis/client/roundtrip_latency",
		Description: "The distribution of milliseconds of the roundtrip latencies for a Redis method invocation",
		Aggregation: defaultMillisecondsDistribution,
		Measure:     MRoundtripLatencyMs,
		TagKeys:     []tag.Key{KeyCommandName},
	},
	{
		Name:        "redis/client/writes",
		Description: "The number of write operations",
		Aggregation: view.Count(),
		Measure:     MWrites,
		TagKeys:     []tag.Key{KeyCommandName},
	},
	{
		Name:        "redis/client/reads",
		Description: "The number of read operations",
		Aggregation: view.Count(),
		Measure:     MReads,
		TagKeys:     []tag.Key{KeyCommandName},
	},
	{
		Name:        "redis/client/errors",
		Description: "The number of errors encountered",
		Aggregation: view.Count(),
		Measure:     MErrors,
		TagKeys:     []tag.Key{KeyCommandName, KeyDetail, KeyKind},
	},
	{
		Name:        "redis/client/connections_closed",
		Description: "The number of connections that have been closed, disambiguated by keys such as stale, idle, complete",
		Aggregation: view.Count(),
		Measure:     MConnectionsClosed,
		TagKeys:     []tag.Key{KeyState},
	},
	{
		Name:        "redis/client/connections_open",
		Description: "The number of open connections, but disambiguated by different states e.g. new, reused",
		Aggregation: view.Count(),
		Measure:     MConnectionsOpen,
		TagKeys:     []tag.Key{KeyState},
	},
}

func SinceInMilliseconds(startTime time.Time) float64 {
	return float64(time.Since(startTime).Nanoseconds()) * 1e6
}
