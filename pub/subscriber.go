package pub

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/iiinsomnia/gochat/utils"
	"github.com/tidwall/gjson"
)

// MaxSubscriberListCount 公众号订阅者列表的最大数目
const MaxSubscriberListCount = 10000

// Subscriber 微信公众号订阅者信息
type Subscriber struct {
	Subscribe      int     `json:"subscribe"`
	OpenID         string  `json:"openid"`
	NickName       string  `json:"nickName"`
	Sex            int     `json:"sex"`
	Language       string  `json:"language"`
	City           string  `json:"city"`
	Province       string  `json:"province"`
	Country        string  `json:"country"`
	AvatarURL      string  `json:"headimgurl"`
	SubscribeTime  int64   `json:"subscribe_time"`
	UnionID        string  `json:"unionid"`
	Remark         string  `json:"remark"`
	GroupID        int64   `json:"groupid"`
	TagidList      []int64 `json:"tagid_list"`
	SubscribeScene string  `json:"subscribe_scene"`
	QRScene        int64   `json:"qr_scene"`
	QRSceneStr     string  `json:"qr_scene_str"`
}

// SubscriberList 微信公众号订阅者列表
type SubscriberList struct {
	Total      int                 `json:"total"`
	Count      int                 `json:"count"`
	Data       map[string][]string `json:"data"`
	NextOpenID string              `json:"next_openid"`
}

// GetSubscriberInfo 获取微信公众号订阅者信息
func GetSubscriberInfo(accessToken, openid string) (*Subscriber, error) {
	resp, err := utils.HTTPGet(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN", accessToken, openid))

	if err != nil {
		return nil, err
	}

	r := gjson.ParseBytes(resp)

	if r.Get("errcode").Int() != 0 {
		return nil, errors.New(r.Get("errmsg").String())
	}

	reply := new(Subscriber)

	if err := json.Unmarshal(resp, reply); err != nil {
		return nil, err
	}

	return reply, nil
}

// BatchGetSubscriberInfo 批量获取微信公众号订阅者信息
func BatchGetSubscriberInfo(accessToken string, openid ...string) ([]*Subscriber, error) {
	l := len(openid)

	if l == 0 {
		return []*Subscriber{}, nil
	}

	userList := make([]map[string]string, 0, l)

	for _, v := range openid {
		userList = append(userList, map[string]string{
			"openid": v,
			"lang":   "zh_CN",
		})
	}

	b, err := json.Marshal(map[string][]map[string]string{"user_list": userList})

	if err != nil {
		return nil, err
	}

	resp, err := utils.HTTPPost(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info/batchget?access_token=%s", accessToken), b, utils.WithRequestHeader("Content-Type", "application/json; charset=utf-8"))

	if err != nil {
		return nil, err
	}

	r := gjson.ParseBytes(resp)

	if r.Get("errcode").Int() != 0 {
		return nil, errors.New(r.Get("errmsg").String())
	}

	reply := make(map[string][]*Subscriber)

	if err := json.Unmarshal(resp, reply); err != nil {
		return nil, err
	}

	return reply["user_info_list"], nil
}

// GetSubscriberList 获取微信公众号订阅者列表
func GetSubscriberList(accessToken string, nextOpenID ...string) (*SubscriberList, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/get?access_token=%s", accessToken)

	if len(nextOpenID) > 0 {
		url = fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/get?access_token=%s&next_openid=%s", accessToken, nextOpenID[0])
	}

	resp, err := utils.HTTPGet(url)

	if err != nil {
		return nil, err
	}

	r := gjson.ParseBytes(resp)

	if r.Get("errcode").Int() != 0 {
		return nil, errors.New(r.Get("errmsg").String())
	}

	reply := new(SubscriberList)

	if err := json.Unmarshal(resp, reply); err != nil {
		return nil, err
	}

	return reply, nil
}