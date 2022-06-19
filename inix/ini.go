package inix

import (
	"errors"

	"gopkg.in/ini.v1"
)

// ParseIni returns info of a section of an ini file.
func ParseIni(iniPath, section string) (*ini.Section, error) {
	f, err := ini.Load(iniPath)
	if err != nil {
		return nil, errors.New("faild to load an ini file")
	}
	return f.Section(section), nil
}
