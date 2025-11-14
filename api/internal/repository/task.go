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
            WHERE user_id = $1 AND DATE(created_at) = $2 AND task_id IS NOT NULL
        )`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, today).Scan(&exists)
	return exists, err
}

func (r *TaskRepo) Create(ctx context.Context, task *model.Task) error {
	query := `
		INSERT INTO task (user_id, title, description, difficulty, tag_id, date, done, xp_reward, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		task.UserID, task.Title, task.Description, task.Difficulty, task.TagID,
		task.Date, task.Done, task.XPReward, task.CreatedAt,
	).Scan(&task.ID)
}

func (r *TaskRepo) GetByID(ctx context.Context, id int, userID int64) (*model.Task, error) {
	var task model.Task
	var tagName, tagColor *string

	query := `
		SELECT t.id, t.user_id, t.title, t.description, t.difficulty, t.tag_id, 
		       t.date, t.done, t.xp_reward, t.created_at,
		       tag.name, tag.color
		FROM task t
		LEFT JOIN tag ON t.tag_id = tag.id
		WHERE t.id = $1 AND t.user_id = $2
	`

	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&task.ID, &task.UserID, &task.Title, &task.Description, &task.Difficulty,
		&task.TagID, &task.Date, &task.Done, &task.XPReward, &task.CreatedAt,
		&tagName, &tagColor,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task with id=%d not found", id)
		}
		return nil, err
	}

	// Если есть тег, заполняем структуру Tag
	if tagName != nil && task.TagID != nil {
		task.Tag = &model.Tag{
			ID:    *task.TagID,
			Name:  *tagName,
			Color: *tagColor,
		}
	}

	return &task, nil
}

func (r *TaskRepo) GetAll(ctx context.Context, userID int64) ([]model.Task, error) {
	query := `
		SELECT t.id, t.user_id, t.title, t.description, t.difficulty, t.tag_id, 
		       t.date, t.done, t.xp_reward, t.created_at,
		       tag.name, tag.color
		FROM task t
		LEFT JOIN tag ON t.tag_id = tag.id
		WHERE t.user_id = $1
		ORDER BY t.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []model.Task{}
	for rows.Next() {
		var task model.Task
		var tagName, tagColor *string

		if err := rows.Scan(
			&task.ID, &task.UserID, &task.Title, &task.Description, &task.Difficulty,
			&task.TagID, &task.Date, &task.Done, &task.XPReward, &task.CreatedAt,
			&tagName, &tagColor,
		); err != nil {
			return nil, err
		}

		// Если есть тег, заполняем структуру Tag
		if tagName != nil && task.TagID != nil {
			task.Tag = &model.Tag{
				ID:    *task.TagID,
				Name:  *tagName,
				Color: *tagColor,
			}
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepo) Update(ctx context.Context, task *model.Task) error {
	query := `
		UPDATE task
		SET title = $1, description = $2, difficulty = $3, tag_id = $4, 
		    date = $5, done = $6, xp_reward = $7
		WHERE id = $8 AND user_id = $9
	`
	result, err := r.db.Exec(ctx, query,
		task.Title, task.Description, task.Difficulty, task.TagID, task.Date,
		task.Done, task.XPReward, task.ID, task.UserID,
	)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("task with id=%d not found or does not belong to user", task.ID)
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

func (r *TaskRepo) GetByStatus(ctx context.Context, userID int64, done bool) ([]model.Task, error) {
	query := `
		SELECT t.id, t.user_id, t.title, t.description, t.difficulty, t.tag_id, 
		       t.date, t.done, t.xp_reward, t.created_at,
		       tag.name, tag.color
		FROM task t
		LEFT JOIN tag ON t.tag_id = tag.id
		WHERE t.user_id = $1 AND t.done = $2
		ORDER BY t.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID, done)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []model.Task{}
	for rows.Next() {
		var task model.Task
		var tagName, tagColor *string

		if err := rows.Scan(
			&task.ID, &task.UserID, &task.Title, &task.Description, &task.Difficulty,
			&task.TagID, &task.Date, &task.Done, &task.XPReward, &task.CreatedAt,
			&tagName, &tagColor,
		); err != nil {
			return nil, err
		}

		if tagName != nil && task.TagID != nil {
			task.Tag = &model.Tag{
				ID:    *task.TagID,
				Name:  *tagName,
				Color: *tagColor,
			}
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepo) MarkAsDone(ctx context.Context, id int, userID int64) error {
	query := `UPDATE task SET done = true WHERE id = $1 AND user_id = $2`
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
	query := `UPDATE task SET done = false WHERE id = $1 AND user_id = $2`
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

func (r *TaskRepo) GetByDate(ctx context.Context, userID int64, date time.Time) ([]model.Task, error) {
	query := `
		SELECT t.id, t.user_id, t.title, t.description, t.difficulty, t.tag_id, 
		       t.date, t.done, t.xp_reward, t.created_at,
		       tag.name, tag.color
		FROM task t
		LEFT JOIN tag ON t.tag_id = tag.id
		WHERE t.user_id = $1 AND t.date = $2
		ORDER BY t.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []model.Task{}
	for rows.Next() {
		var task model.Task
		var tagName, tagColor *string

		if err := rows.Scan(
			&task.ID, &task.UserID, &task.Title, &task.Description, &task.Difficulty,
			&task.TagID, &task.Date, &task.Done, &task.XPReward, &task.CreatedAt,
			&tagName, &tagColor,
		); err != nil {
			return nil, err
		}

		if tagName != nil && task.TagID != nil {
			task.Tag = &model.Tag{
				ID:    *task.TagID,
				Name:  *tagName,
				Color: *tagColor,
			}
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepo) GetByTag(ctx context.Context, userID int64, tagID int) ([]model.Task, error) {
	query := `
		SELECT t.id, t.user_id, t.title, t.description, t.difficulty, t.tag_id, 
		       t.date, t.done, t.xp_reward, t.created_at,
		       tag.name, tag.color
		FROM task t
		LEFT JOIN tag ON t.tag_id = tag.id
		WHERE t.user_id = $1 AND t.tag_id = $2
		ORDER BY t.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID, tagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []model.Task{}
	for rows.Next() {
		var task model.Task
		var tagName, tagColor *string

		if err := rows.Scan(
			&task.ID, &task.UserID, &task.Title, &task.Description, &task.Difficulty,
			&task.TagID, &task.Date, &task.Done, &task.XPReward, &task.CreatedAt,
			&tagName, &tagColor,
		); err != nil {
			return nil, err
		}

		if tagName != nil && task.TagID != nil {
			task.Tag = &model.Tag{
				ID:    *task.TagID,
				Name:  *tagName,
				Color: *tagColor,
			}
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *TaskRepo) ResetTag(ctx context.Context, userID int64, tagID int) error {
	query := `UPDATE task SET tag_id = NULL WHERE user_id = $1 AND tag_id = $2`
	_, err := r.db.Exec(ctx, query, userID, tagID)
	return err
}
