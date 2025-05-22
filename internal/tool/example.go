package tool

import (
	"context"
	"fmt"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/trace"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ExampleTool is an example MCP tool.
type ExampleTool struct {
	config *config.Config
}

// NewExampleTool returns a new ExampleTool instance.
func NewExampleTool(cfg *config.Config) *ExampleTool {
	return &ExampleTool{
		config: cfg,
	}
}

// Name returns the ExampleTool name.
func (t *ExampleTool) Name() string {
	return "example-tool"
}

// Options returns the ExampleTool options.
func (t *ExampleTool) Options() []mcp.ToolOption {
	return []mcp.ToolOption{
		mcp.WithDescription("returns the application name"),
	}
}

// Handle returns the ExampleTool request handler.
func (t *ExampleTool) Handle() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx, span := trace.CtxTracer(ctx).Start(ctx, "ExampleTool.Handle")
		defer span.End()

		log.CtxLogger(ctx).Info().Msg("ExampleTool.Handle")

		return mcp.NewToolResultText(fmt.Sprintf("the application name is %s", t.config.AppName())), nil
	}
}
