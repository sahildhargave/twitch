package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

   "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	rtmtokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"
	//rtmtokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"
	"github.com/gin-gonic/gin"
)

var appID, appCertificate string

func main() {

	appIDEnv, appIDExists := os.LookupEnv("APP_ID")
	appCertEnv, appCertExists := os.LookupEnv("APP_CERTIFICATE")

	if !appIDExists || !appCertExists {
		log.Fatal("FATAL ERROR: ENV not properly configured, check APP_ID and APP_CERTIFICATE")
	} else {
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
	api.GET("rtm/:uid/", getRtmToken)
	api.GET("rte/:channelName/:role/:tokenType/:uid/", getBothTokens)

	api.Run(":8080")
}

func getRtcToken(c *gin.Context) {

}

func getRtmToken(c *gin.Context) {

}

func getBothTokens(c *gin.Context) {
log.Printf("dual token\n")
  // get rtc param values
  channelName, tokentype, uidStr, role, expireTimestamp, rtcParamErr := parseRtcParams(c)

  if rtcParamErr != nil {
    c.Error(rtcParamErr)
    c.AbortWithStatusJSON(400, gin.H{
      "message": "Error Generating RTC token: " + rtcParamErr.Error(),
      "status":  400,
    })
    return
  }
  // generate the rtcToken
  rtcToken, rtcTokenErr := generateRtcToken(channelName, uidStr, tokentype, role, expireTimestamp)
  // generate rtmToken.
  rtmToken, rtmTokenErr := rtmtokenbuilder.BuildToken(appID, appCertificate, uidStr,rtmtokenbuilder.Role, expireTimestamp)

  if rtcTokenErr != nil {
    log.Println(rtcTokenErr) // token failed to generate
    c.Error(rtcTokenErr)
    errMsg := "Error Generating RTC token - " + rtcTokenErr.Error()
    c.AbortWithStatusJSON(400, gin.H{
      "status": 400,
      "error":  errMsg,
    })
  } else if rtmTokenErr != nil {
    log.Println(rtmTokenErr) // token failed to generate
    c.Error(rtmTokenErr)
    errMsg := "Error Generating RTC token - " + rtmTokenErr.Error()
    c.AbortWithStatusJSON(400, gin.H{
      "status": 400,
      "error":  errMsg,
    })
  } else {
    log.Println("RTC Token generated")
    c.JSON(200, gin.H{
      "rtcToken": rtcToken,
      "rtmToken": rtmToken,
    })
  }
}

func parseRtcParams(c *gin.Context) (channelName, tokentype, uidStr string, role rtctokenbuilder2.Role, expireTimestamp uint32, err error) {
channelName = c.Param("channelName")
  roleStr := c.Param("role")
  tokentype = c.Param("tokentype")
  uidStr = c.Param("uid")
  expireTime := c.DefaultQuery("expiry", "3600")

  if roleStr == "publisher" {
    role = rtctokenbuilder2.RolePublisher
  } else {
    role = rtctokenbuilder2.RoleSubscriber
  }

  expireTime64, parseErr := strconv.ParseUint(expireTime, 10, 64)
  if parseErr != nil {
    // if string conversion fails return an error
    err = fmt.Errorf("failed to parse expireTime: %s, causing error: %s", expireTime, parseErr)
  }

  // set timestamps
  expireTimeInSeconds := uint32(expireTime64)
  currentTimestamp := uint32(time.Now().UTC().Unix())
  expireTimestamp = currentTimestamp + expireTimeInSeconds

  return channelName, tokentype, uidStr, role, expireTimestamp, err
}

func parseRtmParams(c *gin.Context) (uidStr string, expireTimestamp uint32, err error) {
 uidStr = c.Param("uid")
  expireTime := c.DefaultQuery("expiry", "3600")

  expireTime64, parseErr := strconv.ParseUint(expireTime, 10, 64)
  if parseErr != nil {
    // if string conversion fails return an error
    err = fmt.Errorf("failed to parse expireTime: %s, causing error: %s", expireTime, parseErr)
  }

  // set timestamps
  expireTimeInSeconds := uint32(expireTime64)
  currentTimestamp := uint32(time.Now().UTC().Unix())
  expireTimestamp = currentTimestamp + expireTimeInSeconds

  // check if string conversion fails
  return uidStr, expireTimestamp, err
}

func generateRtcToken(channelName, uidStr, tokentype string, role rtctokenbuilder2.Role, expireTimestamp uint32) (rtcToken string, err error) {

  if tokentype == "userAccount" {
    log.Printf("Building Token with userAccount: %s\n", uidStr)
    rtcToken, err = rtctokenbuilder2.BuildTokenWithAccount(appID, appCertificate, channelName, uidStr, role, expireTimestamp)
    return rtcToken, err

  } else if tokentype == "uid" {
    uid64, parseErr := strconv.ParseUint(uidStr, 10, 64)
    // check if conversion fails
    if parseErr != nil {
      err = fmt.Errorf("failed to parse uidStr: %s, to uint causing error: %s", uidStr, parseErr)
      return "", err
    }

    uid := uint32(uid64) // convert uid from uint64 to uint 32
    log.Printf("Building Token with uid: %d\n", uid)
    rtcToken, err = rtctokenbuilder2.BuildTokenWithUid(appID, appCertificate, channelName, uid, role, expireTimestamp)
    return rtcToken, err

  } else {
    err = fmt.Errorf("failed to generate RTC token for Unknown Tokentype: %s", tokentype)
    log.Println(err)
    return "", err
  }
}
