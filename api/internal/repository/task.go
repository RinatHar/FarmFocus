package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepo struct {
	db *pgxpool.Pool
}

func NewTaskRepo(db *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{db: db}
}

// Create
func (r *TaskRepo) Create(ctx context.Context, t *model.Task) error {
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	query := `INSERT INTO tasks 
        (user_id, type, title, description, importance, category_id, due_date, repeat_interval, reminder_time, status, xp_reward, gold_reward, created_at, created_by, updated_at, updated_by)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) RETURNING id`

	return r.db.QueryRow(ctx, query,
		t.UserID, t.Type, t.Title, t.Description, t.Importance, t.CategoryID,
		t.DueDate, t.RepeatInterval, t.ReminderTime, t.Status, t.XPReward, t.GoldReward,
		t.CreatedAt, t.CreatedBy, t.UpdatedAt, t.UpdatedBy,
	).Scan(&t.ID)
}

// GetAll
func (r *TaskRepo) GetAll(ctx context.Context, userID int) ([]model.Task, error) {
	query := `
        SELECT id, user_id, type, title, description, importance, category_id, due_date, 
               repeat_interval, reminder_time, status, xp_reward, gold_reward, created_at,
               created_by, updated_at, updated_by
        FROM tasks
        WHERE user_id=$1 AND deleted_at IS NULL
    `
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []model.Task{}
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.Type, &t.Title, &t.Description, &t.Importance,
			&t.CategoryID, &t.DueDate, &t.RepeatInterval, &t.ReminderTime,
			&t.Status, &t.XPReward, &t.GoldReward, &t.CreatedAt,
			&t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

// GetByID
func (r *TaskRepo) GetByID(ctx context.Context, id int, userID int) (*model.Task, error) {
	var t model.Task
	query := `
        SELECT id, user_id, type, title, description, importance, category_id, due_date, 
               repeat_interval, reminder_time, status, xp_reward, gold_reward, created_at,
               created_by, updated_at, updated_by
        FROM tasks
        WHERE id=$1 AND user_id=$2 AND deleted_at IS NULL
    `
	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&t.ID, &t.UserID, &t.Type, &t.Title, &t.Description, &t.Importance,
		&t.CategoryID, &t.DueDate, &t.RepeatInterval, &t.ReminderTime,
		&t.Status, &t.XPReward, &t.GoldReward, &t.CreatedAt,
		&t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with id=%d not found or deleted", id)
		}
		return nil, err
	}
	return &t, nil
}

// Update
func (r *TaskRepo) Update(ctx context.Context, t *model.Task) error {
	t.UpdatedAt = time.Now()

	query := `
        UPDATE tasks
        SET title=$1, description=$2, importance=$3, category_id=$4, due_date=$5, 
            repeat_interval=$6, reminder_time=$7, status=$8, updated_at=$9, updated_by=$10
        WHERE id=$11 AND user_id=$12 AND deleted_at IS NULL
        RETURNING id, user_id, type, title, description, importance, category_id, due_date, 
                  repeat_interval, reminder_time, status, xp_reward, gold_reward, created_at, 
                  created_by, updated_at, updated_by
    `

	err := r.db.QueryRow(ctx, query,
		t.Title, t.Description, t.Importance, t.CategoryID, t.DueDate, t.RepeatInterval, t.ReminderTime,
		t.Status, t.UpdatedAt, t.UpdatedBy, t.ID, t.UserID,
	).Scan(
		&t.ID, &t.UserID, &t.Type, &t.Title, &t.Description, &t.Importance, &t.CategoryID,
		&t.DueDate, &t.RepeatInterval, &t.ReminderTime, &t.Status, &t.XPReward, &t.GoldReward,
		&t.CreatedAt, &t.CreatedBy, &t.UpdatedAt, &t.UpdatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("task with id=%d not found or deleted", t.ID)
		}
		return err
	}

	return nil
}

// Delete (soft delete)
func (r *TaskRepo) Delete(ctx context.Context, id, userID, deletedBy int) error {
	deletedAt := time.Now()

	query := `
        UPDATE tasks
        SET deleted_at=$1, deleted_by=$2
        WHERE id=$3 AND user_id=$4 AND deleted_at IS NULL
        RETURNING id
    `

	var deletedID int
	err := r.db.QueryRow(ctx, query, deletedAt, deletedBy, id, userID).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("task with id=%d not found, already deleted, or does not belong to user %d", id, userID)
		}
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
