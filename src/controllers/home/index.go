package home

import (
	"fmt"
	"go-auth/src/models"
	"go-auth/src/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/ylfyt/go_db/go_db"
	"github.com/ylfyt/meta/meta"
)

type HomeController struct {
}

func (me *HomeController) Setup() []meta.EndPoint {
	return []meta.EndPoint{
		{
			Method:      "GET",
			Path:        "/",
			HandlerFunc: me.home,
		},
		{
			Method:      "GET",
			Path:        "/ping",
			HandlerFunc: me.ping,
		},
	}

}

func (me *HomeController) home(c *fiber.Ctx, db *go_db.DB) meta.ResponseDto {
	reqId := utils.GetContext[string](c, "reqId")
	fmt.Println("ReqId:", *reqId)
	var product *models.Product
	err := db.GetFirst(&product, `SELECT * FROM products LIMIT 1`)
	if err != nil {
		fmt.Println("Err", err)
		return utils.GetInternalErrorResponse("Something wrong")
	}

	return utils.GetSuccessResponse(product)
}
