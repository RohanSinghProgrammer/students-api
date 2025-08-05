package storage

import "github.com/rohansinghprogrammer/sudents-api/internals/types"

type Storage interface{
	CreateStudent(name string, email string, age int) (uint64, error)
	GetStudentById(id uint64) (types.Student, error)
	GetStudentsList() ([]types.Student, error)
	UpdateStudent(id uint64, name string, email string, age int) (types.Student, error)
}