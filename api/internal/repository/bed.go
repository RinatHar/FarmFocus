package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/RinatHar/FarmFocus/api/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BedRepo struct {
	db *pgxpool.Pool
}

func NewBedRepo(db *pgxpool.Pool) *BedRepo {
	return &BedRepo{db: db}
}

func (r *BedRepo) Create(ctx context.Context, bed *model.Bed) error {
	query := `
		INSERT INTO bed (user_id, cell_number, is_locked, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		bed.UserID, bed.CellNumber, bed.IsLocked, bed.CreatedAt,
	).Scan(&bed.ID)
}

func (r *BedRepo) GetByID(ctx context.Context, id int) (*model.Bed, error) {
	var bed model.Bed
	query := `
		SELECT id, user_id, cell_number, is_locked, created_at
		FROM bed
		WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&bed.ID, &bed.UserID, &bed.CellNumber, &bed.IsLocked, &bed.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("bed with id=%d not found", id)
		}
		return nil, err
	}
	return &bed, nil
}

func (r *BedRepo) GetByUser(ctx context.Context, userID int64) ([]model.Bed, error) {
	query := `
		SELECT id, user_id, cell_number, is_locked, created_at
		FROM bed
		WHERE user_id = $1
		ORDER BY cell_number
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	beds := []model.Bed{}
	for rows.Next() {
		var bed model.Bed
		if err := rows.Scan(
			&bed.ID, &bed.UserID, &bed.CellNumber, &bed.IsLocked, &bed.CreatedAt,
		); err != nil {
			return nil, err
		}
		beds = append(beds, bed)
	}
	return beds, nil
}

func (r *BedRepo) GetByCellNumber(ctx context.Context, userID int64, cellNumber int) (*model.Bed, error) {
	var bed model.Bed
	query := `
		SELECT id, user_id, cell_number, is_locked, created_at
		FROM bed
		WHERE user_id = $1 AND cell_number = $2
	`
	err := r.db.QueryRow(ctx, query, userID, cellNumber).Scan(
		&bed.ID, &bed.UserID, &bed.CellNumber, &bed.IsLocked, &bed.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("bed with cell_number=%d not found", cellNumber)
		}
		return nil, err
	}
	return &bed, nil
}

func (r *BedRepo) UpdateLockStatus(ctx context.Context, id int, isLocked bool) error {
	query := `
		UPDATE bed
		SET is_locked = $1
		WHERE id = $2
		RETURNING id
	`
	var updatedID int
	err := r.db.QueryRow(ctx, query, isLocked, id).Scan(&updatedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("bed with id=%d not found", id)
		}
		return err
	}
	return nil
}

func (r *BedRepo) Unlock(ctx context.Context, id int) error {
	return r.UpdateLockStatus(ctx, id, false)
}

func (r *BedRepo) Lock(ctx context.Context, id int) error {
	return r.UpdateLockStatus(ctx, id, true)
}

func (r *BedRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM bed WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("bed with id=%d not found", id)
	}
	return nil
}

func (r *BedRepo) GetAvailableBeds(ctx context.Context, userID int64) ([]model.Bed, error) {
	query := `
		SELECT id, user_id, cell_number, is_locked, created_at
		FROM bed
		WHERE user_id = $1 AND is_locked = false
		ORDER BY cell_number
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	beds := []model.Bed{}
	for rows.Next() {
		var bed model.Bed
		if err := rows.Scan(
			&bed.ID, &bed.UserID, &bed.CellNumber, &bed.IsLocked, &bed.CreatedAt,
		); err != nil {
			return nil, err
		}
		beds = append(beds, bed)
	}
	return beds, nil
}

