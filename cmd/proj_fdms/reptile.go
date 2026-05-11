package main

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/mocheer/pluto/pkg/ds"
	"github.com/mocheer/pluto/pkg/ds/ds_json"
	"github.com/mocheer/pluto/pkg/ts/ctp"
	"github.com/mocheer/pluto/pkg/ts/window"
	"github.com/mocheer/vesta/pkg/vesta"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func main() {
	go ReptileQ()
	// go ReptileM()
	// ReptileW()
	select {}
}

// Province
type Province struct {
	Name    string      `json:"name"`
	Adcd    string      `json:"adcd"`
	HasData bool        `json:"hasData"`
	Items   []*Province `json:"items"`
}

// Options
type Options struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// QMenu
type QMenu struct {
	Provinces []*Province `json:"provinces"`
	Options   []*Options  `json:"options"`
}

//	if (adcd.substring(6, 9) != '000') {
//	    window.parent.$("#mapframe").contents().find("#mapframe").attr("src", "/fdms/pro/table/IA_ADC_ADINFO_" + adcd);
//	} else {
//	    window.parent.$("#mapframe").contents().find("#mapframe").attr("src", "/fdms/pro/table/STS_IA_COUNTY_" + adcd);
//	}

//	 $("#myselect").change(function () {
//			var tableName = $('.selectpicker').selectpicker('val');
//			let url = "http://10.135.6.98:80/fdms//pro/table/" + tableName + $("#mapframe").attr("class");
//			$("#mapframe").attr("src", url);
//	});
func ReptileQ() {

	vm := vesta.New().Head()
	// vm.Nav("http://10.135.6.98/")
	// vm.WaitID("username")
	// vm.Eval(fmt.Sprintf(`
	// 	 document.querySelector("#username").value = "福建省操作员"
	// 	 document.querySelector("#password").value = "czyyh!123"

	// 	 function loadAndRunTesseract() {
	// 		// 1. 创建type="module"的script标签
	// 		const script = document.createElement('script');
	// 		script.type = 'module';
	// 		// 2. 设置script内容（内联模块代码）
	// 		script.innerHTML = %s;
	// 		// 3. 添加到文档中执行
	// 		document.head.appendChild(script);
	// 		return null;
	// 	}
	// 	loadAndRunTesseract()
	// 	// loginFun()
	// `, "`"+`
	// 		import tesseract from 'https://cdn.jsdelivr.net/npm/tesseract.js@5/dist/tesseract.esm.min.js';
	// 		const { createWorker } = tesseract;
	// 		const rec = (async () => {
	// 			const worker = await createWorker('eng');
	// 			// const ret = await worker.recognize('https://tesseract.projectnaptha.com/img/eng_bw.png');
	// 			const t = new Date().toISOString().slice(0,-1)
	// 			const ret = await worker.recognize('http://10.135.6.98/Account/GetValidateCode?time='+t+'&browserName=Safari')
	// 			document.querySelector("#CheckCode").value = ret.data.text
	// 			console.log(ret.data.text);
	// 			loginFun()
	// 			// await worker.terminate();
	// 		})();
	//   `+"`"), nil).Run()
	// time.Sleep(10 * time.Second)
	vm.Nav("http://10.135.6.98/fdms/q")
	qdata := new(QMenu)
	if !ds.IsExist("./q/provinces.json") {
		getProvincesScript := `(function(){
				// 省份 dom
				let items = Array.from(window.frames[0].document.querySelectorAll('.provinceElement'))
				// 省份
				let provinces = items.map(e=>{
					let capitalA = e.children[0].querySelector('a')
					let citylistA = Array.from(e.children[1].querySelectorAll('a.btn'))
					return {
						// 省份名称
						name:capitalA.innerText,
						adcd:capitalA.id || capitalA.getAttribute('adcdcode'),
						hasData:true,
						// 城市名称
						items:citylistA.map(e=>{
							return {
								name:e.innerText,
								hasData:true,
								adcd:e.id || e.getAttribute('adcdcode'),
								items:[]
							}
						}).filter(e=>e.adcd)
					}
				})
				// 选项dom
				// let options = Array.from(window.frames[1].document.querySelector("body>header>div>div>ul").querySelectorAll('a'))
				let options = Array.from(window.frames[1].document.querySelector("#myselect").querySelectorAll('option'))
				let config =  options.map(e=>{return {value:e.value, name:e.innerText}})

				let qdata =  {
					provinces,
					options:config
				};
				return qdata
		})()`
		vm.Eval(getProvincesScript, qdata).Run()
		log.Println(len(qdata.Provinces))
		//
		for _, pp := range qdata.Provinces {
			for _, p2 := range pp.Items {
				if !p2.HasData {
					log.Println(p2.Adcd, p2.Name, "无数据")
					continue
				}
				qdata2 := new(QMenu)
				itemUrl := fmt.Sprintf(`window.frames[0].$("#cityList").attr("src", "/fdms/pro/dist/%s");1;`, string(p2.Adcd[0:6]))
				log.Println(itemUrl)
				err := vm.Eval(itemUrl, nil).Run()
				if err != nil {
					log.Println(p2.Adcd, p2.Name, err)
					continue
				}
				time.Sleep(1 * time.Second)
				vm.Eval(getProvincesScript, qdata2).Run()
				p2.Items = append(p2.Items, qdata2.Provinces...)
				log.Println(p2.Adcd, p2.Name, len(p2.Items))
			}
		}

		if len(qdata.Provinces) > 0 {
			ds_json.Save("./q/provinces.json", qdata)
		}
	} else {
		ds_json.ReadFile("./q/provinces.json", qdata)
	}

	var r func(qdata *QMenu)
	r = func(qdata *QMenu) {
		for _, p := range qdata.Provinces {
			for _, o := range qdata.Options {
				typeID := o.Value
				// 调查评价总体汇总表需要额外处理
				if typeID == "STS_IA_PROV_" {
					typeID = "STS_IA_COUNTY_"
				}
				fname := fmt.Sprintf("./q/%s/%s_data.json", typeID, p.Adcd)
				log.Println(fname)
				total := 800000
				if ds.IsExist(fname) {
					data := ds_json.ReadGJSON(fname)
					ftotal := int(data.Get("total").Int())
					if ftotal != 0 {
						total = ftotal
						if total > 4e5 {
							for i := range 100 {
								log.Println("数据量非常庞大,暂不抓取", i, total)
							}
							continue
						}
						rows := data.Get("rows").Array()
						// 他们的数据总数和具体的行数经常有误
						if total-len(rows) > 0 && total > 4000 {
							log.Println("没采集完", total, len(rows))
						} else {
							log.Println("文件已存在")
							continue
						}
					}

				}
				// 城市级表格数据
				url := fmt.Sprintf("http://10.135.6.98/fdms/pro/table/%s%s/data", typeID, p.Adcd)
				log.Println(url)
				var data string

				vm.Eval(fmt.Sprintf(`fetch("%s",{
						"headers": {
							"accept": "application/json, text/javascript, */*; q=0.01",
							"accept-language": "zh-CN,zh;q=0.9",
							"cache-control": "no-cache",
							"content-type": "application/x-www-form-urlencoded; charset=UTF-8",
							"pragma": "no-cache",
							"x-requested-with": "XMLHttpRequest"
						},
						"referrer": "http://10.135.6.98/fdms/q",
						"body": "page=1&rows=%d",
						"method": "POST",
						"mode": "cors",
						"credentials": "include"
					}).then(res=>res.text())`, url, total), &data).Run()
				t := int(gjson.Get(data, "total").Num)
				if data != "" && t == 0 {
					if len(data) > 100 {
						log.Println(total, string(data[0:100]))
					} else {
						log.Println(total, data)
					}
				}
				rows := gjson.Get(data, "rows").Array()
				nt := len(rows)
				if nt < t && t <= total {
					result, err := sjson.Set(data, "total", nt)
					if err != nil {
						log.Println(err)
					}
					data = result
				}
				log.Println("采集数据量", nt, t, total)

				if data != "" {
					ds.Save(fname, []byte(data))
				}
			}

			if len(p.Items) > 0 {
				r(&QMenu{
					Provinces: p.Items,
					Options:   qdata.Options,
				})
			}
		}
	}
	// 暂不采集
	// r(qdata)
	// 小流域数据
	var r2 func(qdata *QMenu)
	r2 = func(qdata *QMenu) {
		for _, p := range qdata.Provinces {
			fname := fmt.Sprintf("./q/小流域(补充)/%s.json", p.Adcd)
			if ds.IsExist(fname) {
				continue
			}
			// http://10.135.6.98/fdms/pages/search/query/getVillageInfo/350111100212000
			var data string
			url := fmt.Sprintf("http://10.135.6.98/fdms/pages/search/query/getVillageInfo/%s", p.Adcd)
			vm.Eval(fmt.Sprintf(`fetch("%s", {
  "headers": {
    "accept": "*/*",
    "accept-language": "zh-CN,zh;q=0.9",
    "cache-control": "no-cache",
    "pragma": "no-cache",
    "x-requested-with": "XMLHttpRequest"
  },
  "referrer": "http://10.135.6.98/fdms/pages/frameui/modules/xianjitongji/huizongbiao/model/village.jsp?adcd=%s",
  "body": null,
  "method": "POST",
  "mode": "cors",
  "credentials": "include"
}).then(res=>res.text())`, url, p.Adcd), &data).Run()

			log.Println(fname)
			if data != "" {
				ds.Save(fname, []byte(data))
			}
			if len(p.Items) > 0 {
				r2(&QMenu{
					Provinces: p.Items,
					Options:   qdata.Options,
				})
			}

		}

	}
	r2(qdata)
}

