package api

import (
	"context"
	"fmt"
	"gotrains/train_webs/train_web/global"
	"gotrains/train_webs/train_web/models"
	"gotrains/train_webs/train_web/proto"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err == nil {
		return
	}
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "参数错误",
			})
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"message": e.Message(),
			})
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "内部错误",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "其他错误",
			})
		}
		return
	}
	zap.S().Errorw("grpc 请求失败", "msg", err.Error())
	c.JSON(500, gin.H{
		"message": err.Error(),
	})
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": RemoveTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func getClaimes(c *gin.Context) *models.CustomClaims {
	claims, _ := c.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	return currentUser
}

// 获取列车信息
func GetTrains(c *gin.Context) {
	cUser := getClaimes(c)
	zap.S().Infof("访问用户: %d", cUser.ID)
	pn := c.DefaultQuery("page", "0")
	pnInt, _ := strconv.Atoi(pn)
	ps := c.DefaultQuery("size", "10")
	psInt, _ := strconv.Atoi(ps)
	// 获取所有列车信息,此处并不需要参数太多
	resp, err := global.TrainClient.GetTrainList(context.Background(), &proto.TrainPageInfo{
		Pn: uint32(pnInt),
		Ps: uint32(psInt),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	var trains []models.TrainInfo
	for _, value := range resp.Data {
		train := models.TrainInfo{
			Code:         value.Data.Code,
			Type:         value.Data.Type,
			StartStation: value.Data.StartStation,
			EndStation:   value.Data.EndStation,
			StartTime:    time.Unix(value.Data.StartTime, 0).Format("2006-01-02 15:04:05"),
			EndTime:      time.Unix(value.Data.EndTime, 0).Format("2006-01-02 15:04:05"),
			StartPinyin:  value.Data.StartPinyin,
			EndPinyin:    value.Data.EndPinyin,
		}
		trains = append(trains, train)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    trains,
	})
}

// 添加列车信息
func AddTrains(c *gin.Context) {
	cc := getClaimes(c)
	zap.S().Infof("访问用户: %d", cc.ID)
	var taininfo models.TrainInfo
	err2 := c.ShouldBindJSON(&taininfo)
	if err2 != nil {
		zap.S().Errorw("AddTrains 参数错误", "msg", err2.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	// 没时间了，暂时不检查了
	startTime, _ := time.Parse("2006-01-02 15:04:05", taininfo.StartTime)
	endTime, _ := time.Parse("2006-01-02 15:04:05", taininfo.EndTime)
	resp, err := global.TrainClient.CreateTrain(context.Background(), &proto.TrainRequest{
		Code:         taininfo.Code,
		Type:         taininfo.Type,
		StartStation: taininfo.StartStation,
		EndStation:   taininfo.EndStation,
		StartTime:    startTime.Unix(),
		EndTime:      endTime.Unix(),
		StartPinyin:  taininfo.StartPinyin,
		EndPinyin:    taininfo.EndPinyin,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "添加成功",
		"train_code": resp.Data.Code,
	})
}

// 获取某一个车次的所有车厢信息
func GetCarriages(c *gin.Context) {
	// TODO: 下次在写吧
}

// 添加车厢信息
func AddCarriages(c *gin.Context) {
	// TODO: 下次在写吧
}

func Decimal(value float32, prec int) float64 {
	format := "%." + strconv.Itoa(prec) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, value), 64)
	return res
}

// 获取余票信息
func GetTickets(c *gin.Context) {
	cc := getClaimes(c)
	zap.S().Infof("访问用户: %d", cc.ID)
	// 只通过日期查询
	// trainCode := c.Query("train_code")
	// 暂且不管
	startStation := c.Query("start")
	endStation := c.Query("end")
	startTime := c.Query("date")
	t, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		zap.S().Errorw("GetTickets 参数错误", "msg", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
			"success": false,
		})
		return
	}
	startTime = t.Format("2006-01-02")
	zap.S().Infof("start: %s, end: %s, date: %s", startStation, endStation, startTime)
	if startStation == "" || endStation == "" || startTime == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
			"success": false,
		})
		return
	}
	resp, err := global.TicketClient.GetTicketList(context.Background(), &proto.TicketRequest{TrainCode: "", Date: startTime, StartStation: startStation, EndStation: endStation})
	if err != nil {
		zap.S().Errorw("GetTickets 参数错误", "msg", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	var tickets []models.TicketInfo
	for _, value := range resp.Data {
		dt, _ := time.Parse("2006-01-02 15:04:05", value.Data.Date)
		st := time.Unix(value.Data.StartTime, 0)
		et := time.Unix(value.Data.EndTime, 0)
		ticket := models.TicketInfo{
			ID:               int64(value.Data.Id),
			TrainCode:        value.Data.TrainCode,
			Date:             dt.Format("2006-01-02"),
			Start:            value.Data.StartStation,
			StartPinyin:      value.Data.StartPinyin,
			StartTime:        st.Format("2006-01-02 15:04:05"),
			StartIndex:       int32(value.Data.StartIndex),
			End:              value.Data.EndStation,
			EndPinyin:        value.Data.EndPinyin,
			EndTime:          et.Format("2006-01-02 15:04:05"),
			EndIndex:         int32(value.Data.EndIndex),
			FirstClassLast:   int32(value.Data.FirstClassLast),
			SecondClassLast:  int32(value.Data.SecondClassLast),
			FirstClassPrice:  Decimal(value.Data.FirstClassPrice, 1),
			SecondClassPrice: Decimal(value.Data.SecondClassPrice, 1),
			SoftBerthLast:    int32(value.Data.SoftberthLast),
			SoftBerthPrice:   Decimal(value.Data.SoftberthPrice, 1),
			HardBerthLast:    int32(value.Data.HardberthLast),
			HardBerthPrice:   Decimal(value.Data.HardberthPrice, 1),
			Duration:         et.Sub(st).Hours(),
		}
		tickets = append(tickets, ticket)
	}
	msg := ""
	if len(tickets) == 0 {
		msg = "没有查询到余票信息"
	} else {
		msg = "查询成功"
	}
	c.JSON(http.StatusOK, gin.H{
		"message": msg,
		"content": gin.H{
			"total": len(tickets),
			"list":  tickets,
		},
		"success": true,
	})
}

// 获取为某一列车的座位信息
func GetSeatsByTrain(c *gin.Context) {
	cc := getClaimes(c)
	zap.S().Infof("访问用户: %d", cc.ID)
	trainCode := c.Query("train_code")
	if trainCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	resp, err := global.SeatClient.GetSeatList(context.Background(), &proto.SeatPageInfo{
		Seat: &proto.SeatInfo{
			TrainCode: trainCode,
		},
		Date: time.Now().Format("2006-01-02"),
	})

	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	var seats []models.SeatInfo
	for _, value := range resp.Data {
		zap.S().Infof("座位信息: %v", value)
		seat := models.SeatInfo{
			TrainCode:     value.Data.TrainCode,
			CarriageIndex: value.Data.CarriageIndex,
			SeatType:      value.Data.SeatType,
			SeatIndex:     value.Data.SeatIndex,
			Row:           value.Data.Row,
			Column:        value.Data.Column,
		}
		seats = append(seats, seat)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取成功",
		"data":    seats,
	})
}

// 获取车站信息

func GetStations(c *gin.Context) {
	cc := getClaimes(c)
	zap.S().Infof("访问用户: %d", cc.ID)
	resp, err := global.StationClient.GetStationList(context.Background(), &proto.StationPageInfo{Pn: 1, Ps: 100})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	var resplist []models.StationInfo
	for _, value := range resp.Data {
		stationinfo := models.StationInfo{
			Id:     int32(value.Data.Id),
			Name:   value.Data.Name,
			Pinyin: value.Data.Pinyin,
			Py:     value.Data.FirstLetter,
		}
		resplist = append(resplist, stationinfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "station",
		"success": true,
		"content": resplist,
	})
}
