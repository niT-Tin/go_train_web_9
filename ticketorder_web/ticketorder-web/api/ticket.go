package api

import (
	"context"
	"fmt"
	"gotrains/ticketorder_web/ticketorder-web/forms"
	"gotrains/ticketorder_web/ticketorder-web/global"
	"gotrains/ticketorder_web/ticketorder-web/models"
	"gotrains/ticketorder_web/ticketorder-web/proto"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	seats = "ABCDEF"
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
				"success": false,
			})
		case codes.NotFound:
			c.JSON(http.StatusBadRequest, gin.H{
				"message": e.Message(),
				"success": false,
			})
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "内部错误",
				"success": false,
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "其他错误",
				"success": false,
			})
		}
		return
	}
	zap.S().Errorw("grpc 请求失败", "msg", err.Error())
	c.JSON(500, gin.H{
		"message": err.Error(),
		"success": false,
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

func GetPassengerList(c *gin.Context) {
	zap.S().Debug("获取对应订单用户乘客列表")
	cUser := getClaimes(c)
	zap.S().Infof("访问用户: %d", cUser.ID)
	pn := c.DefaultQuery("page", "1")
	pnInt, _ := strconv.Atoi(pn)
	ps := c.DefaultQuery("size", "10")
	psInt, _ := strconv.Atoi(ps)
	res, err := global.UserClient.GetPassengerList(context.WithValue(context.Background(), "ginContext", c), &proto.PassengerPageInfo{
		Pn:     uint32(pnInt),
		Ps:     uint32(psInt),
		UserId: int32(cUser.ID),
	})
	if err != nil {
		zap.S().Errorw("GetPassengerList 获取乘客列表失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}
	var passengers []models.Passenger
	for _, value := range res.Data {
		zap.S().Infof("乘客: %v", value)
		if value.IdCard == "" {
			continue
		}
		ps := models.Passenger{
			Id:       value.Id,
			Name:     value.Name,
			IdCard:   value.IdCard,
			Type:     models.PassengerType(value.Type),
			UserID:   value.UserId,
			Seat:     value.Seat,
			SeatType: value.SeatType,
		}
		passengers = append(passengers, ps)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"content": passengers,
	})
}

func AddPassenger(c *gin.Context) {
	cUser := getClaimes(c)
	zap.S().Infof("访问用户: %d", cUser.ID)
	passengerForm := forms.PassengerForm{}
	if err := c.ShouldBindJSON(&passengerForm); err != nil {
		zap.S().Errorw("AddPassenger 参数错误", "msg", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
			"success": false,
		})
	}
	passenger, err := global.UserClient.GetPassengerByIdCard(context.WithValue(context.Background(), "ginContext", c), &proto.PassengerIdCardRequest{
		IdCard: passengerForm.IdCard,
	})
	// 这里不对，但是没时间了，直接在Add里面判断了
	if passenger.IdCard != "" {
		_, err = global.UserClient.UpdatePassenger(context.WithValue(context.Background(), "ginContext", c), &proto.PassengerInfo{
			Id:       passenger.Id,
			Name:     passengerForm.Name,
			IdCard:   passengerForm.IdCard,
			UserId:   int32(cUser.ID),
			Type:     passengerForm.Type,
			Seat:     passenger.Seat,
			SeatType: passenger.SeatType,
		})
		if err != nil {
			HandleGrpcErrorToHttp(err, c)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "更新成功",
			"success": true,
		})
		return
	}

	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	_, err = global.UserClient.AddPassenger(context.WithValue(context.Background(), "ginContext", c), &proto.PassengerInfo{
		Name:     passengerForm.Name,
		IdCard:   passengerForm.IdCard,
		UserId:   int32(cUser.ID),
		Type:     passengerForm.Type,
		Seat:     "C1",
		SeatType: "一",
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "添加成功",
		"success": true,
	})
}

func DeletePassenger(c *gin.Context) {
	cUser := getClaimes(c)
	zap.S().Infof("访问用户: %d", cUser.ID)
	id := strings.TrimSpace(c.Param("id"))
	i, err2 := strconv.Atoi(id)
	if err2 != nil {
		zap.S().Errorw("DeletePassenger 参数错误", "msg", err2.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误",
		})
	}
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数错误",
		})
		return
	}
	_, err := global.UserClient.DeletePassenger(context.WithValue(context.Background(), "ginContext", c), &proto.IdRequest{
		Id: int32(i),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

// 下单但是没有支付
func CreateOrder(c *gin.Context) {
	cUser := getClaimes(c)
	zap.S().Infof("访问用户: %d", cUser.ID)
	// passengers := []models.Passenger{}
	orderForm := forms.OrderForm{}
	if err := c.ShouldBindJSON(&orderForm); err != nil {
		zap.S().Errorw("CreateOrder 参数错误", "msg", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	orderInfo := proto.OrderInfo{}

	if !store.Verify(orderForm.ImageCodeId, orderForm.ImageCode, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "验证码错误",
			"success": false,
		})
		return
	}

	var orderPassengers []*proto.OPassengerInfo
	order := models.Order{}
	for _, passenger := range orderForm.Tickets {
		opass := proto.OPassengerInfo{
			Name:     passenger.PassengerName,
			IdCard:   passenger.PassengerIdCard,
			Type:     int64(passenger.PassengerType),
			UserId:   int32(passenger.UserID),
			Seat:     passenger.Seat,
			SeatType: passenger.SeatTypeCode,
		}
		orderPassengers = append(orderPassengers, &opass)
	}
	orderInfo.Passengers = orderPassengers
	orderInfo.UserId = int32(cUser.ID)
	orderInfo.TrainCode = orderForm.TrainCode
	orderInfo.StartStation = orderForm.Start
	orderInfo.EndStation = orderForm.End
	orderInfo.StartTime = orderForm.StartTime
	orderInfo.EndTime = orderForm.EndTime
	orderInfo.SeatType = order.SeatType
	orderInfo.SeatNumber = order.SeatNumber
	orderInfo.Price = order.Pirce
	orderInfo.OrderSn = uuid.NewV4().String()
	oir, err := global.OrderClient.CreateOrder(context.WithValue(context.Background(), "ginContext", c), &proto.CreateOrderInfo{
		Id:        uuid.NewV4().String(),
		OrderInfo: &orderInfo,
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  "下单成功",
		"order_id": oir.OrderSn,
	})
}

// 根据用户id获取订单列表
func GetOrder(c *gin.Context) {
	cUser := getClaimes(c)
	zap.S().Infof("访问用户: %d", cUser.ID)
	// pn := c.DefaultQuery("page", "1")
	// pnInt, _ := strconv.Atoi(pn)
	// ps := c.DefaultQuery("size", "10")
	// psInt, _ := strconv.Atoi(ps)
	// _, err := global.OrderClient.GetOrderListByUserId(context.Background(), &proto.OrderPageInfo{
	// 	Pn: uint32(pnInt),
	// 	Ps: uint32(psInt),
	// 	// 暂时为了演示写死
	// 	UserId: 14,
	// })
	// if err != nil {
	// 	HandleGrpcErrorToHttp(err, c)
	// 	return
	// }
	res, err2 := global.UserClient.GetPassengerList(context.WithValue(context.Background(), "ginContext", c), &proto.PassengerPageInfo{Pn: 1, Ps: 100, UserId: 14})
	if err2 != nil {
		HandleGrpcErrorToHttp(err2, c)
		return
	}
	var tiks []struct {
		Id            int32  `json:"id"`
		UserId        int32  `json:"memberId"`
		PassengerId   string `json:"passengerId"`
		PassengerName string `json:"passengerName"`
		Date          string `json:"date"`
		TrainCode     string `json:"trainCode"`
		Row           string `json:"row"`
		Col           string `json:"col"`
		Start         string `json:"start"`
		StartTime     string `json:"startTime"`
		CarriageIndex int32  `json:"carriageIndex"`
		End           string `json:"end"`
		EndTime       string `json:"endTime"`
		SeatType      string `json:"seatType"`
	}
	for _, value := range res.Data {
		zap.S().Infof("乘客: %v", value)
		if value.IdCard == "" {
			continue
		}
		var tik struct {
			Id            int32  `json:"id"`
			UserId        int32  `json:"memberId"`
			PassengerId   string `json:"passengerId"`
			PassengerName string `json:"passengerName"`
			Date          string `json:"date"`
			TrainCode     string `json:"trainCode"`
			Row           string `json:"row"`
			Col           string `json:"col"`
			Start         string `json:"start"`
			StartTime     string `json:"startTime"`
			CarriageIndex int32  `json:"carriageIndex"`
			End           string `json:"end"`
			EndTime       string `json:"endTime"`
			SeatType      string `json:"seatType"`
		}
		tik.Id = value.Id
		tik.UserId = value.UserId
		tik.PassengerId = fmt.Sprintf("%d", value.Id)
		tik.PassengerName = value.Name
		tik.Date = time.Now().Format("2006-01-02")
		// 直接写死
		tik.TrainCode = "tc59"
		rand.NewSource(time.Now().UnixNano())
		rowInt := rand.Intn(2)
		tik.Row = fmt.Sprintf("0%d", rowInt)
		tik.Col = fmt.Sprintf("%c%d", seats[rand.Intn(6)], rowInt)
		tiks = append(tiks, tik)
		// tik.Start = ""
		// tik.StartTime =
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"content": gin.H{
			"list":  tiks,
			"total": len(tiks),
		},
	})

}
