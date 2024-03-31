package usecase

import (
	"context"
	"filmoteka/configs"
	utils "filmoteka/pkg"
	"filmoteka/pkg/models"
	"filmoteka/repository/psx"
	"filmoteka/repository/session"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

//go:generate mockgen -source=core.go -destination=../mocks/core_mock.go -package=mocks
type ICore interface {
	GetUserName(ctx context.Context, sid string) (string, error)
	CreateSession(ctx context.Context, login string) (models.Session, error)
	FindActiveSession(ctx context.Context, sid string) (bool, error)
	KillSession(ctx context.Context, sid string) error
	GetUserId(ctx context.Context, sid string) (uint64, error)

	CreateUserAccount(login string, password string) error
	FindUserAccount(login string, password string) (*models.UserItem, bool, error)
	FindUserByLogin(login string) (bool, error)

	GetAnnouncements(page uint64, pageSize uint64) ([]models.Announcement, error)
	GetAnnouncement(id uint64) (*models.Announcement, error)
	CreateAnnouncement(announcement *models.Announcement, userId uint64) error
	SearchAnnouncements(page, pageSize, minCost, maxCost uint64, order string) ([]models.Announcement, error)
}

type Core struct {
	log           *logrus.Logger
	profiles      psx.IRepository
	announcements psx.IRepository
	sessions      session.ISessionRepo
}

func GetCore(psxCfg *configs.DbPsxConfig, redisCfg *configs.DbRedisCfg, log *logrus.Logger) (*Core, error) {
	filmRepo, err := psx.GetPsxRepo(psxCfg)
	if err != nil {
		return nil, fmt.Errorf("get psx error error: %s", err.Error())
	}
	log.Info("Psx created successful")

	authRepo, err := session.GetAuthRepo(redisCfg)
	if err != nil {
		return nil, fmt.Errorf("get auth repo error: %s", err.Error())
	}
	log.Info("Redis created successful")

	core := &Core{
		log:           log,
		profiles:      filmRepo,
		sessions:      authRepo,
		announcements: filmRepo,
	}

	return core, nil
}

func (c *Core) GetUserId(ctx context.Context, sid string) (uint64, error) {
	login, err := c.sessions.GetUserLogin(ctx, sid)
	if err != nil {
		return 0, fmt.Errorf("get user login error: %s", err.Error())
	}

	id, err := c.profiles.GetUserId(login)
	if err != nil {
		return 0, fmt.Errorf("get user id error: %s", err.Error())
	}

	return id, nil
}

func (c *Core) GetUserName(ctx context.Context, sid string) (string, error) {
	login, err := c.sessions.GetUserLogin(ctx, sid)
	if err != nil {
		return "", fmt.Errorf("get user name error: %s", err.Error())
	}

	return login, nil
}

func (c *Core) CreateSession(ctx context.Context, login string) (models.Session, error) {
	sid := utils.RandStringRunes(32)

	newSession := models.Session{
		Login:     login,
		SID:       sid,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	sessionAdded, err := c.sessions.AddSession(ctx, newSession)

	if !sessionAdded && err != nil {
		return models.Session{}, err
	}

	if !sessionAdded {
		return models.Session{}, nil
	}

	return newSession, nil
}

func (c *Core) FindActiveSession(ctx context.Context, sid string) (bool, error) {
	login, err := c.sessions.CheckActiveSession(ctx, sid)

	if err != nil {
		return false, fmt.Errorf("find active session error: %s", err.Error())
	}

	return login, nil
}

func (c *Core) KillSession(ctx context.Context, sid string) error {
	_, err := c.sessions.DeleteSession(ctx, sid)

	if err != nil {
		return fmt.Errorf("delete session error: %s", err.Error())
	}

	return nil
}

func (c *Core) CreateUserAccount(login string, password string) error {
	hashPassword := utils.HashPassword(password)
	err := c.profiles.CreateUser(login, hashPassword)
	if err != nil {
		return fmt.Errorf("create user account error: %s", err.Error())
	}

	return nil
}

func (c *Core) FindUserAccount(login string, password string) (*models.UserItem, bool, error) {
	hashPassword := utils.HashPassword(password)
	user, found, err := c.profiles.GetUser(login, hashPassword)
	if err != nil {
		return nil, false, fmt.Errorf("find user account error: %s", err.Error())
	}
	return user, found, nil
}

func (c *Core) FindUserByLogin(login string) (bool, error) {
	found, err := c.profiles.FindUser(login)
	if err != nil {
		return false, fmt.Errorf("find user by login error: %s", err.Error())
	}

	return found, nil
}

func (c *Core) GetAnnouncements(page uint64, pageSize uint64) ([]models.Announcement, error) {
	announcements, err := c.announcements.GetAnnouncements(page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("get announcements error: %s", err.Error())
	}

	return announcements, nil
}

func (c *Core) GetAnnouncement(id uint64) (*models.Announcement, error) {
	announcement, err := c.announcements.GetAnnouncement(id)
	if err != nil {
		return nil, fmt.Errorf("get announcement: %s", err.Error())
	}

	return announcement, nil
}

func (c *Core) CreateAnnouncement(announcement *models.Announcement, userId uint64) error {

	err := c.announcements.CreateAnnouncement(announcement, userId)
	if err != nil {
		return fmt.Errorf("create announcement error: %s", err.Error())
	}

	return err
}

func (c *Core) SearchAnnouncements(page, pageSize, minCost, maxCost uint64, order string) ([]models.Announcement, error) {
	announcements, err := c.announcements.SearchAnnouncements(page, pageSize, minCost, maxCost, order)
	if err != nil {
		return nil, fmt.Errorf("search announcements error: %s", err.Error())
	}

	return announcements, nil
}
