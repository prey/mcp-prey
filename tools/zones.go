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

type ZoneNotificationParams struct {
	WhenIn  string `json:"when_in,omitempty" jsonschema:"description=on|off"`
	WhenOut string `json:"when_out,omitempty" jsonschema:"description=on|off"`
}

type ZoneTriggerParams struct {
	Context    string         `json:"context" jsonschema:"description=when_in|when_out"`
	ActionName string         `json:"action_name" jsonschema:"description=alarm|alert|lock|missing"`
	Options    map[string]any `json:"options,omitempty" jsonschema:"description=Action options"`
}

type ZonesListParams struct {
	Page     int `json:"page,omitempty" jsonschema:"default=1"`
	PageSize int `json:"page_size,omitempty" jsonschema:"default=20,minimum=1,maximum=100"`
}

type ZonesGetParams struct {
	ZoneID string `json:"zoneId" jsonschema:"description=ID of the zone"`
}

type ZonesCreateParams struct {
	Name          string                  `json:"name" jsonschema:"description=Zone name"`
	Lat           float64                 `json:"lat,omitempty" jsonschema:"description=Latitude"`
	Lng           float64                 `json:"lng,omitempty" jsonschema:"description=Longitude"`
	Radius        int64                   `json:"radius,omitempty" jsonschema:"description=Radius in meters"`
	Color         string                  `json:"color,omitempty" jsonschema:"description=Hex color"`
	Devices       []string                `json:"devices,omitempty" jsonschema:"description=Device IDs to assign"`
	Actions       []ZoneTriggerParams     `json:"actions,omitempty" jsonschema:"description=Zone triggers"`
	Notifications *ZoneNotificationParams `json:"notifications,omitempty" jsonschema:"description=Notification settings"`
}

type ZonesUpdateParams struct {
	ZoneID        string                  `json:"zoneId" jsonschema:"description=ID of the zone"`
	Name          string                  `json:"name,omitempty" jsonschema:"description=Zone name"`
	Lat           float64                 `json:"lat,omitempty" jsonschema:"description=Latitude"`
	Lng           float64                 `json:"lng,omitempty" jsonschema:"description=Longitude"`
	Radius        int64                   `json:"radius,omitempty" jsonschema:"description=Radius in meters"`
	Color         string                  `json:"color,omitempty" jsonschema:"description=Hex color"`
	AddDevices    []string                `json:"add_devices,omitempty" jsonschema:"description=Device IDs to add"`
	RemoveDevices []string                `json:"remove_devices,omitempty" jsonschema:"description=Device IDs to remove"`
	Actions       []ZoneTriggerParams     `json:"actions,omitempty" jsonschema:"description=Zone triggers"`
	RemoveActions []string                `json:"remove_actions,omitempty" jsonschema:"description=when_in|when_out"`
	Notifications *ZoneNotificationParams `json:"notifications,omitempty" jsonschema:"description=Notification settings"`
}

func validateZoneTrigger(t ZoneTriggerParams) error {
	if err := internal.RequireOneOf(t.Context, "context", "when_in", "when_out"); err != nil {
		return err
	}
	if err := internal.RequireOneOf(t.ActionName, "action_name", "alarm", "alert", "lock", "missing"); err != nil {
		return err
	}
	return nil
}

func validateZoneNotifications(n *ZoneNotificationParams) error {
	if n == nil {
		return nil
	}
	if n.WhenIn != "" {
		if err := internal.RequireOneOf(n.WhenIn, "when_in", "on", "off"); err != nil {
			return err
		}
	}
	if n.WhenOut != "" {
		if err := internal.RequireOneOf(n.WhenOut, "when_out", "on", "off"); err != nil {
			return err
		}
	}
	return nil
}

