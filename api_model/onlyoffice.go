package api_model

type ReqSample struct {
	SampleName    string `json:"sample_name" binding:"required"`
	FileName      string `json:"file_name" binding:"required"`
	FileExtension string `json:"file_extension" binding:"required"`
}

type ReqDownload struct {
	FileName string `query:"file_name" binding:"required"`
}
