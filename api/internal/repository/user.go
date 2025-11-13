package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetLastLogin(ctx context.Context, userID int64) (*time.Time, error) {
	var lastLogin *time.Time
	query := `SELECT last_login FROM user_info WHERE id = $1`
	err := r.db.QueryRow(ctx, query, userID).Scan(&lastLogin)
	if err != nil {
		return nil, err
	}
	return lastLogin, nil
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO "user_info" (max_id, username, created_at, last_login, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		user.MaxID, user.Username, user.CreatedAt, user.LastLogin, user.IsActive,
	).Scan(&user.ID)
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	query := `
		SELECT id, max_id, username, created_at, last_login, is_active
		FROM "user_info"
		WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.MaxID, &user.Username, &user.CreatedAt, &user.LastLogin, &user.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with id=%d not found", id)
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByMaxID(ctx context.Context, maxID int64) (*model.User, error) {
	var user model.User
	query := `
		SELECT id, max_id, username, created_at, last_login, is_active
		FROM "user_info"
		WHERE max_id = $1
	`
	err := r.db.QueryRow(ctx, query, maxID).Scan(
		&user.ID, &user.MaxID, &user.Username, &user.CreatedAt, &user.LastLogin, &user.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // пользователь не найден - это нормально
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE "user_info"
		SET username = $1, last_login = $2, is_active = $3
		WHERE id = $4
		RETURNING id, max_id, username, created_at, last_login, is_active
	`
	err := r.db.QueryRow(ctx, query,
		user.Username, user.LastLogin, user.IsActive, user.ID,
	).Scan(
		&user.ID, &user.MaxID, &user.Username, &user.CreatedAt, &user.LastLogin, &user.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user with id=%d not found", user.ID)
		}
		return err
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM "user_info" WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user with id=%d not found", id)
	}

	return nil
}

func (r *UserRepo) GetAllActiveUsers(ctx context.Context) ([]model.User, error) {
	query := `
		SELECT id, max_id, username, created_at, last_login, is_active 
		FROM user_info 
		WHERE is_active = true
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.MaxID,
			&user.Username,
			&user.CreatedAt,
			&user.LastLogin,
			&user.IsActive,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
