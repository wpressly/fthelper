package metrics

import (
	"github.com/kamontat/fthelper/metric/src/collectors"
	"github.com/kamontat/fthelper/metric/src/connection"
	"github.com/kamontat/fthelper/metric/src/freqtrade"
	"github.com/kamontat/fthelper/shared/commandline/commands"
	"github.com/prometheus/client_golang/prometheus"
)

var Info = collectors.NewMetrics(
	collectors.NewMetric(prometheus.NewDesc(
		prometheus.BuildFQName("fthelper", "build", "info"),
		"Information relate to fthelper, get information from commandline metadata",
		[]string{"version", "commit", "date"},
		nil,
	), func(desc *prometheus.Desc, conn connection.Http, param *commands.ExecutorParameter) []prometheus.Metric {
		return []prometheus.Metric{prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			float64(1),
			param.Meta.Version,
			param.Meta.Commit,
			param.Meta.Date,
		)}
	}),
	collectors.NewMetric(prometheus.NewDesc(
		prometheus.BuildFQName("freqtrade", "build", "info"),
		"Information relate to freqtrade, 0 meaning server is down",
		freqtrade.SummaryLabel(),
		nil,
	), func(desc *prometheus.Desc, conn connection.Http, param *commands.ExecutorParameter) []prometheus.Metric {
		var connection = freqtrade.ToConnection(conn)
		var up = freqtrade.NewPingI(connection)

		return []prometheus.Metric{prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			float64(up),
			freqtrade.NewSummary(connection, param.Cache)...,
		)}
	}),
	collectors.NewMetric(prometheus.NewDesc(
		prometheus.BuildFQName("freqtrade", "", "state"),
		"Freqtrade run state",
		append(freqtrade.SummaryLabel(), "text"),
		nil,
	), func(desc *prometheus.Desc, conn connection.Http, param *commands.ExecutorParameter) []prometheus.Metric {
		var connection = freqtrade.ToConnection(conn)
		var status = freqtrade.NewStatus(connection)

		var labels = append(freqtrade.NewSummary(connection, param.Cache), status.StateStr())
		return []prometheus.Metric{prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			float64(status.StateInt()),
			labels...,
		)}
	}),
	collectors.NewMetric(prometheus.NewDesc(
		prometheus.BuildFQName("freqtrade", "", "mode"),
		"Freqtrade mode",
		append(freqtrade.SummaryLabel(), "text"),
		nil,
	), func(desc *prometheus.Desc, conn connection.Http, param *commands.ExecutorParameter) []prometheus.Metric {
		var connection = freqtrade.ToConnection(conn)
		var status = freqtrade.NewStatus(connection)

		var labels = append(freqtrade.NewSummary(connection, param.Cache), status.RunMode)
		return []prometheus.Metric{prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			float64(status.ModeInt()),
			labels...,
		)}
	}),
)
