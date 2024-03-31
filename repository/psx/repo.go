package psx

import (
	"database/sql"
	"errors"
	"filmoteka/configs"
	"filmoteka/pkg/models"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"strings"
)

//go:generate mockgen -source=repo.go -destination=../../mocks/repo_mock.go -package=mocks
type IRepository interface {
	GetUser(login string, password []byte) (*models.UserItem, bool, error)
	FindUser(login string) (bool, error)
	CreateUser(login string, password []byte) error
	GetUserId(login string) (uint64, error)

	GetAnnouncements(page, pageSize uint64) ([]models.Announcement, error)
	GetAnnouncement(id uint64) (*models.Announcement, error)
	SearchAnnouncements(page, pageSize, minCost, maxCost uint64, order string) ([]models.Announcement, error)
	CreateAnnouncement(announcement *models.Announcement, userId uint64) error
}

type Repository struct {
	db *sql.DB
}

func GetPsxRepo(config *configs.DbPsxConfig) (*Repository, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s password= %s host=%s port=%d sslmode=%s",
		config.User, config.Dbname, config.Password, config.Host, config.Port, config.Sslmode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open error: %s", err.Error())
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("sql ping error: %s", err.Error())
	}
	db.SetMaxOpenConns(config.MaxOpenConns)

	return &Repository{db: db}, nil
}

func (repo *Repository) GetUser(login string, password []byte) (*models.UserItem, bool, error) {
	post := &models.UserItem{}

	err := repo.db.QueryRow("SELECT profile.id, profile.login FROM profile "+
		"WHERE profile.login = $1 AND profile.password = $2 ", login, password).Scan(&post.Id, &post.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("get query user error: %s", err.Error())
	}

	return post, true, nil
}

func (repo *Repository) FindUser(login string) (bool, error) {
	post := &models.UserItem{}

	err := repo.db.QueryRow(
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

func (repo *Repository) CreateUser(login string, password []byte) error {
	var userID uint64
	err := repo.db.QueryRow("INSERT INTO profile(login, password) VALUES($1, $2) RETURNING id", login, password).Scan(&userID)
	if err != nil {
		return fmt.Errorf("create user error: %s", err.Error())
	}

	return nil
}

func (repo *Repository) GetUserId(login string) (uint64, error) {
	var userID uint64

	err := repo.db.QueryRow(
		"SELECT profile.id FROM profile WHERE profile.login = $1", login).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user not found for login: %s", login)
		}
		return 0, fmt.Errorf("select user profile id error: %s", err.Error())
	}

	return userID, nil
}

func (repo *Repository) GetAnnouncements(page uint64, pageSize uint64) ([]models.Announcement, error) {
	var announcements []models.Announcement

	rows, err := repo.db.Query("SELECT announcement.header, announcement.photo_href, announcement.info, announcement.cost FROM announcement OFFSET $1 LIMIT $2", page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("get announcements in repo error: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var announcement models.Announcement

		err = rows.Scan(&announcement.Header, &announcement.Info, &announcement.Photo, &announcement.Cost)
		if err != nil {
			return nil, fmt.Errorf("get announcements scan error: %s", err.Error())
		}

		announcements = append(announcements, announcement)
	}

	return announcements, nil
}

func (repo *Repository) GetAnnouncement(id uint64) (*models.Announcement, error) {
	var announcement models.Announcement

	err := repo.db.QueryRow("SELECT announcement.header, announcement.photo_href, announcement.info, announcement.cost FROM announcement WHERE announcement.id = $1", id).Scan(&announcement.Header, &announcement.Info, &announcement.Photo, &announcement.Cost)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("select announcement error: %s", err.Error())
	}

	return &announcement, nil
}

func (repo *Repository) CreateAnnouncement(announcement *models.Announcement, userId uint64) error {
	_, err := repo.db.Exec("INSERT INTO announcement(id_profile, header, photo_href, info, cost) VALUES ($1, $2, $3, $4, $5)", userId, announcement.Header, announcement.Photo, announcement.Info, announcement.Cost)
	if err != nil {
		return fmt.Errorf("exec create announcement error: %s", err.Error())
	}

	return nil
}

func (repo *Repository) SearchAnnouncements(page, pageSize, minCost, maxCost uint64, order string) ([]models.Announcement, error) {
	var announcements []models.Announcement
	var str strings.Builder
	var params []interface{}

	str.WriteString("SELECT announcement.header, announcement.photo_href, announcement.info, announcement.cost FROM announcement ")

	switch maxCost {
	case 0:
		str.WriteString(fmt.Sprintf("WHERE announcement.cost > $1 ORDER BY announcement.%s DESC OFFSET $2 LIMIT $3 ", order))
		params = append(params, minCost, page, pageSize)
	default:
		str.WriteString(fmt.Sprintf("WHERE announcement.cost > $1 AND announcement.cost < $2 ORDER BY announcement.%s DESC OFFSET $3 LIMIT $4 ", order))
		params = append(params, minCost, maxCost, page, pageSize)
	}

	rows, err := repo.db.Query(str.String(), params...)
	if err != nil {
		return nil, fmt.Errorf("query error: %s", err.Error())
	}

	for rows.Next() {
		var announcement models.Announcement

		err = rows.Scan(&announcement.Header, &announcement.Photo, &announcement.Info, &announcement.Cost)
		if err != nil {
			return nil, fmt.Errorf("scan error: %s", err.Error())
		}

		announcements = append(announcements, announcement)
	}

	fmt.Println(announcements)

	return announcements, nil
}
