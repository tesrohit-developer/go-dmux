package unsideline

import unsideline_models "github.com/flipkart-incubator/go-dmux/unsideline-models"

type Unsideline interface {
	UnsidelineByKey(request unsideline_models.UnsidelineByKeyRequest) (string, error)
}
