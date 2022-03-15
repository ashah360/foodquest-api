package mock

import (
	"fmt"
	"github.com/ashah360/foodquest-api/internal/api/model"
	"github.com/brianvoe/gofakeit"
	"math/rand"
	"strings"
)

// 12345
var mockPassword = "$2a$10$bCCaYlBpW5xwn/LkFHsdy.r9JpJLwqY.EX/T3mYE.O1Suoc5W8nG."

func NewUser() *model.User {
	p := gofakeit.Person()

	return &model.User{
		Email:       fmt.Sprintf("%s%s%d@gmail.com", strings.ToLower(p.FirstName[0:1]), strings.ToLower(p.LastName), rand.Intn(100)),
		Password:    mockPassword,
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		PhoneNumber: strings.NewReplacer("(", "", ")", "", "-", "").Replace(p.Contact.Phone),
	}
}
