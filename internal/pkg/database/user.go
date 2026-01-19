package database

import (
	"database/sql"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func (db *DB) GetUsers(limit, offset int) ([]User, error) {
	rows, err := db.Conn.Query(
		`SELECT id, name, email, age FROM users ORDER BY id LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (db *DB) GetUserByID(id int64) (*User, error) {
	u := &User{}
	err := db.Conn.QueryRow(
		`SELECT id, name, email, age FROM users WHERE id = ?`,
		id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Age)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}

func (db *DB) CreateUser(u *User) error {
	res, err := db.Conn.Exec(
		`INSERT INTO users (name, email, age) VALUES (?, ?, ?)`,
		u.Name, u.Email, u.Age,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id
	return nil

}

func (db *DB) UpdateUser(u *User) error {
	res, err := db.Conn.Exec(
		`UPDATE users SET name = ?, email = ?, age = ? WHERE id = ?`,
		u.Name, u.Email, u.Age, u.ID,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (db *DB) DeleteUser(id int64) error {
	res, err := db.Conn.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrUserNotFound
	}

	return nil
}