// function shengChange(itemId, src) {
//     var name = getPicName(itemId)
//     if (name === "") {
//         $("#title").html("");
//         $("#qgt").attr("src", src + getPicName(itemId) + ".bmp");//pic/11/
//     } else {
//         $("#qgt").attr("src", src + getPicName(itemId) + ".bmp");//pic/11/
//     }
// }

//	function getPicName(layer) {
//	    var name = "";
//	    switch (layer) {
//	        case "layer7":
//	            name = "1山洪灾害防治村分布图";
//	            break;
//	        case "layer5":
//	            name = "2山洪灾害防治区企事业单位分布图";
//	            break;
//	        case "layer0":
//	            name = "3山洪灾害防治区分布图";
//	            break;
//	        case "layer1":
//	            name = "4山洪灾害危险区分布图";
//	            break;
//	        case "layer6":
//	            name = "5历史山洪灾害分布图";
//	            break;
//	        case "layer2":
//	            name = "6自动监测站分布图";
//	            break;
//	        case "layer4":
//	            name = "7山洪灾害监测预警设施分布图";
//	            break;
//	        case "layer3":
//	            name = "8山洪灾害防治区涉水工程分布图";
//	            break;
//	        case "layer8":
//	            name = "9需防洪治理山洪沟分布图";
//	            break;
//	        case "layer11":
//	            name = "10分析评价村现状防洪能力成果分布图";
//	            break;
//	        case "layer10":
//	            name = "11二十年一遇设计暴雨（1小时）分布图";
//	            break;
//	        case "layer14":
//	            name = "12临界雨量（1小时）重现期分布图";
//	            break;
//	        case "layer9":
//	            name = "13临界雨量（1小时）分布图";
//	            break;
//	        case "layer12":
//	            name = "14小流域洪峰模数分布图";
//	            break;
//	        case "layer13":
//	            name = "15小流域汇流时间分布图";
//	            break;
//	    }
//	    return name;
//	}

