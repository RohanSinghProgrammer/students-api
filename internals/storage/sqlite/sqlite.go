package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rohansinghprogrammer/sudents-api/internals/config"
	"github.com/rohansinghprogrammer/sudents-api/internals/types"
)

type Sqlite struct {
	DB *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error)  {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER UNIQUE PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		DB: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (uint64, error) {
	stmt, err := s.DB.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	results, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := results.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(lastId), nil
}

func (s *Sqlite) GetStudentById(id uint64) (types.Student, error) {
	stmt, err := s.DB.Prepare(`SELECT id, name, email, age FROM students WHERE id=? LIMIT 1`)
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, err
		}
		return types.Student{}, err
	}

	return student, nil
}

func (s *Sqlite) GetStudentsList() ([]types.Student, error) {
	stmt, err := s.DB.Prepare(`SELECT id, name, email, age FROM students`)
	if err != nil {
		return []types.Student{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return []types.Student{}, err
	}
	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}
	return students, nil
}

func (s *Sqlite) UpdateStudent(id uint64, name string, email string, age int) (types.Student, error) {

    // This is a good practice to prevent silent failures.
    existingStudent, err := s.GetStudentById(id)
    if err != nil {
        return types.Student{}, err // Will return sql.ErrNoRows if student not found
    }

    // Prepare the update statement.
    stmt, err := s.DB.Prepare("UPDATE students SET name=?, email=?, age=? WHERE id=?")
    if err != nil {
        return types.Student{}, err
    }
    defer stmt.Close()

    // Execute the update query.
    result, err := stmt.Exec(name, email, age, id)
    if err != nil {
        return types.Student{}, err
    }

    // Check how many rows were affected.
    // If RowsAffected is 0, it means the update didn't happen (though in this
    // case, our previous check for the student's existence makes this less likely).
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return types.Student{}, err
    }

    if rowsAffected == 0 {
        // This case might be rare due to the existence check, but it's good practice.
        return types.Student{}, sql.ErrNoRows
    }

    // Construct the updated student object from the parameters and return it.
    // This avoids a second, unnecessary database query.
    updatedStudent := types.Student{
        ID:    existingStudent.ID, // Keep the original ID
        Name:  name,
        Email: email,
        Age:   age,
    }

    return updatedStudent, nil
}