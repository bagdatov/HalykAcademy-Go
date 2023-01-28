package main

import (
	"fmt"

	"main/flatjson"
)

type OutStruct struct {
	Name                                string
	Age                                 int64
	EducationDegree                     string
	educationAverageGrade               float32 // Приватное поле должно остаться пустым
	EducationUniversity                 string
	EducationFacultyName                string
	EducationFacultyDepartment          string
	EducationFacultyAdviserFirstName    string
	EducationFacultyAdviserLastName     string
	EducationFacultyAdviserDegree       string
	EducationFacultyAdviserArticleCount int32
}

func main() {
	data := `
	{
		"Name": "asd",
		"Age": 21,
		"Education": {
			"Degree": "bachelor",
			"AverageGrade": 4.4,
			"University": "ENU",
			"Faculty": {
				"Name": "Mechmath",
				"Department": "Mathematical and computer modeling",
				"Adviser": {
					"FirstName": "Ivanov",
					"LastName": "Ivan",
					"Degree": "PhD",
					"ArticleCount": 30
				}
			}
		}
	}
	`

	out := &OutStruct{}
	err := flatjson.Unmarshal([]byte(data), out)
	if err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		return
	}
	fmt.Printf("%+v\n", out)

	out2 := &OutStruct{}
	err = flatjson.UnmarshalAlternative([]byte(data), out2)
	if err != nil {
		fmt.Printf("Unmarshal error: %v\n", err)
		return
	}
	fmt.Printf("%+v\n", out2)
}
