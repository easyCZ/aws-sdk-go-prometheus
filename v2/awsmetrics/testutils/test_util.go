package testutils

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/smithy-go/middleware"
	"github.com/jonathan-innis/aws-sdk-go-prometheus/v2/awsmetrics"
)

type MetricDataRecorderPublisher struct {
	Data *awsmetrics.MetricData
}

func (mdrp *MetricDataRecorderPublisher) PostRequestMetrics(data *awsmetrics.MetricData) error {
	mdrp.Data = data
	return nil
}

func (mdrp *MetricDataRecorderPublisher) PostStreamMetrics(data *awsmetrics.MetricData) error {
	mdrp.Data = data
	return nil
}

type NoopPublisher struct{}

func (np *NoopPublisher) PostRequestMetrics(data *awsmetrics.MetricData) error {
	return nil
}

func (np *NoopPublisher) PostStreamMetrics(data *awsmetrics.MetricData) error {
	return nil
}

type ErrorPublisher struct{}

func (tp *ErrorPublisher) PostRequestMetrics(data *awsmetrics.MetricData) error {
	return fmt.Errorf("publisher error")
}

func (tp *ErrorPublisher) PostStreamMetrics(data *awsmetrics.MetricData) error {
	return fmt.Errorf("publisher error")
}

type NoopInitializeHandler struct{}
type ErrorInitializeHandler struct{}
type NoopSerializeHandler struct{}
type NoopFinalizeHandler struct{}
type NoopDeserializeHandler struct{}
type StreamingBodyBuildHandler struct {
	Result interface{}
}

func NoopRequestCloner(i interface{}) interface{} {
	return i
}

func (NoopInitializeHandler) HandleInitialize(ctx context.Context, in middleware.InitializeInput) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	return middleware.InitializeOutput{}, middleware.Metadata{}, nil
}

func (ErrorInitializeHandler) HandleInitialize(ctx context.Context, in middleware.InitializeInput) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	return middleware.InitializeOutput{}, middleware.Metadata{}, fmt.Errorf("init error")
}

func (NoopFinalizeHandler) HandleFinalize(ctx context.Context, in middleware.FinalizeInput) (
	out middleware.FinalizeOutput, metadata middleware.Metadata, err error,
) {
	return middleware.FinalizeOutput{}, middleware.Metadata{}, nil
}

func (NoopDeserializeHandler) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	return middleware.DeserializeOutput{}, middleware.Metadata{}, nil
}

func (NoopSerializeHandler) HandleSerialize(ctx context.Context, in middleware.SerializeInput) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	return middleware.SerializeOutput{}, middleware.Metadata{}, nil
}

func (s *StreamingBodyBuildHandler) HandleBuild(ctx context.Context, in middleware.BuildInput) (
	out middleware.BuildOutput, metadata middleware.Metadata, err error,
) {
	return middleware.BuildOutput{Result: s.Result}, middleware.Metadata{}, nil
}

type TestReadCloser struct {
	Data   []byte
	offset int
}

func (m *TestReadCloser) Read(p []byte) (int, error) {
	if m.offset >= len(m.Data) {
		return 0, io.EOF
	}
	n := copy(p, m.Data[m.offset:])
	m.offset += n
	return n, nil
}

func (m *TestReadCloser) Close() error {
	return nil
}
