package usecase

import (
	"anncouncement/configs"
	"anncouncement/pkg/models"
	annoucement_repo "anncouncement/services/announcement/repository/annoucement"
	auth "anncouncement/services/authorization/delivery/proto"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:generate mockgen -source=core.go -destination=../mocks/core_mock.go -package=mocks
type ICore interface {
	GetAnnouncements(page uint64, pageSize uint64) ([]models.Announcement, error)
	GetAnnouncement(id uint64) (*models.Announcement, error)
	CreateAnnouncement(announcement *models.Announcement, userId uint64) error
	SearchAnnouncements(page, pageSize, minCost, maxCost uint64, order string) ([]models.Announcement, error)

	GetUserId(ctx context.Context, sid string) (uint64, error)
}

type Core struct {
	log           *logrus.Logger
	announcements annoucement_repo.IRepository
	client        auth.AuthorizationClient
}

func GetClient(address string) (auth.AuthorizationClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc connect err: %w", err)
	}
	client := auth.NewAuthorizationClient(conn)

	return client, nil
}

func GetCore(grpcCfg *configs.GrpcConfig, psxCfg *configs.DbPsxConfig, log *logrus.Logger) (*Core, error) {
	filmRepo, err := annoucement_repo.GetPsxRepo(psxCfg)
	if err != nil {
		return nil, fmt.Errorf("get psx error error: %s", err.Error())
	}
	log.Info("Psx created successful")

	authRepo, err := GetClient(grpcCfg.Addr + ":" + grpcCfg.Port)
	if err != nil {
		return nil, fmt.Errorf("get auth repo error: %s", err.Error())
	}

	core := &Core{
		log:           log,
		announcements: filmRepo,
		client:        authRepo,
	}

	return core, nil
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

func (c *Core) GetUserId(ctx context.Context, sid string) (uint64, error) {
	response, err := c.client.GetId(ctx, &auth.FindIdRequest{Sid: sid})
	if err != nil {
		return 0, fmt.Errorf("get user id err: %w", err)
	}
	return response.Value, nil
}
