package schemas

import "strings"

func UpperStartCamel(camel string) string {
	names := strings.Split(camel, "_")
	t := make([]string, 0, len(names)*2)
	for i := range names {
		t = append(t, strings.ToUpper(names[i][:1]), strings.ToLower(names[i][1:]))
	}

	return strings.Join(t, "")
}

func LowerStartCamel(camel string) string {
	names := strings.Split(camel, "_")
	t := make([]string, 0, len(names))
	for i := range names {
		t = append(t, strings.ToLower(names[i]))
	}

	return strings.Join(t, "")
}
