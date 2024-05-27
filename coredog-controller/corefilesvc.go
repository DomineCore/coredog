package coredogcontroller

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/DomineCore/coredog/internal/notice"
	"github.com/DomineCore/coredog/pb"
	"github.com/sirupsen/logrus"
)

type CorefileService struct {
}

func (s *CorefileService) Sub(ctx context.Context, r *pb.Corefile) (*pb.Corefile, error) {
	logrus.Infof("recevier a newfile:%s from host:%s", r.Filepath, r.Ip)
	cfg = getCfg()
	for _, noticeChan := range cfg.NoticeChannel {
		_, corefilename := filepath.Split(r.Filepath)
		msg := buildMessage(cfg.MessageTemplate, corefilename, cfg.MessageLabels, r.Url)
		if noticeChan.Chan == "wechat" {
			c := notice.NewWechatWebhookMsg(noticeChan.Webhookurl)
			c.Notice(msg)
		} else if noticeChan.Chan == "slack" {
			c := notice.NewSlackWebhookMsg(noticeChan.Webhookurl)
			c.Notice(msg)
		} else {
			logrus.Warnf("unsupported notice channel:%s", noticeChan.Chan)
		}
	}

	return &pb.Corefile{}, nil
}

func buildMessage(template string, corefilename string, labels map[string]string, url string) string {
	// 1 replace the labels into template
	msg := template
	for key, val := range labels {
		msg = strings.ReplaceAll(msg, fmt.Sprintf("{%v}", key), val)
	}
	msg = strings.ReplaceAll(msg, CorefileName, corefilename)
	msg = strings.ReplaceAll(msg, CorefileUrl, url)
	return msg
}
