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
	Slices        []utils.UserSlice
	MemoryCurrent *prometheus.GaugeVec
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
		Help: "Memory consumed by unit slices on this node",
	},
		[]string{"nodename", "uid", "username"},
	)
}

func (self *cgrpUserExporter) RecordMemoryCurrent() {
	go func() {
		for {
			self.MemoryCurrent.Reset()
			for _, slice := range self.Slices {
				self.MemoryCurrent.With(prometheus.Labels{
					"username": slice.Username,
					"uid":      slice.UID,
					"nodename": self.Hostname,
				}).Set(float64(slice.MemoryCurrent))
			}

			time.Sleep(time.Duration(self.Timeout) * time.Second)
		}
	}()
}

func CgroupUserExporter(userSlicePath string, timeout int) *cgrpUserExporter {
	exporter := cgrpUserExporter{
		Hostname: utils.GetHostname(),
		Slices:   utils.GetUserSlices(userSlicePath),
		Timeout:  timeout,
	}

	exporter.initGauges()

	return &exporter
}
