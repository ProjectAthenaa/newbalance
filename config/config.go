package config

import (
	"github.com/ProjectAthenaa/sonic-core/sonic"
	"os"
	"strings"
)

var Module *sonic.Module

func init() {
	var name = "NewBalance"

	if podName := os.Getenv("POD_NAME"); podName != "" {
		name = strings.Split(podName, "-")[0]
	}

	pidkey := "LOOKUP_pid"
	sizekey := "SIZE"
	widthkey := "WIDTH"
	genderkey := "GENDER"

	Module = &sonic.Module{
		Name: name,
		Fields: []*sonic.ModuleField{
			{
				Validation: "/\\w+?-\\d+?/",
				Type:       sonic.FieldTypeText,
				Label:      "PID",
				FieldKey:   &pidkey,
			},
			{
				Type:       sonic.FieldTypeShoeSize,
				Label:      "Size",
				FieldKey:   &sizekey,
			},
			{
				Type:       sonic.FieldTypeWidth,
				Label:      "Width",
				FieldKey:   &widthkey,
			},
			{
				Type:       sonic.FieldTypeGender,
				Label:      "Gender",
				FieldKey:   &genderkey,
			},
		},
	}
}
