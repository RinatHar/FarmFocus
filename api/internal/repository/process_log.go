package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProgressLogRepo struct {
	db *pgxpool.Pool
}

func NewProgressLogRepo(db *pgxpool.Pool) *ProgressLogRepo {
	return &ProgressLogRepo{db: db}
}

func (r *ProgressLogRepo) Create(ctx context.Context, log *model.ProgressLog) error {
	query := `
		INSERT INTO progress_log (user_id, task_id, habit_id, xp_earned, gold_earned, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRow(ctx, query,
		log.UserID,
		log.TaskID,
		log.HabitID,
		log.XPEarned,
		log.GoldEarned,
		log.CreatedAt,
	).Scan(&log.ID)

	return err
}

func (r *ProgressLogRepo) GetByID(ctx context.Context, id int) (*model.ProgressLog, error) {
	var log model.ProgressLog

	query := `
		SELECT id, user_id, task_id, habit_id, xp_earned, gold_earned, created_at
		FROM progress_log
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&log.ID,
		&log.UserID,
		&log.TaskID,
		&log.HabitID,
		&log.XPEarned,
		&log.GoldEarned,
		&log.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("progress log with id=%d not found", id)
		}
		return nil, err
	}

	return &log, nil
}

func (r *ProgressLogRepo) GetByUser(ctx context.Context, userID int64, fromDate, toDate *time.Time) ([]model.ProgressLog, error) {
	query := `
		SELECT id, user_id, task_id, habit_id, xp_earned, gold_earned, created_at
		FROM progress_log
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argCount := 1

	if fromDate != nil {
		argCount++
		query += " AND created_at >= $" + fmt.Sprint(argCount)
		args = append(args, *fromDate)
	}

	if toDate != nil {
		argCount++
		query += " AND created_at <= $" + fmt.Sprint(argCount)
		args = append(args, *toDate)
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []model.ProgressLog{}
	for rows.Next() {
		var log model.ProgressLog

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.TaskID,
			&log.HabitID,
			&log.XPEarned,
			&log.GoldEarned,
			&log.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (r *ProgressLogRepo) GetByUserWithDetails(ctx context.Context, userID int64, fromDate, toDate *time.Time) ([]model.ProgressLogWithDetails, error) {
	query := `
		SELECT 
			pl.id, 
			pl.user_id, 
			pl.task_id, 
			pl.habit_id, 
			pl.xp_earned, 
			pl.gold_earned, 
			pl.created_at,
			t.title as task_title,
			h.title as habit_title,
			CASE 
				WHEN pl.task_id IS NOT NULL THEN 'task'
				WHEN pl.habit_id IS NOT NULL THEN 'habit'
			END as type
		FROM progress_log pl
		LEFT JOIN task t ON pl.task_id = t.id
		LEFT JOIN habit h ON pl.habit_id = h.id
		WHERE pl.user_id = $1
	`

	args := []interface{}{userID}
	argCount := 1

	if fromDate != nil {
		argCount++
		query += " AND pl.created_at >= $" + fmt.Sprint(argCount)
		args = append(args, *fromDate)
	}

	if toDate != nil {
		argCount++
		query += " AND pl.created_at <= $" + fmt.Sprint(argCount)
		args = append(args, *toDate)
	}

	query += " ORDER BY pl.created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []model.ProgressLogWithDetails{}
	for rows.Next() {
		var log model.ProgressLogWithDetails
		var taskTitle, habitTitle *string

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.TaskID,
			&log.HabitID,
			&log.XPEarned,
			&log.GoldEarned,
			&log.CreatedAt,
			&taskTitle,
			&habitTitle,
			&log.Type,
		)

		if err != nil {
			return nil, err
		}

		log.TaskTitle = taskTitle
		log.HabitTitle = habitTitle

		logs = append(logs, log)
	}

	return logs, nil
}

func (r *ProgressLogRepo) GetByTask(ctx context.Context, taskID int) ([]model.ProgressLog, error) {
	query := `
		SELECT id, user_id, task_id, habit_id, xp_earned, gold_earned, created_at
		FROM progress_log
		WHERE task_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []model.ProgressLog{}
	for rows.Next() {
		var log model.ProgressLog

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.TaskID,
			&log.HabitID,
			&log.XPEarned,
			&log.GoldEarned,
			&log.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (r *ProgressLogRepo) GetByHabit(ctx context.Context, habitID int) ([]model.ProgressLog, error) {
	query := `
		SELECT id, user_id, task_id, habit_id, xp_earned, gold_earned, created_at
		FROM progress_log
		WHERE habit_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, habitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []model.ProgressLog{}
	for rows.Next() {
		var log model.ProgressLog

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.TaskID,
			&log.HabitID,
			&log.XPEarned,
			&log.GoldEarned,
			&log.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (r *ProgressLogRepo) GetUserProgressForDate(ctx context.Context, userID int64, date time.Time) ([]model.ProgressLog, error) {
	query := `
		SELECT id, user_id, task_id, habit_id, xp_earned, gold_earned, created_at
		FROM progress_log
		WHERE user_id = $1 AND DATE(created_at) = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID, date.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []model.ProgressLog{}
	for rows.Next() {
		var log model.ProgressLog

		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.TaskID,
			&log.HabitID,
			&log.XPEarned,
			&log.GoldEarned,
			&log.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (r *ProgressLogRepo) GetTotalXP(ctx context.Context, userID int64, fromDate, toDate *time.Time) (int, error) {
	query := `
		SELECT COALESCE(SUM(xp_earned), 0)
		FROM progress_log
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argCount := 1

	if fromDate != nil {
		argCount++
		query += " AND created_at >= $" + fmt.Sprint(argCount)
		args = append(args, *fromDate)
	}

	if toDate != nil {
		argCount++
		query += " AND created_at <= $" + fmt.Sprint(argCount)
		args = append(args, *toDate)
	}

	var totalXP int
	err := r.db.QueryRow(ctx, query, args...).Scan(&totalXP)
	return totalXP, err
}

func (r *ProgressLogRepo) GetTotalGold(ctx context.Context, userID int64, fromDate, toDate *time.Time) (int, error) {
	query := `
		SELECT COALESCE(SUM(gold_earned), 0)
		FROM progress_log
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argCount := 1

	if fromDate != nil {
		argCount++
		query += " AND created_at >= $" + fmt.Sprint(argCount)
		args = append(args, *fromDate)
	}

	if toDate != nil {
		argCount++
		query += " AND created_at <= $" + fmt.Sprint(argCount)
		args = append(args, *toDate)
	}

	var totalGold int
	err := r.db.QueryRow(ctx, query, args...).Scan(&totalGold)
	return totalGold, err
}

func (r *ProgressLogRepo) DeleteByID(ctx context.Context, id int) error {
	query := `DELETE FROM progress_log WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("progress log with id=%d not found", id)
	}

	return nil
}

func (r *ProgressLogRepo) HasUserCompletedTaskToday(ctx context.Context, userID int64) (bool, error) {
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

func (r *ProgressLogRepo) HasUserCompletedHabitToday(ctx context.Context, userID int64) (bool, error) {
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

// GetLastByTaskID возвращает последнюю запись progress_log для задачи
func (r *ProgressLogRepo) GetLastByTaskID(ctx context.Context, taskID int) (*model.ProgressLog, error) {
	query := `
		SELECT id, user_id, task_id, habit_id, xp_earned, gold_earned, created_at
		FROM progress_log
		WHERE task_id = $1 AND xp_earned > 0
		ORDER BY created_at DESC
		LIMIT 1
	`

	var log model.ProgressLog
	err := r.db.QueryRow(ctx, query, taskID).Scan(
		&log.ID,
		&log.UserID,
		&log.TaskID,
		&log.HabitID,
		&log.XPEarned,
		&log.GoldEarned,
		&log.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("progress log for task_id=%d not found", taskID)
		}
		return nil, err
	}

	return &log, nil
}

// GetLastByHabitID возвращает последнюю запись progress_log для привычки
func (r *ProgressLogRepo) GetLastByHabitID(ctx context.Context, habitID int) (*model.ProgressLog, error) {
	query := `
		SELECT id, user_id, task_id, habit_id, xp_earned, gold_earned, created_at
		FROM progress_log
		WHERE habit_id = $1 AND xp_earned > 0
		ORDER BY created_at DESC
		LIMIT 1
	`

	var log model.ProgressLog
	err := r.db.QueryRow(ctx, query, habitID).Scan(
		&log.ID,
		&log.UserID,
		&log.TaskID,
		&log.HabitID,
		&log.XPEarned,
		&log.GoldEarned,
		&log.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("progress log for habit_id=%d not found", habitID)
		}
		return nil, err
	}

	return &log, nil
}

// DeleteByTaskID удаляет все записи progress_log для указанной задачи
func (r *ProgressLogRepo) DeleteByTaskID(ctx context.Context, taskID int) error {
	query := `DELETE FROM progress_log WHERE task_id = $1`
	result, err := r.db.Exec(ctx, query, taskID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	fmt.Printf("Deleted %d progress logs for task %d\n", rowsAffected, taskID)

	return nil
}

// DeleteByHabitID удаляет все записи progress_log для указанной привычки
func (r *ProgressLogRepo) DeleteByHabitID(ctx context.Context, habitID int) error {
	query := `DELETE FROM progress_log WHERE habit_id = $1`
	result, err := r.db.Exec(ctx, query, habitID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	fmt.Printf("Deleted %d progress logs for habit %d\n", rowsAffected, habitID)

	return nil
}

// В ProgressRepo добавьте метод:
func (r *ProgressLogRepo) GetLastActivityDate(ctx context.Context, userID int64) (time.Time, error) {
	var lastActivity time.Time
	query := `
		SELECT COALESCE(MAX(created_at), '0001-01-01'::timestamp)
		FROM progress_log 
		WHERE user_id = $1 AND xp_earned > 0
	`
	err := r.db.QueryRow(ctx, query, userID).Scan(&lastActivity)
	if err != nil {
		return time.Time{}, err
	}
	return lastActivity, nil
}
