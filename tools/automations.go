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

type AutomationsListParams struct {
	Page     int `json:"page,omitempty" jsonschema:"default=1"`
	PageSize int `json:"page_size,omitempty" jsonschema:"default=20,minimum=1,maximum=100"`
}

type AutomationsGetParams struct {
	AutomationID string `json:"automationId" jsonschema:"description=ID of the automation"`
}

func automationsList(ctx context.Context, args AutomationsListParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.automations.list", false); err != nil {
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
	req, err := client.NewRequest(http.MethodGet, "/automations", q, nil)
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

func automationsGet(ctx context.Context, args AutomationsGetParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.automations.get", false); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.AutomationID, "automationId"); err != nil {
		return nil, err
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	var payload any
	req, err := client.NewRequest(http.MethodGet, "/automations/"+args.AutomationID, url.Values{}, nil)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), nil), nil
}

var AutomationsList = mcprey.MustTool(
	"prey.automations.list",
	"List automations.",
	automationsList,
	mcp.WithTitleAnnotation("List automations"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

var AutomationsGet = mcprey.MustTool(
	"prey.automations.get",
	"Get automation details.",
	automationsGet,
	mcp.WithTitleAnnotation("Get automation"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

func AddAutomationTools(m *server.MCPServer) {
	AutomationsList.Register(m)
	AutomationsGet.Register(m)
}
