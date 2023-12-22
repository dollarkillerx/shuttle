package pkg

import (
	"fmt"
	"testing"
	"time"

	"github.com/rs/xid"
	"google.dev/google/shuttle/utils"
)

func TestVerifyToken_FromJSON(t *testing.T) {
	var vto = &VerifyToken{
		Expiration: time.Now().Add(time.Hour * 5).Unix(),
		NodeID:     "xxxxx",
		UserJWT:    "user jwt",
	}

	aesKey := utils.PadKeyString(xid.New().String(), 32)

	toString := vto.ToToken(aesKey)
	fmt.Println(toString)

	var vto2 = new(VerifyToken)
	vto2.FromToken(toString, aesKey)
	fmt.Println(vto2)
}

func TestVp2(t *testing.T) {
	rx := "UxJV+Y/6j7QdD8eAe+oW643qPVq+Qyt5aNDAfEDlF7SC/B+QtcBWV6Mv6Fte1HAc7ufHXXF15JwGiLQiyr+ZBRyRczuPPEvGI675SII4CloQfcWf+d3Zmwa/u29FPZl3UboAApvrJSY+I4ll3atV8RnF2mZXGcMLN/qLnoGsA2OiZxUtNz0d9kZ0z/+NTOgHZNOMdUYHsgi65G4UhsSEvUBcOp3HQPc1kKwOedXzpPiYJs9EqVkmGRG0f687bSEF"

	aesKey := utils.PadKeyString("cgbsk003c01nkbq0ahgg", 32)

	var vto2 = new(VerifyToken)
	vto2.FromToken(rx, aesKey)
	fmt.Println(vto2)
}
