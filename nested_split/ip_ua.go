package nested_split

import "strings"

func GetIpAndSplitUA(ip, ua string) string {
	if ip == "" || ua == "" {
		return ""
	}
	ua = strings.Split(strings.Split(strings.Split(strings.Split(strings.Split(strings.Split(strings.Split(
		strings.Split(strings.Split(strings.Split(strings.Split(strings.Split(strings.Split(ua, " - HuabenApp")[0], " MeetYouClient")[0], " CSDNApp")[0], " mztapp")[0], " fezpet")[0], " DWD_HSQ")[0], " avmPlus")[0], " QMNovel")[0], " Weibo")[0],
		" Html5Plus")[0], "TESHUBIAOSHI")[0], " motor")[0], " Safari")[0]
	ipUa := ip + "," + ua
	return ipUa
}
