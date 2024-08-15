package chat

const (
	LEVEL_NONE  = "NONE"
	LEVEL_DEBUG = "DEBUG"
	LEVEL_TRACE = "TRACE"
)

// if debugMode = trace, then it will log all the debug messages
// if debugMode = info, then it will log only the general logs
func shouldLog(debugMode string, logLevel string) bool {
	switch debugMode {
	case LEVEL_NONE:
		return false
	case LEVEL_TRACE:
		return true
	default:
		return debugMode == logLevel
	}
}
