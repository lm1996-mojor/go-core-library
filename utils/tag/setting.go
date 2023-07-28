package tag

import "strings"

// ParseTagSetting 分割目标转化为带:分割符的map
func ParseTagSetting(str string, sep string) map[string]string {
	settings := map[string]string{}
	names := strings.Split(str, sep)
	for i := 0; i < len(names); i++ {
		values := strings.Split(names[i], ":")
		key := strings.TrimSpace(values[0])
		if len(values) >= 2 {
			settings[key] = strings.Join(values[1:], ":")
		} else if key != "" {
			settings[key] = values[1]
		}
	}
	return settings
}
