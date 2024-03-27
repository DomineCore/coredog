package coredogwatcher

import (
	"context"
	"os"

	"github.com/DomineCore/coredog/internal/store"
	"github.com/DomineCore/coredog/internal/watcher"
	"github.com/DomineCore/coredog/pb"
	"github.com/sirupsen/logrus"
)

func getHostip() string {
	return os.Getenv("HOST_IP")
}

func WatchCorefile() {
	cfg = getCfg()
	for _, path := range cfg.ScrapePaths {
		recevier := make(chan string)
		w := watcher.NewFileWatcher(recevier)
		w.Watch(path)

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
			default:
				continue
			}
		}
	}
}
