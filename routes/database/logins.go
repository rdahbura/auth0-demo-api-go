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
	client, err := mongodb.GetMongoClient()
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

	user := models.User{}
	collection := client.Database(config.MongoDb).Collection("users")
	err = collection.FindOne(context.TODO(), filter, &opts).Decode(&user)
	if httppkg.HandleError(c, err) {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password))
	if httppkg.HandleError(c, err) {
		return
	}

	user.PasswordHash = ""

	c.Header("Content-Type", config.MimeApplicationJson)
	c.JSON(http.StatusOK, user)
}
