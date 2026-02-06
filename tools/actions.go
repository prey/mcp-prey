package tools

import (
	"context"
	"net/http"
	"net/url"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	mcprey "mcp-prey"
	"mcp-prey/internal"
	"mcp-prey/prey"
)

type DeviceActionTriggerParams struct {
	DeviceID   string         `json:"deviceId" jsonschema:"description=ID of the device"`
	Command    string         `json:"command" jsonschema:"description=Command to execute (start)"`
	ActionName string         `json:"action_name" jsonschema:"description=Action name (alarm|alert|lock)"`
	Options    map[string]any `json:"options,omitempty" jsonschema:"description=Action options"`
}

type DeviceStatusSetParams struct {
	DeviceID string `json:"deviceId" jsonschema:"description=ID of the device"`
	Missing  bool   `json:"missing" jsonschema:"description=true to mark missing, false to recover"`
}

func deviceActionTrigger(ctx context.Context, args DeviceActionTriggerParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.devices.action.trigger", true); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.DeviceID, "deviceId"); err != nil {
		return nil, err
	}
	if err := internal.RequireOneOf(args.Command, "command", "start"); err != nil {
		return nil, err
	}
	if err := internal.RequireOneOf(args.ActionName, "action_name", "alarm", "alert", "lock"); err != nil {
		return nil, err
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	body := map[string]any{
		"command":     args.Command,
		"action_name": args.ActionName,
	}
	if args.Options != nil {
		body["options"] = args.Options
	}
	var payload any
	req, err := client.NewRequest(http.MethodPut, "/devices/"+args.DeviceID+"/action", url.Values{}, body)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), nil), nil
}

func deviceStatusSet(ctx context.Context, args DeviceStatusSetParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.devices.status.set", true); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.DeviceID, "deviceId"); err != nil {
		return nil, err
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	body := map[string]any{"missing": args.Missing}
	var payload any
	req, err := client.NewRequest(http.MethodPut, "/devices/"+args.DeviceID+"/missing", url.Values{}, body)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), nil), nil
}

var DeviceActionTrigger = mcprey.MustTool(
	"prey.devices.action.trigger",
	"Trigger a device action (write).",
	deviceActionTrigger,
	mcp.WithTitleAnnotation("Trigger device action"),
)

var DeviceStatusSet = mcprey.MustTool(
	"prey.devices.status.set",
	"Set device status (write).",
	deviceStatusSet,
	mcp.WithTitleAnnotation("Set device status"),
)

func AddActionTools(m *server.MCPServer) {
	DeviceActionTrigger.Register(m)
	DeviceStatusSet.Register(m)
}