// 基础图+图例+arcgis服务json叠加+arcgis服务
func ReptileT() {

}

// <div class="utils">
//
//	<button class="btn btn-success btn-circle btn-lg " type="button" title="村貌" id="01">村貌
//	</button>
//	<button class="btn btn-success btn-circle btn-lg" type="button" title="房屋分类" id="03">房屋分类
//	</button>
//	<button class="btn btn-success btn-circle btn-lg" type="button" title="重要城集镇" id="63">重要城集镇
//	</button>
//	<button class="btn btn-success btn-circle btn-lg" type="button" title="沿河村落居民" id="62">沿河村落居民
//	</button>
//	<button class="btn btn-success btn-circle btn-lg" type="button" title="企事业单位" id="53">企事业单位
//	</button>
//	<button class="btn btn-success btn-circle btn-lg" type="button" title="历史洪痕" id="71">历史洪痕
//	</button>
//	<button class="btn btn-success btn-circle btn-lg" type="button" title="水闸工程" id="34">水闸工程
//	</button>
//	<button class="btn btn-success btn-circle btn-lg" type="button" title="水库工程" id="33">水库工程
//	</button>
//	<button class="btn btn-success btn-circle btn-lg" type="button" title="桥梁工程" id="38">桥梁工程
//	</button>
//	<button class="btn btn-success btn-circle btn-lg" type="button" title="路涵工程" id="37">路涵工程
//	</button>
//	<button class="btn btn-success btn-circle btn-lg" type="button" title="塘堰坝工程" id="36">塘堰坝工程
//	</button>
//	<button class="btn btn-success btn-circle btn-lg" type="button" title="横断面" id="10">横断面
//	</button>
//	<button class="btn btn-success btn-circle btn-lg" type="button" title="纵断面" id="07">纵断面
//	</button>
//
// </div>
// Array.from(temp1.querySelectorAll('button')).map(e=>({type:e.id,name:e.innerText}))
//
// [{"type":"01","name":"村貌"},{"type":"03","name":"房屋分类"},{"type":"63","name":"重要城集镇"},{"type":"62","name":"沿河村落居民"},{"type":"53","name":"企事业单位"},{"type":"71","name":"历史洪痕"},{"type":"34","name":"水闸工程"},{"type":"33","name":"水库工程"},{"type":"38","name":"桥梁工程"},{"type":"37","name":"路涵工程"},{"type":"36","name":"塘堰坝工程"},{"type":"10","name":"横断面"},{"type":"07","name":"纵断面"}]
type MImages struct {
	Total int    `json:"total"`
	Rows  []Rows `json:"rows"`
}
type Rows struct {
	Pid       int    `json:"PID"`
	Objpid    string `json:"OBJPID"`
	Objtp     string `json:"OBJTP"`
	Objnm     any    `json:"OBJNM"`
	Adcd      string `json:"ADCD"`
	Fpath     string `json:"FPATH"`
	Lgtd      any    `json:"LGTD"`
	Lttd      any    `json:"LTTD"`
	Ptime     any    `json:"PTIME"`
	Fname     string `json:"FNAME"`
	Multitype any    `json:"MULTITYPE"`
	Flag      int    `json:"FLAG"`
	Comments  any    `json:"COMMENTS"`
	Moditime  any    `json:"MODITIME"`
	Status    string `json:"STATUS"`
	GUID      string `json:"GUID"`
	Signer    any    `json:"SIGNER"`
	Audbatch  string `json:"AUDBATCH"`
	Cadcd     string `json:"CADCD"`
	Rn        int    `json:"RN"`
	Name      string `json:"NAME"`
}

