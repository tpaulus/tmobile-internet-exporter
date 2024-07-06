package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	metricNamespace = "tmobile_gateway"
)

var (
	gatewayUp = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: metricNamespace,
		Name:      "gateway_up",
		Help:      "Gateway is able to be scraped",
	})

	// Device Metrics
	uptimeMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricNamespace,
		Name:      "uptime_seconds",
		Help:      "How many seconds a gateway has been continuously powered on.",
	}, []string{"serial"})
	softwareInfoMetric = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: metricNamespace,
		Name:      "software_info",
		Help:      "Gateway software information",
	}, []string{"serial", "version"})
	hardwareInfoMetric = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: metricNamespace,
		Name:      "hardware_info",
		Help:      "Gateway hardware information",
	}, []string{"serial", "name", "model", "version"})

	// Radio Metrics
	channelIdMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricNamespace,
		Name:      "channel_id",
		Help:      "Cellular channel number",
	}, []string{"type", "band"})
	rsrpMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricNamespace,
		Name:      "rsrp",
		Help:      "Cellular reference signal RX power",
	}, []string{"type", "band"})
	rsrqMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricNamespace,
		Name:      "rsrq",
		Help:      "Cellular reference signal RX quality",
	}, []string{"type", "band"})
	rssiMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricNamespace,
		Name:      "rssi",
		Help:      "Cellular received signal strength indicator",
	}, []string{"type", "band"})
	snrMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricNamespace,
		Name:      "snr",
		Help:      "Cellular signal to noise ratio",
	}, []string{"type", "band"})
	barsMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricNamespace,
		Name:      "bars",
		Help:      "Cellular signal strength in bars",
	}, []string{"type", "band"})
)

type Collector struct {
	target string
}

func (c *Collector) ScrapeEvery(ctx context.Context, clock clockwork.Clock, d time.Duration) error {
	t := clock.NewTicker(d)
	defer t.Stop()

	// Scrape immediatly
	err := c.scrape()
	if err != nil {
		gatewayUp.Set(0)
		fmt.Printf("%s\n", err.Error())
	} else {
		gatewayUp.Set(1)
	}

	// Scrape periodically
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.Chan():
			err := c.scrape()
			if err != nil {
				gatewayUp.Set(0)
				fmt.Println(err.Error())
			} else {
				gatewayUp.Set(1)
			}
		}
	}
}

func (c *Collector) scrape() error {
	resp, err := http.Get(fmt.Sprintf("http://%s/TMI/v1/gateway?get=all", c.target))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 status code received - %d", resp.StatusCode)
	}

	var status GatewayStatus
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		return err
	}

	// Device Metrics
	uptimeMetric.WithLabelValues(status.Device.Serial).Set(float64(status.Time.UpTime))
	softwareInfoMetric.WithLabelValues(status.Device.Serial, status.Device.SoftwareVersion).Observe(1)
	hardwareInfoMetric.WithLabelValues(status.Device.Serial, status.Device.Name, status.Device.Model, status.Device.HardwareVersion).Observe(1)

	// Radio Metrics
	// - 4G
	{
		d := status.Signal.FourG
		for _, v := range d.Bands {
			channelIdMetric.WithLabelValues("4g", v).Set(float64(d.Cid))
			rsrpMetric.WithLabelValues("4g", v).Set(float64(d.Rsrp))
			rsrqMetric.WithLabelValues("4g", v).Set(float64(d.Rsrq))
			rssiMetric.WithLabelValues("4g", v).Set(float64(d.Rssi))
			snrMetric.WithLabelValues("4g", v).Set(float64(d.Sinr))
			barsMetric.WithLabelValues("4g", v).Set(float64(d.Bars))
		}
	}

	// - 5G
	{
		d := status.Signal.FiveG
		for _, v := range d.Bands {
			channelIdMetric.WithLabelValues("5g", v).Set(float64(d.Cid))
			rsrpMetric.WithLabelValues("5g", v).Set(float64(d.Rsrp))
			rsrqMetric.WithLabelValues("5g", v).Set(float64(d.Rsrq))
			rssiMetric.WithLabelValues("5g", v).Set(float64(d.Rssi))
			snrMetric.WithLabelValues("5g", v).Set(float64(d.Sinr))
			barsMetric.WithLabelValues("5g", v).Set(float64(d.Bars))
		}
	}

	return nil
}
