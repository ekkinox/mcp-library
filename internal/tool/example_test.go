package tool_test

import (
	"context"
	"strings"
	"testing"

	"github.com/ankorstore/yokai-mcp-template/internal"
	"github.com/ankorstore/yokai/fxmcpserver/fxmcpservertest"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/ankorstore/yokai/trace/tracetest"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
)

func TestExampleTool(t *testing.T) {
	var testServer *fxmcpservertest.MCPSSETestServer
	var logBuffer logtest.TestLogBuffer
	var traceExporter tracetest.TestTraceExporter
	var metricsRegistry *prometheus.Registry

	internal.RunTest(t, fx.Populate(&testServer, &logBuffer, &traceExporter, &metricsRegistry))

	defer testServer.Close()

	testClient, err := testServer.StartClient(context.Background())
	assert.NoError(t, err)

	defer testClient.Close()

	req := mcp.CallToolRequest{}
	req.Params.Name = "example-tool"

	res, err := testClient.CallTool(context.Background(), req)
	assert.NoError(t, err)
	assert.False(t, res.IsError)

	resContents, ok := res.Content[0].(mcp.TextContent)
	assert.True(t, ok)
	assert.Equal(t, "the application name is mcp-app", resContents.Text)

	logtest.AssertHasLogRecord(t, logBuffer, map[string]interface{}{
		"level":        "info",
		"mcpMethod":    "tools/call",
		"mcpTool":      "example-tool",
		"mcpTransport": "sse",
		"message":      "MCP request success",
	})

	tracetest.AssertHasTraceSpan(
		t,
		traceExporter,
		"MCP tools/call example-tool",
		attribute.String("mcp.method", "tools/call"),
		attribute.String("mcp.tool", "example-tool"),
		attribute.String("mcp.transport", "sse"),
	)

	expectedMetric := `
        # HELP mcp_server_requests_total Number of processed MCP requests
        # TYPE mcp_server_requests_total counter
		mcp_server_requests_total{method="initialize",status="success",target=""} 1
        mcp_server_requests_total{method="tools/call",status="success",target="example-tool"} 1
    `

	err = testutil.GatherAndCompare(
		metricsRegistry,
		strings.NewReader(expectedMetric),
		"mcp_server_requests_total",
	)
	assert.NoError(t, err)
}
