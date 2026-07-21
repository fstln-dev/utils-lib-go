package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/google/uuid"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"time"
)

// InteractiveCard 代表整个消息卡片对象
type InteractiveCard struct {
	Schema string      `json:"schema"`
	Header *CardHeader `json:"header,omitempty"`
	Body   *CardBody   `json:"body,omitempty"`
	Config *CardConfig `json:"config,omitempty"`
}

func (ic *InteractiveCard) ToJSONString() (string, error) {
	b, err := json.Marshal(ic)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CardHeader 定义了卡片的标题栏
type CardHeader struct {
	Title    *CardText `json:"title"`
	Subtitle *CardText `json:"subtitle,omitempty"`
	Template string    `json:"template,omitempty"`
	Icon     *CardIcon `json:"icon,omitempty"`
	Padding  string    `json:"padding,omitempty"`
}

// CardIcon 定义了标题栏的图标
type CardIcon struct {
	Tag   string `json:"tag"`   // "standard_icon"
	Token string `json:"token"` // e.g., "warning-hollow_filled"
}

// CardBody 定义了卡片的主体部分
type CardBody struct {
	Direction string        `json:"direction,omitempty"`
	Padding   string        `json:"padding,omitempty"`
	Elements  []interface{} `json:"elements"`
}

// CardConfig 定义了卡片的配置项
type CardConfig struct {
	UpdateMulti bool     `json:"update_multi,omitempty"`
	Locales     []string `json:"locales,omitempty"`
}

// CardText 是卡片中通用的文本单元
type CardText struct {
	Tag         string            `json:"tag"` // "plain_text" or "lark_md"
	Content     string            `json:"content,omitempty"`
	I18nContent map[string]string `json:"i18n_content,omitempty"`
}

// CardElementMarkdown 代表一个 Markdown 文本元素
type CardElementMarkdown struct {
	Tag         string            `json:"tag"` // "markdown"
	Content     string            `json:"content,omitempty"`
	I18nContent map[string]string `json:"i18n_content,omitempty"`
	TextAlign   string            `json:"text_align,omitempty"`
	Margin      string            `json:"margin,omitempty"`
}

// CardElementHR 代表一条分割线
type CardElementHR struct {
	Tag    string `json:"tag"` // "hr"
	Margin string `json:"margin,omitempty"`
}

func FeishuAlertMiddleware(client *lark.Client, chatID string, serverName string, ignore map[string]interface{}) middleware.Middleware {
	// 1. 定义一个常量来控制最大并发告警数
	const maxConcurrentAlerts = 10
	// 2. 创建一个带缓冲的 channel 作为信号量，限制并发的 goroutine 数量
	// 这个信号量对于所有通过此中间件的请求是共享的
	var alertSemaphore = make(chan struct{}, maxConcurrentAlerts)

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			reply, err = handler(ctx, req)
			if err == nil {
				return
			}

			// 从 error 中解析出结构化的错误信息 (kratos/errors)
			se := errors.FromError(err)
			if se == nil {
				return
			}

			// 检查这个错误 reason 是否在我们的忽略列表里
			if _, ok := ignore[se.Reason]; ok {
				return
			}

			// 3. 在启动 goroutine 之前，先获取一个信号量。
			// 如果已达到最大并发数，这里会阻塞，直到有其他 goroutine 完成并释放信号量
			alertSemaphore <- struct{}{}

			// 异步发送飞书告警，避免阻塞主流程
			go func() {
				// 4. 使用 defer 确保在 goroutine 退出时，无论成功或失败，都释放信号量
				defer func() {
					<-alertSemaphore
				}()

				// 创建一个独立的后台上下文，它的生命周期与原始请求 ctx 无关
				// 这样即使原始请求超时或被取消，也不会影响告警的发送
				alertCtx := context.Background()

				alertTime := time.Now().Format("2006-01-02 15:04:05")
				card := &InteractiveCard{
					Schema: "2.0",
					Header: &CardHeader{
						Title: &CardText{
							Tag:     "plain_text",
							Content: fmt.Sprintf("🚨 服务告警: %s", serverName),
						},
						Template: "red", // 使用红色模板，醒目地表示这是一个错误
						Icon: &CardIcon{
							Tag:   "standard_icon",
							Token: "warning-hollow_filled", // 使用警告图标
						},
					},
					Body: &CardBody{
						Elements: []interface{}{
							// 使用 Markdown 元素分块展示信息，比纯文本更清晰
							&CardElementMarkdown{
								Tag:     "markdown",
								Content: fmt.Sprintf("**错误码 (Code):**\n`%d`", se.Code),
							},
							&CardElementMarkdown{
								Tag:     "markdown",
								Content: fmt.Sprintf("**错误原因 (Reason):**\n`%s`", se.Reason),
							},
							&CardElementMarkdown{
								Tag: "markdown",
								// 使用代码块来展示详细的错误堆栈，保持格式
								Content: fmt.Sprintf("**详细信息 (Detail):**\n%s", err.Error()),
							},
							// 添加分割线，让布局更美观
							&CardElementHR{Tag: "hr"},
							&CardElementMarkdown{
								Tag:     "markdown",
								Content: fmt.Sprintf("<font color='grey'>告警时间: %s</font>", alertTime),
							},
						},
					},
				}

				contentStr, _ := card.ToJSONString()
				fReq := larkim.NewCreateMessageReqBuilder().
					ReceiveIdType(larkim.ReceiveIdTypeChatId).
					Body(larkim.NewCreateMessageReqBodyBuilder().
						ReceiveId(chatID).
						MsgType(larkim.MsgTypeInteractive).
						Content(contentStr).
						Uuid(uuid.New().String()).
						Build()).
					Build()

				resp, err := client.Im.V1.Message.Create(alertCtx, fReq)
				if err != nil {
					fmt.Printf("发送飞书告警失败: %v", err)
					return
				}
				if !resp.Success() {
					fmt.Printf("飞书机器人发送消息失败, code: %d, msg: %s", resp.Code, resp.Msg)
				}
			}()
			return
		}
	}
}
