package tools

import "github.com/mark3labs/mcp-go/server"

func AddAllTools(m *server.MCPServer) {
	AddAccountTools(m)
	AddUserTools(m)
	AddDeviceTools(m)
	AddLabelTools(m)
	AddZoneTools(m)
	AddAutomationTools(m)
	AddMassActionTools(m)
	AddActionTools(m)
}
