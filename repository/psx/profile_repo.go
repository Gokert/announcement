package psx

import (
	"database/sql"
	"errors"
	"filmoteka/configs"
	"filmoteka/pkg/models"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/sirupsen/logrus"
)

type IProfileRepo interface {
	GetUser(login string, password []byte) (*models.UserItem, bool, error)
	FindUser(login string) (bool, error)
	CreateUser(login string, password []byte) error
	GetUserId(login string) (uint64, error)
}

type ProfileRepo struct {
	DB *sql.DB
}

func GetFilmRepo(config *configs.DbPsxConfig, log *logrus.Logger) (*ProfileRepo, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s password= %s host=%s port=%d sslmode=%s",
		config.User, config.Dbname, config.Password, config.Host, config.Port, config.Sslmode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Errorf("sql open error: %s", err.Error())
		return nil, fmt.Errorf("get user repo err: %s", err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Errorf("sql ping error: %s", err.Error())
		return nil, fmt.Errorf("get user repo error: %s", err.Error())
	}
	db.SetMaxOpenConns(config.MaxOpenConns)

	log.Info("Psx created successful")
	return &ProfileRepo{DB: db}, nil
}

func (repo *ProfileRepo) GetUser(login string, password []byte) (*models.UserItem, bool, error) {
	post := &models.UserItem{}

	err := repo.DB.QueryRow("SELECT profile.id, profile.login FROM profile "+
		"WHERE profile.login = $1 AND profile.password = $2 ", login, password).Scan(&post.Id, &post.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("get query user error: %s", err.Error())
	}

	err = repo.DB.QueryRow("SELECT profile.id, profile.login, profile.password FROM profile ").Scan(&post.Id, &post.Login, &post.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("get query user error: %s", err.Error())
	}

	return post, true, nil
}

func (repo *ProfileRepo) FindUser(login string) (bool, error) {
	post := &models.UserItem{}

	err := repo.DB.QueryRow(
		"SELECT login FROM profile "+
			"WHERE login = $1", login).Scan(&post.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("find user query error: %s", err.Error())
	}

	return true, nil
}

func (repo *ProfileRepo) CreateUser(login string, password []byte) error {
	var userID uint64
	err := repo.DB.QueryRow("INSERT INTO profile(login, password) VALUES($1, $2) RETURNING id", login, password).Scan(&userID)
	if err != nil {
		return fmt.Errorf("create user error: %s", err.Error())
	}

	return nil
}

func (repo *ProfileRepo) GetUserId(login string) (uint64, error) {
	var userID uint64

	err := repo.DB.QueryRow(
		"SELECT profile.id FROM profile WHERE profile.login = $1", login).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user not found for login: %s", login)
		}
		return 0, fmt.Errorf("get userpro file id error: %s", err.Error())
	}

	return userID, nil
}
