package main

import (
	"github.com/mocheer/vesta/pkg/vesta"
)

// 省委
// https://map.baidu.com/search/%E4%B8%AD%E5%9B%BD%E5%85%B1%E4%BA%A7%E5%85%9A%E5%B1%B1%E4%B8%9C%E7%9C%81%E5%A7%94%E5%91%98%E4%BC%9A/@13026044.01,4365806.55,21z,87t,-90.04h?querytype=s&da_src=shareurl&wd=%E4%B8%AD%E5%9B%BD%E5%85%B1%E4%BA%A7%E5%85%9A%E5%B1%B1%E4%B8%9C%E7%9C%81%E5%A7%94%E5%91%98%E4%BC%9A&c=288&src=0&wd2=%E6%B5%8E%E5%8D%97%E5%B8%82%E5%B8%82%E4%B8%AD%E5%8C%BA&pn=0&sug=1&l=17&b=(13016182.15497,4374781.86557;13019397.70503,4376398.01443)&from=webmap&biz_forward=%7B%22scaler%22:1,%22styles%22:%22pl%22%7D&device_ratio=1#panoid=09016200122207211646188987E&panotype=street&heading=90.04&pitch=0&l=21&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=09016200122207211646188987E

// 省政府顶部
// https://map.baidu.com/search/%E4%B9%90%E5%B1%B1%E5%B0%8F%E5%8C%BA%E5%8D%97%E5%8C%BA-6%E5%8F%B7%E6%A5%BC/@13027326.82,4368558.76,21z,87t,-169.07h#panoid=09016200121902151455253091F&panotype=street&heading=169.07&pitch=0&l=13&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=09016200121902151455253091F

// 省政府中间
// https://map.baidu.com/poi/%E5%B1%B1%E4%B8%9C%E9%81%93%E9%BD%90%E6%8A%95%E8%B5%84%E5%92%A8%E8%AF%A2%E6%9C%89%E9%99%90%E5%85%AC%E5%8F%B8/@13027543.99,4367994.2,21z,87t,112.72h?uid=37c4ba6c19f9c9091ce99927&ugc_type=3&ugc_ver=1&device_ratio=1&compat=1&pcevaname=pc4.1&querytype=detailConInfo&da_src=shareurl#panoid=09016200011704151120090806S&panotype=street&heading=213.49&pitch=0.94&l=21&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=09016200011704151120090806S

// 济南龙奥大厦
// https://map.baidu.com/poi/%E5%B1%B1%E4%B8%9C%E9%81%93%E9%BD%90%E6%8A%95%E8%B5%84%E5%92%A8%E8%AF%A2%E6%9C%89%E9%99%90%E5%85%AC%E5%8F%B8/@13037866.18,4366134.32,21z,87t,-103.49h?uid=37c4ba6c19f9c9091ce99927&ugc_type=3&ugc_ver=1&device_ratio=1&compat=1&pcevaname=pc4.1&querytype=detailConInfo&da_src=shareurl#panoid=09016200122207181032082355D&panotype=street&heading=103.49&pitch=0&l=21&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=09016200122207181032082355D

// 泉城广场
// https://map.baidu.com/poi//@13027198.1,4367253.3,21z,87t,16.44h#panoid=09016200122207221118541117E&panotype=street&heading=343.56&pitch=0&l=13&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=09016200122207221118541117E

// 长途总站
// https://map.baidu.com/search//@13024161.07,4370579.45,21z,87t,11.24h#panoid=09016200121902171652198687I&panotype=street&heading=95.22&pitch=-11.45&l=13&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=09016200121902171652198687I

// 济南广场汽车站 济南站
// https://map.baidu.com/search/%E6%B5%8E%E5%8D%97%E5%B9%BF%E5%9C%BA%E6%B1%BD%E8%BD%A6%E7%AB%99/@13024303.09,4367988.05,21z,87t,-110.57h?querytype=s&da_src=shareurl&wd=%E6%B5%8E%E5%8D%97%E5%B9%BF%E5%9C%BA%E6%B1%BD%E8%BD%A6%E7%AB%99&c=288&src=0&wd2=%E6%B5%8E%E5%8D%97%E5%B8%82%E5%A4%A9%E6%A1%A5%E5%8C%BA&pn=0&sug=1&l=18&b=(13023125.85663,4367898.60719;13024769.12236,4368724.51939)&from=webmap&biz_forward=%7B%22scaler%22:1,%22styles%22:%22pl%22%7D&sug_forward=3f303e6a63813ce2b54f5444&device_ratio=1#panoid=09016200122207241356458925A&panotype=street&heading=110.57&pitch=0&l=21&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=09016200122207241356458925A

// 济南遥城机场
// https://map.baidu.com/search/%E6%B5%8E%E5%8D%97%E5%B9%BF%E5%9C%BA%E6%B1%BD%E8%BD%A6%E7%AB%99/@13048778.57,4393059.56,21z,87t,-76.02h?querytype=s&da_src=shareurl&wd=%E6%B5%8E%E5%8D%97%E5%B9%BF%E5%9C%BA%E6%B1%BD%E8%BD%A6%E7%AB%99&c=288&src=0&wd2=%E6%B5%8E%E5%8D%97%E5%B8%82%E5%A4%A9%E6%A1%A5%E5%8C%BA&pn=0&sug=1&l=18&b=(13023125.85663,4367898.60719;13024769.12236,4368724.51939)&from=webmap&biz_forward=%7B%22scaler%22:1,%22styles%22:%22pl%22%7D&sug_forward=3f303e6a63813ce2b54f5444&device_ratio=1#panoid=0901620012210105120527314HC&panotype=street&heading=76.01&pitch=0&l=21&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=0901620012210105120527314HC

