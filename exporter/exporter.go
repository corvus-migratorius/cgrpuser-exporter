package exporter

import (
	"fmt"
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

func (self *cgrpUserExporter) setSliceNames() {
	go func() {
		for {
			fmt.Println("Exporter is alive")
			time.Sleep(time.Duration(self.Timeout) * time.Second)
		}
	}()
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

func CgroupUserExporter(userSlicePath string, timeout int) *cgrpUserExporter {
	exporter := cgrpUserExporter{
		Hostname:  utils.GetHostname(),
		SlicePath: userSlicePath,
		Timeout:   timeout,
	}

	exporter.initGauges()

	return &exporter
}
