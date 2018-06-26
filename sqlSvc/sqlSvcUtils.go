package sqlSvc

import (
	"database/sql"
)

func convertToString(input []sql.NullString) []string {
	colValues := make([]string, len(input))
	for i := range input {
		if input[i].Valid {
			colValues[i] = input[i].String
		} else {
			colValues[i] = ""
		}
	}
	return colValues
}
