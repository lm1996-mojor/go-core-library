package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/lm1996-mojor/go-core-library/log"
	"github.com/lm1996-mojor/go-core-library/rest/req"
	"github.com/lm1996-mojor/go-core-library/utils/repo"

	"github.com/hashicorp/go-uuid"
)

type CodeGenerationRules struct {
	req.CommonModel
	PrefixStr          string `gorm:"column:prefix_str;type:string" json:"prefixStr"`                                // 编号前缀字符串
	Remark             string `gorm:"column:remark;type:string" json:"remark"`                                       // 备注
	DateTimeStr        string `gorm:"column:date_time_str;type:string" json:"dateTimeStr,omitempty"`                 // 日期时间串（y-m-d-h-M-s，顺序组合）
	OtherStrType       int8   `gorm:"column:other_str_type;type:tinyint" json:"otherStrType,omitempty"`              // 其他字符串类型（1、随机字符串 2、UUID字符串 3、MD5字符串）
	OtherStrDigits     int    `gorm:"column:other_str_digits;type:int" json:"otherStrDigits,omitempty"`              // 其他字符串位数（-1 全部 其他自由输入）
	OtherStrFormat     int8   `gorm:"column:other_str_format;type:tinyint" json:"otherStrFormat,omitempty"`          // 其他字符串格式（1 数字 2 英文 3 英文+数字）
	OtherStrFormatCase int8   `gorm:"column:other_str_format_case;type:tinyint" json:"otherStrFormatCase,omitempty"` // 其他字符串格式大小写（1 全部大写 2 全部小写）
	Status             int8   `gorm:"column:status;type:tinyint" json:"status,omitempty"`                            // 状态（1 启用 2 停用）
}

func (c *CodeGenerationRules) TableName() string {
	return "code_generation_rules"
}

func (c *CodeGenerationRules) allColumn() []string {
	columns := []string{"prefix_str", "remark", "date_time_str", "other_str_type", "other_str_digits", "other_str_format", "other_str_format_case"}
	commonMdl := req.CommonModel{}
	columns = append(columns, commonMdl.GetCommonModelColumns()...)
	return columns
}

// 获取编码前缀
func obtainCodePrefixText(codeType int) CodeGenerationRules {
	var code CodeGenerationRules
	db := repo.ObtainCustomDbByDbName("platform_management")
	if db != nil {
		db.Table(code.TableName()).Where("id = ?", codeType).Select(code.allColumn()).Scan(&code)
	} else {
		log.Error("请配置platform_management数据库")
		panic("服务器错误")
	}
	return code
}
func getDateTimeStr(format string) string {
	dataTimeStr := ""
	formatList := strings.Split(format, "-")
	now := time.Now()
	for _, formatStr := range formatList {
		switch formatStr {
		case "y":
			dataTimeStr += fmt.Sprintf("%d", now.Year())
		case "m":
			dataTimeStr += fmt.Sprintf("%d", now.Month())
		case "d":
			dataTimeStr += fmt.Sprintf("%d", now.Day())
		case "h":
			dataTimeStr += fmt.Sprintf("%d", now.Hour())
		case "M":
			dataTimeStr += fmt.Sprintf("%d", now.Minute())
		case "s":
			dataTimeStr += fmt.Sprintf("%d", now.Second())
		default:
			dataTimeStr += "00"
		}
	}
	return dataTimeStr
}

func getRandomStr(digit int, format int8, formatCase int8) string {
	rand.Seed(time.Now().UnixNano()) // 初始化随机数种子
	var charset string
	switch format {
	case 1:
		charset = "0123456789"
	case 2:
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	default:
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}
	result := make([]byte, digit)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	str := string(result)
	strResult := ""
	if format != 1 {
		if formatCase == 1 {
			strResult = strings.ToUpper(str)
		} else {
			strResult = strings.ToLower(str)
		}
	}

	return strResult
}

func getUUIDStr(digit int, format int8, formatCase int8) string {
	uuId, _ := uuid.GenerateUUID()
	srcStr := strings.ReplaceAll(uuId, "-", "")
	return getDigitStr(digit, srcStr, format, formatCase)
}

func getMd5Str(digit int, format int8, formatCase int8) string {
	hasher := md5.New()
	hasher.Write([]byte(getRandomStr(10, 2, 1)))
	hashBytes := hasher.Sum(nil)
	srcStr := hex.EncodeToString(hashBytes)
	return getDigitStr(digit, srcStr, format, formatCase)
}

func getDigitStr(digit int, srcStr string, format int8, formatCase int8) string {
	rand.Seed(time.Now().UnixNano()) // 初始化随机数种子
	if format != 1 {
		if formatCase == 1 {
			srcStr = strings.ToUpper(srcStr)
		} else {
			srcStr = strings.ToLower(srcStr)
		}
	}
	var charset string
	switch format {
	case 1:
		charset = "0123456789"
	case 2:
		charset = "abcdefghijklmnopqrstuvwxyz"
		if formatCase == 1 {
			charset = strings.ToUpper(charset)
		}
	default:
		charset = "abcdefghijklmnopqrstuvwxyz0123456789"
		if formatCase == 1 {
			charset = strings.ToUpper(charset)
		}
	}
	// 过滤掉源字符串中不在字符集中的字符
	filteredSource := ""
	for _, char := range srcStr {
		if strings.ContainsRune(charset, char) {
			filteredSource += string(char)
		}
	}
	if filteredSource == "" {
		filteredSource = charset
	}
	if digit != -1 {
		result := make([]byte, digit)
		for i := range result {
			result[i] = filteredSource[rand.Intn(len(filteredSource))]
		}
		return string(result)
	} else {
		return filteredSource
	}
}

// GenerateCodeBySearchId 根据搜索ID生成编码(该id是平台管理端中编码生成规则的ID,平台管理端-工具库-编码生成规则)
func GenerateCodeBySearchId(codeType int) (code string) {
	codePrefix := obtainCodePrefixText(codeType)
	code += codePrefix.PrefixStr
	if codePrefix.DateTimeStr != "" {
		code += getDateTimeStr(codePrefix.DateTimeStr)
	}
	if codePrefix.OtherStrFormat <= 0 || codePrefix.OtherStrFormat > 3 {
		codePrefix.OtherStrFormat = 3
	}
	if codePrefix.OtherStrDigits < -1 {
		codePrefix.OtherStrDigits = -1
	}
	if codePrefix.Status != 1 {
		code += "CODE"
		uuId, _ := uuid.GenerateUUID()
		code += strings.ReplaceAll(uuId, "-", "")
	} else {
		switch codePrefix.OtherStrType {
		case 1:
			code += getRandomStr(codePrefix.OtherStrDigits, codePrefix.OtherStrFormat, codePrefix.OtherStrFormatCase)
		case 2:
			code += getUUIDStr(codePrefix.OtherStrDigits, codePrefix.OtherStrFormat, codePrefix.OtherStrFormatCase)
		case 3:
			code += getMd5Str(codePrefix.OtherStrDigits, codePrefix.OtherStrFormat, codePrefix.OtherStrFormatCase)
		default:
			uuId, _ := uuid.GenerateUUID()
			code += strings.ReplaceAll(uuId, "-", "")
		}
	}
	return code
}
