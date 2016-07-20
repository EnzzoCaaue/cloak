package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var (
	captchaVerificationURL = "https://www.google.com/recaptcha/api/siteverify"
)

type captcha struct {
	Success bool
}

// Validate validates de given struct with its tags
func Validate(f interface{}) []error {
	errs := []error{}
	val := reflect.ValueOf(f).Elem()
	for i := 0; i < val.NumField(); i++ {
		p := val.Type().Field(i).Tag
		fieldAlias := p.Get("alias")
		fieldValue := val.Field(i).String()
		fullTag := p.Get("validate")
		tags := strings.Split(fullTag, ",")
		for _, tag := range tags {
			v, err := valid(tag, fieldAlias, fieldValue)
			if !v {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

func valid(tag, alias, value string) (bool, error) {
	parsedTags := strings.Split(strings.TrimSpace(tag), "=")
	if viper.GetString("mode") == "DEV" {
		return true, nil
	}
	switch parsedTags[0] {
	case "min":
		minValue, _ := strconv.Atoi(parsedTags[1])
		if !validateString(value, minValue, 0) {
			return false, fmt.Errorf("%v: %v", alias, "Invalid minimum field size")
		}
	case "max":
		maxValue, _ := strconv.Atoi(parsedTags[1])
		if !validateString(value, 0, maxValue) {
			return false, fmt.Errorf("%v: %v", alias, "Invalid maximum field size")
		}
	case "regexp":
		if !validateRegexp(value, parsedTags[1]) {
			return false, fmt.Errorf("%v: %v", alias, "Contains invalid characters")
		}
	case "validGender":
		if !validateGender(value) {
			return false, fmt.Errorf("%v: %v", alias, "Invalid gender")
		}
	case "validVocation":
		if !validateVocation(value) {
			return false, fmt.Errorf("%v: %v", alias, "Invalid vocation")
		}
	case "validCaptcha":
		if !validateCaptcha(value) {
			return false, fmt.Errorf("%v: %v", alias, "Wrong captcha response")
		}
	}
	return true, nil
}

func validateRegexp(val, pattern string) bool {
	match, err := regexp.MatchString(pattern, val)
	if err != nil {
		return false
	}
	return match
}

func validateString(val string, min, max int) bool {
	if min == 0 && len(val) > max {
		return false
	}
	if max == 0 && len(val) < min {
		return false
	}
	/*if max > 0 && min > 0 && len(val) < min || len(val) > max {
		return false
	}*/
	return true
}

func validateGender(val string) bool {
	if _, e := genderList[val]; !e {
		return false
	}
	return true
}

func validateVocation(val string) bool {
	if _, e := vocationList[val]; !e {
		return false
	}
	return true
}

func validateCaptcha(val string) bool {
	resp, err := http.PostForm(captchaVerificationURL, url.Values{
		"secret": {
			viper.GetString("captcha.secret"),
		},
		"response": {
			val,
		},
	})
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	captchaResponse := &captcha{}
	err = json.Unmarshal(body, captchaResponse)
	if err != nil {
		return false
	}
	return captchaResponse.Success
}
