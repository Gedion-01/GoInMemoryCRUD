package types

import (
	"strconv"

	"github.com/google/uuid"
)

const (
	minNameLength    = 3
	maxNameLength    = 30
	minAgeLength     = 1
	maxAgeLength     = 3
	minHobbiesLength = 1
)

type CreatePersonParams struct {
	Name    string
	Age     string
	Hobbies []string
}

func (params CreatePersonParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.Name) < minNameLength || len(params.Name) > maxNameLength {
		errors["name"] = "name must be between 3 and 30 characters"
	}
	if len(params.Age) < minAgeLength || len(params.Age) > maxAgeLength && params.Age != "" && params.Age != "0" {
		errors["age"] = "age must be between 1 and 3 characters"
	} else if _, err := strconv.Atoi(params.Age); err != nil {
		errors["age"] = "age must be a number"
	}
	if len(params.Hobbies) < minHobbiesLength {
		errors["hobbies"] = "hobbies must have at least 1 item"
	}
	return errors
}

type Person struct {
	ID      string
	Name    string
	Age     string
	Hobbies []string
}

func NewPersonFromParams(params CreatePersonParams) *Person {
	id := uuid.New().String()

	return &Person{
		ID:      id,
		Name:    params.Name,
		Age:     params.Age,
		Hobbies: params.Hobbies,
	}
}

func UpdatedPersonFromParams(params CreatePersonParams) *Person {
	return &Person{
		Name:    params.Name,
		Age:     params.Age,
		Hobbies: params.Hobbies,
	}
}
