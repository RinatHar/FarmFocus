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

type UserPlantRepo struct {
	db *pgxpool.Pool
}

func NewUserPlantRepo(db *pgxpool.Pool) *UserPlantRepo {
	return &UserPlantRepo{db: db}
}

func (r *UserPlantRepo) MarkPlantsAsWithered(ctx context.Context, userID int64) error {
	query := `UPDATE user_plants SET is_withered = true WHERE user_id = $1 AND is_withered = false`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *UserPlantRepo) RemoveWitheredPlants(ctx context.Context, userID int64) error {
	query := `DELETE FROM user_plants WHERE user_id = $1 AND is_withered = true`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *UserPlantRepo) ResetWitheredStatus(ctx context.Context, userID int64) error {
	query := `UPDATE user_plants SET is_withered = false WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}

func (r *UserPlantRepo) Create(ctx context.Context, plant *model.UserPlant) error {
	query := `
		INSERT INTO user_plant (user_id, seed_id, bed_id, current_growth, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		plant.UserID, plant.SeedID, plant.BedID, plant.CurrentGrowth, plant.CreatedAt,
	).Scan(&plant.ID)
}

func (r *UserPlantRepo) GetByID(ctx context.Context, id int) (*model.UserPlant, error) {
	var plant model.UserPlant
	query := `
		SELECT id, user_id, seed_id, bed_id, current_growth, created_at
		FROM user_plant
		WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&plant.ID, &plant.UserID, &plant.SeedID, &plant.BedID, &plant.CurrentGrowth, &plant.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user plant with id=%d not found", id)
		}
		return nil, err
	}
	return &plant, nil
}

func (r *UserPlantRepo) GetByUser(ctx context.Context, userID int64) ([]model.UserPlant, error) {
	query := `
		SELECT id, user_id, seed_id, bed_id, current_growth, created_at
		FROM user_plant
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plants := []model.UserPlant{}
	for rows.Next() {
		var plant model.UserPlant
		if err := rows.Scan(
			&plant.ID, &plant.UserID, &plant.SeedID, &plant.BedID, &plant.CurrentGrowth, &plant.CreatedAt,
		); err != nil {
			return nil, err
		}
		plants = append(plants, plant)
	}
	return plants, nil
}

func (r *UserPlantRepo) GetByBed(ctx context.Context, bedID int) (*model.UserPlant, error) {
	var plant model.UserPlant
	query := `
		SELECT id, user_id, seed_id, bed_id, current_growth, created_at
		FROM user_plant
		WHERE bed_id = $1
	`
	err := r.db.QueryRow(ctx, query, bedID).Scan(
		&plant.ID, &plant.UserID, &plant.SeedID, &plant.BedID, &plant.CurrentGrowth, &plant.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &plant, nil
}

func (r *UserPlantRepo) UpdateGrowth(ctx context.Context, id int, growth int) error {
	query := `
		UPDATE user_plant
		SET current_growth = $1
		WHERE id = $2
		RETURNING id
	`
	var updatedID int
	err := r.db.QueryRow(ctx, query, growth, id).Scan(&updatedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("user plant with id=%d not found", id)
		}
		return err
	}
	return nil
}

func (r *UserPlantRepo) AddGrowth(ctx context.Context, id int, amount int) (int, error) {
	query := `
		UPDATE user_plant
		SET current_growth = current_growth + $1
		WHERE id = $2
		RETURNING id, current_growth
	`
	var updatedID, newGrowth int
	err := r.db.QueryRow(ctx, query, amount, id).Scan(&updatedID, &newGrowth)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("user plant with id=%d not found", id)
		}
		return 0, err
	}
	return newGrowth, nil
}

func (r *UserPlantRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM user_plant WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("user plant with id=%d not found", id)
	}
	return nil
}

func (r *UserPlantRepo) DeleteByBed(ctx context.Context, bedID int) error {
	query := `DELETE FROM user_plant WHERE bed_id = $1`
	result, err := r.db.Exec(ctx, query, bedID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no plant found in bed_id=%d", bedID)
	}
	return nil
}

