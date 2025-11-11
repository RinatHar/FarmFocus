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

type ProgressLogRepo struct {
	db *pgxpool.Pool
}

func NewProgressLogRepo(db *pgxpool.Pool) *ProgressLogRepo {
	return &ProgressLogRepo{db: db}
}

func (r *ProgressLogRepo) Create(ctx context.Context, log *model.ProgressLog) error {
	query := `
		INSERT INTO progress_log (user_id, task_id, xp_earned, gold_earned, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		log.UserID, log.TaskID, log.XPEarned, log.GoldEarned, log.CreatedAt,
	).Scan(&log.ID)
}

func (r *ProgressLogRepo) GetByID(ctx context.Context, id int) (*model.ProgressLog, error) {
	var log model.ProgressLog
	query := `
		SELECT id, user_id, task_id, xp_earned, gold_earned, created_at
		FROM progress_log
		WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&log.ID, &log.UserID, &log.TaskID, &log.XPEarned, &log.GoldEarned, &log.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("progress log with id=%d not found", id)
		}
		return nil, err
	}
	return &log, nil
}

func (r *ProgressLogRepo) GetByUser(ctx context.Context, userID int64) ([]model.ProgressLog, error) {
	query := `
		SELECT id, user_id, task_id, xp_earned, gold_earned, created_at
		FROM progress_log
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []model.ProgressLog{}
	for rows.Next() {
		var log model.ProgressLog
		if err := rows.Scan(
			&log.ID, &log.UserID, &log.TaskID, &log.XPEarned, &log.GoldEarned, &log.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (r *ProgressLogRepo) GetByTask(ctx context.Context, taskID int) ([]model.ProgressLog, error) {
	query := `
		SELECT id, user_id, task_id, xp_earned, gold_earned, created_at
		FROM progress_log
		WHERE task_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []model.ProgressLog{}
	for rows.Next() {
		var log model.ProgressLog
		if err := rows.Scan(
			&log.ID, &log.UserID, &log.TaskID, &log.XPEarned, &log.GoldEarned, &log.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (r *ProgressLogRepo) GetByUserAndDateRange(ctx context.Context, userID int64, startDate, endDate time.Time) ([]model.ProgressLog, error) {
	query := `
		SELECT id, user_id, task_id, xp_earned, gold_earned, created_at
		FROM progress_log
		WHERE user_id = $1 AND created_at BETWEEN $2 AND $3
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []model.ProgressLog{}
	for rows.Next() {
		var log model.ProgressLog
		if err := rows.Scan(
			&log.ID, &log.UserID, &log.TaskID, &log.XPEarned, &log.GoldEarned, &log.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (r *ProgressLogRepo) GetUserStats(ctx context.Context, userID int64, startDate, endDate time.Time) (totalXP, totalGold int, err error) {
	query := `
		SELECT COALESCE(SUM(xp_earned), 0), COALESCE(SUM(gold_earned), 0)
		FROM progress_log
		WHERE user_id = $1 AND created_at BETWEEN $2 AND $3
	`
	err = r.db.QueryRow(ctx, query, userID, startDate, endDate).Scan(&totalXP, &totalGold)
	return totalXP, totalGold, err
}

func (r *ProgressLogRepo) GetWithTaskDetails(ctx context.Context, userID int64) ([]model.ProgressLogWithTask, error) {
	query := `
		SELECT pl.id, pl.user_id, pl.task_id, pl.xp_earned, pl.gold_earned, pl.created_at,
		       t.title as task_title, t.type as task_type
		FROM progress_log pl
		INNER JOIN task t ON pl.task_id = t.id
		WHERE pl.user_id = $1
		ORDER BY pl.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []model.ProgressLogWithTask{}
	for rows.Next() {
		var log model.ProgressLogWithTask
		if err := rows.Scan(
			&log.ID, &log.UserID, &log.TaskID, &log.XPEarned, &log.GoldEarned, &log.CreatedAt,
			&log.TaskTitle, &log.TaskType,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (r *ProgressLogRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM progress_log WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("progress log with id=%d not found", id)
	}
	return nil
}
