package healthcheck

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/ankorstore/yokai/healthcheck"
)

// DefaultProbeName is the name of the Clickhouse probe.
const DefaultProbeName = "clickhouse"

// ClickhouseProbe is a probe compatible with the [healthcheck] module.
//
// [healthcheck]: https://github.com/ankorstore/yokai/tree/main/healthcheck
type ClickhouseProbe struct {
	name string
	conn driver.Conn
}

// NewClickhouseProbe returns a new [ClickhouseProbe].
func NewClickhouseProbe(conn driver.Conn) *ClickhouseProbe {
	return &ClickhouseProbe{
		name: DefaultProbeName,
		conn: conn,
	}
}

// NewClickhouseProbe returns the name of the [ClickhouseProbe].
func (p *ClickhouseProbe) Name() string {
	return p.name
}

// SetName sets the name of the [ClickhouseProbe].
func (p *ClickhouseProbe) SetName(name string) *ClickhouseProbe {
	p.name = name

	return p
}

// Check returns a successful [healthcheck.CheckerProbeResult] if the database connection can be pinged.
func (p *ClickhouseProbe) Check(ctx context.Context) *healthcheck.CheckerProbeResult {
	err := p.conn.Ping(ctx)
	if err != nil {
		return healthcheck.NewCheckerProbeResult(false, fmt.Sprintf("database ping error: %v", err))
	}

	return healthcheck.NewCheckerProbeResult(true, "database ping success")
}
