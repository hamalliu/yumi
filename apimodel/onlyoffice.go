package apimodel

//ReqSample ...
type ReqSample struct {
	SampleName    string `json:"sample_name" binding:"required"`
	FileName      string `json:"file_name" binding:"required"`
	FileExtension string `json:"file_extension" binding:"required"`
}

//ReqDownload ...
type ReqDownload struct {
	FileName string `query:"file_name" binding:"required"`
}

//ReqDeleteFile ...
type ReqDeleteFile struct {
	FileName string `query:"file_name"`
}

//ReqConvert ...
type ReqConvert struct {
	FileName string `query:"file_name" binding:"required"`
}

//ReqTrack ...
type ReqTrack struct {
	FileName string `query:"file_name" binding:"required"`
	UserID   string `query:"user_id" binding:"required"`
}

//ReqEditor ...
type ReqEditor struct {
	Mode     string `query:"mode" binding:"oneof=view edit"`
	Type     string `query:"type"`
	FileName string `query:"file_name" binding:"required"`
}
