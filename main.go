package main

import (
	"blackHole/Logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var logger = Logger.GetLogger("blackhole")

var countNum = 0

type blackHoleParams struct {
	ReqCode string `json:"reqCode" binding:"required"`
	ReqStartTime int64 `json:"reqStartTime" binding:"required"`
}

type blackHoleData struct {
	ReqCode string `json:"reqCode" binding:"required"`
	ReqStartTime int64 `json:"reqStartTime" binding:"required"`
	ReqArriveTime int64 `json:"reqArriveTime" binding:"required"`
	ResStartTime int64 `json:"resStartTime" binding:"required"`
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "get request!")
	})

	// blackhole for performance test
	r.GET("/blackhole", func(c *gin.Context) {
		countNum += 1
		reqArriveTime := time.Now().UnixNano()
		reqCode := c.DefaultQuery("reqCode", "0")
		reqStartTimeStr := c.DefaultQuery("reqStartTime", "1562454427000")  //默认时间 2019-07-07 07:07:07
		reqStartTime, err := strconv.ParseInt(reqStartTimeStr, 10, 64)
		if err != nil {
			logger.Infof("reqStartTime参数转换失败! err=%v", err)
			c.JSON(500, gin.H{"code": 999999, "message": "reqStartTime参数转换失败!", "data": err})
			return
		}
		//sleepTime := rand.Intn(5)
		//time.Sleep(time.Second * time.Duration(sleepTime))
		resData := blackHoleData{
			ReqCode:reqCode,
			ReqStartTime: reqStartTime,
			ReqArriveTime: reqArriveTime,
			ResStartTime: time.Now().UnixNano(),
		}
		c.JSON(200, gin.H{
			"code": "000000",
			"message": "操作成功!",
			"data": resData,
		})
		defer logger.Infof("累计请求数=%d", countNum)
	})

	r.POST("/blackhole", func(c *gin.Context) {
		countNum += 1
		var reqData blackHoleParams
		reqArriveTime := time.Now().UnixNano()
		err := c.ShouldBind(&reqData)
		if err != nil {
			logger.Infof("Blackhole解析参数失败! err=%v", err)
			c.JSON(500, gin.H{"code": 999999, "message": "Blackhole解析参数失败!", "data": err})
			return
		} else {
			logger.Infof("reqCode=%s, reqStartTime=%v", reqData.ReqCode, reqData.ReqStartTime)
		}
		//sleepTime := rand.Intn(5)
		//time.Sleep(time.Second * time.Duration(sleepTime))
		resData := blackHoleData{
			ReqCode:reqData.ReqCode,
			ReqStartTime: reqData.ReqStartTime,
			ReqArriveTime: reqArriveTime,
			ResStartTime: time.Now().UnixNano(),
		}
		c.JSON(200, gin.H{
			"code": "000000",
			"message": "操作成功!",
			"data": resData,
		})
		defer logger.Infof("累计请求数=%d", countNum)
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	err := r.Run(":9527")
	if err != nil{
		logger.Fatalf("服务器异常退出! err=%v\n", err)
	}
}
