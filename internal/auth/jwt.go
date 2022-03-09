package auth

import (
	"github.com/ashah360/fibertools"
	"github.com/ashah360/foodquest-api/internal/api/cerror"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"os"
	"strings"
	"time"
)

var cookieAndQueryName = "access_token"

// ValidateJWT validates the auth token and extracts the user ID
func ValidateJWT(token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, cerror.ErrMalformedToken
		}

		return []byte(os.Getenv("FOODQUEST_JWT_SECRET")), nil
	})
	if err != nil {
		return "", cerror.ErrMalformedToken
	}

	if !t.Valid {
		return "", cerror.ErrInvalidToken
	}

	claims := t.Claims.(jwt.MapClaims)

	exp, ok := claims["exp"]
	if !ok {
		return "", cerror.ErrMalformedToken
	}

	if time.Since(time.Unix(int64(exp.(float64)), 0)) > 0 {
		return "", cerror.ErrInvalidToken
	}

	uid, ok := claims["id"]
	if !ok {
		return "", cerror.ErrMalformedToken
	}

	return uid.(string), nil
}

// ExtractJWT extracts the token from the passed-in fiber context
func ExtractJWT(c *fiber.Ctx) string {
	a := strings.Split(fibertools.GetHeader(c, "Authorization"), " ")
	if len(a) == 2 && strings.EqualFold(a[0], "bearer") {
		return a[1]
	}

	// if no header, try cookie
	j := utils.UnsafeString([]byte(c.Cookies(cookieAndQueryName)))
	if j != "" {
		return j
	}

	q := utils.UnsafeString([]byte(c.Query(cookieAndQueryName)))
	if len(q) > 0 {
		return q
	}

	return ""
}
