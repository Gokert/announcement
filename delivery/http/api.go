package delivery

import (
	"encoding/json"
	"errors"
	_ "filmoteka/docs"
	utils "filmoteka/pkg"
	"filmoteka/pkg/middleware"
	"filmoteka/pkg/models"
	httpResponse "filmoteka/pkg/response"
	"filmoteka/usecase"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
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

	api.mx.HandleFunc("/signin", api.Signin)
	api.mx.HandleFunc("/signup", api.Signup)
	api.mx.HandleFunc("/logout", api.Logout)
	api.mx.HandleFunc("/authcheck", api.AuthAccept)

	api.mx.HandleFunc("/api/v1/announcements", api.GetAnnouncement)
	api.mx.HandleFunc("/api/v1/announcements/list", api.GetAnnouncements)
	api.mx.HandleFunc("/api/v1/announcements/search", api.SearchAnnouncements)
	api.mx.Handle("/api/v1/announcements/create", middleware.AuthCheck(http.HandlerFunc(api.CreateAnnouncement), core, log))

	return api
}

func (a *Api) ListenAndServe(port string) error {
	err := http.ListenAndServe(":"+port, a.mx)
	if err != nil {
		a.log.Error("ListenAndServer error: ", err.Error())
		return err
	}

	return nil
}

func (a *Api) Signin(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodPost {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	var request models.SigninRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.log.Error("Signin error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		a.log.Error("Signin error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	_, found, err := a.core.FindUserAccount(request.Login, request.Password)
	if err != nil {
		a.log.Error("Signin error: ", err.Error())
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	if !found {
		response.Status = http.StatusUnauthorized
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

	if r.Method != http.MethodPost {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

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
		a.log.Error("Signup error: ", err.Error())
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
		a.log.Error("create user error: ", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	httpResponse.SendResponse(w, r, &response, a.log)
}

func (a *Api) Logout(w http.ResponseWriter, r *http.Request) {
	response := models.Response{Status: http.StatusOK, Body: nil}

	if r.Method != http.MethodDelete {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

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

	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	session, err := r.Cookie("session_id")
	if err == nil && session != nil {
		authorized, _ = a.core.FindActiveSession(r.Context(), session.Value)
	}
	a.log.Warning("API", authorized)
	if !authorized {
		response.Status = http.StatusUnauthorized
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	login, err := a.core.GetUserName(r.Context(), session.Value)
	if err != nil {
		a.log.Error("auth accept error: ", err.Error())
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

	if r.Method != http.MethodPost {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

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

	if handler == nil {
		a.log.Errorf("photo not found")
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	err = utils.ValidateImage(path.Ext(handler.Filename))
	if err != nil {
		a.log.Errorf("validate image error: %s", err.Error())
		response.Status = http.StatusBadRequest
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

	filename := "./images/" + handler.Filename
	if err != nil && handler != nil && photo != nil {
		a.log.Errorf("save photo error: %s", err.Error())
		response.Status = http.StatusBadRequest
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

	filePhoto, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		a.log.Errorf("open file error: %s", err.Error())
		response.Status = http.StatusInternalServerError
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}
	defer filePhoto.Close()

	_, err = io.Copy(filePhoto, photo)
	if err != nil {
		a.log.Errorf("file copy error: %s", err.Error())
		response.Status = http.StatusInternalServerError
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

	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

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
		a.log.Error("get announcements error: ", err.Error())
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

	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

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

	if r.Method != http.MethodGet {
		response.Status = http.StatusMethodNotAllowed
		httpResponse.SendResponse(w, r, &response, a.log)
		return
	}

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
