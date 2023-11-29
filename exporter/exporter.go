/* Package exporter exports data about Systemd user slices as Prometheus gauges.

Data extraction is handled with the `cgrpUserExporter` type.
Use the `CgroupUserExporter` factory function to get a properly initialized
instance of the exporter.

A web server exporting the `/metrics` endpoint is not included in the module,
but is trivial to implement: e.g. see the accompanying `main.go` file. */
package exporter

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"cgrpuser-exporter/utils"
)

type cgrpUserExporter struct {
	Timeout       int
	Hostname      string
	SlicePath     string
	MemoryCurrent *prometheus.GaugeVec
	SwapCurrent   *prometheus.GaugeVec
}

func (self *cgrpUserExporter) initGauges() {
	self.MemoryCurrent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cgrpuser_exporter_memory_current",
		Help: "Memory consumed by user slices on this node",
	},
		[]string{"nodename", "uid", "username"},
	)

	self.SwapCurrent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cgrpuser_exporter_swap_current",
		Help: "Swap space consumed by user slices on this node",
	},
		[]string{"nodename", "uid", "username"},
	)
}

func (self *cgrpUserExporter) RecordMetrics() {
	go func() {
		for {
			self.MemoryCurrent.Reset()
			self.SwapCurrent.Reset()

			slices := utils.GetUserSlices(self.SlicePath)
			for _, slice := range slices {
				username := slice.Username
				uid := slice.UID
				self.MemoryCurrent.With(prometheus.Labels{
					"username": username,
					"uid":      uid,
					"nodename": self.Hostname,
				}).Set(float64(slice.MemoryCurrent))

				self.SwapCurrent.With(prometheus.Labels{
					"username": slice.Username,
					"uid":      slice.UID,
					"nodename": self.Hostname,
				}).Set(float64(slice.SwapCurrent))
			}

			time.Sleep(time.Duration(self.Timeout) * time.Second)
		}
	}()
}

// CgroupUserExporter is factory function serving as constructor for the cgrpUserExporter type
func CgroupUserExporter(userSlicePath string, timeout int) *cgrpUserExporter {
	exporter := cgrpUserExporter{
		Hostname:  utils.GetHostname(),
		SlicePath: userSlicePath,
		Timeout:   timeout,
	}

	exporter.initGauges()

	return &exporter
}
