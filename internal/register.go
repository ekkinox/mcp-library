package internal

import (
	"github.com/ankorstore/yokai-mcp-template/internal/tool"
	"github.com/ankorstore/yokai/fxhealthcheck"
	"github.com/ankorstore/yokai/fxmcpserver"
	"github.com/ankorstore/yokai/fxmcpserver/server"
	"go.uber.org/fx"
)

// Register is used to register the application dependencies.
func Register() fx.Option {
	return fx.Options(
		// MCP registrations
		fxmcpserver.AsMCPServerTools(tool.NewExampleTool),
		// MCP probe
		fxhealthcheck.AsCheckerProbe(server.NewMCPServerProbe),
	)
}
