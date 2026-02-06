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

type LabelsListParams struct {
	Page     int `json:"page,omitempty" jsonschema:"default=1"`
	PageSize int `json:"page_size,omitempty" jsonschema:"default=20,minimum=1,maximum=100"`
}

type LabelsGetParams struct {
	LabelID string `json:"labelId" jsonschema:"description=ID of the label"`
}

type LabelsCreateParams struct {
	Name    string   `json:"name" jsonschema:"description=Label name"`
	Devices []string `json:"devices,omitempty" jsonschema:"description=Device IDs to assign"`
}

func labelsList(ctx context.Context, args LabelsListParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.labels.list", false); err != nil {
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
	req, err := client.NewRequest(http.MethodGet, "/labels", q, nil)
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

func labelsGet(ctx context.Context, args LabelsGetParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.labels.get", false); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.LabelID, "labelId"); err != nil {
		return nil, err
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	var payload any
	req, err := client.NewRequest(http.MethodGet, "/labels/"+args.LabelID, url.Values{}, nil)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), nil), nil
}

func labelsCreate(ctx context.Context, args LabelsCreateParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.labels.create", true); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.Name, "name"); err != nil {
		return nil, err
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	body := map[string]any{"name": args.Name}
	if len(args.Devices) > 0 {
		body["devices"] = args.Devices
	}
	var payload any
	req, err := client.NewRequest(http.MethodPost, "/labels", url.Values{}, body)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), nil), nil
}

var LabelsList = mcprey.MustTool(
	"prey.labels.list",
	"List labels.",
	labelsList,
	mcp.WithTitleAnnotation("List labels"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

var LabelsGet = mcprey.MustTool(
	"prey.labels.get",
	"Get label details.",
	labelsGet,
	mcp.WithTitleAnnotation("Get label"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

var LabelsCreate = mcprey.MustTool(
	"prey.labels.create",
	"Create a new label (write).",
	labelsCreate,
	mcp.WithTitleAnnotation("Create label"),
)

func AddLabelTools(m *server.MCPServer) {
	LabelsList.Register(m)
	LabelsGet.Register(m)
	LabelsCreate.Register(m)
}