func zonesList(ctx context.Context, args ZonesListParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.zones.list", false); err != nil {
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
	req, err := client.NewRequest(http.MethodGet, "/zones", q, nil)
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

func zonesGet(ctx context.Context, args ZonesGetParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.zones.get", false); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.ZoneID, "zoneId"); err != nil {
		return nil, err
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	var payload any
	req, err := client.NewRequest(http.MethodGet, "/zones/"+args.ZoneID, url.Values{}, nil)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), nil), nil
}

func zonesCreate(ctx context.Context, args ZonesCreateParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.zones.create", true); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.Name, "name"); err != nil {
		return nil, err
	}
	if err := validateZoneNotifications(args.Notifications); err != nil {
		return nil, err
	}
	for _, t := range args.Actions {
		if err := validateZoneTrigger(t); err != nil {
			return nil, err
		}
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	body := map[string]any{
		"name": args.Name,
	}
	if args.Lat != 0 {
		body["lat"] = args.Lat
	}
	if args.Lng != 0 {
		body["lng"] = args.Lng
	}
	if args.Radius != 0 {
		body["radius"] = args.Radius
	}
	if args.Color != "" {
		body["color"] = args.Color
	}
	if len(args.Devices) > 0 {
		body["devices"] = args.Devices
	}
	if len(args.Actions) > 0 {
		body["actions"] = args.Actions
	}
	if args.Notifications != nil {
		body["notifications"] = args.Notifications
	}
	var payload any
	req, err := client.NewRequest(http.MethodPost, "/zones", url.Values{}, body)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), nil), nil
}

func zonesUpdate(ctx context.Context, args ZonesUpdateParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.zones.update", true); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.ZoneID, "zoneId"); err != nil {
		return nil, err
	}
	if err := validateZoneNotifications(args.Notifications); err != nil {
		return nil, err
	}
	for _, t := range args.Actions {
		if err := validateZoneTrigger(t); err != nil {
			return nil, err
		}
	}
	for _, r := range args.RemoveActions {
		if err := internal.RequireOneOf(r, "remove_actions", "when_in", "when_out"); err != nil {
			return nil, err
		}
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	body := map[string]any{}
	if args.Name != "" {
		body["name"] = args.Name
	}
	if args.Lat != 0 {
		body["lat"] = args.Lat
	}
	if args.Lng != 0 {
		body["lng"] = args.Lng
	}
	if args.Radius != 0 {
		body["radius"] = args.Radius
	}
	if args.Color != "" {
		body["color"] = args.Color
	}
	if len(args.AddDevices) > 0 {
		body["add_devices"] = args.AddDevices
	}
	if len(args.RemoveDevices) > 0 {
		body["remove_devices"] = args.RemoveDevices
	}
	if len(args.Actions) > 0 {
		body["actions"] = args.Actions
	}
	if len(args.RemoveActions) > 0 {
		body["remove_actions"] = args.RemoveActions
	}
	if args.Notifications != nil {
		body["notifications"] = args.Notifications
	}
	var payload any
	req, err := client.NewRequest(http.MethodPut, "/zones/"+args.ZoneID, url.Values{}, body)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), nil), nil
}

var ZonesList = mcprey.MustTool(
	"prey.zones.list",
	"List zones.",
	zonesList,
	mcp.WithTitleAnnotation("List zones"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

var ZonesGet = mcprey.MustTool(
	"prey.zones.get",
	"Get zone details.",
	zonesGet,
	mcp.WithTitleAnnotation("Get zone"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

var ZonesCreate = mcprey.MustTool(
	"prey.zones.create",
	"Create a new zone (write).",
	zonesCreate,
	mcp.WithTitleAnnotation("Create zone"),
)

var ZonesUpdate = mcprey.MustTool(
	"prey.zones.update",
	"Update a zone (write).",
	zonesUpdate,
	mcp.WithTitleAnnotation("Update zone"),
)

func AddZoneTools(m *server.MCPServer) {
	ZonesList.Register(m)
	ZonesGet.Register(m)
	ZonesCreate.Register(m)
	ZonesUpdate.Register(m)
}
