package entities

type PlayerMetaData struct {
	IP        string
	UserAgent string
	Host      string
	Request   []byte
}

func NewPlayerMetaData(ip, userAgent, host string, request []byte) *PlayerMetaData {
	return &PlayerMetaData{
		IP:        ip,
		UserAgent: userAgent,
		Host:      host,
		Request:   request,
	}
}
