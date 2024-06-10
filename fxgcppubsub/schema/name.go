package schema

import "strings"

func ExtractID(schemaName string) string {
	strParts := strings.Split(schemaName, "/")

	return strParts[len(strParts)-1]
}
