package sideline_models

type SidelineMessageResponse struct {
	Success                     bool
	ConcurrentModificationError bool
	UnknownError                bool
	ErrorMessage                string
}
