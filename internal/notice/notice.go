package notice

import (
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
)

type Notifier interface {
	Notice(content string)
}

type WechatWebhookMsg struct {
	webhookurl string
}

func NewWechatWebhookMsg(webhookurl string) Notifier {
	return WechatWebhookMsg{
		webhookurl: webhookurl,
	}
}

func (wm WechatWebhookMsg) Notice(content string) {
	payload := map[string]interface{}{
		"msgtype": "text",
	}
	payload["text"] = map[string]string{
		"content": content,
	}

	_, _, errs := gorequest.New().Post(wm.webhookurl).Send(payload).End()
	if len(errs) > 0 {
		logrus.Errorf("send notice error:%v", errs)
	}
}

type SlackWebhookMsg struct {
	webhookurl string
}

func NewSlackWebhookMsg(webhookurl string) Notifier {
	return SlackWebhookMsg{
		webhookurl: webhookurl,
	}
}

func (wm SlackWebhookMsg) Notice(content string) {
	payload := map[string]interface{}{
		"text": content,
	}

	_, _, errs := gorequest.New().Post(wm.webhookurl).
		Set("Content-type", "application/json").
		Send(payload).End()
	if len(errs) > 0 {
		logrus.Errorf("send notice error:%v", errs)
	}
}
