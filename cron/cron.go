package cron

import (
	"fmt"
	"gin_demo/models"
	"github.com/robfig/cron"
	"log"
	"time"
)

func main() {
	fmt.Println("starting...")

	cron := cron.New()

	cron.AddFunc("* * * * * *", func() {
		log.Println("Run Models.ClearAllArticles...")
		models.ClearAllArticle()
	})

	cron.AddFunc("* * * * * *", func() {
		log.Println("Run Models.ClearAllTag...")
		models.ClearAllTag()
	})

	cron.Start()

	//这里使用定时器进行阻塞,否则还没等定时任务执行,进程就退出了
	t1 := time.NewTimer(time.Second * 10)

	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}
