package coredogwatcher

import (
	"context"
	"log"
	"os"

	"github.com/DomineCore/coredog/internal/store"
	"github.com/DomineCore/coredog/internal/watcher"
	"github.com/DomineCore/coredog/pb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func getHostip() string {
	return os.Getenv("HOST_IP")
}

func WatchCorefile() {
	cfg = getCfg()
	recevier := make(chan string)
	w := watcher.NewFileWatcher(recevier)
	err := w.Watch(cfg.COREFILE_DIR)
	if err != nil {
		log.Fatal(errors.Wrap(err, "start watch error"))
	}
	s, err := store.NewS3Store(
		cfg.StorageConfig.S3Region,
		cfg.StorageConfig.S3AccessKeyID,
		cfg.StorageConfig.S3SecretAccessKey,
		cfg.StorageConfig.S3Bucket,
		cfg.StorageConfig.S3Endpoint,
		cfg.StorageConfig.StoreDir,
		cfg.StorageConfig.PresignedURLExpireSeconds,
	)
	if err != nil {
		logrus.Fatal(err)
	}
	cli, conn := NewCoreFileServiceClient()
	defer conn.Close()
	for {
		select {
		case corefilePath := <-recevier:
			url, err := s.Upload(context.Background(), corefilePath)
			if err != nil {
				logrus.Errorf("store a corefile error:%v", err)
				return
			}
			if cfg.Gc && (cfg.GcType == "rm") {
				_ = os.Remove(corefilePath)
			} else if cfg.Gc && (cfg.GcType == "truncate") {
				_ = os.Truncate(corefilePath, 0)
			}

			logrus.Infof("corefile upload down: %s", url)
			go func() {
				_, err := cli.Sub(context.Background(), &pb.Corefile{
					Filepath: corefilePath,
					Ip:       getHostip(),
					Url:      url,
				})
				if err != nil {
					logrus.Errorf("pub a corefile error:%v", err)
				}
			}()
		}
	}
}
