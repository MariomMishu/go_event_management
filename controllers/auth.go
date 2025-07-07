package controllers

import (
	"ems/domain"
	"ems/types"
	"ems/utils/msgutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authSvc domain.AuthService
}

func NewAuthController(authSvc domain.AuthService) *AuthController {
	return &AuthController{
		authSvc: authSvc,
	}
}
func (ctrl *AuthController) Login(c echo.Context) error {
	var req types.LoginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, msgutil.InvalidRequestMsg())
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, &types.ValidationError{
			Error: err,
		})
	}

	resp, err := ctrl.authSvc.Login(&req)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, msgutil.SomethingWentWrongMsg())
	}

	return c.JSON(http.StatusOK, resp)
}

func (ctrl *AuthController) Logout(c echo.Context) error {
	return c.JSON(http.StatusOK, msgutil.LogoutSuccessfully())
}
