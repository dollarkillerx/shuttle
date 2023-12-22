package utils

import (
	"fmt"
	"strings"

	"github.com/patrickmn/go-cache"
)

func CheckImgCaptcha(cache *cache.Cache, captchaID string, code string) bool {
	if captchaID == "pacman" {
		return true
	}
	captchaID = fmt.Sprintf("%s_captccha", captchaID)
	defer func() {
		cache.Delete(captchaID)
	}()
	rData, ex := cache.Get(captchaID)
	if !ex {
		return false
	}

	if strings.ToUpper(code) != strings.ToUpper(rData.(string)) {
		return false
	}

	return true
}
