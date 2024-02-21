package router

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/like2foxes/nirlir/pg"
	"golang.org/x/crypto/bcrypt"
)

func getLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "login", nil)
}

func getRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "register", nil)
}

func postLogin(q *pg.Queries, ctx context.Context, jwtSecret string) func(c echo.Context) error {
	return func(c echo.Context) error {
		password := c.FormValue("password")
		email := c.FormValue("email")

		if (email == "") || (password == "") {
			return c.String(http.StatusBadRequest, "invalid credentials")
		}

		log.Println("email: ", email)
		log.Println("password: ", password)

		hashedPass, err := q.GetUserPassword(ctx, email)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		user, err := q.GetUser(ctx, pg.GetUserParams{Email: email, Password: hashedPass})
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		if user.ID == 0 {
			return c.String(http.StatusUnauthorized, "invalid credentials")
		}

		if !checkPasswordHash(password, hashedPass) {
			log.Println("invalid password")
			return c.String(http.StatusUnauthorized, "invalid credentials")
		}

		claims := &jwtCustomClaims{
			user.Email,
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			log.Println("signed string: ", err)
			return c.String(http.StatusInternalServerError, "Error generating token")
		}
		writeCookie(c, t)
		return c.Redirect(http.StatusSeeOther, "/")
	}
}

func writeCookie(c echo.Context, token string) {
	cookie := new(http.Cookie)
	cookie.Name = "jwt_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"

	cookie.HttpOnly = true
	c.SetCookie(cookie)
}

func readCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie("token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func postRegister(q *pg.Queries, ctx context.Context) func(c echo.Context) error {
	return func(c echo.Context) error {
		password := c.FormValue("password")
		email := c.FormValue("email")

		hashed, err := hashPassword(password)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		_, err = q.CreateUser(ctx, pg.CreateUserParams{Email: email, Password: hashed})
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "internal error")
		}

		return c.Redirect(http.StatusSeeOther, "/login")
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type jwtCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