type TypeInfo struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// ReptileM
func ReptileM() {
	qdata := new(QMenu)
	if !ds.IsExist("./m/provinces.json") {
		vm := vesta.New().Head().Nav("http://10.135.6.98")
		vm.Eval(`
		 document.querySelector("#username").value = "福建省操作员"
		 document.querySelector("#password").value = "czyyh!123"
		 // loginFun()
	`, nil).Run()

		time.Sleep(10 * time.Second)
		//
		vm.Nav("http://10.135.6.98/fdms/m").Run()
		time.Sleep(5 * time.Second)
		getProvincesScript := `
  (function(){
			let items = Array.from(window.frames[0].window.frames[0].document.querySelectorAll(".provinceElement"))
			// 省份
			let provinces = items.map(e=>{
				let capitalA = e.children[0].querySelector('a')
				// 这里只包含所有数据
				let citylistA = Array.from(e.children[1].querySelectorAll('li'))
				return {
					// 省份名称
					name:capitalA.innerText,
					adcd:capitalA.id || capitalA.getAttribute('adcdcode'),
					hasData:true,
					// 城市名称
					items:citylistA.map(e=>{
					  let adcd = e.getAttribute('adcdcode')
						let a = e.querySelector('a')
						let hasData = !!a;
						if (!adcd){
							adcd = a.id || a.getAttribute('adcdcode')
						}
						return {
							name:e.innerText,
							hasData,
							adcd,
							items:[]
						}
					})
				}
			})
			let qdata =  {
				provinces,
				options:[],
			};
			return qdata
		})()
	`

		vm.Eval(getProvincesScript, qdata).Run()
		log.Println("城市数量", len(qdata.Provinces))

		for _, pp := range qdata.Provinces {
			for _, p2 := range pp.Items {
				if !p2.HasData {
					log.Println(p2.Adcd, p2.Name, "无数据")
					continue
				}
				qdata2 := new(QMenu)
				vm.Eval(fmt.Sprintf(`window.frames[0].$("#cityList").attr("src", "/fdms/pages/search/query/getSubinfoByPcode1/%s")`, p2.Adcd), nil).Run()
				time.Sleep(2 * time.Second)
				vm.Eval(getProvincesScript, qdata2).Run()
				p2.Items = append(p2.Items, qdata2.Provinces...)
				log.Println(p2.Adcd, p2.Name, len(p2.Items))
			}
		}
		ds_json.Save("./m/provinces.json", qdata)
	} else {
		ds_json.ReadFile("./m/provinces.json", qdata)
	}
	//
	typeInfos := new([]TypeInfo)
	json.Unmarshal([]byte(`[{"type":"01","name":"村貌"},{"type":"03","name":"房屋分类"},{"type":"63","name":"重要城集镇"},{"type":"62","name":"沿河村落居民"},{"type":"53","name":"企事业单位"},{"type":"71","name":"历史洪痕"},{"type":"34","name":"水闸工程"},{"type":"33","name":"水库工程"},{"type":"38","name":"桥梁工程"},{"type":"37","name":"路涵工程"},{"type":"36","name":"塘堰坝工程"},{"type":"10","name":"横断面"},{"type":"07","name":"纵断面"}]`), typeInfos)
	//
	qp := slices.Clone(qdata.Provinces)
	slices.Reverse(qp)
	for _, pppp := range qp[4:] { // 省份
		slices.Reverse(pppp.Items)
		for _, ppp := range pppp.Items { // 城市
			for _, pp := range ppp.Items { // 县区
				for _, p := range pp.Items { // 乡镇，街道
					log.Println(p.Adcd, p.Name, p.HasData)
					if !p.HasData {
						continue
					}
					for _, t := range *typeInfos {
						configName := fmt.Sprintf("./m/多媒体资料/%s/%s/config.json", p.Adcd, t.Type)
						res := new(MImages)
						if ds.IsExist(configName) {
							ds_json.ReadFile(configName, res)
						} else {
							var paSize string = "9000"
							var pageNo string = "1"
							u := fmt.Sprintf("http://10.135.6.98/fdms/pages/search/query/getMediaForPage/%s/%s/%s/%s/", p.Adcd, t.Type, paSize, pageNo)
							data, err := ctp.Get(u)
							if err != nil {
								panic(err)
							}
							err = json.Unmarshal(data, res)
							if err != nil {
								panic(err) //解析失败
							}
							ds.Save(configName, data)
							if res.Total > 0 {
								log.Println(len(res.Rows), res.Total)
							}

						}

						for i, r := range res.Rows {
							fname := strings.ReplaceAll(r.Fpath, "\\", "/")
							src := fmt.Sprintf("./m/多媒体资料/%s/%s/%s.jpg", p.Adcd, t.Type, fname)
							if !ds.IsExist(src) {
								// encodeURI(encodeURI(r.Fpath))
								// 没办法，目标站点太奇葩
								curl := "http://10.135.6.98/fdms/pages/search/query/getImg/" + window.EncodeURI(window.EncodeURI(r.Fpath)) + "/JPG"
								data, err := ctp.Get(curl)
								if err == nil && len(data) > 0 {
									go ds.Save(src, data)
									log.Println("保存成功", i, r.Fpath)
								} else {
									log.Println("失败", curl, err, r.Fpath, len(data))
								}
							}
						}
					}
				}
			}
		}
	}
}

