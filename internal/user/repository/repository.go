package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"user-service/internal/models"
	"user-service/internal/user/service"
)

type userRepository struct {
	Db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) service.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	const op = "userRepository.create"
	query := `INSERT INTO users (id, email, phone, name, last_name, second_name, birth_date, sex) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`
	if err := r.Db.QueryRowxContext(ctx, query, &user.ID, &user.Email, &user.Name, &user.LastName, &user.SecondName,
		&user.BirthDate, &user.Sex).StructScan(user); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	const op = "userRepository.GetByID"
	query := "SELECT * FROM users WHERE id = $1"
	user := &models.User{}
	err := r.Db.QueryRowxContext(ctx, query, id).StructScan(user)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "userRepository.GetByEmail"
	query := "SELECT * FROM users WHERE email = $1"
	user := &models.User{}
	if err := r.Db.QueryRowxContext(ctx, query, email).StructScan(user); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	const op = "userRepository.Update"
	query := `UPDATE users
			  SET email = $1, phone = $2, name = $3, last_name = $4, second_name = $5, birth_date = $6, sex = $7
			  WHERE id = $8 RETURNING *`
	newUser := &models.User{}
	if err := r.Db.QueryRowxContext(ctx, query, user.Email, user.Phone, user.Name, user.LastName, user.SecondName,
		user.BirthDate, user.Sex, user.ID).StructScan(newUser); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return newUser, nil
}
