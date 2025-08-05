package storage

type Storage interface{
	CreateStudent(name string, email string, age int) (uint64, error)
}