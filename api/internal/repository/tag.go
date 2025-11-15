package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TagRepo struct {
	db *pgxpool.Pool
}

func NewTagRepo(db *pgxpool.Pool) *TagRepo {
	return &TagRepo{db: db}
}

func (r *TagRepo) Create(ctx context.Context, tag *model.Tag) error {
	query := `
		INSERT INTO tag (user_id, name, color, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		tag.UserID, tag.Name, tag.Color, tag.CreatedAt,
	).Scan(&tag.ID)
}

func (r *TagRepo) GetByID(ctx context.Context, id int) (*model.Tag, error) {
	var tag model.Tag
	query := `
		SELECT id, user_id, name, color, created_at
		FROM tag
		WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&tag.ID, &tag.UserID, &tag.Name, &tag.Color, &tag.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("tag with id=%d not found", id)
		}
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepo) GetByUser(ctx context.Context, userID int64) ([]model.Tag, error) {
	query := `
		SELECT id, user_id, name, color, created_at
		FROM tag
		WHERE user_id = $1
		ORDER BY name ASC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []model.Tag{}
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(
			&tag.ID, &tag.UserID, &tag.Name, &tag.Color, &tag.CreatedAt,
		); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (r *TagRepo) GetByName(ctx context.Context, name string, userID int64) (*model.Tag, error) {
	var tag model.Tag
	query := `
		SELECT id, user_id, name, color, created_at
		FROM tag
		WHERE name = $1 AND user_id = $2
	`
	err := r.db.QueryRow(ctx, query, name, userID).Scan(
		&tag.ID, &tag.UserID, &tag.Name, &tag.Color, &tag.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("tag with name=%s not found", name)
		}
		return nil, err
	}
	return &tag, nil
}

func (r *TagRepo) Update(ctx context.Context, tag *model.Tag) error {
	query := `
		UPDATE tag
		SET name = $1, color = $2
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, name, color, created_at
	`
	err := r.db.QueryRow(ctx, query,
		tag.Name, tag.Color, tag.ID, tag.UserID,
	).Scan(
		&tag.ID, &tag.UserID, &tag.Name, &tag.Color, &tag.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("tag with id=%d not found", tag.ID)
		}
		return err
	}
	return nil
}

func (r *TagRepo) Delete(ctx context.Context, id int, userID int64) error {
	var taskCount int
	checkQuery := `SELECT COUNT(*) FROM task WHERE tag_id = $1 AND user_id = $2`
	err := r.db.QueryRow(ctx, checkQuery, id, userID).Scan(&taskCount)
	if err != nil {
		return err
	}

	if taskCount > 0 {
		return fmt.Errorf("cannot delete tag: %d task(s) are still using it", taskCount)
	}

	query := `DELETE FROM tag WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("tag with id=%d not found or does not belong to user", id)
	}
	return nil
}

func (r *TagRepo) GetWithTaskCount(ctx context.Context, userID int64) ([]model.TagWithCount, error) {
	query := `
		SELECT t.id, t.user_id, t.name, t.color, t.created_at,
		       COUNT(task.id) as task_count
		FROM tag t
		LEFT JOIN task task ON t.id = task.tag_id AND task.user_id = t.user_id
		WHERE t.user_id = $1
		GROUP BY t.id, t.user_id, t.name, t.color, t.created_at
		ORDER BY t.name ASC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []model.TagWithCount{}
	for rows.Next() {
		var tag model.TagWithCount
		if err := rows.Scan(
			&tag.ID, &tag.UserID, &tag.Name, &tag.Color, &tag.CreatedAt, &tag.TaskCount,
		); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (r *TagRepo) Exists(ctx context.Context, name string, userID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM tag WHERE name = $1 AND user_id = $2)`
	err := r.db.QueryRow(ctx, query, name, userID).Scan(&exists)
	return exists, err
}
