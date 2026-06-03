package entities

type FileContentType byte

const (
	Unown FileContentType = iota
	PDF
	PNG
	JPEG
	MKV
	MP4
)

func (s FileContentType) String() string {
	switch s {
	case PDF:
		return "PDF"
	case PNG:
		return "PNG"
	case JPEG:
		return "Failed"
	case MKV:
		return "MKV"
	case MP4:
		return "MP4"
	default:
		return ""
	}
}
