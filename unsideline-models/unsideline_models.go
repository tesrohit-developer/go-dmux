package unsideline_models

type ScanWithStartRowEndRowRequest struct {
	StartKey string
	EndKey   string
}

type ScanWithStartTimeEndTimeRequest struct {
	StartTime int64
	EndTime   int64
	StartKey  string
	EndKey    string
}

type UnsidelineByKeyRequest struct {
	Key      string
	DmuxItem string
}

type UnsidelineContainerConfig struct {
	Port int64 `json:"port"`
}
