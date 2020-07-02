package api_model

type ReqSample struct {
	FileName      string `json:"file_name" binding:"required"`
	FileExtension string `json:"file_extension" binding:"required"`
}
