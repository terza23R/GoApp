package database

import (
	"context"
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

func (db *DB) GetUsers(ctx context.Context, limit, offset int) ([]User, error) {
	rows, err := db.Conn.QueryContext(
		ctx,
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

func (db *DB) GetUserByID(ctx context.Context, id int64) (*User, error) {
	u := &User{}
	err := db.Conn.QueryRowContext(
		ctx,
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

func (db *DB) CreateUser(ctx context.Context, u *User) error {
	res, err := db.Conn.ExecContext(
		ctx,
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

func (db *DB) UpdateUser(ctx context.Context, u *User) error {
	res, err := db.Conn.ExecContext(
		ctx,
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

func (db *DB) DeleteUser(ctx context.Context, id int64) error {
	res, err := db.Conn.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
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
