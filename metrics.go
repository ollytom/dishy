package dishy

import (
	"fmt"
	"io"
	"text/template"

	"olowe.co/dishy/device"
)

const metricsPage string = `
# HELP dishy_uptime_seconds Seconds since last boot.
# TYPE dishy_uptime_seconds counter
dishy_uptime_seconds {{ .DeviceState.UptimeS }}
# HELP dishy_pop_ping_drop_rate
# TYPE dishy_pop_ping_drop_rate gauge
dishy_pop_ping_drop_rate {{ .PopPingDropRate }}
# HELP dishy_pop_ping_latency_milliseconds
# TYPE dishy_pop_ping_latency gauge
dishy_pop_ping_latency_milliseconds {{ .PopPingLatencyMs }}
# HELP dishy_downlink_throughput Received bytes per second.
# TYPE dishy_downlink_throughput gauge
dishy_downlink_throughput {{ .DownlinkThroughputBps }}
# HELP dishy_uplink_throughput Transmitted bytes per second.
# TYPE dishy_uplink_throughput gauge
dishy_uplink_throughput {{ .UplinkThroughputBps }}
# HELP dishy_obstruction_percentage
# TYPE dishy_obstruction_percentage gauge
dishy_obstruction_percentage {{ .ObstructionStats.FractionObstructed }}
`

var metricsTmpl = template.Must(template.New("metrics").Parse(metricsPage))

// WriteOpenMetrics writes any metrics found in status in [OpenMetrics]
// format to w for use in systems such as Prometheus and VictoriaMetrics.
//
// [OpenMetrics]: https://openmetrics.io/
func WriteOpenMetrics(w io.Writer, status *device.DishGetStatusResponse) error {
	err := metricsTmpl.Execute(w, status)
	if err != nil {
		return fmt.Errorf("execute template: %w", err)
	}
	return nil
}
