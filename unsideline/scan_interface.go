package unsideline

import unsideline_models "github.com/flipkart-incubator/go-dmux/unsideline-models"

type Scan interface {
	ScanWithStartRowEndRow(request unsideline_models.ScanWithStartRowEndRowRequest) ([]string, error)
	ScanWithStartTimeEndTime(request unsideline_models.ScanWithStartTimeEndTimeRequest) ([]string, error)
}
