package entity

type PublishFileRequest struct {
	Classify string `json:"classify"`
	FileName 	string `json:"file_name"`
}

func (p PublishFileRequest) Topic() string {
	return MQPrefix + p.Classify
}
