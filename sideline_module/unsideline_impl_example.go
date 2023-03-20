package sideline_module

type ScanImpl struct {
}

type UnsidelineImpl struct {
}

func (s *ScanImpl) ScanWithStartRowEndRow(request ScanWithStartRowEndRowRequest) ([]string, error) {
	response := []string{"abc", "cde"}
	return response, nil
}

func (s *ScanImpl) ScanWithStartTimeEndTime(request ScanWithStartTimeEndTimeRequest) ([]string, error) {
	response := []string{"abc", "cde"}
	return response, nil
}

func (u *UnsidelineImpl) UnsidelineByKey(request UnsidelineByKeyRequest) (string, error) {
	return "success", nil
}

func unsidelineInitExample() {
	scanImpl := &ScanImpl{}
	unsidelineImpl := &UnsidelineImpl{}
	path := "" // config path
	unsideline.UnsidelineStart(scanImpl, unsidelineImpl, path)
}