func (r *UserPlantRepo) GetWithSeedDetails(ctx context.Context, userID int64) ([]model.UserPlantWithSeed, error) {
	query := `
		SELECT up.id, up.user_id, up.seed_id, up.bed_id, up.current_growth, up.is_withered, up.created_at,
		       s.name as seed_name, s.icon as seed_icon, s.img_plant as seed_img_plant, s.target_growth, s.gold_reward, s.xp_reward
		FROM user_plant up
		INNER JOIN seed s ON up.seed_id = s.id
		WHERE up.user_id = $1
		ORDER BY up.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plants := []model.UserPlantWithSeed{}
	for rows.Next() {
		var plant model.UserPlantWithSeed
		if err := rows.Scan(
			&plant.ID, &plant.UserID, &plant.SeedID, &plant.BedID, &plant.CurrentGrowth, &plant.IsWithered, &plant.CreatedAt,
			&plant.SeedName, &plant.SeedIcon, &plant.SeedImgPlant, &plant.TargetGrowth, &plant.GoldReward, &plant.XPReward,
		); err != nil {
			return nil, err
		}
		plant.GrowthPercent = utils.CalculateGrowthPercent(plant.CurrentGrowth, plant.TargetGrowth)
		plants = append(plants, plant)
	}
	return plants, nil
}

func (r *UserPlantRepo) GetReadyForHarvest(ctx context.Context, userID int64) ([]model.UserPlantWithSeed, error) {
	query := `
		SELECT up.id, up.user_id, up.seed_id, up.bed_id, up.current_growth, up.created_at,
		       s.name as seed_name, s.icon as seed_icon, s.target_growth, s.gold_reward, s.xp_reward
		FROM user_plant up
		INNER JOIN seed s ON up.seed_id = s.id
		WHERE up.user_id = $1 AND up.current_growth >= s.target_growth
		ORDER BY up.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plants := []model.UserPlantWithSeed{}
	for rows.Next() {
		var plant model.UserPlantWithSeed
		if err := rows.Scan(
			&plant.ID, &plant.UserID, &plant.SeedID, &plant.BedID, &plant.CurrentGrowth, &plant.CreatedAt,
			&plant.SeedName, &plant.SeedIcon, &plant.TargetGrowth, &plant.GoldReward, &plant.XPReward,
		); err != nil {
			return nil, err
		}
		plant.GrowthPercent = 100 // готовы к сбору
		plants = append(plants, plant)
	}
	return plants, nil
}

func (r *UserPlantRepo) GetGrowingPlants(ctx context.Context, userID int64) ([]model.UserPlantWithSeed, error) {
	query := `
		SELECT up.id, up.user_id, up.seed_id, up.bed_id, up.current_growth, up.created_at,
		       s.name as seed_name, s.icon as seed_icon, s.target_growth, s.gold_reward, s.xp_reward
		FROM user_plant up
		INNER JOIN seed s ON up.seed_id = s.id
		WHERE up.user_id = $1 AND up.current_growth < s.target_growth
		ORDER BY up.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plants := []model.UserPlantWithSeed{}
	for rows.Next() {
		var plant model.UserPlantWithSeed
		if err := rows.Scan(
			&plant.ID, &plant.UserID, &plant.SeedID, &plant.BedID, &plant.CurrentGrowth, &plant.CreatedAt,
			&plant.SeedName, &plant.SeedIcon, &plant.TargetGrowth, &plant.GoldReward, &plant.XPReward,
		); err != nil {
			return nil, err
		}
		plant.GrowthPercent = utils.CalculateGrowthPercent(plant.CurrentGrowth, plant.TargetGrowth)
		plants = append(plants, plant)
	}
	return plants, nil
}

func (r *UserPlantRepo) MarkAsWithered(ctx context.Context, plantID int) error {
	query := `UPDATE user_plant SET is_withered = true WHERE id = $1`

	result, err := r.db.Exec(ctx, query, plantID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("plant with id=%d not found", plantID)
	}

	return nil
}
