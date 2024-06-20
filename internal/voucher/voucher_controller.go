package voucher

import (
	"log"
	"net/http"
	"time"
	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"github.com/gin-gonic/gin"
	"fmt"
)

type VoucherController struct {
	VoucherService VoucherService
}

func New(voucherService VoucherService) VoucherController {
	return VoucherController{
		VoucherService: voucherService,
	}
}

type ValidateVoucherRequest struct {
    CartAmount float64 `json:"cart_amount"`
}

func (vc *VoucherController) CreateVoucher(ctx *gin.Context) {
	var voucher entity.Voucher
	// voucher.ID = entity.GenerateID()
	istLocation,_:= time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(istLocation)
	formattedTime := now.Format("01-02-2006 15:04:05")
	voucher.CreatedAt = formattedTime
	if err :=  ctx.ShouldBindJSON(&voucher); err != nil {
	ctx.JSON(http.StatusBadRequest,gin.H{"message": err.Error()})
	return
	};

	if(voucher.Type != "restaurant") {
		ctx.JSON(403, gin.H{"status": 403, "message": "cannot create voucher as a user"})
		return
	}
	err := vc.VoucherService.CreateVoucher(&voucher)
	if err != nil{
        ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
        return
	}
	ctx.JSON(http.StatusOK, gin.H{"message":"success","status": http.StatusOK,"voucher_code":voucher.VoucherCode})
	}

func (vc *VoucherController) GetVoucher(ctx *gin.Context) {
	voucherId := ctx.Param("id")
	voucher, err := vc.VoucherService.GetVoucher((&voucherId))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, voucher)
}

func (vc *VoucherController) GetAll(ctx *gin.Context) {
	vouchers, err := vc.VoucherService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, vouchers)

}

func (vc *VoucherController) UpdateVoucher(ctx *gin.Context) {
	var voucher entity.Voucher
	// voucher.ID = ctx.Param("id")
	voucher.VoucherCode = ctx.Param("id")
	if err := ctx.ShouldBindJSON(&voucher); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := vc.VoucherService.UpdateVoucher(&voucher)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	istLocation,_:= time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(istLocation)
	formattedTime := now.Format("01-02-2006 15:04:05")
	voucher.CreatedAt = formattedTime
	ctx.JSON(http.StatusOK, gin.H{"message": "success","status":http.StatusOK})

}

func (vc *VoucherController) DeleteVoucher(ctx *gin.Context) {
	voucherId := ctx.Param("id")
	log.Println(voucherId)
	err := vc.VoucherService.DeleteVoucher(&voucherId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success","status":http.StatusOK})
}

func (vc *VoucherController) ValidateVoucher(ctx *gin.Context) {
	voucherCode := ctx.Param("id")
    log.Println(voucherCode)
    
    var req ValidateVoucherRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }

    cartAmount := req.CartAmount
    min_cart_value,final_amount,amt, err := vc.VoucherService.ValidateVoucher(&voucherCode, &cartAmount)
	if err != nil {
		if err == ErrMinimumAmt {
			message := fmt.Sprintf("Minimum cart amount of %.2f should be reached", min_cart_value)
            ctx.JSON(http.StatusBadRequest, gin.H{"message": message, "status": http.StatusBadRequest})
            return
        }
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success","status":http.StatusOK,"discount_amount":amt,"final_amount":final_amount})
}

func (vc *VoucherController) RegisterVoucherRoutes(rg *gin.RouterGroup) {
	voucherroute := rg.Group("/voucher")
	voucherroute.POST("/create", vc.CreateVoucher)
	voucherroute.GET("/get/:id", vc.GetVoucher)
	voucherroute.GET("/getvouchers", vc.GetAll)
	voucherroute.PATCH("/update/:id", vc.UpdateVoucher)
	voucherroute.DELETE("/delete/:id", vc.DeleteVoucher)
	voucherroute.POST("/validate/:id",vc.ValidateVoucher)
}
