package paginate

import (
	"log/slog"
	"os"
	"strconv"
)

// GlobalConfig holds the global configuration for go-paginate
type GlobalConfig struct {
	DefaultLimit int
	MaxLimit     int
	DebugMode    bool
	logger       *slog.Logger
}

// globalConfig is the singleton instance
var globalConfig = &GlobalConfig{
	DefaultLimit: 10,  // default value
	MaxLimit:     100, // default value
	DebugMode:    false,
	logger:       slog.Default(),
}

// init loads configuration from environment variables
func init() {
	loadFromEnv()
}

// loadFromEnv loads configuration from environment variables
func loadFromEnv() {
	logger := slog.With("component", "go-paginate-config")

	// Load GO_PAGINATE_DEBUG
	if debugStr := os.Getenv("GO_PAGINATE_DEBUG"); debugStr != "" {
		if debug, err := strconv.ParseBool(debugStr); err == nil {
			globalConfig.DebugMode = debug
			logger.Info("Debug mode loaded from environment",
				"GO_PAGINATE_DEBUG", debug)
		} else {
			logger.Warn("Invalid GO_PAGINATE_DEBUG value, using default",
				"value", debugStr,
				"error", err,
				"default", globalConfig.DebugMode)
		}
	}

	// Load GO_PAGINATE_DEFAULT_LIMIT
	if limitStr := os.Getenv("GO_PAGINATE_DEFAULT_LIMIT"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			globalConfig.DefaultLimit = limit
			logger.Info("Default limit loaded from environment",
				"GO_PAGINATE_DEFAULT_LIMIT", limit)
		} else {
			logger.Warn("Invalid GO_PAGINATE_DEFAULT_LIMIT value, using default",
				"value", limitStr,
				"error", err,
				"default", globalConfig.DefaultLimit)
		}
	}

	// Load GO_PAGINATE_MAX_LIMIT
	if maxLimitStr := os.Getenv("GO_PAGINATE_MAX_LIMIT"); maxLimitStr != "" {
		if maxLimit, err := strconv.Atoi(maxLimitStr); err == nil && maxLimit > 0 {
			globalConfig.MaxLimit = maxLimit
			logger.Info("Max limit loaded from environment",
				"GO_PAGINATE_MAX_LIMIT", maxLimit)
		} else {
			logger.Warn("Invalid GO_PAGINATE_MAX_LIMIT value, using default",
				"value", maxLimitStr,
				"error", err,
				"default", globalConfig.MaxLimit)
		}
	}

	logger.Info("Go-paginate configuration initialized",
		"defaultLimit", globalConfig.DefaultLimit,
		"maxLimit", globalConfig.MaxLimit,
		"debugMode", globalConfig.DebugMode)
}

// SetDefaultLimit sets the global default limit
func SetDefaultLimit(limit int) {
	logger := slog.With("component", "go-paginate-config")

	if limit <= 0 {
		logger.Error("Invalid default limit value, must be greater than 0",
			"attempted_value", limit,
			"current_value", globalConfig.DefaultLimit)
		return
	}

	oldValue := globalConfig.DefaultLimit
	globalConfig.DefaultLimit = limit

	logger.Info("Default limit updated",
		"old_value", oldValue,
		"new_value", limit)
}

// SetMaxLimit sets the global maximum limit
func SetMaxLimit(maxLimit int) {
	logger := slog.With("component", "go-paginate-config")

	if maxLimit <= 0 {
		logger.Error("Invalid max limit value, must be greater than 0",
			"attempted_value", maxLimit,
			"current_value", globalConfig.MaxLimit)
		return
	}

	oldValue := globalConfig.MaxLimit
	globalConfig.MaxLimit = maxLimit

	logger.Info("Max limit updated",
		"old_value", oldValue,
		"new_value", maxLimit)
}

// SetDebugMode sets the global debug mode
func SetDebugMode(debug bool) {
	logger := slog.With("component", "go-paginate-config")

	oldValue := globalConfig.DebugMode
	globalConfig.DebugMode = debug

	logger.Info("Debug mode updated",
		"old_value", oldValue,
		"new_value", debug)
}

// GetDefaultLimit returns the global default limit
func GetDefaultLimit() int {
	return globalConfig.DefaultLimit
}

// GetMaxLimit returns the global maximum limit
func GetMaxLimit() int {
	return globalConfig.MaxLimit
}

// IsDebugMode returns the global debug mode status
func IsDebugMode() bool {
	return globalConfig.DebugMode
}

// SetLogger sets a custom logger for the configuration
func SetLogger(logger *slog.Logger) {
	globalConfig.logger = logger
}

// logSQL logs SQL queries when debug mode is enabled
func logSQL(operation, query string, args []any) {
	if globalConfig.DebugMode {
		logger := slog.With("component", "go-paginate-sql")
		logger.Info("Generated SQL query",
			"operation", operation,
			"query", query,
			"args", args,
			"args_count", len(args))
	}
}
