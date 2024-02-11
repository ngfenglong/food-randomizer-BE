package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/ngfenglong/food-randomizer-BE/pkg/models"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	RegisterUser(ctx context.Context, r RegisterUserDto) error
	InsertToken(ctx context.Context, userId int, refreshToken string, expiresAt time.Time) error
	CheckIfUserExists(ctx context.Context, r RegisterUserDto) (usernameCheck, emailCheck bool, err error)
	DeleteToken(ctx context.Context, refreshToken string) error
	IsAdminRequestPending(ctx context.Context, ar AdminRequestDto) (bool, error)
	RegisterRequest(ctx context.Context, teleId, teleUsername string) error
}

type SQLAuthRepository struct {
	db *sql.DB
}

func NewSQLAuthRepository(db *sql.DB) *SQLAuthRepository {
	return &SQLAuthRepository{db: db}
}

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

type AdminRequestDto struct {
	TelegramID       string `json:"telegram_id"`
	TelegramUsername string `json:"telegram_username"`
}

func (repo *SQLAuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	row := repo.db.QueryRowContext(ctx, `
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

func (repo *SQLAuthRepository) RegisterUser(ctx context.Context, r RegisterUserDto) error {
	stmt := `
		Insert into user (username, email, password, role) value (?, ?, ?, ?)
	`
	_, err := repo.db.ExecContext(ctx, stmt, r.Username, r.Email, r.Password, 0)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SQLAuthRepository) InsertToken(ctx context.Context, userId int, refreshToken string, expiresAt time.Time) error {
	stmt := `
		Insert into refresh_token (userId, token, expires_at) value (? , ? , ?)
	`
	_, err := repo.db.ExecContext(ctx, stmt, userId, refreshToken, expiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SQLAuthRepository) CheckIfUserExists(ctx context.Context, r RegisterUserDto) (usernameCheck, emailCheck bool, err error) {
	checkUsernameStmt := `Select Count(*) From user Where username = ?`
	row := repo.db.QueryRowContext(ctx, checkUsernameStmt, r.Username)

	var count int
	err = row.Scan(&count)
	if err != nil {
		return false, false, err
	}

	if count > 0 {
		return true, false, nil
	}

	checkEmailStmt := `Select Count(*) From user Where email = ?`
	row = repo.db.QueryRowContext(ctx, checkEmailStmt, r.Email)
	err = row.Scan(&count)
	if err != nil {
		return false, false, err
	}
	if count > 0 {
		return false, true, nil
	}

	return false, false, nil
}

func (repo *SQLAuthRepository) DeleteToken(ctx context.Context, refreshToken string) error {
	stmt := `Delete from refresh_token where token = ?`
	_, err := repo.db.ExecContext(ctx, stmt, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SQLAuthRepository) IsAdminRequestPending(ctx context.Context, ar AdminRequestDto) (bool, error) {
	var count int
	checkRequestStmt := `SELECT count(*) FROM admin_request WHERE telegram_id = ?`
	row := repo.db.QueryRowContext(ctx, checkRequestStmt, ar.TelegramID)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, err
	}

	return false, nil
}

func (repo *SQLAuthRepository) RegisterRequest(ctx context.Context, teleId, teleUsername string) error {
	stmt := `INSERT INTO admin_request(telegram_id, telegram_username) VALUE (?,?)`
	_, err := repo.db.ExecContext(ctx, stmt, teleId, teleUsername)
	if err != nil {
		return err
	}

	return nil
}
