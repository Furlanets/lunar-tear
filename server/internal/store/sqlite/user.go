package sqlite

import (
	"database/sql"
	"fmt"

	"lunar-tear/server/internal/store"
)

func (s *SQLiteStore) CreateUser(uuid string) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	var existingId int64
	err = tx.QueryRow(`SELECT user_id FROM users WHERE uuid = ?`, uuid).Scan(&existingId)
	if err == nil {
		return existingId, nil
	}

	nowMillis := s.clock().UnixMilli()

	res, err := tx.Exec(`INSERT INTO users (uuid, player_id, os_type, platform_type, user_restriction_type,
		register_datetime, game_start_datetime, latest_version, birth_year, birth_month,
		backup_token, charge_money_this_month) VALUES (?, 0, 2, 2, 0, ?, ?, 0, 2000, 1, 'mock-backup-token', 0)`,
		uuid, nowMillis, nowMillis)
	if err != nil {
		return 0, fmt.Errorf("insert user: %w", err)
	}
	userId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id: %w", err)
	}

	// player_id = user_id
	if _, err := tx.Exec(`UPDATE users SET player_id = ? WHERE user_id = ?`, userId, userId); err != nil {
		return 0, fmt.Errorf("update player_id: %w", err)
	}

	user := store.SeedUserState(userId, uuid, nowMillis)
	if err := writeUserState(tx, userId, user); err != nil {
		return 0, fmt.Errorf("write seed state: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return userId, nil
}

func (s *SQLiteStore) GetUserByUUID(uuid string) (int64, error) {
	var userId int64
	err := s.db.QueryRow(`SELECT user_id FROM users WHERE uuid = ?`, uuid).Scan(&userId)
	if err == sql.ErrNoRows {
		return 0, store.ErrNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("query user: %w", err)
	}
	return userId, nil
}

func (s *SQLiteStore) DefaultUserId() (int64, error) {
	var userId int64
	err := s.db.QueryRow(`SELECT min(user_id) FROM users`).Scan(&userId)
	if err != nil || userId == 0 {
		return 0, store.ErrNotFound
	}
	return userId, nil
}

func (s *SQLiteStore) UpdateUser(userId int64, mutate func(*store.UserState)) (store.UserState, error) {
	before, err := s.LoadUser(userId)
	if err != nil {
		return store.UserState{}, err
	}

	after := store.CloneUserState(before)
	mutate(&after)

	tx, err := s.db.Begin()
	if err != nil {
		return store.UserState{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if err := diffAndSave(tx, userId, &before, &after); err != nil {
		return store.UserState{}, fmt.Errorf("diff and save: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return store.UserState{}, fmt.Errorf("commit: %w", err)
	}

	return after, nil
}
