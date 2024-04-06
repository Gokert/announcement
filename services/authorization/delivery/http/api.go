package delivery

import (
	"anncouncement/pkg/middleware"
	"anncouncement/pkg/models"
	httpResponse "anncouncement/pkg/response"
	"anncouncement/services/authorization/usecase"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
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