// 绕城高速
// 双向车道未完结
// https://map.baidu.com/poi/%E6%B5%8E%E5%8D%97%E7%BB%95%E5%9F%8E%E9%AB%98%E9%80%9F/@13010036.44,4375265.71,21z,87t,-162.38h?uid=fa59b51d8d65c36b82a00c69&ugc_type=3&ugc_ver=1&device_ratio=1&compat=1&pcevaname=pc4.1&querytype=detailConInfo&da_src=shareurl#panoid=01016200001405210813185285X&panotype=street&heading=154.11&pitch=-15.96&l=21&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=01016200001405210813185285X

// 济南西站 西元大厦
// https://map.baidu.com/search/%E6%B5%8E%E5%8D%97%E8%A5%BF%E7%AB%99/@13012866.48,4367130.14,21z,87t,-30.22h?querytype=s&da_src=shareurl&wd=%E6%B5%8E%E5%8D%97%E8%A5%BF%E7%AB%99&c=288&src=0&wd2=%E6%B5%8E%E5%8D%97%E5%B8%82%E6%A7%90%E8%8D%AB%E5%8C%BA&pn=0&sug=1&l=13&b=(13013399.495604469,4368796.19301726;13061957.507735817,4393201.652239527)&from=webmap&biz_forward=%7B%22scaler%22:1,%22styles%22:%22pl%22%7D&sug_forward=578267a89a00ae9cc6effb5b&device_ratio=1#panoid=09016200011506100442250966A&panotype=street&heading=29.22&pitch=0&l=21&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=09016200011506100442250966A

// 济南西站2
// https://map.baidu.com/search/%E6%B5%8E%E5%8D%97%E8%A5%BF%E7%AB%99/@13013533.240000002,4367679.26,21z,87t,-179.5h?querytype=s&da_src=shareurl&wd=%E6%B5%8E%E5%8D%97%E8%A5%BF%E7%AB%99&c=288&src=0&wd2=%E6%B5%8E%E5%8D%97%E5%B8%82%E6%A7%90%E8%8D%AB%E5%8C%BA&pn=0&sug=1&l=13&b=(13013399.495604469,4368796.19301726;13061957.507735817,4393201.652239527)&from=webmap&biz_forward=%7B%22scaler%22:1,%22styles%22:%22pl%22%7D&sug_forward=578267a89a00ae9cc6effb5b&device_ratio=1#panoid=0901620012210105114111201HT&panotype=street&heading=3.33&pitch=-13.95&l=21&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=0901620012210105114111201HT

// 济南东站
// https://map.baidu.com/search//@13042405.84,4378982.92,21z,87t,-84.92h#panoid=09016200011604141536231546O&panotype=street&heading=84.92&pitch=0&l=21&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=09016200011604141536231546O

// 经十路-绕城高速
// https://map.baidu.com/search//@13012437.69,4365516.48,21z,87t,-81.5h#panoid=0901620012210104104253353HT&panotype=street&heading=81.5&pitch=0.07&l=13&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=0901620012210104104253353HT

// 北园高架
// https://map.baidu.com/@13011319.21,4369393.67,21z,87t,111.79h#panoid=01016200001405161350590195M&panotype=street&heading=67.52&pitch=-20&l=21&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=01016200001405161350590195M

// 山东大厦
// https://map.baidu.com/search/@13027084.91,4364527.25,21z,87t,89.44h#panoid=09016200121902280920014335E&panotype=street&heading=71.12&pitch=-19.98&l=13&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=09016200121902280920014335E

// 二环高架路
// https://map.baidu.com/poi/%E4%BA%8C%E7%8E%AF%E5%8D%97%E9%AB%98%E6%9E%B6%E8%B7%AF/@13017691.09,4374790.8,21z,87t,-6.49h?uid=a3209db9de1f0ace70210bdb&ugc_type=3&ugc_ver=1&device_ratio=1&compat=1&pcevaname=pc4.1&querytype=detailConInfo&da_src=shareurl#panoid=09016200122207171257250323O&panotype=street&heading=185.5&pitch=0&l=21&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=09016200122207171257250323O

func main() {
	run()
}

// console.log('$click',683,370)
// allow pasting
// var a = setInterval(()=>{console.log('$click',683,370)},1500);a;
// clearInterval(a)
func run() {

	v := vesta.New().Head().Size(1920, 1020)
	defer v.Cancel()
	v.Nav(`https://map.baidu.com/search//@13012437.69,4365516.48,21z,87t,-81.5h#panoid=09016200122207230757215742R&panotype=street&heading=89.2&pitch=-20&l=21&tn=B_NORMAL_MAP&sc=0&newmap=1&shareurl=1&pid=09016200122207230757215742R`)
	v.InterceptRequestWithBauduPano("经十路")
	v.AddEventMouseClickXY()
	// v.Sleep(100 * time.Hour)
	v.Run()
	//

	select {}
}
