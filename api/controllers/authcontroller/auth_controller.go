package authcontroller

import (
	"encoding/json"
	"example/auth/domain"
	"example/auth/pkg"
	"example/auth/usecase"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"example/auth/config"
)

type UserController struct {
	userUseCase *usecase.UserUseCase
}

func NewController(userUseCase *usecase.UserUseCase) *UserController {
	return &UserController{userUseCase: userUseCase}
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {

	// mengambil inputan json
	var userInput domain.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		pkg.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	username, err := uc.userUseCase.LoginUser(userInput.Username, userInput.Password)

	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response := map[string]string{"message": "Wrong username or password"}
			pkg.ResponseJSON(w, http.StatusUnauthorized, response)
			return
		default:
			response := map[string]string{"message": err.Error()}
			pkg.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		}
	}

	// proses pembuatan token jwt
	expTime := time.Now().Add(time.Minute * 1)
	claims := &config.JWTClaim{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-jwt-mux",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	// medeklarasikan algoritma yang akan digunakan untuk signing
	tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// signed token
	token, err := tokenAlgo.SignedString(config.JWT_KEY)
	if err != nil {
		response := map[string]string{"message": err.Error()}
		pkg.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	// set token yang ke cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	response := map[string]string{"message": "Login Success as " + username}
	pkg.ResponseJSON(w, http.StatusOK, response)

}

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {

	// mengambil inputan json
	var userInput domain.User
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userInput); err != nil {
		response := map[string]string{"message": err.Error()}
		pkg.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}
	defer r.Body.Close()

	user, err := uc.userUseCase.RegisterUser(&userInput)

	if err != nil {
		pkg.ResponseJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]string{"message": "Success", "name": user.Name}
	pkg.ResponseJSON(w, http.StatusOK, response)

}

func (uc *UserController) ProfilUser(w http.ResponseWriter, r *http.Request) {

	user, err := uc.userUseCase.ProfilUser()

	profileUsers := make([]domain.ProfilUser, len(*user))

	for i, item := range *user {
		profileUsers[i] = domain.ProfilUser{
			Id:       item.Id,
			Username: item.Username,
			Name:     item.Name,
		}
	}

	if err != nil {
		pkg.ResponseJSON(w, http.StatusNotFound, err.Error())
		return
	}

	response := map[string][]domain.ProfilUser{"data": profileUsers}
	pkg.ResponseJSON(w, http.StatusOK, response)
}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]string{"message": "Logout Success"}
	pkg.ResponseJSON(w, http.StatusOK, response)
}
