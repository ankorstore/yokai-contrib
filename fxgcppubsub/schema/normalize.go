package schema

import "strings"

// NormalizeSchemaID normalizes a given schemaID to be compatible with the "projects/{projectID}/schemas/{schemaID}" format.
func NormalizeSchemaID(schemaID string) string {
	if strings.Contains(schemaID, "/") {
		strParts := strings.Split(schemaID, "/")

		return strParts[len(strParts)-1]
	}

	return schemaID
}
