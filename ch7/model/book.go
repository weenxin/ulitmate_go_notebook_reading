package model

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

//Catalog Book 的类型
type Catalog int8

const (
	CategoryNovel      Catalog = iota //CategoryNovel 小说
	CategoryShortStory Catalog = iota //CategoryShortStory 短故事
)
const WeightEnvName = "WEIGHT_UNITS"

//MaxShortStoryPages 短故事最大页数
const MaxShortStoryPages = 300

//Book 是测试用例
type Book struct {
	Id     int64  `json:"id,omitempty"`
	Title  string `json:"title,omitempty" json:"title,omitempty"`
	Author string `json:"author,omitempty" json:"author,omitempty"`
	Pages  int32  `json:"pages" json:"pages,omitempty"`
	Weight int32  `json:"weight,omitempty" json:"weight,omitempty"` // 存储时使用g
}

//NewBookFromJSON 通过json创建 Book 对象
func NewBookFromJSON(data string) (*Book, error) {
	var b Book
	if err := json.Unmarshal([]byte(data), &b); err != nil {
		return nil, err
	}
	return &b, nil
}

//Catalog 返回 Book 类型
func (b Book) Catalog() Catalog {
	if b.Pages < MaxShortStoryPages {
		return CategoryShortStory
	} else {
		return CategoryNovel
	}
}

//FirstName 返回作者 FirstName
func (b Book) FirstName() string {
	if !b.IsValid() {
		return ""
	}
	nameParts := strings.Split(b.Author, " ")
	switch len(nameParts) {
	case 1:
		return ""
	default:
		return nameParts[0]
	}
}

//LastName 返回作者 LastName
func (b Book) LastName() string {
	if !b.IsValid() {
		return ""
	}
	nameParts := strings.Split(b.Author, " ")
	return nameParts[len(nameParts)-1]
}

//MiddleName 返回作者 MiddleName
func (b Book) MiddleName() string {
	if !b.IsValid() {
		return ""
	}
	nameParts := strings.Split(b.Author, " ")
	if len(nameParts) < 3 {
		return ""
	}
	return strings.Join(nameParts[1:len(nameParts)-1], " ")
}

//IsValid 返回是否有效
func (b Book) IsValid() bool {
	if len(b.Author) == 0 {
		return false
	}
	return true
}

//HumanReadableWeight 返回可读性高的书本重量
func (b Book) HumanReadableWeight() (string, error) {
	unit := "g"
	content := os.Getenv(WeightEnvName)
	if len(content) != 0 {
		unit = content
	}
	switch unit {
	case "g":
		return fmt.Sprintf("%d%s", b.Weight, unit), nil
	case "kg":
		return fmt.Sprintf("%.3f%s", float64(b.Weight)/1000, unit), nil
	default:
		return "", fmt.Errorf("invalid unit, only support [g,kg]")
	}

}
