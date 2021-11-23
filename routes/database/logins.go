package database

import (
	"context"
	"net/http"

	"dahbura.me/api/config"
	"dahbura.me/api/database/models"
	"dahbura.me/api/database/mongodb"
	httppkg "dahbura.me/api/util/http"
	"dahbura.me/api/util/validation"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/crypto/bcrypt"
)

func Logins(c *gin.Context) {
	db, err := mongodb.GetMongoDb()
	if httppkg.HandleError(c, err) {
		return
	}

	login := models.Login{}
	err = c.ShouldBindJSON(&login)
	if httppkg.HandleError(c, err) {
		return
	}

	validate := validation.GetValidator()

	err = validate.Struct(login)
	if httppkg.HandleError(c, err) {
		return
	}

	filter := bson.M{"email": login.Email}
	projection := bson.M{
		"password": 0,
	}
	opts := options.FindOneOptions{
		Projection: &projection,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultCtxTimeout)
	defer cancel()

	user := models.User{}
	err = db.Collection("users").FindOne(ctx, filter, &opts).Decode(&user)
	if httppkg.HandleError(c, err) {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		c.Status(http.StatusUnauthorized)
		return
	}
	if httppkg.HandleError(c, err) {
		return
	}

	user.PasswordHash = ""

	c.Header("Content-Type", config.MimeApplicationJson)
	c.JSON(http.StatusOK, user)
}
