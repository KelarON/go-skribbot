package model

type DrawingType string

const (
	DRAWING_TYPE_LINE  DrawingType = "line"
	DRAWING_TYPE_POINT DrawingType = "point"
)

var STATUS_LIST = [...]DrawingType{
	DRAWING_TYPE_LINE,
	DRAWING_TYPE_POINT,
}

func GetAllDrawingStatuses() []string {
	var statusList []string
	for _, statusText := range STATUS_LIST {
		statusList = append(statusList, string(statusText))
	}
	return statusList
}
