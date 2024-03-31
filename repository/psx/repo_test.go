package psx

import (
	"database/sql"
	utils "filmoteka/pkg"
	"filmoteka/pkg/models"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"regexp"
	"testing"
)

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "login"})

	hashPsw := utils.HashPassword("p1")
	expect := []*models.UserItem{
		{Id: 1, Login: "l1"},
		{Id: 12, Login: "kek"},
	}

	selectRow := "SELECT profile.id, profile.login FROM profile WHERE profile.login = $1 AND profile.password = $2"

	repo := &Repository{
		db: db,
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Id, item.Login)

		mock.
			ExpectQuery(regexp.QuoteMeta(selectRow)).
			WithArgs(item.Login, hashPsw).
			WillReturnRows(rows)

		user, found, err := repo.GetUser(item.Login, hashPsw)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		if !found {
			t.Errorf("user not found")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

		if !reflect.DeepEqual(user, item) {
			t.Errorf("results not match, want %v, have %v", item, user)
		}
	}

	mock.ExpectQuery(regexp.QuoteMeta(selectRow)).
		WithArgs("unknown_user", hashPsw).
		WillReturnError(sql.ErrNoRows)

	_, found, err := repo.GetUser("unknown_user", hashPsw)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if found {
		t.Errorf("expected user not found")
	}

	mock.
		ExpectQuery(regexp.QuoteMeta(selectRow)).
		WithArgs("error_user", hashPsw).
		WillReturnError(fmt.Errorf("db_error"))

	_, found, err = repo.GetUser("error_user", hashPsw)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestFindUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"login"})
	testUser := models.UserItem{
		Login: "l1",
	}
	expect := []*models.UserItem{&testUser}
	selectRow := "SELECT login FROM profile WHERE"

	for _, item := range expect {
		rows = rows.AddRow(item.Login)
	}

	repo := &Repository{
		db: db,
	}

	mock.
		ExpectQuery(regexp.QuoteMeta(selectRow)).
		WithArgs(expect[0].Login).
		WillReturnRows(rows)

	foundAccount, err := repo.FindUser(expect[0].Login)
	if err != nil {
		t.Errorf("GetUser error: %s", err)
	}
	if !foundAccount {
		t.Errorf("user not found")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.
		ExpectQuery(regexp.QuoteMeta(selectRow)).
		WithArgs(expect[0].Login).
		WillReturnError(fmt.Errorf("db_error"))

	found, err := repo.FindUser(expect[0].Login)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if found {
		t.Errorf("expected not found")
	}

	mock.
		ExpectQuery(regexp.QuoteMeta(selectRow)).
		WithArgs(expect[0].Login).
		WillReturnError(sql.ErrNoRows)

	found, err = repo.FindUser(expect[0].Login)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if found {
		t.Errorf("expected user not found")
	}

}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"login"})

	password := "p1"
	hashPsw := utils.HashPassword(password)

	expect := []*models.UserItem{
		{Login: "l1", Password: password},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Login)
	}

	repo := &Repository{
		db: db,
	}

	selectRow := "INSERT INTO profile(login, password) VALUES($1, $2) RETURNING id"

	mock.
		ExpectQuery(regexp.QuoteMeta(selectRow)).
		WithArgs(expect[0].Login, hashPsw).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = repo.CreateUser(expect[0].Login, hashPsw)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.
		ExpectQuery(regexp.QuoteMeta(selectRow)).
		WithArgs(expect[0].Login).
		WillReturnError(sql.ErrNoRows)

	err = repo.CreateUser(expect[0].Login, hashPsw)
	if err == nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

}

func TestGetUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id"})

	expect := []*models.UserItem{
		{Id: 1, Login: "l1"},
	}

	for _, item := range expect {
		rows = rows.AddRow(item.Id)
	}

	repo := &Repository{
		db: db,
	}

	selectRow := "SELECT profile.id FROM profile WHERE profile.login = $1"

	mock.
		ExpectQuery(regexp.QuoteMeta(selectRow)).
		WithArgs(expect[0].Login).
		WillReturnRows(rows)

	id, err := repo.GetUserId(expect[0].Login)
	if err != nil {
		t.Errorf("get user id error: %s", err.Error())
		return
	}
	if id == 0 {
		t.Errorf("user not found")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	mock.ExpectQuery(regexp.QuoteMeta(selectRow)).
		WithArgs(expect[0].Login).
		WillReturnError(sql.ErrNoRows)

	id, err = repo.GetUserId(expect[0].Login)
	if err == nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	mock.
		ExpectQuery(regexp.QuoteMeta(selectRow)).
		WithArgs(expect[0].Login).
		WillReturnError(fmt.Errorf("db_error"))

	id, err = repo.GetUserId(expect[0].Login)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
