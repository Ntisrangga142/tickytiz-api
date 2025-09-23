package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Auth struct {
	DB  *pgxpool.Pool
	RDB *redis.Client
}

func NewAuthRepo(db *pgxpool.Pool, rdb *redis.Client) *Auth {
	return &Auth{DB: db, RDB: rdb}
}

func (r *Auth) Register(ctx context.Context, email, password string) error {
	// mulai transaction
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction")
	}
	defer func() {
		// kalau ada panic / lupa commit → rollback
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// step 1: insert ke account
	var userID int
	queryAccount := `INSERT INTO account (email, password) VALUES ($1, $2) RETURNING id`
	if err = tx.QueryRow(ctx, queryAccount, email, password).Scan(&userID); err != nil {
		return fmt.Errorf("failed to insert account")
	}

	// step 2: insert ke users (pakai id dari account)
	queryUser := `INSERT INTO users (id) VALUES ($1)`
	if _, err = tx.Exec(ctx, queryUser, userID); err != nil {
		return fmt.Errorf("failed to insert user = %w", err)
	}

	// commit transaction
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction")
	}

	return nil
}

func (r *Auth) Login(ctx context.Context, email string) (int, string, string, error) {
	query := `SELECT id, role, password FROM account WHERE email = $1`
	var id int
	var role string
	var hashedPassword string

	err := r.DB.QueryRow(ctx, query, email).Scan(&id, &role, &hashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, "", "", nil // ✅ biar bisa bedain di handler
		}
		return 0, "", "", err
	}
	return id, role, hashedPassword, nil
}

func (r *Auth) BlacklistToken(ctx context.Context, token string, expiresIn time.Duration) error {
	data := models.BlacklistToken{
		Token:     token,
		ExpiresIn: expiresIn,
	}
	bt, err := json.Marshal(data)
	if err != nil {
		log.Println("❌ Internal Server Error.\nCause:", err)
		return err
	}

	if err := r.RDB.Set(ctx, "blacklist:"+token, bt, expiresIn).Err(); err != nil {
		log.Printf("❌ Redis Error.\nCause: %s\n", err)
		return err
	}

	return nil
}
