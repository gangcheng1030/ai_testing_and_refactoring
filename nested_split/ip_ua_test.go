package nested_split

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIpAndSplitUA(t *testing.T) {
	testCases := []struct {
		name     string
		ip       string
		ua       string
		expected string
	}{
		{
			name:     "empty-ip",
			ip:       "",
			ua:       "Mozilla/5.0",
			expected: "",
		},
		{
			name:     "empty-ua",
			ip:       "192.168.1.1",
			ua:       "",
			expected: "",
		},
		{
			name:     "both-empty",
			ip:       "",
			ua:       "",
			expected: "",
		},
		{
			name:     "normal-ua-without-special-markers",
			ip:       "192.168.1.1",
			ua:       "Mozilla/5.0",
			expected: "192.168.1.1,Mozilla/5.0",
		},
		{
			name:     "ua-with-huabenapp",
			ip:       "10.0.0.1",
			ua:       "Mozilla/5.0 - HuabenApp/1.0",
			expected: "10.0.0.1,Mozilla/5.0",
		},
		{
			name:     "ua-with-meetyouclient",
			ip:       "172.16.0.1",
			ua:       "Mozilla/5.0 MeetYouClient/2.0",
			expected: "172.16.0.1,Mozilla/5.0",
		},
		{
			name:     "ua-with-csdnapp",
			ip:       "192.168.0.100",
			ua:       "Mozilla/5.0 CSDNApp/3.0",
			expected: "192.168.0.100,Mozilla/5.0",
		},
		{
			name:     "ua-with-mztapp",
			ip:       "10.10.10.10",
			ua:       "Mozilla/5.0 mztapp/1.5",
			expected: "10.10.10.10,Mozilla/5.0",
		},
		{
			name:     "ua-with-fezpet",
			ip:       "192.168.2.2",
			ua:       "Mozilla/5.0 fezpet/2.1",
			expected: "192.168.2.2,Mozilla/5.0",
		},
		{
			name:     "ua-with-dwd-hsq",
			ip:       "8.8.8.8",
			ua:       "Mozilla/5.0 DWD_HSQ/1.0",
			expected: "8.8.8.8,Mozilla/5.0",
		},
		{
			name:     "ua-with-avmplus",
			ip:       "1.1.1.1",
			ua:       "Mozilla/5.0 avmPlus/3.0",
			expected: "1.1.1.1,Mozilla/5.0",
		},
		{
			name:     "ua-with-qmnovel",
			ip:       "192.168.100.1",
			ua:       "Mozilla/5.0 QMNovel/1.2",
			expected: "192.168.100.1,Mozilla/5.0",
		},
		{
			name:     "ua-with-weibo",
			ip:       "172.20.0.1",
			ua:       "Mozilla/5.0 Weibo/4.0",
			expected: "172.20.0.1,Mozilla/5.0",
		},
		{
			name:     "ua-with-html5plus",
			ip:       "192.168.50.50",
			ua:       "Mozilla/5.0 Html5Plus/1.0",
			expected: "192.168.50.50,Mozilla/5.0",
		},
		{
			name:     "ua-with-teshubiaoshi",
			ip:       "10.0.1.1",
			ua:       "Mozilla/5.0TESHUBIAOSHI/1.0",
			expected: "10.0.1.1,Mozilla/5.0",
		},
		{
			name:     "ua-with-motor",
			ip:       "192.168.3.3",
			ua:       "Mozilla/5.0 motor/2.0",
			expected: "192.168.3.3,Mozilla/5.0",
		},
		{
			name:     "ua-with-safari",
			ip:       "127.0.0.1",
			ua:       "Mozilla/5.0 Safari/537.36",
			expected: "127.0.0.1,Mozilla/5.0",
		},
		{
			name:     "ua-with-multiple-markers",
			ip:       "192.168.1.100",
			ua:       "Mozilla/5.0 Safari/537.36 motor/2.0 Weibo/4.0",
			expected: "192.168.1.100,Mozilla/5.0",
		},
		{
			name:     "ua-with-all-markers-chain",
			ip:       "10.20.30.40",
			ua:       "BaseUA Safari/1.0 motor/1.0TESHUBIAOSHI Html5Plus/1.0 Weibo/1.0 QMNovel/1.0 avmPlus/1.0 DWD_HSQ/1.0 fezpet/1.0 mztapp/1.0 CSDNApp/1.0 MeetYouClient/1.0 - HuabenApp/1.0",
			expected: "10.20.30.40,BaseUA",
		},
		{
			name:     "ipv6-address",
			ip:       "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			ua:       "Mozilla/5.0 Safari/537.36",
			expected: "2001:0db8:85a3:0000:0000:8a2e:0370:7334,Mozilla/5.0",
		},
		{
			name:     "complex-ua-string",
			ip:       "192.168.1.1",
			ua:       "Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36",
			expected: "192.168.1.1,Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetIpAndSplitUA(tc.ip, tc.ua)
			assert.Equal(t, tc.expected, result)
		})
	}
}
