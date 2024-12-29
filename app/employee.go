package app

import (
	"bytes"
	"encoding/csv"
	"pochtmx/static"
	"sort"
	"strings"
)

type Employee struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	JobTitle  string `json:"job_title"`
}

// GetEmployees reads the embedded CSV file and returns
// a filtered, sorted result set
func GetEmployees(search string) ([]Employee, error) {
	search = strings.ToLower(strings.TrimSpace(search))
	employees := make([]Employee, 0)

	// read the file
	bb, err := static.EmbeddedFS().ReadFile("data/employees.csv")
	if err != nil {
		return employees, err
	}
	reader := csv.NewReader(bytes.NewReader(bb))
	records, err := reader.ReadAll()
	if err != nil {
		return employees, err
	}

	// filter the file
	for _, r := range records {
		e := Employee{r[0], r[1], r[2]}
		if search == "" || strings.Contains(strings.ToLower(e.FirstName), search) || strings.Contains(strings.ToLower(e.LastName), search) {
			employees = append(employees, e)
		}
	}

	// sort the results
	sort.Slice(employees, func(i, j int) bool {
		return employees[i].LastName < employees[j].LastName
	})

	return employees, nil
}
