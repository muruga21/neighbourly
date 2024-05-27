package controllers

import (
	"context"
	"time"

	"github.com/Kishore-kumar-07/neighbourly/backend/config"
	"github.com/Kishore-kumar-07/neighbourly/backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func SetProvider (c *gin.Context) {
	var serviceProviderModel models.ServiceProviderModel;
	serviceProviderModel.SetDefaults();

	client := config.Client;

	if err := c.BindJSON(&serviceProviderModel); err != nil {
		c.JSON(400, gin.H{
			"error": true,
			"message": err.Error(),
		})
		return
	}

	if( serviceProviderModel.Description == "" || serviceProviderModel.Experience == "" || serviceProviderModel.ServiceDescription == ""){
		c.JSON(400, gin.H{
			"error": true,
			"message": "All fields are required",
		})
		return
	}

	serviceProviderModel.Email = c.GetString("email");
	if serviceProviderModel.Email == "" {
		c.JSON(400, gin.H{
			"error": true,
			"message": "Email is missing",
		})
		return
	}

	collection := client.Database("muruga").Collection("serviceProviders");
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, err := collection.InsertOne(ctx, serviceProviderModel);

	if err != nil {
		c.JSON(500, gin.H{
			"error": true,
			"message": "Error while inserting service provider",
		})
		return
	}

	c.JSON(200, gin.H{
		"error": false,
		"message": "Service provider added successfully",
	})
}

func TopRatedProviders (c * gin.Context){
	var serviceProviderModel models.ServiceProviderModel;
	serviceProviderModel.SetDefaults()

	client := config.Client;

	collection := client.Database("muruga").Collection("serviceProviders");
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second);
	defer cancel();

	// findOptions := options.Find().SetSort(map[string]int{"rating": -1}).SetLimit(5)

	doc := collection.FindOne(ctx, bson.M{"rating":0});

	if err:= doc.Decode(&serviceProviderModel); err != nil {
		c.JSON(500, gin.H{"error":true, "message": "couldnt find providers"});
		return
	}

	c.JSON(200, gin.H{
		"error": false,
		"message": serviceProviderModel,
	});
}