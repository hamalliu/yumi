package api_model

type ReqSample struct {
	SampleName    string `json:"sample_name" binding:"required"`
	FileName      string `json:"file_name" binding:"required"`
	FileExtension string `json:"file_extension" binding:"required"`
}

type ReqDownload struct {
	FileName string `query:"file_name" binding:"required"`
}

type ReqDeleteFile struct {
	FileName string `query:"file_name"`
}

type ReqConvert struct {
	FileName string `query:"file_name" binding:"required"`
}

type ReqTrack struct {
	FileName string `query:"file_name" binding:"required"`
	UserId   string `query:"user_id" binding:"required"`
}

type ReqEditor struct {
	Mode     string `query:"mode" binding:"oneof=view edit"`
	Type     string `query:"type"`
	FileName string `query:"file_name" binding:"required"`
}
