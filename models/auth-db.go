package models

import (
	"context"
	"time"
)

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserDto struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	SecretCode string `json:"secret_code"`
}

type LoginResponseDto struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
	UserName     string    `json:"username"`
}

func (m *DBModel) GetUserByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u User
	row := m.DB.QueryRowContext(ctx, `
		Select id, username, email, password, role 
		From user 
		Where email = ?
	`, email)

	err := row.Scan(
		&u.ID,
		&u.UserName,
		&u.Email,
		&u.Password,
		&u.Role,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (m *DBModel) RegisterUser(r RegisterUserDto) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		Insert into user (username, email, password, role) value (?, ?, ?, ?)
	`
	_, err := m.DB.ExecContext(ctx, stmt, r.Username, r.Email, r.Password, 0)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) InsertToken(userId int, refreshToken string, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		Insert into refresh_token (userId, token, expires_at) value (? , ? , ?)
	`
	_, err := m.DB.ExecContext(ctx, stmt, userId, refreshToken, expiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) CheckIfUserExists(r RegisterUserDto) (usernameCheck, emailCheck bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	checkUsernameStmt := `Select Count(*) From user Where username = ?`
	row := m.DB.QueryRowContext(ctx, checkUsernameStmt, r.Username)

	var count int
	err = row.Scan(&count)
	if err != nil {
		return false, false, err
	}

	if count > 0 {
		return true, false, nil
	}

	checkEmailStmt := `Select Count(*) From user Where email = ?`
	row = m.DB.QueryRowContext(ctx, checkEmailStmt, r.Email)
	err = row.Scan(&count)
	if err != nil {
		return false, false, err
	}
	if count > 0 {
		return false, true, nil
	}

	return false, false, nil
}

func (m *DBModel) DeleteToken(refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `Delete from refresh_token where token = ?`
	_, err := m.DB.ExecContext(ctx, stmt, refreshToken)
	if err != nil {
		return err
	}

	return nil
}
