package main

import (
	_ "context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

type authSuccessResponse struct {
	Response struct {
		Params struct {
			Result          string `json:"result"`
			OwnerSteamid    string `json:"ownersteamid"`
			Steamid         string `json:"steamid"`
			Publisherbanned bool   `json:"publisherbanned"`
			Vacbanned       bool   `json:"vacbanned"`
		} `json:"params"`
	} `json:"response"`
}

type authFailureResponse struct {
	Response struct {
		Error struct {
			Errorcode int64  `json:"errorcode"`
			Errordesc string `json:"errordesc"`
		} `json:"error"`
	} `json:"reponse"`
}

func main() {
	router := gin.Default()
	router.GET("/auth", authenticateUser)
	router.Run(":8282")
}

func authenticateUser(c *gin.Context) {
	appid := c.Query("appid")
	key := c.Query("key")
	ticket := c.Query("ticket")
	steamRoute := "https://partner.steam-api.com/ISteamUserAuth/AuthenticateUserTicket/v1?"
	resp, err := http.Get(steamRoute + "key=" + key + "&appid=" + appid + "&ticket=" + ticket)
	if err != nil {
		log.Println("couldn't reach steam auth api")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//log.Println(err)
	}
	var response *authSuccessResponse = &authSuccessResponse{}
	if err = json.Unmarshal(body, response); err != nil || response.Response.Params.Steamid == "" {
		var failResponse *authFailureResponse = &authFailureResponse{}
		if err = json.Unmarshal(body, failResponse); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "FAIL", "message": failResponse.Response.Error.Errordesc})
			return
		}
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "steamID": response.Response.Params.Steamid})
}
