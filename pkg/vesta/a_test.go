package vesta_test

import (
	"strings"
	"testing"
	"time"

	"github.com/mocheer/pluto/pkg/ds"
	"github.com/mocheer/vesta/pkg/vesta"
)

func TestHead(t *testing.T) {

	v := vesta.New().Edge().Head()
	var val1 string
	var val2 string
	var val3 string
	var val4 string
	defer v.Cancel()
	v.Nav("https://baidu.com")
	v.Eval("location.href", &val1)
	v.Sleep(1000 * time.Millisecond)
	v.Nav("https://www.zhihu.com/")
	v.Eval("location.href", &val2)
	v.ClosePage()
	v.Run()

	// 新标签页
	v2 := v.NewContext()
	defer v2.Cancel()
	v2.Nav("http://192.168.118.103:9912/v/studio/login")
	v2.EvaluateAsDevTools("T.appStore.version+''", &val3)
	// v2.Sleep(1000 * time.Millisecond)
	// v2.Eval("T.appStore.version+''", &val4)
	v2.Run()
	// 新浏览器
	// v3 := vesta.NewWithHead()
	// defer v2.Cancel()
	// v3.Nav("https://baidu.com")
	// v3.Eval("location.href", &val1)
	// v3.Sleep(1000 * time.Millisecond)
	// v3.Run()

	t.Log(val1, val2)
	t.Log(val3, val4)

	time.Sleep(2000 * time.Millisecond)
}

func TestInject(t *testing.T) {
	v := vesta.New().Head()

	defer v.Cancel()
	v.Inject(`window.$vesta =  {version:"1.0",body:document.body};`)
	var val1 string
	v.Sleep(2000 * time.Millisecond)

	v.Nav("http://192.168.118.103:9912/v/studio/login")
	v.InterceptRequestWithJS(&vesta.InterceptRequestParams{
		UpdateBody: map[string]string{"http://192.168.118.103:9912/charon/v1/dal/get/config_ww_app?id=ep": `data=>{console.log(data);return {
    "charon": "v11",
    "code": 200,
    "data": {
        "id": "ep",
        "corpid": "ww72c3202458a6eb5d",
        "agentid": "1000071",
        "redirect_uri": "https://sso.istrongcloud.net/",
        "state": "2_713f9fabd0034d2284f28a3011aaff4e"
    },
    "msg": ""
}}`},
	})
	v.Eval("JSON.stringify(Object.keys($vesta))", &val1)
	v.Sleep(3000 * time.Millisecond)
	err := v.Run()
	time.Sleep(20000 * time.Millisecond)

	t.Log(val1)
	t.Log(err)

}

func TestMain(t *testing.T) {
	v := vesta.New().Head()
	defer v.Cancel()
	v.Inject(`originToLowerCase = String.prototype.toLowerCase ;String.prototype.toLowerCase = function(){
	  if(this.length>10){
			  return 'micromessenger'
		}
    return originToLowerCase.call(this)
	}`)
	v.Nav("http://localhost:9912/charon/v1/agent/wskj")

	err := v.Run()
	time.Sleep(time.Hour)
	t.Log(err)
}

func TestCharonReptile(t *testing.T) {
	go func() {
		v := vesta.New().Head()
		defer v.Cancel()

		v.Nav("https://www.msn.cn/zh-cn/weather/maps/precipitation/in-%E7%A6%8F%E5%BB%BA%E7%9C%81,%E5%B0%A4%E6%BA%AA%E5%8E%BF?loc=eyJsIjoi5bCk5rqq5Y6%2FIiwiciI6Iuemj%2BW7uuecgSIsImMiOiLkuK3ljY7kurrmsJHlhbHlkozlm70iLCJpIjoiQ04iLCJnIjoiemgtY24iLCJ4IjoiMTE4LjA0MjYwMjU0IiwieSI6IjI1Ljk0MzIyNjc5In0%3D&weadegreetype=C&cvid=802a9fdd09554d6ba1e7e6f9db9fc737&zoom=7&lightning=0&thdhai=0")

		v.Sleep(time.Second * 10)
		v.Size(1900, 1000)
		v.Viewport(1080, 1080)
		v.Sleep(time.Second * 10)
		data, err := v.GetValue("document.title")

		t.Log(data, err)
	}()
	select {}

}

func TestSave(t *testing.T) {
	// vesta.SaveAllResource("https://webshare2.shanhaibi.com/erwkq60victe/")
	// vesta.SaveAllResource("https://webshare2.shanhaibi.com/hp0hfoh0hxcp/")
	vesta.SaveAllResource("https://webshare2.shanhaibi.com/db339wa89rho/") //大坝水电站 https://www.shanhaibi.com/market/theme/801.html

	// vesta.SaveAllResource("http://192.168.118.103:9912/v/studio/login")
}

func TestA(t *testing.T) {
	v := vesta.New().Edge().Head()

	defer v.Cancel()
	v.Nav("http://221.13.83.50:50403/shzhfz/prePlatform/index")
	v.InterceptRequestWithJS(&vesta.InterceptRequestParams{
		UpdateBody: map[string]string{
			"http://221.13.83.50:50403/shzhfz/shzhfz-province-api/baseplat-auth/login": `data=>{
			   data.code = 200
				 data.data = true
				 data.status = 200
				 return data
			}`,
			"http://221.13.83.50:50403/shzhfz/shzhfz-province-api/baseplat-system/user/getInfo": `data=>{
				 data.code = 200
				 data.data = true
				 data.status = 200
				 return data
			}`,
		},
	})
	v.Sleep(100 * time.Hour)
	v.Run()
}

func TestC(t *testing.T) {

	ds.EachFilesToRename("./baidu_pano", func(name string) string {
		if strings.HasSuffix(name, ".png") {
			return strings.ReplaceAll(name, ".png", ".jpg")
		}
		return name
	})

}
