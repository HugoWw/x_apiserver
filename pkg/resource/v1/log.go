package v1

type DebugLogReq struct {
	DebugModule string `json:"debug_module" enum:"all,apiserver"`
}
