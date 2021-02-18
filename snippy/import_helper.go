package snippy

import (
	"github.com/mpvl/unique"
	"strings"
)

type Exported interface {
	Exported() map[string][]string
	ResetExported()
	Name() string
}

func (s Snippy) Exported() map[string][]string {
	return s.ExporteData
}

func (s *Snippy) ResetExported() {
	s.ExporteData = map[string][]string{}
}

func (s Snippy) Import(value string) string {
	if value == "" {
		return ""
	}
	s.ExporteData["Imports"] = append(s.ExporteData["Imports"], value)
	return ""
}

func (s Snippy) Imports() string {
	elements := s.ExporteData["Imports"]
	unique.Strings(&elements)
	sb := strings.Builder{}
	for _, elem := range elements {
		if elem == "" {
			continue
		}
		sb.WriteRune('"')
		sb.WriteString(elem)
		sb.WriteRune('"')
		sb.WriteRune('\n')

	}
	return sb.String()
}
