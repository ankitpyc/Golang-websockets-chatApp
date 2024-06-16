package Repository

import (
	databases "TCPServer/internal/database/models"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db             *gorm.DB
	ctx            context.Context
	defaultTimeout time.Duration
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{ctx: context.Background(), db: db}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (userRepo *UserRepository) CreateUser(user *databases.User) (*databases.User, error) {
	ctx, cancel := context.WithTimeout(userRepo.ctx, time.Second*15)
	defer cancel()
	hashedPassword, _ := hashPassword(user.Password)
	user.Password = hashedPassword
	result := userRepo.db.WithContext(ctx).Create(&user)
	if result.Error != nil {
		_ = fmt.Errorf("error while creating user %v \n", result.Error.Error())
		return nil, result.Error
	}
	log.Printf("Created user %v , Rows Affected : %d \n", user, result.RowsAffected)
	return user, nil
}

func (userRepo *UserRepository) GetUserByEmailId(email string) (*databases.User, error) {
	ctx, cancel := context.WithTimeout(userRepo.ctx, time.Second*5)
	defer cancel()
	var user databases.User
	result := userRepo.db.WithContext(ctx).Model(databases.User{Email: email}).First(&user)
	if result.Error != nil {
		_ = fmt.Errorf("error while creating user %v \n", result.Error.Error())
		return nil, result.Error
	}
	log.Printf("Fetched user %v , Rows Affected : %d \n", user, result.RowsAffected)
	return &user, nil
}

func (userRepo *UserRepository) Login(user *databases.User) (*databases.User, error) {
	ctx, cancel := context.WithTimeout(userRepo.ctx, time.Second*30)
	dbUser := &databases.User{}
	defer cancel()
	result := userRepo.db.WithContext(ctx).Where("email = ?", user.Email).Limit(1).First(&dbUser)
	// handle specific error while logging in
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("invalid User Name and Password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.GetPassword()), []byte(user.GetPassword())); err != nil {
		log.Println("Password does not match")
		return nil, fmt.Errorf("invalid User Name and Password")
	}

	if result.Error != nil {
		_ = fmt.Errorf("error while creating user %v \n", result.Error.Error())
		return nil, result.Error
	}

	log.Printf("Fetched user %v , Rows Affected : %d \n", user, result.RowsAffected)
	return dbUser, nil
}
