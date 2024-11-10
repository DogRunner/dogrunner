package identity

import (
	"github.com/labstack/echo/v4"
	jwt "github.com/wanrun-develop/wanrun/internal/auth/middleware"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
)

// GetClaims : contextから検証済みのclaims情報を取得する共通関数
//
// args:
//   - echo.MiddlewareFunc: JWT検証のためのミドルウェア設定
//
// return:
//   - *AccountClaims: 検証済みのclaims情報
func GetVerifiedClaims(c echo.Context) (*jwt.AccountClaims, error) {
	logger := log.GetLogger(c).Sugar()

	claims, ok := c.Get(jwt.VERIFIED_CONTEXT_KEY).(*jwt.AccountClaims)
	if !ok || claims == nil {
		wrErr := errors.NewWRError(
			nil,
			"クレーム情報が見つかりません。",
			errors.NewDogownerClientErrorEType(),
		)
		logger.Error(wrErr)
		return nil, wrErr
	}

	return claims, nil
}