type WData struct {
	Total int      `json:"total"`
	Rows  []*WRows `json:"rows"`
}
type Modtime struct {
	Date           int   `json:"date"`
	Day            int   `json:"day"`
	Hours          int   `json:"hours"`
	Minutes        int   `json:"minutes"`
	Month          int   `json:"month"`
	Nanos          int   `json:"nanos"`
	Seconds        int   `json:"seconds"`
	Time           int64 `json:"time"`
	TimezoneOffset int   `json:"timezoneOffset"`
	Year           int   `json:"year"`
}
type WRows struct {
	Cadcd   string   `json:"CADCD"`
	Adnm    string   `json:"ADNM"`
	Dnm     string   `json:"DNM"`
	Dpath   string   `json:"DPATH"`
	Dtype   any      `json:"DTYPE"`
	Modtime *Modtime `json:"MODTIME"`
	Sunit   string   `json:"SUNIT"`
	Rn      int      `json:"RN"`
}

func ReptileW() {
	vm := vesta.New().Head().Nav("http://10.135.6.98")
	vm.Eval(fmt.Sprintf(`
		 document.querySelector("#username").value = "福建省操作员"
		 document.querySelector("#password").value = "czyyh!123"

		 function loadAndRunTesseract() {
					// 1. 创建type="module"的script标签
					const script = document.createElement('script');
					script.type = 'module';
					// 2. 设置script内容（内联模块代码）
					script.innerHTML = %s;
			// 3. 添加到文档中执行
			document.head.appendChild(script);
		}
		loadAndRunTesseract()
		// loginFun()
	`, "`"+`
			import tesseract from 'https://cdn.jsdelivr.net/npm/tesseract.js@5/dist/tesseract.esm.min.js';
			const { createWorker } = tesseract;
			const rec = (async () => {
				const worker = await createWorker('eng');
				// const ret = await worker.recognize('https://tesseract.projectnaptha.com/img/eng_bw.png');
				const t = new Date().toISOString().slice(0,-1)
				const ret = await worker.recognize('http://10.135.6.98/Account/GetValidateCode?time='+t+'&browserName=Safari')
				document.querySelector("#CheckCode").value = ret.data.text
				console.log(ret.data.text);
				loginFun()
				// 
				await worker.terminate();
				setTimeout(()=>{
					console.log('重新识别')
					rec()
				},100)
			})();
    `+"`"), nil).Run()
	time.Sleep(10 * time.Second)
	vm.Nav("http://10.135.6.98/fdms/w")
	qdata := new(QMenu)

	getProvincesScript := `
  (function(){
			let items = Array.from(window.frames[0].window.frames[0].document.querySelectorAll(".provinceElement"))
			// 省份
			let provinces = items.map(e=>{
				let capitalA = e.children[0].querySelector('a')
				// 这里只包含所有数据
				let citylistA = Array.from(e.children[1].querySelectorAll('li'))
				return {
					// 省份名称
					name:capitalA.innerText,
					adcd:capitalA.id || capitalA.getAttribute('adcdcode'),
					hasData:true,
					// 城市名称
					items:citylistA.map(e=>{
					  let adcd = e.getAttribute('adcdcode')
						let a = e.querySelector('a')
						let hasData = !!a;
						if (!adcd){
							adcd = a.id || a.getAttribute('adcdcode')
						}
						return {
							name:e.innerText,
							hasData,
							adcd,
							items:[]
						}
					})
				}
			})
			let qdata =  {
				provinces,
				options:[],
			};
			return qdata
		})()
	`

	vm.Eval(getProvincesScript, qdata).Run()

	log.Println(len(qdata.Provinces))
	if len(qdata.Provinces) > 0 {
		ds_json.Save("./w/config.json", qdata)
	}
	var r func(p *Province)
	r = func(p *Province) {
		fname := fmt.Sprintf("./w/%s/config.json", p.Adcd)
		data := new(WData)
		if ds.IsExist(fname) {
			ds_json.ReadFile(fname, data)

		} else {
			url := fmt.Sprintf("http://10.135.6.98/fdms/pages/search/query/getDoc/%s", p.Adcd)

			// 只要城市就够了，城市有统计所有
			vm.Eval(fmt.Sprintf(`fetch("%s",{
  "headers": {
    "accept": "application/json, text/javascript, */*; q=0.01",
    "accept-language": "zh-CN,zh;q=0.9",
    "cache-control": "no-cache",
    "content-type": "application/x-www-form-urlencoded; charset=UTF-8",
    "pragma": "no-cache",
    "x-requested-with": "XMLHttpRequest"
  },
  "referrer": "http://10.135.6.98",
  "body": "page=1&rows=9000",
  "method": "POST",
  "mode": "cors",
  "credentials": "include"
}).then(res=>res.json())`, url), data).Run()

			if len(data.Rows) > 0 {
				ds_json.Save(fname, data)
			}
		}
		log.Println(fname)
		for _, row := range data.Rows {
			fname := strings.ReplaceAll(row.Dpath, "\\", "/")
			src := fmt.Sprintf("./w/%s/文档%s.pdf", p.Adcd, fname)
			log.Println(fname)
			if ds.IsExist(src) {
				continue
			}
			pname := strings.ReplaceAll(row.Dpath, "\\", "~")
			uname := window.EncodeURI(window.EncodeURI(pname))
			// http://10.135.6.98/fdms/pages/search/query/readDocs/~35-%25E7%25A6%258F%25E5%25BB%25BA~%25E9%25A9%25AC%25E5%25B0%25BE%25E5%258C%25BA~350105_%25E9%25A9%25AC%25E5%25B0%25BE%25E5%258C%25BA2013-2015%25E5%25B9%25B4%25E5%25BA%25A6%25E5%25B1%25B1%25E6%25B4%25AA%25E7%2581%25BE%25E5%25AE%25B3%25E5%2588%2586%25E6%259E%2590%25E8%25AF%2584%25E4%25BB%25B7%25E6%258A%25A5%25E5%2591%258A(%25E6%258A%25A5%25E6%2589%25B9%25E7%25A8%25BF)_%25E7%25A6%258F%25E5%25B7%259E%25E5%25B8%2582%25E6%25B0%25B4%25E5%2588%25A9%25E6%25B0%25B4%25E7%2594%25B5%25E5%25BC%2580%25E5%258F%2591%25E5%2585%25AC%25E5%258F%25B8_201611.doc/PDF
			curl := "http://10.135.6.98/fdms/pages/search/query/readDocs/" + uname + "/PDF"

			data, err := ctp.Get(curl)
			if err != nil || len(data) == 0 {
				log.Println("请求失败", curl, err, pname, len(data))
			} else {
				ds.Save(src, data)
			}

		}

	}
	for _, p := range qdata.Provinces { //省份
		r(p)
	}

}
