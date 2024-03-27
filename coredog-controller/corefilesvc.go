package coredogcontroller

import (
	"context"
	"fmt"

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
		msg := fmt.Sprintf("A new corefile is generated. download: %s", r.Url)
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
