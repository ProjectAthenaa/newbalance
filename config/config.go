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
				Type:           sonic.FieldTypeDropDown,
				Label:          "Size",
				FieldKey:       &sizekey,
				DropdownValues: []string{"0", "1", "1.5", "2", "2.5", "3", "3.5", "4", "4.5", "5", "5.5", "6", "6.5", "7", "7.5", "8", "8.5", "9", "9.5", "10", "10.5", "11", "11.5", "12", "12.5", "13", "13.5", "14", "15", "16", "18", "19", "20"},
			},
			{
				Type:           sonic.FieldTypeDropDown,
				Label:          "Width",
				FieldKey:       &widthkey,
				DropdownValues: []string{"2A", "B", "D", "2E", "4E", "6E", "4A", "B", "M", "W", "XW"},
			},
		},
	}
}
