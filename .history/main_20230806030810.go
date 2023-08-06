//package main
//
//import (
//	"log"
//	"os"
//
//	"github.com/gin-gonic/gin"
//)
//
//var appID, appCertificate string
//
//func main() {
//
//	appIDEnv, appIDExists := os.LookupEnv("9fc6de75e3a14ff7a77c16c7eb6bb767")
//	appCertEnv, appCertExists := os.LookupEnv("b3fab562898646baad8a7a9e5923b0b0")
//
//	if !appIDExists || !appCertExists {
//		log.Fatal("FATAL ERROR: ENV not properly configured, check APP_ID and APP_CERTIFICATE")
//	} else {
//		appID = appIDEnv
//		appCertificate = appCertEnv
//	}
//
//	api := gin.Default()
//
//	api.GET("/ping", func(c *gin.Context) {
//		c.JSON(200, gin.H{
//			"message": "pong",
//		})
//	})
//
//	api.Run(":8080")
//}

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	rtctokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	rtmtokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"
	"github.com/gin-gonic/gin"
)


var appID ,appCertificate string
func main() {


	appIDEnv , appIDExists := os.LookupEnv("APP_ID")
	appCertEnv , appCertExists := os.LookupEnv("APP_CERTIFICATE")

	if !appIDExists || !appCertExists{
		log.Fatal("FATAL ERROR: ENV not properly configured , check APP_ID and APP_CERTIFICATE")

	}else{
		appID = appIDEnv
		appCertificate = appCertEnv
	}
	api := gin.Default()

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api.GET("rtc/:channelName/:role/:tokenType/:uid/", getRtcToken)
	api.GET("rtm/:uid", getRtmToken)
	api.GET("rte/:channelName/:role/:tokenType/:uid/", getBothTokens)
	api.Run(":8080")

}

func getRtcToken(c *gin.Context) {

	channelName, tokenType, uidStr, role, expireTimestamp, err := parseRtcParams(c)

	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Error Generating RTC token :" + err.Error(),
			"status":  400,
		})
		return
	}

	rtcToken, tokenErr := generateRtcToken(channelName, uidStr, tokenType, role, expireTimestamp)

	if tokenErr != nil {
		log.Println(tokenErr)
		c.Error(err)
		c.AbortWithStatusJSON(400, gin.H{
			"status": 400,
			"error":  "Error Generating RTC token:" + tokenErr.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"rtcToken": rtcToken,
		})
	}

}

func getRtmToken(c *gin.Context) {
  uidStr, expireTimestamp, err:= parseRtmParams(c)

  if err!= nil{
	c.Error(err)
	c.AbortWithStatusJSON(400,gin.H{
		"status": 400,
		"message": "Error generating RTM token",
	})
	return
  }

  rtmToken, tokenErr:= rtmtokenbuilder2.BuildToken(appID,appCertificate,uidStr,rtmtokenbuilder2.RoleRtmUser,expireTimestamp)

  if tokenErr != nil{
	log.Println(err)
	c.Error(err)


	errMsg := "Error generating RTM token :" + tokenErr.Error()
	c.AbortWithStatusJSON(400,gin.H{
		"status":400,
		"error": errMsg
	})
  }else{
	c.JSON(200,gin.H{
		"rtmToken ": rtmToken,
	})
  }
}

func getBothTokens(c *gin.Context) {

}

func parseRtcParams(c *gin.Context) (channelName, tokenType, uidStr string, role rtctokenbuilder2.Role, expireTimestamp uint32, err error) {
	channelName = c.Param("channelName")
	roleStr := c.Param("role")
	tokenType = c.Param("tokenType")
	uidStr = c.Param("uid")
	expireTime := c.DefaultQuery("expiry", "3600")

	if roleStr == "publisher" {
		role = rtctokenbuilder2.RolePublisher
	} else {
		role = rtctokenbuilder2.RoleSubscriber
	}

	expireTime64, parseErr := strconv.ParseUint(expireTime, 10, 64)
	if parseErr != nil {
		err = fmt.Errorf("failed to parse expireTime")
	}
	expireTimeInSeconds := uint32(expireTime64)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp = currentTimestamp + expireTimeInSeconds

	return channelName, tokenType, uidStr, role, expireTimestamp, err

}

func parseRtmParams(c *gin.Context) (uidStr string, expireTimestamp uint32, err error) {

	  uidStr = c.Param("uid")
	  expireTime := c.DefaultQuery("expiry","3600")

	  expireTime64, parseErr := strconv.ParseUint(expireTime,10,64)
	  if parseErr != null{
		err = fmt.Errorf("failed to parse expireTime: %s,causing error:%s",expireTime,parseErr)
	  }
	  expireTimeInSeconds!= uint32(expireTime64)
	  currentTimestamp := uint32(time.Now().Unix())
	  expireTimestamp = currentTimestamp + expireTimeInSeconds


	  return uidStr,expireTimestamp,err
}

func generateRtcToken(channelName, uidStr, tokenType string, role rtctokenbuilder2.Role, expireTimestamp uint32) (rtcToken string, err error) {

	if tokenType == "userAccount" {
		rtcToken, err = rtctokenbuilder2.BuildTokenWithAccount(appID, appCertificate, channelName, uidStr, role, expireTimestamp)
	  
		 return rtcToken,err
	
	}else if tokenType == "uid"{
		uid64, parseErr := strconv.ParseUint(uidStr,10,64)
		if parseErr != nil{
			err = fmt.Errorf("failed to parse uidStr:%s, to uint causing error:%s",uidStr,parseErr)
		    return " ",err
		}
		uid := uint32(uid64)
		log.Printf("Building Token with uid: %d\n",uid)
		rtcToken, err = rtctokenbuilder2.BuildTokenWithUid(appID,appCertificate,channelName,uid,role,expireTimestamp)
		return rtcToken,err
	}else{
		err = fmt.Errorf("failed to generate RTC token for Unknown Tokentype: %s",tokenType)
		log.Println(err)
		return " ",err
	}
	
}
