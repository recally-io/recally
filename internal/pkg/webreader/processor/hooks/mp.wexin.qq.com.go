package hooks

import (
	"net/url"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type WeixinMP struct {
	Host string `json:"host"`
}

func NewWeixinMP() *WeixinMP {
	return &WeixinMP{
		Host: "mp.weixin.qq.com",
	}
}

func (p *WeixinMP) RegisterHooks() {
	mdBeforeHooks[p.Host] = append(mdBeforeHooks[p.Host], p.beforeImgHook)
	readabilityBeforeHooks[p.Host] = append(readabilityBeforeHooks[p.Host], p.readabilityParsePublishedTimeHook, p.readabilityCoverHook)
}

func (p *WeixinMP) beforeImgHook(selec *goquery.Selection) {
	selec.Find("img").Each(func(i int, s *goquery.Selection) {
		// src="https://mmbiz.qpic.cn/mmbiz_png/VWpZENjIo5uQjCCFia5oz2y6PTBHwhS77SR5RwSExhic2lhIMicqeiatOorMn5Q1H2PFn1xpiarJtNlRsHRlsiaomnkQ/640?wx_fmt=png&from=appmsg&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1"

		src := s.AttrOr("src", "")
		dataSrc := s.AttrOr("data-src", "")
		if dataSrc != "" && strings.HasPrefix(src, "data:image") {
			src = dataSrc
		}
		imgUrl := p.transferImg(src)
		s.SetAttr("src", imgUrl)
	})
}

func (p *WeixinMP) readabilityParsePublishedTimeHook(doc *goquery.Document, content *webreader.Content) {
	// <em id="publish_time" class="rich_media_meta rich_media_meta_text">2024年12月31日 18:51</em>
	publishTime := doc.Find("#publish_time").First().Text()
	if publishTime == "" {
		return
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, err := time.ParseInLocation("2006年01月02日 15:04", publishTime, loc)
	if err != nil {
		logger.Default.Warn("Failed to parse published time", "err", err)
		return
	}

	content.PublishedTime = &t
}

func (p *WeixinMP) readabilityCoverHook(doc *goquery.Document, content *webreader.Content) {
	content.Cover = p.transferImg(content.Cover)
}

func (p *WeixinMP) transferImg(src string) string {
	if !strings.Contains(src, "mmbiz.qpic.cn") {
		return src
	}
	imgUrl, err := url.Parse(src)
	if err != nil {
		logger.Default.Warn("Failed to parse image url", "err", err)
		return src
	}
	imgType := imgUrl.Query().Get("wx_fmt")
	if imgType != "" {
		newUrl := imgUrl.Scheme + "://" + imgUrl.Host + imgUrl.Path + "." + imgType
		return newUrl
	}
	return src
}
