package server
import (
	"github.com/gin-gonic/gin"
	"golangdynamocrud/services"
	"golangdynamocrud/model"
	"net/http"
)

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "server up and runnning",
	})
}


func addNewUser(c *gin.Context) {
	var newUserDetail model.Detalhe
	if err := c.ShouldBindJSON(&newUserDetail); err != nil {
		errMsg := "Fields idade, nome and profissao are required"
		c.JSON(http.StatusBadRequest, gin.H{
			"data": errMsg,
		})
		return
	}

	dynamoMapForInclusion, err := services.MarshalMapForAttributes(newUserDetail)
	if err != nil {
		errMsg := "Error mapping to dynamo attributes"
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": errMsg,
		})
		return
	}
	input := services.AssembleDynamoItem(dynamoMapForInclusion)
	
	_, err = services.GetDynamoInstance().PutItem(input)
	if err != nil {
		errMsg := "There was an error to send data to dynamodb instance"
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": errMsg,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": "Item created with success",
	})
}

func getAllUsers(c *gin.Context) {
	lista, err := services.AssembleUsersList()
	if err != nil {
		errMsg := "There was an error to send data to dynamodb instance"
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": errMsg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": lista,
	})
}

func getItemById(c *gin.Context) {
	identifier := c.Param("id")
	result, err := services.GetItemById(identifier)
	if err != nil {
		errMsg := "There was an error to send data to dynamodb instance"
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": errMsg,
		})
		return
	}
	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"data": "Not found",
		})
		return
	}
	item, err := services.AssembleUserItem(result)
	if err != nil {
		errMsg := "There was an error to unmarshall item to model instance"
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": errMsg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

func deleteItemById(c *gin.Context) {
	err := services.DeleteItemById(c.Param("id"))
	if err != nil {
		errMsg := "There was an error to unmarshall item to model instance"
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": errMsg,
		})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func updateItemById(c *gin.Context) {
	var detalhe model.Detalhe
	err := c.ShouldBindJSON(&detalhe)
	if err != nil {
		errMsg := "The fields idade, nome and profissao are required"
		c.JSON(http.StatusBadRequest, gin.H{
			"data": errMsg,
		})
		return
	}
	err = services.UpdateItemById(c.Param("id"), detalhe)
	if err != nil {
		errMsg := "There was an error to unmarshall item to model instance"
		c.JSON(http.StatusInternalServerError, gin.H{
			"data": errMsg,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": "Data updated with success",
	})
}