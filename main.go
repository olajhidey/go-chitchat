package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

const API_URL = "https://YOUR_SPACE.signalwire.com/api/chat/tokens"

type ApiResponse struct {
	Token string `json:"token"`
}

type Channels struct {
	Read  bool `json:"read"`
	Write bool `json:"write,omitempty"`
}

type ApiRequestBody struct {
	Ttl      int            `json:"ttl"`
	MemberId string         `json:"member_id"`
	Channels string         `json:"channels"`
	State    map[string]any `json:"state,omitempty"`
}

type RequestBody struct {
	Ttl      int                 `json:"ttl"`
	MemberId string              `json:"member_id"`
	State    map[string]any      `json:"state,omitempty"`
	Channels map[string]Channels `json:"channels"`
}

// Helper function to get environment variable data
func getEnv(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading env file")
	}
	return os.Getenv(key)
}

func TokenHandler(ctx *gin.Context) {

	var apiRequest ApiRequestBody

	// Get environmental variables
	spaceUrl := getEnv("SPACE_URL")
	apiToken := getEnv("API_TOKEN")
	projectId := getEnv("PROJECT_ID")

	baseUrl := "https://" + spaceUrl + "/api/chat/tokens"

	//initiate resty
	client := resty.New()

	if err := ctx.ShouldBindBodyWithJSON(&apiRequest); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	fmt.Printf("Request body: %v\n", apiRequest)

	channels := make(map[string]Channels)

	channels[apiRequest.Channels] = Channels{
		Read:  true,
		Write: true,
	}

	request := RequestBody{
		Ttl:      apiRequest.Ttl,
		MemberId: apiRequest.MemberId,
		Channels: channels,
		State: apiRequest.State,
	}

	response, err := client.R().
		SetBasicAuth(projectId, apiToken).
		SetBody(request).
		SetResult(&ApiResponse{}).
		Post(baseUrl)

	if err != nil {
		fmt.Printf("Error making POST request: %v", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	apiResponse := response.Result().(*ApiResponse)

	ctx.JSON(http.StatusOK, gin.H{
		"token": apiResponse.Token,
	})

}

func main() {

	// Initiate Gin
	router := gin.Default()

	router.Static("/static", "./www")

	// HTML Rendering
	router.GET("/", func(ctx *gin.Context) {
		ctx.File("./www/index.html")
	})

	api := router.Group("/api")

	api.POST("/token", TokenHandler)

	router.Run(":8080")

}
