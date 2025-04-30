package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/himanshukandari14/go-server/internal/config"
	"github.com/himanshukandari14/go-server/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err

	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	res, err := s.Db.Exec("INSERT INTO students (name, email, age) VALUES (?, ?, ?)", name, email, age)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	var student types.Student

	err := s.Db.QueryRow("SELECT id, name, email, age FROM students WHERE id = ?", id).
		Scan(&student.Id, &student.Name, &student.Email, &student.Age)

	return student, err
}

func (s *Sqlite) GetById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM STUDENTS WHERE ID = ? LIMIT 1")

	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Email, &student.Age)

	if err != nil {

		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %w", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error %w", err)
	}
	return student, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student

		err = rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)

		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	return students, nil
}
