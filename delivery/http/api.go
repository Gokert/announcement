package delivery

import (
	"encoding/json"
	"errors"
	utils "filmoteka/pkg"
	"filmoteka/pkg/middleware"
	"filmoteka/pkg/models"
	httpResponse "filmoteka/pkg/response"
	"filmoteka/usecase"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"time"
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

	api.mx.Handle("/signin", middleware.MethodCheck(http.HandlerFunc(api.Signin), http.MethodPost, log))
	api.mx.Handle("/signup", middleware.MethodCheck(http.HandlerFunc(api.Signup), http.MethodPost, log))
	api.mx.Handle("/logout", middleware.MethodCheck(http.HandlerFunc(api.Logout), http.MethodDelete, log))
	api.mx.Handle("/authcheck", middleware.MethodCheck(http.HandlerFunc(api.AuthAccept), http.MethodGet, log))

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

func (a *Api) Signin(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}
	var request models.SigninRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.log.Error("Read all error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		a.log.Error("Unmarshal error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	_, found, err := a.core.FindUserAccount(request.Login, request.Password)
	if err != nil {
		a.log.Error("Find user account error: ", err.Error())
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	if !found {
		response.Status = http.StatusNotFound
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	session, _ := a.core.CreateSession(r.Context(), request.Login)
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.SID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	httpResponse.SendResponse(w, r, &response, a.log)
}

func (a *Api) Signup(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}
	var request models.SignupRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	found, err := a.core.FindUserByLogin(request.Login)
	if err != nil {
		a.log.Errorf("Find user by login error: %s", err.Error())
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	if found {
		response.Status = http.StatusConflict
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = a.core.CreateUserAccount(request.Login, request.Password)
	if err != nil {
		a.log.Error("Create user error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

func (a *Api) Logout(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = a.core.KillSession(r.Context(), cookie.Value)
	if err != nil {
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)

	httpResponse.SendResponse(w, r, &response, a.log)
}

func (a *Api) AuthAccept(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}
	var authorized bool

	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized, _ = a.core.FindActiveSession(r.Context(), session.Value)
	}

	if !authorized {
		response.Status = http.StatusUnauthorized
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	login, err := a.core.GetUserName(r.Context(), session.Value)
	if err != nil {
		a.log.Errorf("Get user name error: %s", err.Error())
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	response.Body = models.AuthCheckResponse{
		Login: login,
	}

	httpResponse.SendResponse(w, r, &response, a.log)
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
