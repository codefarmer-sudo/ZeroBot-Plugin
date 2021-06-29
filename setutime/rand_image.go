package setutime

import (
	"strings"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	RANDOM_API_URL = "https://api.pixivweb.com/anime18r.php?return=img"
	BLOCK_REQUEST  = false
)

func init() { // 插件主体
	zero.OnRegex(`^设置随机图片网址(.*)$`, zero.SuperUserPermission).SetBlock(true).SetPriority(20).
		Handle(func(ctx *zero.Ctx) {
			url := ctx.State["regex_matched"].([]string)[1]
			if !strings.HasPrefix(url, "http") {
				ctx.Send("URL非法!")
			} else {
				RANDOM_API_URL = url
			}
			return
		})
	// 有保护的随机图片
	zero.OnFullMatch("随机图片").SetBlock(true).SetPriority(24).
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.GroupID > 0 {
				Classify(ctx, RANDOM_API_URL, false)
			}
			return
		})
	// 直接随机图片，危险，仅管理员可用
	zero.OnFullMatch("直接随机", zero.AdminPermission).SetBlock(true).SetPriority(24).
		Handle(func(ctx *zero.Ctx) {
			if ctx.Event.GroupID > 0 {
				if BLOCK_REQUEST {
					ctx.Send("请稍后再试哦")
				} else {
					BLOCK_REQUEST = true
					last_message_id := ctx.SendChain(message.Image(RANDOM_API_URL).Add("no_cache", "0"))
					last_group_id := ctx.Event.GroupID
					MsgofGrp[last_group_id] = last_message_id
					BLOCK_REQUEST = false
				}
			}
			return
		})
	// 撤回最后的随机图片
	zero.OnFullMatch("不许好").SetBlock(true).SetPriority(24).
		Handle(func(ctx *zero.Ctx) {
			Vote(ctx, 5)
		})
	// 撤回最后的随机图片
	zero.OnFullMatch("太涩了").SetBlock(true).SetPriority(24).
		Handle(func(ctx *zero.Ctx) {
			Vote(ctx, 6)
		})
	/*
		// 上传一张图进行评价
		zero.OnFullMatch("评价图片", MustHasPicture()).SetBlock(true).SetPriority(24).
			Handle(func(ctx *zero.Ctx) {
				if ctx.Event.GroupID > 0 {
					ctx.Send("少女祈祷中......")
					for _, pic := range ctx.State["image_url"].([]string) {
						fmt.Println(pic)
						Classify(ctx, pic, true)
					}
				}
				return
			})
	*/
}
