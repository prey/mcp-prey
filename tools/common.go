package tools

import (
	"context"
	"fmt"

	"mcp-prey/prey"
)

func ensureToolAllowed(ctx context.Context, toolName string, write bool) error {
	cfg := prey.ConfigFromContext(ctx)
	if !prey.IsToolAllowed(cfg, toolName) {
		return fmt.Errorf("tool not allowed: %s", toolName)
	}
	if write && !cfg.AllowWrite {
		return prey.ErrWriteDisabled
	}
	return nil
}
