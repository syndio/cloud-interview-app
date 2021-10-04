package employeesdb

import (
	"context"
	"database/sql"
	"fmt"
)

type Employee struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type Database struct {
	*sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{db}
}

func (db *Database) List(ctx context.Context) ([]*Employee, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, title FROM employees")
	if err != nil {
		return nil, fmt.Errorf("selecting employees: %w", err)
	}
	defer rows.Close()

	emps := make([]*Employee, 0)
	for rows.Next() {
		var emp Employee
		if err := rows.Scan(&emp.ID, &emp.Title); err != nil {
			return nil, fmt.Errorf("scanning employee: %w", err)
		}
		emps = append(emps, &emp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scanning employees: %w", err)
	}

	return emps, nil
}

func (db *Database) Create(ctx context.Context, employee *Employee) error {
	err := db.QueryRowContext(ctx, "INSERT INTO employees (title) VALUES ($1) RETURNING id", employee.Title).
		Scan(&employee.ID)
	if err != nil {
		return fmt.Errorf("inserting employee: %w", err)
	}
	return nil
}

func (db *Database) Delete(ctx context.Context, employeeID int64) error {
	res, err := db.ExecContext(ctx, "DELETE FROM employees WHERE id=$1", employeeID)
	if err != nil {
		return fmt.Errorf("deleting employee: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting affected count: %w", err)
	}
	if count == 0 {
		return sql.ErrNoRows
	}

	return nil
}
