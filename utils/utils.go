package utils

import "strings"

type Utils interface {
	ConvertVideoNameToProperFileName(videoTitle string) string
}

type UtilsImpl struct {
}

func NewUtilsImpl() Utils {
	return UtilsImpl{}
}

func (ut UtilsImpl) ConvertVideoNameToProperFileName(videoTitle string) string {

	return strings.ReplaceAll(videoTitle, " ", "_") + ".mp4"
}
