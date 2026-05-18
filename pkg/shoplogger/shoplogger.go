package shoplogger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/WhoaWicked/board-game-shop/pkg/utils"
	"github.com/gofiber/fiber/v3"
)

type IShopLogger interface {
	Print() IShopLogger
	Save()
	SetQuery(c fiber.Ctx)
	SetBody(c fiber.Ctx)
	SetResponse(data any)
}

type shopLogger struct {
	Time       string `json:"time"`
	Ip         string `json:"ip"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	Path       string `json:"path"`
	Query      any    `json:"query"`
	Body       any    `json:"body"`
	Response   any    `json:"response"`
}

func InitShopLogger(c fiber.Ctx, res any) IShopLogger {
	log := &shopLogger{
		Time:       time.Now().Local().Format("2006-01-02 15:04:05"),
		Ip:         c.IP(),
		Method:     c.Method(),
		StatusCode: c.Response().StatusCode(),
		Path:       c.Path(),
	}
	log.SetQuery(c)
	log.SetBody(c)
	log.SetResponse(res)
	return log
}

func (l *shopLogger) Print() IShopLogger {
	utils.Debug(l)
	return l
}
func (l *shopLogger) Save() {
	data := utils.Output(l)
	filename := fmt.Sprintf("./assets/logs/shoplogger%v.txt", strings.ReplaceAll(time.Now().Format("2006-01-02"), "-", ""))
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	file.WriteString(string(data) + "\n")
}

func (l *shopLogger) SetQuery(c fiber.Ctx) {
	query := make(map[string]string)
	if err := c.Bind().Query(&query); err != nil {
		log.Printf("query parser error: %v", err)
	}
	l.Query = query
}

func (l *shopLogger) SetBody(c fiber.Ctx) {
	body := make(map[string]any)
	if err := c.Bind().Body(&body); err != nil {
		log.Printf("body parser error: %v", err)
	}
	l.Body = body
}

func (l *shopLogger) SetResponse(res any) {
	l.Response = res
}
