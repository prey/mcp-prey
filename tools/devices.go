package tools

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	mcprey "mcp-prey"
	"mcp-prey/internal"
	"mcp-prey/prey"
)

type DevicesListParams struct {
	Page     int `json:"page,omitempty" jsonschema:"default=1,description=Page number"`
	PageSize int `json:"page_size,omitempty" jsonschema:"default=20,minimum=1,maximum=100,description=Number of records per page"`
}

type DevicesGetParams struct {
	DeviceID string `json:"deviceId" jsonschema:"description=ID of the device"`
}

type DevicesDeleteParams struct {
	DeviceID string `json:"deviceId" jsonschema:"description=ID of the device"`
}

type DevicesReportsListParams struct {
	DeviceID string `json:"deviceId" jsonschema:"description=ID of the device"`
	Page     int    `json:"page,omitempty" jsonschema:"default=1,description=Page number"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"default=20,minimum=1,maximum=100,description=Number of records per page"`
}

type DevicesReportsGetParams struct {
	DeviceID string `json:"deviceId" jsonschema:"description=ID of the device"`
	ReportID string `json:"reportId" jsonschema:"description=ID of the report"`
}

type DevicesLocationHistoryParams struct {
	DeviceID string `json:"deviceId" jsonschema:"description=ID of the device"`
	Format   string `json:"format,omitempty" jsonschema:"default=json,description=Response format: json or csv"`
}

func devicesList(ctx context.Context, args DevicesListParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.devices.list", false); err != nil {
		return nil, err
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	q, err := internal.AddPagination(url.Values{}, args.Page, args.PageSize)
	if err != nil {
		return nil, err
	}
	var payload any
	req, err := client.NewRequest(http.MethodGet, "/devices", q, nil)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	meta, err := internal.Meta(args.Page, args.PageSize)
	if err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), meta), nil
}

func devicesGet(ctx context.Context, args DevicesGetParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.devices.get", false); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.DeviceID, "deviceId"); err != nil {
		return nil, err
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	var payload any
	req, err := client.NewRequest(http.MethodGet, "/devices/"+args.DeviceID, url.Values{}, nil)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), nil), nil
}

func devicesDelete(ctx context.Context, args DevicesDeleteParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.devices.delete", true); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.DeviceID, "deviceId"); err != nil {
		return nil, err
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	var payload any
	req, err := client.NewRequest(http.MethodDelete, "/devices/"+args.DeviceID, url.Values{}, nil)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), nil), nil
}

func devicesReportsList(ctx context.Context, args DevicesReportsListParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.devices.reports.list", false); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.DeviceID, "deviceId"); err != nil {
		return nil, err
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	q, err := internal.AddPagination(url.Values{}, args.Page, args.PageSize)
	if err != nil {
		return nil, err
	}
	var payload any
	req, err := client.NewRequest(http.MethodGet, "/devices/"+args.DeviceID+"/reports", q, nil)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	meta, err := internal.Meta(args.Page, args.PageSize)
	if err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), meta), nil
}

func devicesReportsGet(ctx context.Context, args DevicesReportsGetParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.devices.reports.get", false); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.DeviceID, "deviceId"); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.ReportID, "reportId"); err != nil {
		return nil, err
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	var payload any
	req, err := client.NewRequest(http.MethodGet, "/devices/"+args.DeviceID+"/reports/"+args.ReportID, url.Values{}, nil)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), nil), nil
}

func devicesLocationHistory(ctx context.Context, args DevicesLocationHistoryParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.devices.location_history.get", false); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.DeviceID, "deviceId"); err != nil {
		return nil, err
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	format := strings.ToLower(strings.TrimSpace(args.Format))
	if format == "csv" {
		req, err := client.NewRequest(http.MethodGet, "/devices/"+args.DeviceID+"/location_activity.csv", url.Values{}, nil)
		if err != nil {
			return nil, err
		}
		b, contentType, err := client.DoRaw(req.WithContext(ctx))
		if err != nil {
			return nil, err
		}
		encoded := base64.StdEncoding.EncodeToString(b)
		return internal.Wrap(map[string]any{
			"content_type": contentType,
			"base64":       encoded,
		}, nil), nil
	}
	var payload any
	req, err := client.NewRequest(http.MethodGet, "/devices/"+args.DeviceID+"/location_activity", url.Values{}, nil)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), nil), nil
}

var DevicesList = mcprey.MustTool(
	"prey.devices.list",
	"List devices in the account.",
	devicesList,
	mcp.WithTitleAnnotation("List devices"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

var DevicesGet = mcprey.MustTool(
	"prey.devices.get",
	"Get device details by ID.",
	devicesGet,
	mcp.WithTitleAnnotation("Get device"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

var DevicesDelete = mcprey.MustTool(
	"prey.devices.delete",
	"Delete a device (write).",
	devicesDelete,
	mcp.WithTitleAnnotation("Delete device"),
)

var DevicesReportsList = mcprey.MustTool(
	"prey.devices.reports.list",
	"List reports for a device.",
	devicesReportsList,
	mcp.WithTitleAnnotation("List device reports"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

var DevicesReportsGet = mcprey.MustTool(
	"prey.devices.reports.get",
	"Get a device report by ID.",
	devicesReportsGet,
	mcp.WithTitleAnnotation("Get device report"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

var DevicesLocationHistory = mcprey.MustTool(
	"prey.devices.location_history.get",
	"Get device location history (JSON or CSV).",
	devicesLocationHistory,
	mcp.WithTitleAnnotation("Get device location history"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

func AddDeviceTools(m *server.MCPServer) {
	DevicesList.Register(m)
	DevicesGet.Register(m)
	DevicesDelete.Register(m)
	DevicesReportsList.Register(m)
	DevicesReportsGet.Register(m)
	DevicesLocationHistory.Register(m)
}
