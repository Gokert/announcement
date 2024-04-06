package delivery

import (
	utils "anncouncement/pkg"
	"anncouncement/pkg/middleware"
	"anncouncement/pkg/models"
	httpResponse "anncouncement/pkg/response"
	"anncouncement/services/announcement/usecase"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Api struct {
	log  *logrus.Logger
	mx   *http.ServeMux
	core usecase.ICore
}

func GetApi(core *usecase.Core, log *logrus.Logger) *Api {
	api := &Api{
		core: core,
		log:  log,
		mx:   http.NewServeMux(),
	}

	api.mx.Handle("/api/v1/announcements", middleware.MethodCheck(http.HandlerFunc(api.GetAnnouncement), http.MethodGet, log))
	api.mx.Handle("/api/v1/announcements/list", middleware.MethodCheck(http.HandlerFunc(api.GetAnnouncements), http.MethodGet, log))
	api.mx.Handle("/api/v1/announcements/search", middleware.MethodCheck(http.HandlerFunc(api.SearchAnnouncements), http.MethodGet, log))
	api.mx.Handle("/api/v1/announcements/create", middleware.MethodCheck(middleware.AuthCheck(http.HandlerFunc(api.CreateAnnouncement), core, log), http.MethodPost, log))

	return api
}

func (a *Api) ListenAndServe(port string) error {
	err := http.ListenAndServe(":"+port, a.mx)
	if err != nil {
		a.log.Errorf("listen error: %s", err.Error())
		return err
	}

	return nil
}

func (a *Api) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		a.log.Error("parse form-data error", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	header := r.FormValue("header")
	info := r.FormValue("info")
	cost, err := strconv.ParseUint(r.FormValue("cost"), 10, 64)
	if err != nil {
		a.log.Errorf("parse cost error: %s", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	photo, handler, err := r.FormFile("photo")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		a.log.Errorf("parse photo error: %s", err.Error())
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	filename, status, err := utils.SaveImage(&photo, handler, "./images/")
	if err != nil {
		a.log.Errorf("save photo error: %s", err.Error())
		response.Status = status
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = utils.ValidateSize(header, info, cost)
	if err != nil {
		a.log.Errorf("validate size error: %s", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	userId, _ := r.Context().Value(middleware.UserIDKey).(uint64)

	err = a.core.CreateAnnouncement(&models.Announcement{Header: header, Info: info, Cost: cost, Photo: filename}, userId)
	if err != nil {
		a.log.Errorf("create announcement error: %s", err.Error())
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

func (a *Api) GetAnnouncements(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		page = 0
	}

	pageSize, err := strconv.ParseUint(r.URL.Query().Get("per_page"), 10, 64)
	if err != nil {
		pageSize = 8
	}

	announcements, err := a.core.GetAnnouncements(page, pageSize)
	if err != nil {
		a.log.Errorf("get announcements error: %s", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	response.Body = &models.Announcements{
		Count:         uint64(len(announcements)),
		Announcements: announcements,
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

func (a *Api) GetAnnouncement(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	annId, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		a.log.Error("Parse id error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	announcement, err := a.core.GetAnnouncement(annId)
	if err != nil {
		a.log.Error("Get announcement error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	response.Body = announcement

	httpResponse.SendResponse(w, r, &response, a.log)
}

func (a *Api) SearchAnnouncements(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 64)
	if err != nil {
		page = 0
	}

	pageSize, err := strconv.ParseUint(r.URL.Query().Get("per_page"), 10, 64)
	if err != nil {
		pageSize = 8
	}

	minCost, err := strconv.ParseUint(r.URL.Query().Get("min_cost"), 10, 64)
	if err != nil {
		minCost = 0
	}

	maxCost, err := strconv.ParseUint(r.URL.Query().Get("max_cost"), 10, 64)
	if err != nil {
		maxCost = 0
	}

	order := r.URL.Query().Get("sort_by")
	if order == "" {
		order = "date"
	}

	announcements, err := a.core.SearchAnnouncements(page, pageSize, minCost, maxCost, order)
	if err != nil {
		a.log.Error("Search announcements error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	response.Body = &models.Announcements{
		Count:         uint64(len(announcements)),
		Announcements: announcements,
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}
