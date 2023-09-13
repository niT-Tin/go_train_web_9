package api

import (
	"context"
	"gotrains/ticketorder_web/ticketorder-web/forms"
	"gotrains/ticketorder_web/ticketorder-web/global"
	"gotrains/ticketorder_web/ticketorder-web/models"
	"gotrains/ticketorder_web/ticketorder-web/proto"
	"net/http"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"

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
	res, err := global.UserClient.GetPassengerList(context.Background(), &proto.PassengerPageInfo{
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
	passenger, err := global.UserClient.GetPassengerByIdCard(context.Background(), &proto.PassengerIdCardRequest{
		IdCard: passengerForm.IdCard,
	})
	// 这里不对，但是没时间了，直接在Add里面判断了
	if passenger.IdCard != "" {
		_, err = global.UserClient.UpdatePassenger(context.Background(), &proto.PassengerInfo{
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
	_, err = global.UserClient.AddPassenger(context.Background(), &proto.PassengerInfo{
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
	_, err := global.UserClient.DeletePassenger(context.Background(), &proto.IdRequest{
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
	passengers := []models.Passenger{}
	if err := c.ShouldBindJSON(&passengers); err != nil {
		zap.S().Errorw("CreateOrder 参数错误", "msg", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "参数错误",
		})
		return
	}
	orderInfo := proto.OrderInfo{}
	var orderPassengers []*proto.OPassengerInfo
	order := models.Order{}
	for _, passenger := range passengers {
		opass := proto.OPassengerInfo{
			Name:     passenger.Name,
			IdCard:   passenger.IdCard,
			Type:     int64(passenger.Type),
			UserId:   int32(passenger.UserID),
			Seat:     passenger.Seat,
			SeatType: passenger.SeatType,
		}
		orderPassengers = append(orderPassengers, &opass)
	}
	orderInfo.Passengers = orderPassengers
	orderInfo.UserId = int32(cUser.ID)
	orderInfo.TrainId = order.TrainID
	orderInfo.StartStation = order.StartStation
	orderInfo.EndStation = order.EndStation
	orderInfo.StartTime = order.StartTime.Format("2006-01-02 15:04:05")
	orderInfo.EndTime = order.EndTime.Format("2006-01-02 15:04:05")
	orderInfo.SeatType = order.SeatType
	orderInfo.SeatNumber = order.SeatNumber
	orderInfo.Price = order.Pirce
	oir, err := global.OrderClient.CreateOrder(context.Background(), &proto.CreateOrderInfo{
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
