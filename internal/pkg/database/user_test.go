package database

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func newMockDB(t *testing.T) (*DB, sqlmock.Sqlmock, func()) {
	t.Helper()

	conn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	cleanup := func() { _ = conn.Close() }
	return &DB{Conn: conn}, mock, cleanup
}

func TestGetUsers(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "email", "age"}).
		AddRow(int64(1), "A", "a@test.com", 20).
		AddRow(int64(2), "B", "b@test.com", 30)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, name, email, age FROM users ORDER BY id LIMIT ? OFFSET ?`,
	)).
		WithArgs(10, 0).
		WillReturnRows(rows)

	got, err := db.GetUsers(context.Background(), 10, 0)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 users, got %d", len(got))
	}
	if got[0].Email != "a@test.com" {
		t.Fatalf("expected first email a@test.com, got %s", got[0].Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, name, email, age FROM users WHERE id = ?`,
	)).
		WithArgs(int64(999)).
		WillReturnError(sql.ErrNoRows)

	_, err := db.GetUserByID(context.Background(), 999)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestCreateUser_SetsID(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta(
		`INSERT INTO users (name, email, age) VALUES (?, ?, ?)`,
	)).
		WithArgs("Mahir", "mahir@test.com", 24).
		WillReturnResult(sqlmock.NewResult(7, 1))

	u := &User{Name: "Mahir", Email: "mahir@test.com", Age: 24}
	err := db.CreateUser(context.Background(), u)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if u.ID != 7 {
		t.Fatalf("expected ID=7, got %d", u.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE users SET name = ?, email = ?, age = ? WHERE id = ?`,
	)).
		WithArgs("X", "x@test.com", 10, int64(123)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	u := &User{ID: 123, Name: "X", Email: "x@test.com", Age: 10}
	err := db.UpdateUser(context.Background(), u)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	db, mock, cleanup := newMockDB(t)
	defer cleanup()

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM users WHERE id = ?`,
	)).
		WithArgs(int64(123)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := db.DeleteUser(context.Background(), 123)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}
