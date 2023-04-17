package models

type Program struct {
	ProgramId      string
	SpecId         string
	VuzId          string
	HasProfessions bool
	Description    string
	Form           string
	Exams          []string
	Base           Basic
}