func (r *BedRepo) GetWithPlants(ctx context.Context, userID int64) ([]model.BedWithUserPlant, error) {
	query := `
		SELECT b.id, b.user_id, b.cell_number, b.is_locked, b.created_at,
		       up.id as plant_id, up.seed_id, up.current_growth, up.created_at as plant_created_at,
		       s.name as seed_name, s.icon as seed_icon, s.target_growth, s.gold_reward, s.xp_reward
		FROM bed b
		LEFT JOIN user_plant up ON b.id = up.bed_id
		LEFT JOIN seed s ON up.seed_id = s.id
		WHERE b.user_id = $1
		ORDER BY b.cell_number
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	beds := []model.BedWithUserPlant{}
	for rows.Next() {
		var bed model.BedWithUserPlant
		var plantID, seedID, currentGrowth *int
		var plantCreatedAt *string
		var seedName, seedIcon *string
		var targetGrowth, goldReward, xpReward *int

		err := rows.Scan(
			&bed.ID, &bed.UserID, &bed.CellNumber, &bed.IsLocked, &bed.CreatedAt,
			&plantID, &seedID, &currentGrowth, &plantCreatedAt,
			&seedName, &seedIcon, &targetGrowth, &goldReward, &xpReward,
		)
		if err != nil {
			return nil, err
		}

		// Если есть растение, создаем структуру UserPlantWithSeed
		if plantID != nil {
			bed.UserPlant = &model.UserPlantWithSeed{
				UserPlant: model.UserPlant{
					ID:            *plantID,
					UserID:        userID,
					SeedID:        *seedID,
					BedID:         bed.ID,
					CurrentGrowth: *currentGrowth,
				},
				SeedName:      *seedName,
				SeedIcon:      getStringPtr(seedIcon),
				TargetGrowth:  *targetGrowth,
				GoldReward:    *goldReward,
				XPReward:      *xpReward,
				GrowthPercent: utils.CalculateGrowthPercent(*currentGrowth, *targetGrowth),
			}
		}

		beds = append(beds, bed)
	}
	return beds, nil
}

func (r *BedRepo) GetEmptyBeds(ctx context.Context, userID int64) ([]model.Bed, error) {
	query := `
		SELECT b.id, b.user_id, b.cell_number, b.is_locked, b.created_at
		FROM bed b
		LEFT JOIN user_plant up ON b.id = up.bed_id
		WHERE b.user_id = $1 AND up.id IS NULL AND b.is_locked = false
		ORDER BY b.cell_number
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	beds := []model.Bed{}
	for rows.Next() {
		var bed model.Bed
		if err := rows.Scan(
			&bed.ID, &bed.UserID, &bed.CellNumber, &bed.IsLocked, &bed.CreatedAt,
		); err != nil {
			return nil, err
		}
		beds = append(beds, bed)
	}
	return beds, nil
}

func (r *BedRepo) CreateInitialBeds(ctx context.Context, userID int64, count int) error {
	for i := 1; i <= count; i++ {
		// Первые 1 грядка разблокирована, остальные заблокированы
		isLocked := i > 1
		query := `
			INSERT INTO bed (user_id, cell_number, is_locked, created_at)
			VALUES ($1, $2, $3, NOW())
			ON CONFLICT (user_id, cell_number) DO NOTHING
		`
		_, err := r.db.Exec(ctx, query, userID, i, isLocked)
		if err != nil {
			return err
		}
	}
	return nil
}

// Вспомогательные функции
func getStringPtr(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func (r *BedRepo) GetAll(ctx context.Context) ([]model.Bed, error) {
	query := `
		SELECT id, user_id, cell_number, is_locked, created_at
		FROM bed
		ORDER BY cell_number
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	beds := []model.Bed{}
	for rows.Next() {
		var bed model.Bed
		err := rows.Scan(
			&bed.ID,
			&bed.UserID,
			&bed.CellNumber,
			&bed.IsLocked,
			&bed.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		beds = append(beds, bed)
	}
	return beds, nil
}

// UnlockNextBed разблокирует следующую закрытую грядку пользователя (по cellNumber)
func (r *BedRepo) UnlockNextBed(ctx context.Context, userID int64) (*model.Bed, error) {
	// Находим первую заблокированную грядку с минимальным cell_number
	query := `
		SELECT id, user_id, cell_number, is_locked, created_at
		FROM bed
		WHERE user_id = $1 AND is_locked = true
		ORDER BY cell_number
		LIMIT 1
	`

	var bed model.Bed
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&bed.ID, &bed.UserID, &bed.CellNumber, &bed.IsLocked, &bed.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no locked beds available for user %d", userID)
		}
		return nil, err
	}

	// Разблокируем грядку
	if err := r.Unlock(ctx, bed.ID); err != nil {
		return nil, fmt.Errorf("failed to unlock bed %d: %w", bed.ID, err)
	}

	// Получаем обновленную грядку
	unlockedBed, err := r.GetByID(ctx, bed.ID)
	if err != nil {
		return nil, err
	}

	return unlockedBed, nil
}
