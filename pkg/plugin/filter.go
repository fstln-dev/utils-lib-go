package plugin

import (
	"errors"
	"fmt"
	"regexp"
)

// 大小写字母，数字，下划线字符串检测
const pwRegexp string = `^[a-zA-Z0-9_]{%d,%d}$`

func CheckWord(word string, lower, upper int) (bool, error) {
	return regexp.MatchString(fmt.Sprintf(pwRegexp, lower, upper), word)
}

//
const iRegexpModel string = `^[a-zA-Z0-9]{%d}$`

func CheckInvCode(code string, charCount int) (bool, error) {
	return regexp.MatchString(fmt.Sprintf(iRegexpModel, charCount), code)
}

//
const emailRegexpModel string = `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`

func CheckEmail(email string) (bool, error) {
	return regexp.MatchString(emailRegexpModel, email)
}

const shopNameFilter string = `^(\S+).myshopify.com$`

func SubmatchShopName(shop string) (string, error) {
	res := regexp.MustCompile(shopNameFilter).FindStringSubmatch(shop)

	if len(res) < 2 {
		return "", errors.New("Not Find ShopName")
	}
	return res[1], nil
}

const emptyFilter string = `\s+`

func ReplaceEmpty(str, rep string) (string, error) {
	reg, err := regexp.Compile(emptyFilter)
	if err != nil {
		return "", err
	}
	return reg.ReplaceAllString(str, rep), nil
}
