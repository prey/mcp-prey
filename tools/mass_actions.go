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

type MassActionsListParams struct {
	Page     int `json:"page,omitempty" jsonschema:"default=1"`
	PageSize int `json:"page_size,omitempty" jsonschema:"default=20,minimum=1,maximum=100"`
}

type MassActionsGetParams struct {
	MassActionID string `json:"massActionId" jsonschema:"description=ID of the mass action"`
}

func massActionsList(ctx context.Context, args MassActionsListParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.mass_actions.list", false); err != nil {
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
	req, err := client.NewRequest(http.MethodGet, "/mass_actions", q, nil)
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

func massActionsGet(ctx context.Context, args MassActionsGetParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.mass_actions.get", false); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.MassActionID, "massActionId"); err != nil {
		return nil, err
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	var payload any
	req, err := client.NewRequest(http.MethodGet, "/mass_actions/"+args.MassActionID, url.Values{}, nil)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), nil), nil
}

var MassActionsList = mcprey.MustTool(
	"prey.mass_actions.list",
	"List mass actions.",
	massActionsList,
	mcp.WithTitleAnnotation("List mass actions"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

var MassActionsGet = mcprey.MustTool(
	"prey.mass_actions.get",
	"Get mass action details.",
	massActionsGet,
	mcp.WithTitleAnnotation("Get mass action"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

func AddMassActionTools(m *server.MCPServer) {
	MassActionsList.Register(m)
	MassActionsGet.Register(m)
}
