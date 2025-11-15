package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserSeedRepo struct {
	db *pgxpool.Pool
}

func NewUserSeedRepo(db *pgxpool.Pool) *UserSeedRepo {
	return &UserSeedRepo{db: db}
}

func (r *UserSeedRepo) CreateOrUpdate(ctx context.Context, userSeed *model.UserSeed) error {
	query := `
		INSERT INTO user_seed (user_id, seed_id, quantity, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, seed_id) 
		DO UPDATE SET quantity = user_seed.quantity + EXCLUDED.quantity
		RETURNING id, quantity
	`
	return r.db.QueryRow(ctx, query,
		userSeed.UserID, userSeed.SeedID, userSeed.Quantity, userSeed.CreatedAt,
	).Scan(&userSeed.ID, &userSeed.Quantity)
}

func (r *UserSeedRepo) GetByUserAndSeed(ctx context.Context, userID int64, seedID int) (*model.UserSeed, error) {
	var userSeed model.UserSeed
	query := `
		SELECT id, user_id, seed_id, quantity, created_at
		FROM user_seed
		WHERE user_id = $1 AND seed_id = $2
	`
	err := r.db.QueryRow(ctx, query, userID, seedID).Scan(
		&userSeed.ID, &userSeed.UserID, &userSeed.SeedID, &userSeed.Quantity, &userSeed.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &userSeed, nil
}

func (r *UserSeedRepo) GetByUser(ctx context.Context, userID int64) ([]model.UserSeed, error) {
	query := `
		SELECT id, user_id, seed_id, quantity, created_at
		FROM user_seed
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userSeeds := []model.UserSeed{}
	for rows.Next() {
		var userSeed model.UserSeed
		if err := rows.Scan(
			&userSeed.ID, &userSeed.UserID, &userSeed.SeedID, &userSeed.Quantity, &userSeed.CreatedAt,
		); err != nil {
			return nil, err
		}
		userSeeds = append(userSeeds, userSeed)
	}
	return userSeeds, nil
}

func (r *UserSeedRepo) UpdateQuantity(ctx context.Context, userID int64, seedID int, quantity int64) error {
	query := `
		UPDATE user_seed
		SET quantity = $1
		WHERE user_id = $2 AND seed_id = $3
		RETURNING id
	`
	var id int
	err := r.db.QueryRow(ctx, query, quantity, userID, seedID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user seed not found for user_id=%d and seed_id=%d", userID, seedID)
		}
		return err
	}
	return nil
}

func (r *UserSeedRepo) AddQuantity(ctx context.Context, userID int64, seedID int, amount int64) error {
	query := `
		UPDATE user_seed
		SET quantity = quantity + $1
		WHERE user_id = $2 AND seed_id = $3
		RETURNING id, quantity
	`
	var id int
	var newQuantity int64
	err := r.db.QueryRow(ctx, query, amount, userID, seedID).Scan(&id, &newQuantity)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user seed not found for user_id=%d and seed_id=%d", userID, seedID)
		}
		return err
	}
	return nil
}

func (r *UserSeedRepo) SubtractQuantity(ctx context.Context, userID int64, seedID int, amount int64) (bool, error) {
	query := `
		UPDATE user_seed
		SET quantity = quantity - $1
		WHERE user_id = $2 AND seed_id = $3 AND quantity >= $1
		RETURNING id, quantity
	`
	var id int
	var newQuantity int64
	err := r.db.QueryRow(ctx, query, amount, userID, seedID).Scan(&id, &newQuantity)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *UserSeedRepo) Delete(ctx context.Context, userID int64, seedID int) error {
	query := `DELETE FROM user_seed WHERE user_id = $1 AND seed_id = $2`
	result, err := r.db.Exec(ctx, query, userID, seedID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user seed not found for user_id=%d and seed_id=%d", userID, seedID)
	}
	return nil
}

func (r *UserSeedRepo) GetUserSeedsWithDetails(ctx context.Context, userID int64) ([]model.UserSeedWithDetails, error) {
	query := `
		SELECT us.id, us.user_id, us.seed_id, us.quantity, us.created_at,
		       s.name as seed_name, s.icon as seed_icon, s.target_growth as seed_target_growth, 
		       s.rarity as seed_rarity, s.img_plant as seed_imgPlant
		FROM user_seed us
		INNER JOIN seed s ON us.seed_id = s.id
		WHERE us.user_id = $1
		ORDER BY us.quantity DESC, s.name
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userSeeds := []model.UserSeedWithDetails{}
	for rows.Next() {
		var userSeed model.UserSeedWithDetails
		if err := rows.Scan(
			&userSeed.ID, &userSeed.UserID, &userSeed.SeedID, &userSeed.Quantity, &userSeed.CreatedAt,
			&userSeed.SeedName, &userSeed.SeedIcon, &userSeed.TargetGrowth, &userSeed.Rarity, &userSeed.ImgPlant,
		); err != nil {
			return nil, err
		}
		userSeeds = append(userSeeds, userSeed)
	}
	return userSeeds, nil
}

func (r *UserSeedRepo) GetAvailableSeedsForUser(ctx context.Context, userID int64, userLevel int) ([]model.SeedWithUserData, error) {
	if userLevel <= 0 {
		userLevel = 1
	}

	query := `
        SELECT s.id, s.name, s.icon, s.level_required, s.target_growth, s.rarity, 
               s.modification, s.gold_reward, s.xp_reward, s.created_at,
               COALESCE(us.quantity, 0) as user_quantity
        FROM seed s
        LEFT JOIN user_seed us ON s.id = us.seed_id AND us.user_id = $1
        WHERE s.level_required <= $2
        ORDER BY s.level_required, s.rarity DESC, s.name
    `
	rows, err := r.db.Query(ctx, query, userID, userLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seeds := []model.SeedWithUserData{}
	for rows.Next() {
		var seed model.SeedWithUserData
		var userQuantity int64
		if err := rows.Scan(
			&seed.ID, &seed.Name, &seed.Icon, &seed.LevelRequired, &seed.TargetGrowth,
			&seed.Rarity, &seed.Modification, &seed.GoldReward, &seed.XPReward, &seed.CreatedAt,
			&userQuantity,
		); err != nil {
			return nil, err
		}
		seed.UserQuantity = userQuantity
		seed.IsOwned = userQuantity > 0
		seeds = append(seeds, seed)
	}
	return seeds, nil
}

func (r *UserSeedRepo) GetTotalSeedCount(ctx context.Context, userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM user_seed WHERE user_id = $1 AND quantity > 0`
	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

// AddOrUpdateQuantity добавляет или обновляет количество семян у пользователя
func (r *UserSeedRepo) AddOrUpdateQuantity(ctx context.Context, userID int64, seedID int, quantity int) error {
	query := `
		INSERT INTO user_seed (user_id, seed_id, quantity, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (user_id, seed_id) 
		DO UPDATE SET quantity = user_seed.quantity + EXCLUDED.quantity
	`
	_, err := r.db.Exec(ctx, query, userID, seedID, quantity)
	return err
}
