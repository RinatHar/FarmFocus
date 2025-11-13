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

type TaskRepo struct {
	db *pgxpool.Pool
}

func NewTaskRepo(db *pgxpool.Pool) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) HasCompletedTaskToday(ctx context.Context, userID int64) (bool, error) {
	today := time.Now().Format("2006-01-02")
	query := `
        SELECT EXISTS(
            SELECT 1 FROM progress_log 
            WHERE user_id = $1 AND DATE(created_at) = $2
        )`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, today).Scan(&exists)
	return exists, err
}

func (r *TaskRepo) Create(ctx context.Context, task *model.Task) error {
	query := `
		INSERT INTO task (user_id, type, title, description, difficulty, tag_id, due_date, repeat_interval, is_done, xp_reward, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		task.UserID, task.Type, task.Title, task.Description, task.Difficulty, task.TagID,
		task.DueDate, task.RepeatInterval, task.IsDone, task.XPReward, task.CreatedAt,
	).Scan(&task.ID)
}

func (r *TaskRepo) GetByID(ctx context.Context, id int, userID int64) (*model.Task, error) {
	var task model.Task
	query := `
		SELECT id, user_id, type, title, description, difficulty, tag_id, due_date, 
		       repeat_interval, is_done, xp_reward, created_at
		FROM task
		WHERE id = $1 AND user_id = $2
	`
	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&task.ID, &task.UserID, &task.Type, &task.Title, &task.Description, &task.Difficulty,
		&task.TagID, &task.DueDate, &task.RepeatInterval, &task.IsDone, &task.XPReward, &task.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task with id=%d not found", id)
		}
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepo) GetAll(ctx context.Context, userID int64) ([]model.Task, error) {
	query := `
		SELECT id, user_id, type, title, description, difficulty, tag_id, due_date, 
		       repeat_interval, is_done, xp_reward, created_at
		FROM task
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []model.Task{}
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(
			&task.ID, &task.UserID, &task.Type, &task.Title, &task.Description, &task.Difficulty,
			&task.TagID, &task.DueDate, &task.RepeatInterval, &task.IsDone, &task.XPReward, &task.CreatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepo) Update(ctx context.Context, task *model.Task) error {
	query := `
		UPDATE task
		SET title = $1, description = $2, difficulty = $3, tag_id = $4, due_date = $5, 
		    repeat_interval = $6, is_done = $7, xp_reward = $8
		WHERE id = $9 AND user_id = $10
		RETURNING id, user_id, type, title, description, difficulty, tag_id, due_date, 
		          repeat_interval, is_done, xp_reward, created_at
	`
	err := r.db.QueryRow(ctx, query,
		task.Title, task.Description, task.Difficulty, task.TagID, task.DueDate,
		task.RepeatInterval, task.IsDone, task.XPReward, task.ID, task.UserID,
	).Scan(
		&task.ID, &task.UserID, &task.Type, &task.Title, &task.Description, &task.Difficulty,
		&task.TagID, &task.DueDate, &task.RepeatInterval, &task.IsDone, &task.XPReward, &task.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("task with id=%d not found", task.ID)
		}
		return err
	}
	return nil
}

func (r *TaskRepo) Delete(ctx context.Context, id int, userID int64) error {
	query := `DELETE FROM task WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("task with id=%d not found or does not belong to user", id)
	}
	return nil
}

func (r *TaskRepo) GetByStatus(ctx context.Context, userID int64, isDone bool) ([]model.Task, error) {
	query := `
		SELECT id, user_id, type, title, description, difficulty, tag_id, due_date, 
		       repeat_interval, is_done, xp_reward, created_at
		FROM task
		WHERE user_id = $1 AND is_done = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID, isDone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []model.Task{}
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(
			&task.ID, &task.UserID, &task.Type, &task.Title, &task.Description, &task.Difficulty,
			&task.TagID, &task.DueDate, &task.RepeatInterval, &task.IsDone, &task.XPReward, &task.CreatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepo) MarkAsDone(ctx context.Context, id int, userID int64) error {
	query := `UPDATE task SET is_done = true WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("task with id=%d not found or does not belong to user", id)
	}
	return nil
}

func (r *TaskRepo) MarkAsUndone(ctx context.Context, id int, userID int64) error {
	query := `UPDATE task SET is_done = false WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("task with id=%d not found or does not belong to user", id)
	}
	return nil
}

func (r *TaskRepo) GetOverdue(ctx context.Context, userID int64) ([]model.Task, error) {
	query := `
		SELECT id, user_id, type, title, description, difficulty, tag_id, due_date, 
		       repeat_interval, is_done, xp_reward, created_at
		FROM task
		WHERE user_id = $1 AND is_done = false AND due_date < CURRENT_DATE
		ORDER BY due_date ASC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []model.Task{}
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(
			&task.ID, &task.UserID, &task.Type, &task.Title, &task.Description, &task.Difficulty,
			&task.TagID, &task.DueDate, &task.RepeatInterval, &task.IsDone, &task.XPReward, &task.CreatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepo) GetByTag(ctx context.Context, userID int64, tagID int) ([]model.Task, error) {
	query := `
		SELECT id, user_id, type, title, description, difficulty, tag_id, due_date, 
		       repeat_interval, is_done, xp_reward, created_at
		FROM task
		WHERE user_id = $1 AND tag_id = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID, tagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []model.Task{}
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(
			&task.ID, &task.UserID, &task.Type, &task.Title, &task.Description, &task.Difficulty,
			&task.TagID, &task.DueDate, &task.RepeatInterval, &task.IsDone, &task.XPReward, &task.CreatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
