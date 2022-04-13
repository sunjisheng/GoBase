package Utility

import (
"bufio"
"fmt"
"io"
"os"
"regexp"
	"strconv"
	"strings"
)

type Section map[string]string

type IniFile struct {
	sections map[string]Section
}

type ErrSyntax struct {
	Line   int
	Source string
}

var (
	sectionRegex = regexp.MustCompile(`^\[(.*)\]$`)
	assignRegex  = regexp.MustCompile(`^([^=]+)=(.*)$`)
)

func (e ErrSyntax) Error() string {
	return fmt.Sprintf("invalid INI syntax on line %d: %s", e.Line, e.Source)
}

func (this *IniFile) Section(name string) Section {
	section := this.sections[name]
	if section == nil {
		section = make(Section)
		this.sections[name] = section
	}
	return section
}

func (this *IniFile) ExistSection(name string) bool {
	_, ok := this.sections[name]
	if ok {
		return true
	} else {
		return false
	}
}

func (this *IniFile) GetString(section, key string) string {
	if s := this.sections[section]; s != nil {
		value, ok := s[key]
		if ok {
			return value
		}
	}
	return ""
}

func (this *IniFile) GetInt(section, key string) uint32 {
	if s := this.sections[section]; s != nil {
		value, ok := s[key]
		if ok {
			intVal, err := strconv.Atoi(value)
			if err == nil {
				return uint32(intVal)
			}
		}
	}
	return 0
}

func (this *IniFile) LoadFile(file string) (err error) {
	in, err := os.Open(file)
	if err != nil {
		return
	}
	this.sections = make(map[string]Section)
	defer in.Close()
	bufin := bufio.NewReader(in)
	return parseFile(bufin, this)
}

func parseFile(in *bufio.Reader, file *IniFile) (err error) {
	section := ""
	lineNum := 0
	for done := false; !done; {
		var line string
		if line, err = in.ReadString('\n'); err != nil {
			if err == io.EOF {
				done = true
			} else {
				return
			}
		}
		lineNum++
		line = strings.TrimSpace(line)
		if len(line) == 0 { //空行
			continue
		}
		if line[0] == ';' || line[0] == '#' { //注释
			continue
		}

		if groups := assignRegex.FindStringSubmatch(line); groups != nil {
			key, val := groups[1], groups[2]
			key, val = strings.TrimSpace(key), strings.TrimSpace(val)
			file.Section(section)[key] = val
		} else if groups := sectionRegex.FindStringSubmatch(line); groups != nil {
			name := strings.TrimSpace(groups[1])
			section = name
			// Create the section if it does not exist
			file.Section(section)
		} else {
			return ErrSyntax{lineNum, line}
		}

	}
	return nil
}

