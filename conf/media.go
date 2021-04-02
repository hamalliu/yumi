package conf

import "yumi/pkg/types"

// Media 媒体配置
type Media struct {
	StoragePath                string          // 附件路径
	MultipleFileUploadsMaxSize types.SpaceSize // 多媒体上传最大限制
	SingleFileUploadsMaxSize   types.SpaceSize // 单媒体上传最大限制
}
