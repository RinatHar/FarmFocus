package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStatRepo struct {
	db *pgxpool.Pool
}

func NewUserStatRepo(db *pgxpool.Pool) *UserStatRepo {
	return &UserStatRepo{db: db}
}

func (r *UserStatRepo) GetByUserID(ctx context.Context, userID int64) (*model.UserStat, error) {
	var stat model.UserStat
	query := `
		SELECT id, user_id, experience, gold, streak, total_plant_harvested, total_task_completed, updated_at
		FROM user_stat
		WHERE user_id = $1
	`
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&stat.ID, &stat.UserID, &stat.Experience, &stat.Gold, &stat.Streak,
		&stat.TotalPlantHarvested, &stat.TotalTaskCompleted, &stat.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user stat for user_id=%d not found", userID)
		}
		return nil, err
	}
	return &stat, nil
}

func (r *UserStatRepo) Create(ctx context.Context, stat *model.UserStat) error {
	query := `
		INSERT INTO user_stat (user_id, experience, gold, streak, total_plant_harvested, total_task_completed, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		stat.UserID, stat.Experience, stat.Gold, stat.Streak,
		stat.TotalPlantHarvested, stat.TotalTaskCompleted, stat.UpdatedAt,
	).Scan(&stat.ID)
}

func (r *UserStatRepo) Update(ctx context.Context, stat *model.UserStat) error {
	query := `
		UPDATE user_stat
		SET experience = $1, gold = $2, streak = $3, total_plant_harvested = $4, total_task_completed = $5, updated_at = $6
		WHERE user_id = $7
		RETURNING id
	`
	err := r.db.QueryRow(ctx, query,
		stat.Experience, stat.Gold, stat.Streak, stat.TotalPlantHarvested,
		stat.TotalTaskCompleted, stat.UpdatedAt, stat.UserID,
	).Scan(&stat.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user stat for user_id=%d not found", stat.UserID)
		}
		return err
	}
	return nil
}

func (r *UserStatRepo) AddExperience(ctx context.Context, userID int64, amount int64) error {
	query := `
		UPDATE user_stat
		SET experience = experience + $1, updated_at = NOW()
		WHERE user_id = $2
	`
	result, err := r.db.Exec(ctx, query, amount, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user stat for user_id=%d not found", userID)
	}

	return nil
}

func (r *UserStatRepo) AddGold(ctx context.Context, userID int64, amount int64) error {
	query := `
		UPDATE user_stat
		SET gold = gold + $1, updated_at = NOW()
		WHERE user_id = $2
	`
	result, err := r.db.Exec(ctx, query, amount, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user stat for user_id=%d not found", userID)
	}

	return nil
}

func (r *UserStatRepo) IncrementStreak(ctx context.Context, userID int64) error {
	query := `
		UPDATE user_stat
		SET streak = streak + 1, updated_at = NOW()
		WHERE user_id = $1
	`
	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user stat for user_id=%d not found", userID)
	}

	return nil
}

func (r *UserStatRepo) ResetStreak(ctx context.Context, userID int64) error {
	query := `
		UPDATE user_stat
		SET streak = 0, updated_at = NOW()
		WHERE user_id = $1
	`
	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user stat for user_id=%d not found", userID)
	}

	return nil
}
