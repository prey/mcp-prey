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

type AccountGetParams struct{}

func accountGet(ctx context.Context, _ AccountGetParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.account.get", false); err != nil {
		return nil, err
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	var payload any
	req, err := client.NewRequest(http.MethodGet, "/account", url.Values{}, nil)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), nil), nil
}

var AccountGet = mcprey.MustTool(
	"prey.account.get",
	"Retrieve Prey account information.",
	accountGet,
	mcp.WithTitleAnnotation("Get account"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

func AddAccountTools(m *server.MCPServer) {
	AccountGet.Register(m)
}
