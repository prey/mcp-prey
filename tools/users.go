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

type UsersListParams struct {
	Page     int `json:"page,omitempty" jsonschema:"default=1"`
	PageSize int `json:"page_size,omitempty" jsonschema:"default=20,minimum=1,maximum=100"`
}

type UsersGetParams struct {
	UserID string `json:"userId" jsonschema:"description=ID of the user"`
}

func usersList(ctx context.Context, args UsersListParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.users.list", false); err != nil {
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
	req, err := client.NewRequest(http.MethodGet, "/users", q, nil)
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

func usersGet(ctx context.Context, args UsersGetParams) (any, error) {
	if err := ensureToolAllowed(ctx, "prey.users.get", false); err != nil {
		return nil, err
	}
	if err := internal.RequireID(args.UserID, "userId"); err != nil {
		return nil, err
	}
	client := prey.ClientFromContext(ctx)
	if client == nil {
		return nil, &mcprey.HardError{Err: prey.ErrMissingClient}
	}
	var payload any
	req, err := client.NewRequest(http.MethodGet, "/users/"+args.UserID, url.Values{}, nil)
	if err != nil {
		return nil, err
	}
	if err := client.DoJSON(req.WithContext(ctx), &payload); err != nil {
		return nil, err
	}
	return internal.Wrap(internal.MaskSensitive(payload), nil), nil
}

var UsersList = mcprey.MustTool(
	"prey.users.list",
	"List users in the account.",
	usersList,
	mcp.WithTitleAnnotation("List users"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

var UsersGet = mcprey.MustTool(
	"prey.users.get",
	"Get user details by ID.",
	usersGet,
	mcp.WithTitleAnnotation("Get user"),
	mcp.WithIdempotentHintAnnotation(true),
	mcp.WithReadOnlyHintAnnotation(true),
)

func AddUserTools(m *server.MCPServer) {
	UsersList.Register(m)
	UsersGet.Register(m)
}
