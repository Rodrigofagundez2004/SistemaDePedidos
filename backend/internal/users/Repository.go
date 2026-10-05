package users

import (
	"errors"

	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailAlreadyExists = errors.New("email already exists")
var NameAlreadyExists = errors.New("name already exists")

type Repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) *Repository {

	return &Repository{db: db}
}
func (r *Repository) Create(u *User) error {
	//verificar que el emial no exista
	var count int64
	if err := r.db.Model(&User{}).Where("email = ?", u.Email).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrEmailAlreadyExists
	}
	return r.db.Create(u).Error
}
func (r *Repository) FindByEmail(email string) (*User, error) {
	var u User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}
func (r *Repository) FindByName(fullName string) (*User, error) {
	var u User
	if err := r.db.Where("full_name = ?", fullName).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}
func (r *Repository) FindByID(id string) (*User, error) {
	var u User
	err := r.db.Where("id = ?", id).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
