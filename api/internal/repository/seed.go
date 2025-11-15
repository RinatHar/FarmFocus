package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeedRepo struct {
	db *pgxpool.Pool
}

func NewSeedRepo(db *pgxpool.Pool) *SeedRepo {
	return &SeedRepo{db: db}
}

func (r *SeedRepo) Create(ctx context.Context, seed *model.Seed) error {
	query := `
		INSERT INTO seed (name, icon, img_plant, level_required, target_growth, rarity, modification, gold_reward, xp_reward, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		seed.Name, seed.Icon, seed.ImgPlant, seed.LevelRequired, seed.TargetGrowth, seed.Rarity,
		seed.Modification, seed.GoldReward, seed.XPReward, seed.CreatedAt,
	).Scan(&seed.ID)
}

func (r *SeedRepo) GetByID(ctx context.Context, id int) (*model.Seed, error) {
	var seed model.Seed
	query := `
		SELECT id, name, icon, img_plant, level_required, target_growth, rarity, modification, gold_reward, xp_reward, created_at
		FROM seed
		WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&seed.ID, &seed.Name, &seed.Icon, &seed.ImgPlant, &seed.LevelRequired, &seed.TargetGrowth,
		&seed.Rarity, &seed.Modification, &seed.GoldReward, &seed.XPReward, &seed.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("seed with id=%d not found", id)
		}
		return nil, err
	}
	return &seed, nil
}

func (r *SeedRepo) GetAll(ctx context.Context) ([]model.Seed, error) {
	query := `
		SELECT id, name, icon, img_plant, level_required, target_growth, rarity, modification, gold_reward, xp_reward, created_at
		FROM seed
		ORDER BY level_required, name
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seeds := []model.Seed{}
	for rows.Next() {
		var seed model.Seed
		if err := rows.Scan(
			&seed.ID, &seed.Name, &seed.Icon, &seed.ImgPlant, &seed.LevelRequired, &seed.TargetGrowth,
			&seed.Rarity, &seed.Modification, &seed.GoldReward, &seed.XPReward, &seed.CreatedAt,
		); err != nil {
			return nil, err
		}
		seeds = append(seeds, seed)
	}
	return seeds, nil
}

func (r *SeedRepo) GetByLevel(ctx context.Context, level int) ([]model.Seed, error) {
	query := `
		SELECT id, name, icon, img_plant, level_required, target_growth, rarity, modification, gold_reward, xp_reward, created_at
		FROM seed
		WHERE level_required <= $1
		ORDER BY level_required, rarity DESC, name
	`
	rows, err := r.db.Query(ctx, query, level)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seeds := []model.Seed{}
	for rows.Next() {
		var seed model.Seed
		if err := rows.Scan(
			&seed.ID, &seed.Name, &seed.Icon, &seed.ImgPlant, &seed.LevelRequired, &seed.TargetGrowth,
			&seed.Rarity, &seed.Modification, &seed.GoldReward, &seed.XPReward, &seed.CreatedAt,
		); err != nil {
			return nil, err
		}
		seeds = append(seeds, seed)
	}
	return seeds, nil
}

func (r *SeedRepo) GetByRarity(ctx context.Context, rarity string) ([]model.Seed, error) {
	query := `
		SELECT id, name, icon, img_plant, level_required, target_growth, rarity, modification, gold_reward, xp_reward, created_at
		FROM seed
		WHERE rarity = $1
		ORDER BY level_required, name
	`
	rows, err := r.db.Query(ctx, query, rarity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seeds := []model.Seed{}
	for rows.Next() {
		var seed model.Seed
		if err := rows.Scan(
			&seed.ID, &seed.Name, &seed.Icon, &seed.ImgPlant, &seed.LevelRequired, &seed.TargetGrowth,
			&seed.Rarity, &seed.Modification, &seed.GoldReward, &seed.XPReward, &seed.CreatedAt,
		); err != nil {
			return nil, err
		}
		seeds = append(seeds, seed)
	}
	return seeds, nil
}

func (r *SeedRepo) Update(ctx context.Context, seed *model.Seed) error {
	query := `
		UPDATE seed
		SET name = $1, icon = $2, img_plant = $3, level_required = $4, target_growth = $5, rarity = $6,
			modification = $7, gold_reward = $8, xp_reward = $9
		WHERE id = $10
		RETURNING id, name, icon, img_plant, level_required, target_growth, rarity, modification, gold_reward, xp_reward, created_at
	`
	err := r.db.QueryRow(ctx, query,
		seed.Name, seed.Icon, seed.ImgPlant, seed.LevelRequired, seed.TargetGrowth, seed.Rarity,
		seed.Modification, seed.GoldReward, seed.XPReward, seed.ID,
	).Scan(
		&seed.ID, &seed.Name, &seed.Icon, &seed.ImgPlant, &seed.LevelRequired, &seed.TargetGrowth,
		&seed.Rarity, &seed.Modification, &seed.GoldReward, &seed.XPReward, &seed.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("seed with id=%d not found", seed.ID)
		}
		return err
	}
	return nil
}

func (r *SeedRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM seed WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("seed with id=%d not found", id)
	}

	return nil
}
