package server
import (
	"github.com/gin-gonic/gin"
	"golangdynamocrud/services"
	"golangdynamocrud/model"
	"github.com/google/uuid"
	"net/http"
)

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "server up and runnning",
	})
}
func generateIndexForDynamo() string {
	return "users:"+uuid.New().String()
}

func addNewUser(c *gin.Context) {
	var newUserDetail model.Detalhe
	err := c.ShouldBindJSON(&newUserDetail)
	if err != nil {
		errMsg := "The fields idade, nome and profissao are required"
		c.JSON(http.StatusBadRequest, gin.H{
			"data": errMsg,
		})
		return
	}

	newUser := model.Usuario{
		Identificador: generateIndexForDynamo(),
		Detalhes: newUserDetail,
	}

	dynamoMapForInclusion, err := services.MarshalMapForAttributes(newUser)
	if err != nil {
		errMsg := "The fields idade, nome and profissao are required"
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
	c.Header("URI", newUser.Identificador)
	c.JSON(http.StatusCreated, gin.H{
		"data": "Item created with success",
		"Uri": "/users/"+newUser.Identificador,
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
