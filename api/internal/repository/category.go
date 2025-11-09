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

type CategoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{db: db}
}

// Create
func (r *CategoryRepo) Create(ctx context.Context, c *model.Category) error {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	return r.db.QueryRow(ctx,
		"INSERT INTO categories (user_id, name, color, created_at, created_by, updated_at, updated_by) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id",
		c.UserID, c.Name, c.Color, c.CreatedAt, c.CreatedBy, c.UpdatedAt, c.UpdatedBy).Scan(&c.ID)
}

// GetAll
func (r *CategoryRepo) GetAll(ctx context.Context, userID int) ([]model.Category, error) {
	rows, err := r.db.Query(ctx, "SELECT id, user_id, name, color, created_at FROM categories WHERE user_id=$1 AND deleted_at IS NULL", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cats := []model.Category{}
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Color, &c.CreatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, nil
}

// Update
func (r *CategoryRepo) Update(ctx context.Context, c *model.Category) error {
	c.UpdatedAt = time.Now()

	query := `
        UPDATE categories
        SET name=$1, color=$2, updated_at=$3, updated_by=$4
        WHERE id=$5 AND user_id=$6 AND deleted_at IS NULL
        RETURNING id, user_id, name, color, created_at, created_by, updated_at, updated_by
    `

	err := r.db.QueryRow(ctx, query,
		c.Name, c.Color, c.UpdatedAt, c.UpdatedBy, c.ID, c.UserID,
	).Scan(
		&c.ID, &c.UserID, &c.Name, &c.Color, &c.CreatedAt, &c.CreatedBy, &c.UpdatedAt, &c.UpdatedBy,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("category with id=%d not found or deleted", c.ID)
		}
		return err
	}

	return nil
}

// Delete
func (r *CategoryRepo) Delete(ctx context.Context, id, userID, deletedBy int) error {
	deletedAt := time.Now()

	query := `
        UPDATE categories
        SET deleted_at=$1, deleted_by=$2
        WHERE id=$3 AND user_id=$4 AND deleted_at IS NULL
        RETURNING id
    `

	var deletedID int
	err := r.db.QueryRow(ctx, query, deletedAt, deletedBy, id, userID).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("category with id=%d not found, already deleted, or does not belong to user %d", id, userID)
		}
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}
