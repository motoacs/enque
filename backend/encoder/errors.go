package encoder

// Error codes (design doc 6.3).
const (
	ErrValidation             = "E_VALIDATION"
	ErrToolNotFound           = "E_TOOL_NOT_FOUND"
	ErrToolVersionUnsupported = "E_TOOL_VERSION_UNSUPPORTED"
	ErrEncoderNotImplemented  = "E_ENCODER_NOT_IMPLEMENTED"
	ErrSessionRunning         = "E_SESSION_RUNNING"
	ErrIO                     = "E_IO"
	ErrInternal               = "E_INTERNAL"
)
