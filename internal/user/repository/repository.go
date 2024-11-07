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
	query := `INSERT INTO users (email, name, last_name, second_name, birth_date, sex) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`
	if err := r.Db.QueryRowxContext(ctx, query, &user.Email, &user.Name, &user.LastName, &user.SecondName,
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
	//TODO implement me
	panic("implement me")
}
