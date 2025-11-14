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

type HabitRepo struct {
	db *pgxpool.Pool
}

func NewHabitRepo(db *pgxpool.Pool) *HabitRepo {
	return &HabitRepo{db: db}
}

func (r *HabitRepo) HasCompletedHabitToday(ctx context.Context, userID int64) (bool, error) {
	today := time.Now().Format("2006-01-02")
	query := `
        SELECT EXISTS(
            SELECT 1 FROM progress_log 
            WHERE user_id = $1 AND DATE(created_at) = $2 AND habit_id IS NOT NULL
        )`
	var exists bool
	err := r.db.QueryRow(ctx, query, userID, today).Scan(&exists)
	return exists, err
}

func (r *HabitRepo) Create(ctx context.Context, habit *model.Habit) error {
	query := `
		INSERT INTO habit (user_id, title, description, difficulty, tag_id, 
		                  done, count, period, every, start_date, xp_reward, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		habit.UserID, habit.Title, habit.Description, habit.Difficulty, habit.TagID,
		habit.Done, habit.Count, habit.Period, habit.Every, habit.StartDate,
		habit.XPReward, habit.CreatedAt,
	).Scan(&habit.ID)
}

func (r *HabitRepo) GetByID(ctx context.Context, id int, userID int64) (*model.Habit, error) {
	var habit model.Habit
	var tagName, tagColor *string

	query := `
		SELECT h.id, h.user_id, h.title, h.description, h.difficulty, h.tag_id, 
		       h.done, h.count, h.period, h.every, h.start_date, h.xp_reward, h.created_at,
		       tag.name, tag.color
		FROM habit h
		LEFT JOIN tag ON h.tag_id = tag.id
		WHERE h.id = $1 AND h.user_id = $2
	`

	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&habit.ID, &habit.UserID, &habit.Title, &habit.Description, &habit.Difficulty,
		&habit.TagID, &habit.Done, &habit.Count, &habit.Period, &habit.Every,
		&habit.StartDate, &habit.XPReward, &habit.CreatedAt,
		&tagName, &tagColor,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("habit with id=%d not found", id)
		}
		return nil, err
	}

	// Если есть тег, заполняем структуру Tag
	if tagName != nil && habit.TagID != nil {
		habit.Tag = &model.Tag{
			ID:    *habit.TagID,
			Name:  *tagName,
			Color: *tagColor,
		}
	}

	return &habit, nil
}

func (r *HabitRepo) GetAll(ctx context.Context, userID int64) ([]model.Habit, error) {
	query := `
		SELECT h.id, h.user_id, h.title, h.description, h.difficulty, h.tag_id, 
		       h.done, h.count, h.period, h.every, h.start_date, h.xp_reward, h.created_at,
		       tag.name, tag.color
		FROM habit h
		LEFT JOIN tag ON h.tag_id = tag.id
		WHERE h.user_id = $1
		ORDER BY h.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habits := []model.Habit{}
	for rows.Next() {
		var habit model.Habit
		var tagName, tagColor *string

		if err := rows.Scan(
			&habit.ID, &habit.UserID, &habit.Title, &habit.Description, &habit.Difficulty,
			&habit.TagID, &habit.Done, &habit.Count, &habit.Period, &habit.Every,
			&habit.StartDate, &habit.XPReward, &habit.CreatedAt,
			&tagName, &tagColor,
		); err != nil {
			return nil, err
		}

		// Если есть тег, заполняем структуру Tag
		if tagName != nil && habit.TagID != nil {
			habit.Tag = &model.Tag{
				ID:    *habit.TagID,
				Name:  *tagName,
				Color: *tagColor,
			}
		}

		habits = append(habits, habit)
	}
	return habits, nil
}

func (r *HabitRepo) Update(ctx context.Context, habit *model.Habit) error {
	query := `
		UPDATE habit
		SET title = $1, description = $2, difficulty = $3, tag_id = $4, 
		    done = $5, count = $6, period = $7, every = $8, start_date = $9, xp_reward = $10
		WHERE id = $11 AND user_id = $12
	`
	result, err := r.db.Exec(ctx, query,
		habit.Title, habit.Description, habit.Difficulty, habit.TagID,
		habit.Done, habit.Count, habit.Period, habit.Every, habit.StartDate,
		habit.XPReward, habit.ID, habit.UserID,
	)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("habit with id=%d not found or does not belong to user", habit.ID)
	}
	return nil
}

func (r *HabitRepo) Delete(ctx context.Context, id int, userID int64) error {
	query := `DELETE FROM habit WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("habit with id=%d not found or does not belong to user", id)
	}
	return nil
}

func (r *HabitRepo) GetByStatus(ctx context.Context, userID int64, done bool) ([]model.Habit, error) {
	query := `
		SELECT h.id, h.user_id, h.title, h.description, h.difficulty, h.tag_id, 
		       h.done, h.count, h.period, h.every, h.start_date, h.xp_reward, h.created_at,
		       tag.name, tag.color
		FROM habit h
		LEFT JOIN tag ON h.tag_id = tag.id
		WHERE h.user_id = $1 AND h.done = $2
		ORDER BY h.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID, done)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habits := []model.Habit{}
	for rows.Next() {
		var habit model.Habit
		var tagName, tagColor *string

		if err := rows.Scan(
			&habit.ID, &habit.UserID, &habit.Title, &habit.Description, &habit.Difficulty,
			&habit.TagID, &habit.Done, &habit.Count, &habit.Period, &habit.Every,
			&habit.StartDate, &habit.XPReward, &habit.CreatedAt,
			&tagName, &tagColor,
		); err != nil {
			return nil, err
		}

		if tagName != nil && habit.TagID != nil {
			habit.Tag = &model.Tag{
				ID:    *habit.TagID,
				Name:  *tagName,
				Color: *tagColor,
			}
		}

		habits = append(habits, habit)
	}
	return habits, nil
}

