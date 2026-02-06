package prey

import (
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultPreyURL = "https://api.preyproject.com/v1"

	preyAPIKeyEnvVar       = "PREY_API_KEY"
	preyAPIBaseEnvVar      = "PREY_API_BASE"
	preyTimeoutMsEnvVar    = "PREY_TIMEOUT_MS"
	preyAllowWriteEnvVar   = "PREY_ALLOW_WRITE"
	preyAllowedToolsEnvVar = "PREY_ALLOWED_TOOLS"
	preyDebugEnvVar        = "PREY_DEBUG"
	preyDisableRateLimit   = "PREY_RATE_LIMIT_DISABLE"

	preyURLHeader    = "X-Prey-URL"
	preyAPIKeyHeader = "X-Prey-API-Key"
)

type Config struct {
	Debug                   bool
	IncludeArgumentsInSpans bool
	URL                     string
	APIKey                  string
	AllowWrite              bool
	AllowedTools            map[string]struct{}
	Timeout                 time.Duration
	DisableRateLimit        bool
}

func envBool(key string) bool {
	val := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	return val == "1" || val == "true" || val == "yes"
}

func allowedToolsFromEnv() map[string]struct{} {
	val := strings.TrimSpace(os.Getenv(preyAllowedToolsEnvVar))
	if val == "" {
		return nil
	}
	set := make(map[string]struct{})
	for _, t := range strings.Split(val, ",") {
		name := strings.TrimSpace(t)
		if name != "" {
			set[name] = struct{}{}
		}
	}
	return set
}

func timeoutFromEnv() time.Duration {
	val := strings.TrimSpace(os.Getenv(preyTimeoutMsEnvVar))
	if val == "" {
		return 30 * time.Second
	}
	ms, err := strconv.Atoi(val)
	if err != nil || ms <= 0 {
		return 30 * time.Second
	}
	return time.Duration(ms) * time.Millisecond
}

func baseURLFromEnv() string {
	u := strings.TrimRight(os.Getenv(preyAPIBaseEnvVar), "/")
	if u == "" {
		return defaultPreyURL
	}
	return u
}

func apiKeyFromEnv() string {
	return strings.TrimSpace(os.Getenv(preyAPIKeyEnvVar))
}
