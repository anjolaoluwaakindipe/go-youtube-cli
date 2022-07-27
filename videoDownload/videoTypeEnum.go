package videodownload

// Download type enum
type DownloadType int

const (
	SingleVideo DownloadType = iota
	PlayList
)

func (dt DownloadType) String() string {
	switch dt {
	case SingleVideo:
		return "SingleVidel"
	case PlayList:
		return "PlayList"
	}
	return "unknown"
}