func (r *HabitRepo) MarkAsDone(ctx context.Context, id int, userID int64) error {
	query := `UPDATE habit SET done = true WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("habit with id=%d not found or does not belong to user", id)
	}
	return nil
}

func (r *HabitRepo) MarkAsUndone(ctx context.Context, id int, userID int64) error {
	query := `UPDATE habit SET done = false WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("habit with id=%d not found or does not belong to user", id)
	}
	return nil
}

func (r *HabitRepo) IncrementCount(ctx context.Context, id int, userID int64) error {
	query := `UPDATE habit SET count = count + 1 WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("habit with id=%d not found or does not belong to user", id)
	}
	return nil
}

func (r *HabitRepo) ResetCount(ctx context.Context, id int, userID int64) error {
	query := `UPDATE habit SET count = 0 WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("habit with id=%d not found or does not belong to user", id)
	}
	return nil
}

func (r *HabitRepo) GetByPeriod(ctx context.Context, userID int64, period string) ([]model.Habit, error) {
	query := `
		SELECT h.id, h.user_id, h.title, h.description, h.difficulty, h.tag_id, 
		       h.done, h.count, h.period, h.every, h.start_date, h.xp_reward, h.created_at,
		       tag.name, tag.color
		FROM habit h
		LEFT JOIN tag ON h.tag_id = tag.id
		WHERE h.user_id = $1 AND h.period = $2
		ORDER BY h.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID, period)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habits := []model.Habit{}
	for rows.Next() {
		var habit model.Habit
		var tagName, tagColor *string

		if err := rows.Scan(
			&habit.ID, &habit.UserID, &habit.Title, &habit.Description, &habit.Difficulty,
			&habit.TagID, &habit.Done, &habit.Count, &habit.Period, &habit.Every,
			&habit.StartDate, &habit.XPReward, &habit.CreatedAt,
			&tagName, &tagColor,
		); err != nil {
			return nil, err
		}

		if tagName != nil && habit.TagID != nil {
			habit.Tag = &model.Tag{
				ID:    *habit.TagID,
				Name:  *tagName,
				Color: *tagColor,
			}
		}

		habits = append(habits, habit)
	}
	return habits, nil
}

func (r *HabitRepo) GetByTag(ctx context.Context, userID int64, tagID int) ([]model.Habit, error) {
	query := `
		SELECT h.id, h.user_id, h.title, h.description, h.difficulty, h.tag_id, 
		       h.done, h.count, h.period, h.every, h.start_date, h.xp_reward, h.created_at,
		       tag.name, tag.color
		FROM habit h
		LEFT JOIN tag ON h.tag_id = tag.id
		WHERE h.user_id = $1 AND h.tag_id = $2
		ORDER BY h.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID, tagID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habits := []model.Habit{}
	for rows.Next() {
		var habit model.Habit
		var tagName, tagColor *string

		if err := rows.Scan(
			&habit.ID, &habit.UserID, &habit.Title, &habit.Description, &habit.Difficulty,
			&habit.TagID, &habit.Done, &habit.Count, &habit.Period, &habit.Every,
			&habit.StartDate, &habit.XPReward, &habit.CreatedAt,
			&tagName, &tagColor,
		); err != nil {
			return nil, err
		}

		if tagName != nil && habit.TagID != nil {
			habit.Tag = &model.Tag{
				ID:    *habit.TagID,
				Name:  *tagName,
				Color: *tagColor,
			}
		}

		habits = append(habits, habit)
	}
	return habits, nil
}

func (r *HabitRepo) GetAllUsersWithHabits(ctx context.Context) ([]model.User, error) {
	query := `
		SELECT DISTINCT u.id, u.max_id, u.username, u.created_at, u.last_login, u.is_active 
		FROM user_info u 
		JOIN habit h ON u.id = h.user_id 
		WHERE u.is_active = true
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.MaxID,
			&user.Username,
			&user.CreatedAt,
			&user.LastLogin,
			&user.IsActive,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// MarkAsDoneAndIncrementCount помечает привычку как выполненную и увеличивает счетчик
func (r *HabitRepo) MarkAsDoneAndIncrementCount(ctx context.Context, id int, userID int64) error {
	query := `UPDATE habit SET done = true, count = count + 1 WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("habit with id=%d not found or does not belong to user", id)
	}
	return nil
}

// MarkAsUndoneAndDecrementCount помечает привычку как невыполненную и уменьшает счетчик
func (r *HabitRepo) MarkAsUndoneAndDecrementCount(ctx context.Context, id int, userID int64) error {
	query := `UPDATE habit SET done = false, count = GREATEST(0, count - 1) WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("habit with id=%d not found or does not belong to user", id)
	}
	return nil
}

func (r *HabitRepo) ResetTag(ctx context.Context, userID int64, tagID int) error {
	query := `UPDATE habit SET tag_id = NULL WHERE user_id = $1 AND tag_id = $2`
	_, err := r.db.Exec(ctx, query, userID, tagID)
	return err
}
