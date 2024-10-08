package remote

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"runtime"
	"scnu_acm_rank/biz/config"
	"scnu_acm_rank/biz/model"
	"sort"
	"strconv"
	"strings"
)

type VjCrawler struct {
	userName string
	passWord string
	cookie   string
}

type VjRespJson struct {
	ID           int                  `json:"id"`
	Title        string               `json:"title"`
	Begin        int64                `json:"begin"`
	Length       int64                `json:"length"`
	IsReplay     bool                 `json:"isReplay"`
	Participants map[string][3]string `json:"participants"`
	Submissions  [][4]int64           `json:"submissions"`
}

var VJCrawler *VjCrawler

const getRankUrl = "https://vjudge.net/contest/rank/single/"
const loginUrl = "https://vjudge.net/user/login"
const checkLoginStatusUrl = "https://vjudge.net/user/checkLogInStatus"

func init() {
	//config := model.Config{}
	//middle.DB.Model(&model.Config{}).First(&config)

	//VJCrawler.userName = config.VjUserName
	//VJCrawler.passWord = config.VjPassWord
	//VJCrawler.cookie = config.VjCookie
	//if flag, err := VJCrawler.checkLoginStatus(); flag && err != nil {
	//	return
	//}
	//if cookie, err := VJCrawler.Login(); err != nil {
	//	VJCrawler.cookie = cookie
	//}
	VJCrawler = &VjCrawler{}
	runtime.KeepAlive(VJCrawler)
	config.Add(VJCrawler)
}

func (vj *VjCrawler) Update() {
	v := config.Conf
	vj.passWord = v.VjPassWord
	vj.cookie = v.VjCookie
	vj.userName = v.VjUserName
}

func (vj *VjCrawler) checkLoginStatus() (bool, error) {

	req, err := http.NewRequest("POST", checkLoginStatusUrl, nil)

	if err != nil {
		return false, err
	}

	req.Header.Add("cookie", vj.cookie)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}
	return string(body) == "true", nil
}

func (vj *VjCrawler) Login() bool {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("username", vj.userName)
	_ = writer.WriteField("password", vj.passWord)
	_ = writer.WriteField("captcha", "")
	err := writer.Close()
	if err != nil {
		return false
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", loginUrl, payload)

	if err != nil {
		return false
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false
	}
	if string(body) != "success" {
		return false
	}

	cookies := res.Header.Values("Set-Cookie")
	cookie := ""
	for _, s := range cookies {
		ck := strings.Split(s, ";")
		if len(ck) > 0 {
			cookie += ck[0] + ";"
		}
		fmt.Println(ck[0])
	}
	vj.cookie = cookie
	fmt.Println(cookie)
	return true
}

func (vj *VjCrawler) GetTrainRes(contest string) (string, error) {
	f, err := vj.checkLoginStatus()
	if !f || err != nil {
		if !vj.Login() {
			return "", errors.New("login fail")
		}
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://vjudge.net/contest/rank/single/"+contest, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Apipost client Runtime/+https://www.apipost.cn/")
	req.Header.Set("cookie", vj.cookie)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	temp := VjRespJson{}
	err = json.Unmarshal(bodyText, &temp)
	if err != nil {
		return "", err
	}
	return string(bodyText), nil
}

func (vj *VjCrawler) AnalysisRes(v interface{}) (*AnalysisRes, *model.Competition, error) {
	resJson := v.(string)
	res := VjRespJson{}
	err := json.Unmarshal([]byte(resJson), &res)
	if err != nil {
		return nil, nil, err
	}
	mp2name := make(map[string]*PersonalRes, len(res.Participants))
	for k, _ := range res.Participants {
		mp2name[k] = &PersonalRes{
			Name:        res.Participants[k][0],
			Submissions: make(map[int64]*Submission),
		}
	}

	sort.Slice(res.Submissions, func(a, b int) bool {
		return res.Submissions[a][3] < res.Submissions[b][3]
	})

	for _, v := range res.Submissions {
		ts := fmt.Sprintf("%v", v[0])
		if v[3] > res.Length/1000 {
			break
		}
		sub := mp2name[ts].Submissions[v[1]]
		if sub == nil {
			sub = &Submission{
				AcceptTime: -1,
			}
			mp2name[ts].Submissions[v[1]] = sub
		}
		if sub.AcceptTime != -1 {
			continue
		}
		if v[2] == 1 {
			//mp.AcceptTime = v[3]
			sub.AcceptTime = v[3]
			mp2name[ts].SolveCnt++
			mp2name[ts].Penalty += int(v[3]) + sub.SubCnt*20*60
		} else {
			sub.SubCnt++
		}
	}
	sli := make([]*PersonalRes, 0, len(mp2name))
	for k, _ := range mp2name {
		sli = append(sli, mp2name[k])
	}
	sort.Slice(sli, func(a, b int) bool {
		if sli[a].SolveCnt == sli[b].SolveCnt {
			return sli[a].Penalty < sli[b].Penalty
		}
		return sli[a].SolveCnt > sli[b].SolveCnt
	})
	for i, _ := range sli {
		sli[i].Rank = int64(i + 1)
		sli[i].Penalty *= 1000
	}
	r := &AnalysisRes{Result: sli}
	return r, &model.Competition{
		Name:      res.Title,
		CpId:      strconv.Itoa(res.ID),
		Kind:      0,
		StartDate: res.Begin,
		Length:    res.Length,
	}, nil
}
