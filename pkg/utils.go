package utils

import (
	"crypto/sha512"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"unicode/utf8"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

const (
	maxLenHeader = 50
	minLenHeader = 1
	maxLenInfo   = 150
	minLenInfo   = 1
	minCost      = 0
)

func HashPassword(password string) []byte {
	hashPassword := sha512.Sum512([]byte(password))
	passwordByteSlice := hashPassword[:]
	return passwordByteSlice
}

func RandStringRunes(seed int) string {
	symbols := make([]rune, seed)
	for i := range symbols {
		symbols[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(symbols)
}

func ValidateSize(header, info string, cost uint64) error {
	headerLength := utf8.RuneCountInString(header)
	if headerLength < minLenHeader || headerLength > maxLenHeader {
		return fmt.Errorf("header size error")
	}

	infoLength := utf8.RuneCountInString(info)
	if infoLength < minLenInfo || infoLength > maxLenInfo {
		return fmt.Errorf("info size error")
	}

	if cost < minCost {
		return fmt.Errorf("cost error")
	}

	return nil
}

func ValidateImage(filename string) error {
	allowedExtensions := []string{".png", ".jpg", ".jpeg", ".webp"}

	matchFound := false
	for _, allowedExtension := range allowedExtensions {
		if filename == allowedExtension {
			matchFound = true
			break
		}
	}

	if !matchFound {
		return fmt.Errorf("invalid file extension: %s", filename)
	}

	return nil
}

func SaveImage(photo *multipart.File, handler *multipart.FileHeader, pathSave string) (string, int, error) {
	if handler == nil {
		return "", http.StatusBadRequest, fmt.Errorf("photo not found")
	}

	err := ValidateImage(path.Ext(handler.Filename))
	if err != nil {
		return "", http.StatusBadRequest, fmt.Errorf("validate image error: %s", err.Error())

	}

	filename := pathSave + handler.Filename
	if err != nil && handler != nil && photo != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("save photo error: %s", err.Error())
	}

	filePhoto, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("open file error: %s", err.Error())
	}
	defer filePhoto.Close()

	_, err = io.Copy(filePhoto, *photo)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("file copy error: %s", err.Error())
	}

	return filename, 0, nil
}

const (
	InvalidEmailOrPasswordError     = "Invalid email or password"
	SessionRepositoryNotActiveError = "Session repository not active"
	ProfileRepositoryNotActiveError = "Profile repository not active"
	CreateProfileError              = "Create profile failed"
	ProfileNotFoundError            = "Profile not found"
	GetProfileError                 = "Get profile failed"
	GetProfileRoleError             = "Get profile role failed"
	RatingSizeError                 = "Rating must be from 0 to 10"
	TitleSizeError                  = "Title size must be from 1 to 150"
	DescriptionSizeError            = "Description size must be from 1 to 1000"
	FilmsListNotFoundError          = "Films list not found"
	ActorNameSizeError              = "Actor name size must be from 1 to 150"
	GrpcRecievError                 = "gRPC recieve error"
)
