package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GoodRepo struct {
	db *pgxpool.Pool
}

func NewGoodRepo(db *pgxpool.Pool) *GoodRepo {
	return &GoodRepo{db: db}
}

func (r *GoodRepo) Create(ctx context.Context, good *model.Good) error {
	query := `
        INSERT INTO good (user_id, type, id_good, quantity, cost, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `

	now := time.Now()
	err := r.db.QueryRow(ctx, query,
		good.UserID,
		good.Type,
		good.IDGood,
		good.Quantity,
		good.Cost,
		now,
		now,
	).Scan(&good.ID)

	return err
}

func (r *GoodRepo) GetByUser(ctx context.Context, userID int64) ([]model.Good, error) {
	query := `
        SELECT id, user_id, type, id_good, quantity, cost, created_at, updated_at
        FROM good
        WHERE user_id = $1
        ORDER BY type, id_good
    `

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	goods := []model.Good{}
	for rows.Next() {
		var good model.Good
		err := rows.Scan(
			&good.ID,
			&good.UserID,
			&good.Type,
			&good.IDGood,
			&good.Quantity,
			&good.Cost,
			&good.CreatedAt,
			&good.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		goods = append(goods, good)
	}
	return goods, nil
}

func (r *GoodRepo) GetByUserAndType(ctx context.Context, userID int64, goodType string) ([]model.Good, error) {
	query := `
        SELECT id, user_id, type, id_good, quantity, cost, created_at, updated_at
        FROM good
        WHERE user_id = $1 AND type = $2
        ORDER BY id_good
    `

	rows, err := r.db.Query(ctx, query, userID, goodType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	goods := []model.Good{}
	for rows.Next() {
		var good model.Good
		err := rows.Scan(
			&good.ID,
			&good.UserID,
			&good.Type,
			&good.IDGood,
			&good.Quantity,
			&good.Cost,
			&good.CreatedAt,
			&good.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		goods = append(goods, good)
	}
	return goods, nil
}

func (r *GoodRepo) GetByID(ctx context.Context, id int) (*model.Good, error) {
	query := `
        SELECT id, user_id, type, id_good, quantity, cost, created_at, updated_at
        FROM good
        WHERE id = $1
    `

	var good model.Good
	err := r.db.QueryRow(ctx, query, id).Scan(
		&good.ID,
		&good.UserID,
		&good.Type,
		&good.IDGood,
		&good.Quantity,
		&good.Cost,
		&good.CreatedAt,
		&good.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("good with id=%d not found", id)
		}
		return nil, err
	}
	return &good, nil
}

func (r *GoodRepo) UpdateQuantity(ctx context.Context, id int, quantity int) error {
	query := `
        UPDATE good 
        SET quantity = $1, updated_at = $2
        WHERE id = $3
    `

	result, err := r.db.Exec(ctx, query, quantity, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("good with id=%d not found", id)
	}

	return nil
}

func (r *GoodRepo) UpdateCost(ctx context.Context, id int, cost int) error {
	query := `
        UPDATE good 
        SET cost = $1, updated_at = $2
        WHERE id = $3
    `

	result, err := r.db.Exec(ctx, query, cost, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("good with id=%d not found", id)
	}

	return nil
}

func (r *GoodRepo) AddQuantity(ctx context.Context, id int, amount int) error {
	query := `
        UPDATE good 
        SET quantity = quantity + $1, updated_at = $2
        WHERE id = $3
    `

	result, err := r.db.Exec(ctx, query, amount, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("good with id=%d not found", id)
	}

	return nil
}

func (r *GoodRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM good WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("good with id=%d not found", id)
	}

	return nil
}

func (r *GoodRepo) GetByUserTypeAndIDGood(ctx context.Context, userID int64, goodType string, idGood int) (*model.Good, error) {
	query := `
        SELECT id, user_id, type, id_good, quantity, cost, created_at, updated_at
        FROM good
        WHERE user_id = $1 AND type = $2 AND id_good = $3
    `

	var good model.Good
	err := r.db.QueryRow(ctx, query, userID, goodType, idGood).Scan(
		&good.ID,
		&good.UserID,
		&good.Type,
		&good.IDGood,
		&good.Quantity,
		&good.Cost,
		&good.CreatedAt,
		&good.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Not found, but not an error
		}
		return nil, err
	}
	return &good, nil
}

// CreateOrUpdate создает или обновляет товар
func (r *GoodRepo) CreateOrUpdate(ctx context.Context, good *model.Good) error {
	existing, err := r.GetByUserTypeAndIDGood(ctx, good.UserID, good.Type, good.IDGood)
	if err != nil {
		return err
	}

	if existing == nil {
		// Create new
		return r.Create(ctx, good)
	} else {
		// Update existing
		good.ID = existing.ID
		query := `
            UPDATE good 
            SET quantity = $1, cost = $2, updated_at = $3
            WHERE id = $4
        `
		_, err := r.db.Exec(ctx, query, good.Quantity, good.Cost, time.Now(), good.ID)
		return err
	}
}

func (r *GoodRepo) DeleteByUser(ctx context.Context, userID int64) error {
	query := `DELETE FROM good WHERE user_id = $1`

	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	fmt.Printf("Deleted %d goods for user %d\n", rowsAffected, userID)

	return nil
}
