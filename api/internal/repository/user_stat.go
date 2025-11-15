package repository

import (
	"context"
	"fmt"
	"time"

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

// В методе Create
func (r *UserStatRepo) Create(ctx context.Context, stat *model.UserStat) error {
	query := `
		INSERT INTO user_stat 
		(user_id, experience, gold, total_tasks_completed, total_plant_harvested, 
		 current_streak, longest_streak, is_drought, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		stat.UserID, stat.Experience, stat.Gold,
		stat.TotalTasksCompleted, stat.TotalPlantHarvested,
		stat.CurrentStreak, stat.LongestStreak, stat.IsDrought, stat.UpdatedAt,
	).Scan(&stat.ID)
}

// В методе GetByUserID
func (r *UserStatRepo) GetByUserID(ctx context.Context, userID int64) (*model.UserStat, error) {
	var stat model.UserStat
	query := `
		SELECT id, user_id, experience, gold, total_tasks_completed, total_plant_harvested, 
		       current_streak, longest_streak, is_drought, updated_at
		FROM user_stat 
		WHERE user_id = $1
	`
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&stat.ID, &stat.UserID, &stat.Experience, &stat.Gold,
		&stat.TotalTasksCompleted, &stat.TotalPlantHarvested,
		&stat.CurrentStreak, &stat.LongestStreak, &stat.IsDrought, &stat.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user stat not found for user_id=%d", userID)
		}
		return nil, err
	}
	return &stat, nil
}

// В методе Update
func (r *UserStatRepo) Update(ctx context.Context, stat *model.UserStat) error {
	query := `
		UPDATE user_stat 
		SET experience = $1, gold = $2, total_tasks_completed = $3, 
		    total_plant_harvested = $4, current_streak = $5, longest_streak = $6, 
		    is_drought = $7, updated_at = $8
		WHERE id = $9
	`
	_, err := r.db.Exec(ctx, query,
		stat.Experience, stat.Gold,
		stat.TotalTasksCompleted, stat.TotalPlantHarvested,
		stat.CurrentStreak, stat.LongestStreak, stat.IsDrought, stat.UpdatedAt,
		stat.ID,
	)
	return err
}

// RemoveExperience удаляет опыт у пользователя
func (r *UserStatRepo) RemoveExperience(ctx context.Context, userID int64, amount int64) error {
	query := `
		UPDATE user_stat 
		SET experience = GREATEST(0, experience - $1), updated_at = $2
		WHERE user_id = $3
	`

	_, err := r.db.Exec(ctx, query, amount, time.Now(), userID)
	return err
}

func (r *UserStatRepo) AddExperience(ctx context.Context, userID int64, amount int64) error {
	query := `UPDATE user_stat SET experience = experience + $1, updated_at = NOW() WHERE user_id = $2`
	_, err := r.db.Exec(ctx, query, amount, userID)
	return err
}

func (r *UserStatRepo) AddGold(ctx context.Context, userID int64, amount int64) error {
	query := `UPDATE user_stat SET gold = gold + $1, updated_at = NOW() WHERE user_id = $2`
	_, err := r.db.Exec(ctx, query, amount, userID)
	return err
}

func (r *UserStatRepo) IncrementStreak(ctx context.Context, userID int64) error {
	query := `
		UPDATE user_stat 
		SET current_streak = current_streak + 1, 
		    longest_streak = GREATEST(longest_streak, current_streak + 1),
		    updated_at = NOW() 
		WHERE user_id = $1
	`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *UserStatRepo) ResetStreak(ctx context.Context, userID int64) error {
	query := `UPDATE user_stat SET current_streak = 0, updated_at = NOW() WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

// Новые методы для работы с засухой

// SetDrought устанавливает состояние засухи для пользователя
func (r *UserStatRepo) SetDrought(ctx context.Context, userID int64, isDrought bool) error {
	query := `UPDATE user_stat SET is_drought = $1, updated_at = NOW() WHERE user_id = $2`
	_, err := r.db.Exec(ctx, query, isDrought, userID)
	return err
}

// ResetDrought сбрасывает засуху (устанавливает is_drought = false)
func (r *UserStatRepo) ResetDrought(ctx context.Context, userID int64) error {
	query := `UPDATE user_stat SET is_drought = false, updated_at = NOW() WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

// GetDroughtStatus возвращает статус засухи пользователя
func (r *UserStatRepo) GetDroughtStatus(ctx context.Context, userID int64) (bool, error) {
	var isDrought bool
	query := `SELECT is_drought FROM user_stat WHERE user_id = $1`
	err := r.db.QueryRow(ctx, query, userID).Scan(&isDrought)
	if err != nil {
		return false, err
	}
	return isDrought, nil
}

// CheckAndSetDrought проверяет и устанавливает засуху на основе последней активности
func (r *UserStatRepo) CheckAndSetDrought(ctx context.Context, userID int64, lastActivity time.Time) error {
	// Если с последней активности прошло больше 3 дней - устанавливаем засуху
	query := `
		UPDATE user_stat 
		SET is_drought = (NOW() - $1::timestamp) > INTERVAL '3 days'
		WHERE user_id = $2
	`
	_, err := r.db.Exec(ctx, query, lastActivity, userID)
	return err
}
