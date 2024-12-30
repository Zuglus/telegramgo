package utils

import "strings"

// SplitMonths разделяет строку с месяцами на срез строк
func SplitMonths(monthsStr string) []string {
	if monthsStr == "" {
		return []string{}
	}
	return strings.Split(monthsStr, ",")
}