package ocppj

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const (
	requestsInboundMetric  = "ocpp_requests_inbound"
	requestsOutboundMetric = "ocpp_requests_outbound"
	requestQueueMetric     = "ocpp_request_queue_size"
)

const (
	attributeChargePointId = "charge_point_id"
	attributeOcppVersion   = "ocpp_version"
	attributeFeature       = "feature"
	attributeError         = "error"
)

type ocppMetricsError string

var (
	chargePointError     = ocppMetricsError("charge_point_error")
	metricsInternalError = ocppMetricsError("internal_error")
	metricsNetworkError  = ocppMetricsError("network_error")
	payloadError         = ocppMetricsError("payload_error")
	validationError      = ocppMetricsError("validation_error")
)

type ocppMetrics struct {
	requestsIn         metric.Int64Histogram
	requestsOut        metric.Int64Histogram
	requestQueueMetric metric.Int64ObservableGauge
	meter              metric.Meter
}

// newOcppMetrics Creates a new metrics instance
func newOcppMetrics(meterProvider metric.MeterProvider, ocppVersion string) (*ocppMetrics, error) {
	if meterProvider == nil {
		return nil, errors.New("meterProvider is required")
	}

	meter := meterProvider.Meter(
		"ocpp",
		metric.WithInstrumentationAttributes(attribute.String(attributeOcppVersion, ocppVersion)),
	)

	requestsIn, err := meter.Int64Histogram(
		requestsInboundMetric,
		metric.WithDescription("Number of inbound requests"),
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create %s metric", requestsInboundMetric))
	}

	requestsOut, err := meter.Int64Histogram(
		requestsOutboundMetric,
		metric.WithDescription("Number of outbound requests"),
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create %s metric", requestsOutboundMetric))
	}

	queueSize, err := meter.Int64ObservableGauge(
		requestQueueMetric,
		metric.WithDescription("Current queue length"),
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to create %s metric", requestQueueMetric))
	}

	metrics := &ocppMetrics{
		requestsIn:         requestsIn,
		requestsOut:        requestsOut,
		requestQueueMetric: queueSize,
		meter:              meter,
	}
	return metrics, nil
}

func (m *ocppMetrics) IncrementInboundRequests(ctx context.Context, chargePointId, requestName string, error *ocppMetricsError) {
	attrs := []attribute.KeyValue{
		attribute.String(attributeChargePointId, chargePointId),
	}

	// Optionally add a request name. Should be present most of the time, except when we cannot unmarshal the request.
	if requestName != "" {
		attribute.String(attributeFeature, requestName)
	}

	if error != nil {
		attrs = append(attrs, attribute.String(attributeError, string(*error)))
	}

	metricAttrs := metric.WithAttributes(attrs...)
	m.requestsIn.Record(ctx, 1, metricAttrs)
}

func (m *ocppMetrics) IncrementOutboundRequests(ctx context.Context, chargePointId, requestName string, error *ocppMetricsError) {
	attrs := []attribute.KeyValue{
		attribute.String(attributeChargePointId, chargePointId),
	}

	// Optionally add a request name. Should be present most of the time, except when we cannot unmarshal the request.
	if requestName != "" {
		attribute.String(attributeFeature, requestName)
	}

	if error != nil {
		attrs = append(attrs, attribute.String(attributeError, string(*error)))
	}

	metricAttrs := metric.WithAttributes(attrs...)
	m.requestsOut.Record(ctx, 1, metricAttrs)
}

func (m *ocppMetrics) ObserveRequestQueue(ctx context.Context, queueMap ServerQueueMap) {
	_, err := m.meter.RegisterCallback(
		func(ctx context.Context, observer metric.Observer) error {
			observer.ObserveInt64(m.requestQueueMetric, int64(queueMap.Size()))
			return nil
		},
		m.requestQueueMetric,
	)
	if err != nil {
		log.Error(err)
	}
}
