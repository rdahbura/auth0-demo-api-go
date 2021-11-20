package database

import (
	"context"
	"net/http"
	"time"

	"dahbura.me/api/config"
	"dahbura.me/api/database/models"
	"dahbura.me/api/database/mongodb"
	httppkg "dahbura.me/api/util/http"
	"dahbura.me/api/util/validation"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	client, err := mongodb.GetMongoClient()
	if httppkg.HandleError(c, err) {
		return
	}

	user := models.User{}
	err = c.ShouldBindJSON(&user)
	if httppkg.HandleError(c, err) {
		return
	}

	validate := validation.GetValidator()

	err = validate.Struct(user)
	if httppkg.HandleError(c, err) {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if httppkg.HandleError(c, err) {
		return
	}

	user.Password = ""
	user.PasswordHash = string(hash)

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultCtxTimeout)
	defer cancel()

	collection := client.Database(config.MongoDb).Collection("users")
	result, err := collection.InsertOne(ctx, user)
	if httppkg.HandleError(c, err) {
		return
	}

	id, _ := result.InsertedID.(primitive.ObjectID)
	user.Id = id
	user.PasswordHash = ""

	c.Header("Content-Type", config.MimeApplicationJson)
	c.JSON(http.StatusCreated, user)
}

func DeleteUser(c *gin.Context) {
	client, err := mongodb.GetMongoClient()
	if httppkg.HandleError(c, err) {
		return
	}

	objectId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if httppkg.HandleError(c, err) {
		return
	}

	filter := bson.M{"_id": objectId}

	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultCtxTimeout)
	defer cancel()

	collection := client.Database(config.MongoDb).Collection("users")
	result, err := collection.DeleteOne(ctx, filter)
	if httppkg.HandleError(c, err) {
		return
	}
	if result.DeletedCount == 0 {
		c.Status(http.StatusNotFound)
		return
	}

	c.Header("Content-Type", config.MimeApplicationJson)
	c.Status(http.StatusNoContent)
}

func GetUsers(c *gin.Context) {
	client, err := mongodb.GetMongoClient()
	if httppkg.HandleError(c, err) {
		return
	}

	var filter primitive.M
	if email := c.Query("email"); email != "" {
		filter = bson.M{"email": email}
	} else {
		filter = bson.M{}
	}

	projection := bson.M{
		"password":      0,
		"password_hash": 0,
	}
	opts := options.FindOptions{
		Projection: &projection,
	}

	ctxFind, cancelFind := context.WithTimeout(context.Background(), config.DefaultCtxTimeout)
	defer cancelFind()

	collection := client.Database(config.MongoDb).Collection("users")
	cursor, err := collection.Find(ctxFind, filter, &opts)
	if httppkg.HandleError(c, err) {
		return
	}

	ctxCursor, cancelCursor := context.WithTimeout(context.Background(), config.DefaultCtxTimeout)
	defer cancelCursor()

	users := []models.User{}
	err = cursor.All(ctxCursor, &users)
	if httppkg.HandleError(c, err) {
		return
	}

	c.Header("Content-Type", config.MimeApplicationJson)
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	client, err := mongodb.GetMongoClient()
	if httppkg.HandleError(c, err) {
		return
	}

	objectId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if httppkg.HandleError(c, err) {
		return
	}

	filter := bson.M{"_id": objectId}
	projection := bson.M{
		"password":      0,
		"password_hash": 0,
	}
	opts := options.FindOneOptions{
		Projection: &projection,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultCtxTimeout)
	defer cancel()

	user := models.User{}
	collection := client.Database(config.MongoDb).Collection("users")
	err = collection.FindOne(ctx, filter, &opts).Decode(&user)
	if err == mongo.ErrNoDocuments {
		c.Status(http.StatusNotFound)
		return
	}
	if httppkg.HandleError(c, err) {
		return
	}

	c.Header("Content-Type", config.MimeApplicationJson)
	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	client, err := mongodb.GetMongoClient()
	if httppkg.HandleError(c, err) {
		return
	}

	objectId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if httppkg.HandleError(c, err) {
		return
	}

	user := models.User{}
	err = c.ShouldBindJSON(&user)
	if httppkg.HandleError(c, err) {
		return
	}

	validate := validation.GetValidator()

	err = validate.Struct(user)
	if httppkg.HandleError(c, err) {
		return
	}

	if user.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if httppkg.HandleError(c, err) {
			return
		}

		user.Password = ""
		user.PasswordHash = string(hash)
	}

	now := time.Now()
	user.UpdatedAt = now

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": user}
	projection := bson.M{
		"password":      0,
		"password_hash": 0,
	}
	after := options.After
	opts := options.FindOneAndUpdateOptions{
		Projection:     &projection,
		ReturnDocument: &after,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultCtxTimeout)
	defer cancel()

	collection := client.Database(config.MongoDb).Collection("users")
	result := collection.FindOneAndUpdate(ctx, filter, update, &opts)
	if httppkg.HandleError(c, result.Err()) {
		return
	}

	updatedUser := models.User{}
	err = result.Decode(&updatedUser)
	if httppkg.HandleError(c, err) {
		return
	}

	c.Header("Content-Type", config.MimeApplicationJson)
	c.JSON(http.StatusOK, updatedUser)
}
