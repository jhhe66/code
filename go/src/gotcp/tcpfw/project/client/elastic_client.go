package main

import (
	// "github.com/bitly/go-simplejson"
	// "database/sql"

	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
	"gopkg.in/olivere/elastic.v3"
)

type Options struct {
	// Example of verbosity with level
	Verbose []bool `short:"v" long:"verbose" description:"Verbose output"`

	// Example of optional value
	ClientConf string `short:"c" long:"conf" description:"Clinet Config" optional:"no"`
}

type DistributeIterm struct {
	TotalTryCount    int64
	TotalFailedCount int64
	CountryName      string
}

var options Options
var parser = flags.NewParser(&options, flags.Default)
var staticCount = 50

var infoLog *log.Logger

var NickToName map[string]string = map[string]string{
	"AD": "安道尔共和国",
	"AE": "阿拉伯联合酋长国",
	"AF": "阿富汗",
	"AG": "安提瓜和巴布达",
	"AI": "安圭拉岛",
	"AL": "阿尔巴尼亚",
	"AM": "亚美尼亚",
	"AO": "安哥拉",
	"AR": "阿根廷",
	"AT": "奥地利",
	"AU": "澳大利亚",
	"AZ": "阿塞拜疆",
	"BB": "巴巴多斯",
	"BD": "孟加拉国",
	"BE": "比利时",
	"BF": "布基纳法索",
	"BG": "保加利亚",
	"BH": "巴林",
	"BI": "布隆迪",
	"BJ": "贝宁",
	"BL": "巴勒斯坦",
	"BM": "百慕大群岛",
	"BN": "文莱",
	"BO": "玻利维亚",
	"BR": "巴西",
	"BS": "巴哈马",
	"BW": "博茨瓦纳",
	"BY": "白俄罗斯",
	"BZ": "伯利兹",
	"CA": "加拿大",
	"CF": "中非共和国",
	"CG": "刚果",
	"CH": "瑞士",
	"CK": "库克群岛",
	"CL": "智利",
	"CM": "喀麦隆",
	"CN": "中国",
	"CO": "哥伦比亚",
	"CR": "哥斯达黎加",
	"CS": "捷克",
	"CU": "古巴",
	"CY": "塞浦路斯",
	"CZ": "捷克",
	"DE": "德国",
	"DJ": "吉布提",
	"DK": "丹麦",
	"DO": "多米尼加共和国",
	"DZ": "阿尔及利亚",
	"EC": "厄瓜多尔",
	"EE": "爱沙尼亚",
	"EG": "埃及",
	"ES": "西班牙",
	"ET": "埃塞俄比亚",
	"FI": "芬兰",
	"FJ": "斐济",
	"FR": "法国",
	"GA": "加蓬",
	"GB": "英国",
	"GD": "格林纳达",
	"GE": "格鲁吉亚",
	"GF": "法属圭亚那",
	"GH": "加纳",
	"GI": "直布罗陀",
	"GM": "冈比亚",
	"GN": "几内亚",
	"GR": "希腊",
	"GT": "危地马拉",
	"GU": "关岛",
	"GY": "圭亚那",
	"HK": "香港特别行政区",
	"HN": "洪都拉斯",
	"HT": "海地",
	"HU": "匈牙利",
	"ID": "印度尼西亚",
	"IE": "爱尔兰",
	"IL": "以色列",
	"IN": "印度",
	"IQ": "伊拉克",
	"IR": "伊朗",
	"IS": "冰岛",
	"IT": "意大利",
	"JM": "牙买加",
	"JO": "约旦",
	"JP": "日本",
	"KE": "肯尼亚",
	"KG": "吉尔吉斯坦",
	"KH": "柬埔寨",
	"KP": "朝鲜",
	"KR": "韩国",
	"KT": "科特迪瓦共和国",
	"KW": "科威特",
	"KZ": "哈萨克斯坦",
	"LA": "老挝",
	"LB": "黎巴嫩",
	"LC": "圣卢西亚",
	"LI": "列支敦士登",
	"LK": "斯里兰卡",
	"LR": "利比里亚",
	"LS": "莱索托",
	"LT": "立陶宛",
	"LU": "卢森堡",
	"LV": "拉脱维亚",
	"LY": "利比亚",
	"MA": "摩洛哥",
	"MC": "摩纳哥",
	"MD": "摩尔多瓦",
	"MG": "马达加斯加",
	"ML": "马里",
	"MM": "缅甸",
	"MN": "蒙古",
	"MO": "澳门",
	"MS": "蒙特塞拉特岛",
	"MT": "马耳他",
	"MU": "毛里求斯",
	"MV": "马尔代夫",
	"MW": "马拉维",
	"MX": "墨西哥",
	"MY": "马来西亚",
	"MZ": "莫桑比克",
	"NA": "纳米比亚",
	"NE": "尼日尔",
	"NG": "尼日利亚",
	"NI": "尼加拉瓜",
	"NL": "荷兰",
	"NO": "挪威",
	"NP": "尼泊尔",
	"NR": "瑙鲁",
	"NZ": "新西兰",
	"OM": "阿曼",
	"PA": "巴拿马",
	"PE": "秘鲁",
	"PF": "法属玻利尼西亚",
	"PG": "巴布亚新几内亚",
	"PH": "菲律宾",
	"PK": "巴基斯坦",
	"PL": "波兰",
	"PR": "波多黎各",
	"PT": "葡萄牙",
	"PY": "巴拉圭",
	"QA": "卡塔尔",
	"RO": "罗马尼亚",
	"RU": "俄罗斯",
	"SA": "沙特阿拉伯",
	"SB": "所罗门群岛",
	"SC": "塞舌尔",
	"SD": "苏丹",
	"SE": "瑞典",
	"SG": "新加坡",
	"SI": "斯洛文尼亚",
	"SK": "斯洛伐克",
	"SL": "塞拉利昂",
	"SM": "圣马力诺",
	"SN": "塞内加尔",
	"SO": "索马里",
	"SR": "苏里南",
	"ST": "圣多美和普林西比",
	"SV": "萨尔瓦多",
	"SY": "叙利亚",
	"SZ": "斯威士兰",
	"TD": "乍得",
	"TG": "多哥",
	"TH": "泰国",
	"TJ": "塔吉克斯坦",
	"TM": "土库曼斯坦",
	"TN": "突尼斯",
	"TO": "汤加",
	"TR": "土耳其",
	"TT": "特立尼达和多巴哥",
	"TW": "台湾省",
	"TZ": "坦桑尼亚",
	"UA": "乌克兰",
	"UG": "乌干达",
	"US": "美国",
	"UY": "乌拉圭",
	"UZ": "乌兹别克斯坦",
	"VC": "圣文森特岛",
	"VE": "委内瑞拉",
	"VN": "越南",
	"YE": "也门",
	"YU": "南斯拉夫",
	"ZA": "南非",
	"ZM": "赞比亚",
	"ZR": "扎伊尔",
	"ZW": "津巴布韦",
}

type ResultIterm struct {
	key   string `json:"key"`
	count string `json:"doc_count"`
}

func GetDayOfMonth(year, month int) int {
	switch month {
	case 1:
		return 31
	case 2:
		if (year%4 == 0 && year%100 != 0) || (year%400 == 0) {
			return 29
		} else {
			return 28
		}
	case 3:
		return 31
	case 4:
		return 30
	case 5:
		return 31
	case 6:
		return 30
	case 7:
		return 31
	case 8:
		return 31
	case 9:
		return 30
	case 10:
		return 31
	case 11:
		return 30
	case 12:
		return 31
	default:
		return 30
	}
}

func GetFirstReconnectData(client *elastic.Client, year, month, day int) {
	strTargetDay := fmt.Sprintf("%4d.%02d.%02d", year, month, day)
	tsBegin := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Unix()
	tsBegin = tsBegin * 1000 // change to ms
	tsEnd := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.Local).Unix()
	tsEnd = tsEnd * 1000 // change to ms
	infoLog.Printf("tsBegin=%v tsEnd=%v", tsBegin, tsEnd)

	indexBegin := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day-1)

	indexEnd := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day)
	infoLog.Printf("indexBegegin=%s indexEnd=%s", indexBegin, indexEnd)
	rangeQuery := elastic.NewRangeQuery("@timestamp")
	rangeQuery = rangeQuery.Gte(tsBegin)
	rangeQuery = rangeQuery.Lte(tsEnd)
	rangeQuery = rangeQuery.Format("epoch_millis")
	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(rangeQuery)

	strTotalQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND retry:0")
	strTotalQuery = strTotalQuery.AnalyzeWildcard(true)
	totalQuery := elastic.NewBoolQuery()
	totalQuery = totalQuery.Must(strTotalQuery)
	totalQuery = totalQuery.Filter(boolQuery)
	totalResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(totalQuery). // return all results, but ...
		Pretty(true).      // pretty print request and response JSON
		Do()               // execute
	if err != nil {
		infoLog.Println("client.Search failed err=", err)
		return
	}
	totalCount := totalResult.TotalHits()

	strAggQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND retry:0 AND NOT wns_code:0")
	strAggQuery = strAggQuery.AnalyzeWildcard(true)
	aggQuery := elastic.NewBoolQuery()
	aggQuery = aggQuery.Must(strAggQuery)
	aggQuery = aggQuery.Filter(boolQuery)

	termsAgg := elastic.NewTermsAggregation()
	termsAgg = termsAgg.Field("geoip.country_code3.raw")
	termsAgg = termsAgg.Size(30)
	termsAgg = termsAgg.OrderByCountDesc()

	searchResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(aggQuery).            // return all results, but ...
		Aggregation("2", termsAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	failTotalCount := searchResult.TotalHits()
	infoLog.Printf("================================================================================================")
	infoLog.Printf("%s 首次重练总次数: %d 首次重练失败总次数:%d 失败占比:%f%%", strTargetDay, totalCount, failTotalCount, float32(failTotalCount)/float32(totalCount)*100)
	infoLog.Printf("------------------------------------------------------------------------------------------------")
	bucketKeyIterm, ok := searchResult.Aggregations.Terms("2")
	if ok {
		for i, iterm := range bucketKeyIterm.Buckets {
			infoLog.Printf("index:%2v 国家:%s 失败次数:%5v 占比:%f%%", i, iterm.Key.(string), iterm.DocCount, float32(iterm.DocCount)/float32(totalCount)*100)
		}
	} else {
		infoLog.Printf("searchResult.Aggregations.Terms trans err")
	}
	infoLog.Printf("================================================================================================")
}

func GetIosFirstReconnectData(client *elastic.Client, year, month, day int) {
	strTargetDay := fmt.Sprintf("%4d.%02d.%02d", year, month, day)
	tsBegin := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Unix()
	tsBegin = tsBegin * 1000 // change to ms
	tsEnd := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.Local).Unix()
	tsEnd = tsEnd * 1000 // change to ms
	infoLog.Printf("tsBegin=%v tsEnd=%v", tsBegin, tsEnd)

	indexBegin := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day-1)

	indexEnd := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day)
	infoLog.Printf("indexBegegin=%s indexEnd=%s", indexBegin, indexEnd)
	rangeQuery := elastic.NewRangeQuery("@timestamp")
	rangeQuery = rangeQuery.Gte(tsBegin)
	rangeQuery = rangeQuery.Lte(tsEnd)
	rangeQuery = rangeQuery.Format("epoch_millis")
	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(rangeQuery)

	strTotalQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND retry:0 AND terminaltype:0")
	strTotalQuery = strTotalQuery.AnalyzeWildcard(true)
	totalQuery := elastic.NewBoolQuery()
	totalQuery = totalQuery.Must(strTotalQuery)
	totalQuery = totalQuery.Filter(boolQuery)
	totalResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(totalQuery). // return all results, but ...
		Pretty(true).      // pretty print request and response JSON
		Do()               // execute
	if err != nil {
		infoLog.Println("client.Search failed err=", err)
		return
	}
	totalCount := totalResult.TotalHits()

	strAggQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND retry:0 AND terminaltype:0 AND NOT wns_code:0")
	strAggQuery = strAggQuery.AnalyzeWildcard(true)
	aggQuery := elastic.NewBoolQuery()
	aggQuery = aggQuery.Must(strAggQuery)
	aggQuery = aggQuery.Filter(boolQuery)

	termsAgg := elastic.NewTermsAggregation()
	termsAgg = termsAgg.Field("geoip.country_code3.raw")
	termsAgg = termsAgg.Size(30)
	termsAgg = termsAgg.OrderByCountDesc()

	searchResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(aggQuery).            // return all results, but ...
		Aggregation("2", termsAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	failTotalCount := searchResult.TotalHits()
	infoLog.Printf("=====================================================================================================")
	infoLog.Printf("%s iOS 首次重练总次数: %d iOS 首次重练失败总次数:%d 失败占比:%f%%", strTargetDay, totalCount, failTotalCount, float32(failTotalCount)/float32(totalCount)*100)
	infoLog.Printf("-----------------------------------------------------------------------------------------------------")
	bucketKeyIterm, ok := searchResult.Aggregations.Terms("2")
	if ok {
		for i, iterm := range bucketKeyIterm.Buckets {
			infoLog.Printf("index:%2v 国家:%s 失败次数:%5v 占比:%f%%", i, iterm.Key.(string), iterm.DocCount, float32(iterm.DocCount)/float32(totalCount)*100)
		}
	} else {
		infoLog.Printf("searchResult.Aggregations.Terms trans err")
	}
	infoLog.Printf("=====================================================================================================")
}

func GetAndroidFirstReconnectData(client *elastic.Client, year, month, day int) {
	strTargetDay := fmt.Sprintf("%4d.%02d.%02d", year, month, day)
	tsBegin := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Unix()
	tsBegin = tsBegin * 1000 // change to ms
	tsEnd := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.Local).Unix()
	tsEnd = tsEnd * 1000 // change to ms
	infoLog.Printf("tsBegin=%v tsEnd=%v", tsBegin, tsEnd)

	indexBegin := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day-1)

	indexEnd := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day)
	infoLog.Printf("indexBegegin=%s indexEnd=%s", indexBegin, indexEnd)
	rangeQuery := elastic.NewRangeQuery("@timestamp")
	rangeQuery = rangeQuery.Gte(tsBegin)
	rangeQuery = rangeQuery.Lte(tsEnd)
	rangeQuery = rangeQuery.Format("epoch_millis")
	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(rangeQuery)

	strTotalQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND retry:0 AND NOT terminaltype:0")
	strTotalQuery = strTotalQuery.AnalyzeWildcard(true)
	totalQuery := elastic.NewBoolQuery()
	totalQuery = totalQuery.Must(strTotalQuery)
	totalQuery = totalQuery.Filter(boolQuery)
	totalResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(totalQuery). // return all results, but ...
		Pretty(true).      // pretty print request and response JSON
		Do()               // execute
	if err != nil {
		infoLog.Println("client.Search failed err=", err)
		return
	}
	totalCount := totalResult.TotalHits()

	strAggQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND retry:0 AND NOT terminaltype:0 AND NOT wns_code:0")
	strAggQuery = strAggQuery.AnalyzeWildcard(true)
	aggQuery := elastic.NewBoolQuery()
	aggQuery = aggQuery.Must(strAggQuery)
	aggQuery = aggQuery.Filter(boolQuery)

	termsAgg := elastic.NewTermsAggregation()
	termsAgg = termsAgg.Field("geoip.country_code3.raw")
	termsAgg = termsAgg.Size(30)
	termsAgg = termsAgg.OrderByCountDesc()

	searchResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(aggQuery).            // return all results, but ...
		Aggregation("2", termsAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	failTotalCount := searchResult.TotalHits()
	infoLog.Printf("============================================================================================================")
	infoLog.Printf("%s Android 首次重练总次数: %d Android 首次重练失败总次数:%d 失败占比:%f%%", strTargetDay, totalCount, failTotalCount, float32(failTotalCount)/float32(totalCount)*100)
	infoLog.Printf("------------------------------------------------==============================------------------------------")
	bucketKeyIterm, ok := searchResult.Aggregations.Terms("2")
	if ok {
		for i, iterm := range bucketKeyIterm.Buckets {
			infoLog.Printf("index:%2v 国家:%s 失败次数:%5v 占比:%f%%", i, iterm.Key.(string), iterm.DocCount, float32(iterm.DocCount)/float32(totalCount)*100)
		}
	} else {
		infoLog.Printf("searchResult.Aggregations.Terms trans err")
	}
	infoLog.Printf("============================================================================================================")
}

func GetFirstReconnectDistribute(client *elastic.Client, year, month, day int) {
	strTargetDay := fmt.Sprintf("%4d.%02d.%02d", year, month, day)
	tsBegin := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Unix()
	tsBegin = tsBegin * 1000 // change to ms
	tsEnd := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.Local).Unix()
	tsEnd = tsEnd * 1000 // change to ms
	infoLog.Printf("tsBegin=%v tsEnd=%v", tsBegin, tsEnd)

	indexBegin := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day-1)

	indexEnd := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day)
	infoLog.Printf("indexBegegin=%s indexEnd=%s", indexBegin, indexEnd)
	rangeQuery := elastic.NewRangeQuery("@timestamp")
	rangeQuery = rangeQuery.Gte(tsBegin)
	rangeQuery = rangeQuery.Lte(tsEnd)
	rangeQuery = rangeQuery.Format("epoch_millis")
	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(rangeQuery)

	termsAgg := elastic.NewTermsAggregation()
	termsAgg = termsAgg.Field("geoip.country_code2.raw")
	termsAgg = termsAgg.Size(staticCount)
	termsAgg = termsAgg.OrderByCountDesc()

	strTotalQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND retry:0")
	strTotalQuery = strTotalQuery.AnalyzeWildcard(true)
	totalQuery := elastic.NewBoolQuery()
	totalQuery = totalQuery.Must(strTotalQuery)
	totalQuery = totalQuery.Filter(boolQuery)
	totalResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(totalQuery).          // return all results, but ...
		Aggregation("2", termsAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	if err != nil {
		infoLog.Println("client.Search failed err=", err)
		return
	}
	totalCount := totalResult.TotalHits()
	mapTotalTry := make(map[string]int64, staticCount)
	totalBucketKeyIterm, ok := totalResult.Aggregations.Terms("2")
	if ok {
		for _, iterm := range totalBucketKeyIterm.Buckets {
			mapTotalTry[iterm.Key.(string)] = iterm.DocCount
		}
	} else {
		infoLog.Printf("searchResult.Aggregations.Terms trans err")
	}

	strAggQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND retry:0 AND NOT wns_code:0")
	strAggQuery = strAggQuery.AnalyzeWildcard(true)
	aggQuery := elastic.NewBoolQuery()
	aggQuery = aggQuery.Must(strAggQuery)
	aggQuery = aggQuery.Filter(boolQuery)

	searchResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(aggQuery).            // return all results, but ...
		Aggregation("2", termsAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	failTotalCount := searchResult.TotalHits()
	infoLog.Printf("================================================================================================")
	infoLog.Printf("%s 首次重练总次数: %d 首次重练失败总次数:%d 失败占比:%f%%", strTargetDay, totalCount, failTotalCount, float32(failTotalCount)/float32(totalCount)*100)
	infoLog.Printf("------------------------------------------------------------------------------------------------")
	mapItermTry := make(map[string]int64, staticCount)
	mapRateToCountry := make(map[float64]DistributeIterm, staticCount)
	sliceRate := make([]float64, staticCount)
	bucketKeyIterm, ok := searchResult.Aggregations.Terms("2")
	if ok {
		for _, iterm := range bucketKeyIterm.Buckets {
			mapItermTry[iterm.Key.(string)] = iterm.DocCount
		}
	} else {
		infoLog.Printf("searchResult.Aggregations.Terms trans err")
	}
	var index int = 0
	for k, v := range mapTotalTry {
		if count, ok := mapItermTry[k]; ok {
			rate := float64(count) / float64(v)
			mapRateToCountry[rate] = DistributeIterm{TotalTryCount: v, TotalFailedCount: count, CountryName: k}
			sliceRate[index] = rate
			index++
		} else {
			//infoLog.Printf("国家:%s 总的尝试次数:%5v 没找到总的失败次数", k, v)
		}
	}
	infoLog.Println()
	sort.Float64s(sliceRate)
	for _, k := range sliceRate {
		infoLog.Printf("国家:%s-%-10s\t 总的尝试次数:%7v 总的失败次数:%5v 占比:%8.4f%%",
			mapRateToCountry[k].CountryName,
			NickToName[mapRateToCountry[k].CountryName],
			mapRateToCountry[k].TotalTryCount,
			mapRateToCountry[k].TotalFailedCount,
			k*100)
	}
	infoLog.Printf("================================================================================================")
}

func GetIosFirstReconnectDistribute(client *elastic.Client, year, month, day int) {
	strTargetDay := fmt.Sprintf("%4d.%02d.%02d", year, month, day)
	tsBegin := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Unix()
	tsBegin = tsBegin * 1000 // change to ms
	tsEnd := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.Local).Unix()
	tsEnd = tsEnd * 1000 // change to ms
	infoLog.Printf("tsBegin=%v tsEnd=%v", tsBegin, tsEnd)

	indexBegin := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day-1)

	indexEnd := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day)
	infoLog.Printf("indexBegegin=%s indexEnd=%s", indexBegin, indexEnd)
	rangeQuery := elastic.NewRangeQuery("@timestamp")
	rangeQuery = rangeQuery.Gte(tsBegin)
	rangeQuery = rangeQuery.Lte(tsEnd)
	rangeQuery = rangeQuery.Format("epoch_millis")
	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(rangeQuery)

	termsAgg := elastic.NewTermsAggregation()
	termsAgg = termsAgg.Field("geoip.country_code2.raw")
	termsAgg = termsAgg.Size(staticCount)
	termsAgg = termsAgg.OrderByCountDesc()

	strTotalQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND retry:0 AND terminaltype:0")
	strTotalQuery = strTotalQuery.AnalyzeWildcard(true)
	totalQuery := elastic.NewBoolQuery()
	totalQuery = totalQuery.Must(strTotalQuery)
	totalQuery = totalQuery.Filter(boolQuery)
	totalResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(totalQuery).          // return all results, but ...
		Aggregation("2", termsAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	if err != nil {
		infoLog.Println("client.Search failed err=", err)
		return
	}
	totalCount := totalResult.TotalHits()
	mapTotalTry := make(map[string]int64, staticCount)
	totalBucketKeyIterm, ok := totalResult.Aggregations.Terms("2")
	if ok {
		for _, iterm := range totalBucketKeyIterm.Buckets {
			mapTotalTry[iterm.Key.(string)] = iterm.DocCount
		}
	} else {
		infoLog.Printf("searchResult.Aggregations.Terms trans err")
	}

	strAggQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND retry:0 AND terminaltype:0 AND NOT wns_code:0")
	strAggQuery = strAggQuery.AnalyzeWildcard(true)
	aggQuery := elastic.NewBoolQuery()
	aggQuery = aggQuery.Must(strAggQuery)
	aggQuery = aggQuery.Filter(boolQuery)

	searchResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(aggQuery).            // return all results, but ...
		Aggregation("2", termsAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	failTotalCount := searchResult.TotalHits()
	infoLog.Printf("================================================================================================")
	infoLog.Printf("%s iOS 首次重练总次数: %d iOS 首次重练失败总次数:%d 失败占比:%f%%", strTargetDay, totalCount, failTotalCount, float32(failTotalCount)/float32(totalCount)*100)
	infoLog.Printf("------------------------------------------------------------------------------------------------")
	mapItermTry := make(map[string]int64, staticCount)
	mapRateToCountry := make(map[float64]DistributeIterm, staticCount)
	sliceRate := make([]float64, staticCount)
	bucketKeyIterm, ok := searchResult.Aggregations.Terms("2")
	if ok {
		for _, iterm := range bucketKeyIterm.Buckets {
			mapItermTry[iterm.Key.(string)] = iterm.DocCount
		}
	} else {
		infoLog.Printf("searchResult.Aggregations.Terms trans err")
	}
	var index int = 0
	for k, v := range mapTotalTry {
		if count, ok := mapItermTry[k]; ok {
			rate := float64(count) / float64(v)
			mapRateToCountry[rate] = DistributeIterm{TotalTryCount: v, TotalFailedCount: count, CountryName: k}
			sliceRate[index] = rate
			index++
		} else {
			//infoLog.Printf("国家:%s 总的尝试次数:%5v 没找到总的失败次数", k, v)
		}
	}
	infoLog.Println()
	sort.Float64s(sliceRate)
	for _, k := range sliceRate {
		infoLog.Printf("国家:%s-%-10s\t 总的尝试次数:%7v 总的失败次数:%5v 占比:%8.4f%%",
			mapRateToCountry[k].CountryName,
			NickToName[mapRateToCountry[k].CountryName],
			mapRateToCountry[k].TotalTryCount,
			mapRateToCountry[k].TotalFailedCount,
			k*100)
	}
	infoLog.Printf("================================================================================================")
}

func GetAndroidFirstReconnectDistribute(client *elastic.Client, year, month, day int) {
	strTargetDay := fmt.Sprintf("%4d.%02d.%02d", year, month, day)
	tsBegin := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Unix()
	tsBegin = tsBegin * 1000 // change to ms
	tsEnd := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.Local).Unix()
	tsEnd = tsEnd * 1000 // change to ms
	infoLog.Printf("tsBegin=%v tsEnd=%v", tsBegin, tsEnd)

	indexBegin := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day-1)

	indexEnd := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day)
	infoLog.Printf("indexBegegin=%s indexEnd=%s", indexBegin, indexEnd)
	rangeQuery := elastic.NewRangeQuery("@timestamp")
	rangeQuery = rangeQuery.Gte(tsBegin)
	rangeQuery = rangeQuery.Lte(tsEnd)
	rangeQuery = rangeQuery.Format("epoch_millis")
	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(rangeQuery)

	termsAgg := elastic.NewTermsAggregation()
	termsAgg = termsAgg.Field("geoip.country_code2.raw")
	termsAgg = termsAgg.Size(staticCount)
	termsAgg = termsAgg.OrderByCountDesc()

	strTotalQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND retry:0 AND NOT terminaltype:0")
	strTotalQuery = strTotalQuery.AnalyzeWildcard(true)
	totalQuery := elastic.NewBoolQuery()
	totalQuery = totalQuery.Must(strTotalQuery)
	totalQuery = totalQuery.Filter(boolQuery)
	totalResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(totalQuery).          // return all results, but ...
		Aggregation("2", termsAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	if err != nil {
		infoLog.Println("client.Search failed err=", err)
		return
	}
	totalCount := totalResult.TotalHits()
	mapTotalTry := make(map[string]int64, staticCount)
	totalBucketKeyIterm, ok := totalResult.Aggregations.Terms("2")
	if ok {
		for _, iterm := range totalBucketKeyIterm.Buckets {
			mapTotalTry[iterm.Key.(string)] = iterm.DocCount
		}
	} else {
		infoLog.Printf("searchResult.Aggregations.Terms trans err")
	}

	strAggQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND retry:0 AND NOT terminaltype:0 AND NOT wns_code:0")
	strAggQuery = strAggQuery.AnalyzeWildcard(true)
	aggQuery := elastic.NewBoolQuery()
	aggQuery = aggQuery.Must(strAggQuery)
	aggQuery = aggQuery.Filter(boolQuery)

	searchResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(aggQuery).            // return all results, but ...
		Aggregation("2", termsAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	failTotalCount := searchResult.TotalHits()
	infoLog.Printf("================================================================================================")
	infoLog.Printf("%s Android 首次重练总次数: %d Android 首次重练失败总次数:%d 失败占比:%f%%", strTargetDay, totalCount, failTotalCount, float32(failTotalCount)/float32(totalCount)*100)
	infoLog.Printf("------------------------------------------------------------------------------------------------")
	mapItermTry := make(map[string]int64, staticCount)
	mapRateToCountry := make(map[float64]DistributeIterm, staticCount)
	sliceRate := make([]float64, staticCount)
	bucketKeyIterm, ok := searchResult.Aggregations.Terms("2")
	if ok {
		for _, iterm := range bucketKeyIterm.Buckets {
			mapItermTry[iterm.Key.(string)] = iterm.DocCount
		}
	} else {
		infoLog.Printf("searchResult.Aggregations.Terms trans err")
	}
	var index int = 0
	for k, v := range mapTotalTry {
		if count, ok := mapItermTry[k]; ok {
			rate := float64(count) / float64(v)
			mapRateToCountry[rate] = DistributeIterm{TotalTryCount: v, TotalFailedCount: count, CountryName: k}
			sliceRate[index] = rate
			index++
		} else {
			//infoLog.Printf("国家:%s 总的尝试次数:%5v 没找到总的失败次数", k, v)
		}
	}
	infoLog.Println()
	sort.Float64s(sliceRate)
	for _, k := range sliceRate {
		infoLog.Printf("国家:%s-%-10s\t 总的尝试次数:%7v 总的失败次数:%5v 占比:%8.4f%%",
			mapRateToCountry[k].CountryName,
			NickToName[mapRateToCountry[k].CountryName],
			mapRateToCountry[k].TotalTryCount,
			mapRateToCountry[k].TotalFailedCount,
			k*100)
	}

	infoLog.Printf("================================================================================================")
}

func GetReconnectData(client *elastic.Client, year, month, day int) {
	strTargetDay := fmt.Sprintf("%4d.%02d.%02d", year, month, day)
	tsBegin := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Unix()
	tsBegin = tsBegin * 1000 // change to ms
	tsEnd := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.Local).Unix()
	tsEnd = tsEnd * 1000 // change to ms
	infoLog.Printf("tsBegin=%v tsEnd=%v", tsBegin, tsEnd)

	indexBegin := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day-1)

	indexEnd := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day)
	infoLog.Printf("indexBegegin=%s indexEnd=%s", indexBegin, indexEnd)
	rangeQuery := elastic.NewRangeQuery("@timestamp")
	rangeQuery = rangeQuery.Gte(tsBegin)
	rangeQuery = rangeQuery.Lte(tsEnd)
	rangeQuery = rangeQuery.Format("epoch_millis")
	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(rangeQuery)

	strTotalQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\"")
	strTotalQuery = strTotalQuery.AnalyzeWildcard(true)
	totalQuery := elastic.NewBoolQuery()
	totalQuery = totalQuery.Must(strTotalQuery)
	totalQuery = totalQuery.Filter(boolQuery)
	totalResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(totalQuery). // return all results, but ...
		Pretty(true).      // pretty print request and response JSON
		Do()               // execute
	if err != nil {
		infoLog.Println("client.Search failed err=", err)
		return
	}
	totalCount := totalResult.TotalHits()

	strAggQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND NOT wns_code:0")
	strAggQuery = strAggQuery.AnalyzeWildcard(true)
	aggQuery := elastic.NewBoolQuery()
	aggQuery = aggQuery.Must(strAggQuery)
	aggQuery = aggQuery.Filter(boolQuery)

	termsAgg := elastic.NewTermsAggregation()
	termsAgg = termsAgg.Field("geoip.country_code3.raw")
	termsAgg = termsAgg.Size(30)
	termsAgg = termsAgg.OrderByCountDesc()

	searchResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(aggQuery).            // return all results, but ...
		Aggregation("2", termsAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	failTotalCount := searchResult.TotalHits()
	infoLog.Printf("================================================================================================")
	infoLog.Printf("%s 重练总次数: %d 重练失败总次数:%d 失败占比:%f%%", strTargetDay, totalCount, failTotalCount, float32(failTotalCount)/float32(totalCount)*100)
	infoLog.Printf("------------------------------------------------------------------------------------------------")
	bucketKeyIterm, ok := searchResult.Aggregations.Terms("2")
	if ok {
		for i, iterm := range bucketKeyIterm.Buckets {
			infoLog.Printf("index:%2v 国家:%s 失败次数:%5v 占比:%f%%", i, iterm.Key.(string), iterm.DocCount, float32(iterm.DocCount)/float32(totalCount)*100)
		}
	} else {
		infoLog.Printf("searchResult.Aggregations.Terms trans err")
	}
	infoLog.Printf("================================================================================================")
}

func GetIosReconnectData(client *elastic.Client, year, month, day int) {
	strTargetDay := fmt.Sprintf("%4d.%02d.%02d", year, month, day)
	tsBegin := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Unix()
	tsBegin = tsBegin * 1000 // change to ms
	tsEnd := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.Local).Unix()
	tsEnd = tsEnd * 1000 // change to ms
	infoLog.Printf("tsBegin=%v tsEnd=%v", tsBegin, tsEnd)

	indexBegin := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day-1)

	indexEnd := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day)
	infoLog.Printf("indexBegegin=%s indexEnd=%s", indexBegin, indexEnd)
	rangeQuery := elastic.NewRangeQuery("@timestamp")
	rangeQuery = rangeQuery.Gte(tsBegin)
	rangeQuery = rangeQuery.Lte(tsEnd)
	rangeQuery = rangeQuery.Format("epoch_millis")
	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(rangeQuery)

	strTotalQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND terminaltype:0")
	strTotalQuery = strTotalQuery.AnalyzeWildcard(true)
	totalQuery := elastic.NewBoolQuery()
	totalQuery = totalQuery.Must(strTotalQuery)
	totalQuery = totalQuery.Filter(boolQuery)
	totalResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(totalQuery). // return all results, but ...
		Pretty(true).      // pretty print request and response JSON
		Do()               // execute
	if err != nil {
		infoLog.Println("client.Search failed err=", err)
		return
	}
	totalCount := totalResult.TotalHits()

	strAggQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND terminaltype:0 AND NOT wns_code:0")
	strAggQuery = strAggQuery.AnalyzeWildcard(true)
	aggQuery := elastic.NewBoolQuery()
	aggQuery = aggQuery.Must(strAggQuery)
	aggQuery = aggQuery.Filter(boolQuery)

	termsAgg := elastic.NewTermsAggregation()
	termsAgg = termsAgg.Field("geoip.country_code3.raw")
	termsAgg = termsAgg.Size(30)
	termsAgg = termsAgg.OrderByCountDesc()

	searchResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(aggQuery).            // return all results, but ...
		Aggregation("2", termsAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	failTotalCount := searchResult.TotalHits()
	infoLog.Printf("=====================================================================================================")
	infoLog.Printf("%s iOS 重练总次数: %d iOS 重练失败总次数:%d 失败占比:%f%%", strTargetDay, totalCount, failTotalCount, float32(failTotalCount)/float32(totalCount)*100)
	infoLog.Printf("-----------------------------------------------------------------------------------------------------")
	bucketKeyIterm, ok := searchResult.Aggregations.Terms("2")
	if ok {
		for i, iterm := range bucketKeyIterm.Buckets {
			infoLog.Printf("index:%2v 国家:%s 失败次数:%5v 占比:%f%%", i, iterm.Key.(string), iterm.DocCount, float32(iterm.DocCount)/float32(totalCount)*100)
		}
	} else {
		infoLog.Printf("searchResult.Aggregations.Terms trans err")
	}
	infoLog.Printf("=====================================================================================================")
}

func GetAndroidReconnectData(client *elastic.Client, year, month, day int) {
	strTargetDay := fmt.Sprintf("%4d.%02d.%02d", year, month, day)
	tsBegin := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Unix()
	tsBegin = tsBegin * 1000 // change to ms
	tsEnd := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.Local).Unix()
	tsEnd = tsEnd * 1000 // change to ms
	infoLog.Printf("tsBegin=%v tsEnd=%v", tsBegin, tsEnd)

	indexBegin := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day-1)

	indexEnd := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day)
	infoLog.Printf("indexBegegin=%s indexEnd=%s", indexBegin, indexEnd)
	rangeQuery := elastic.NewRangeQuery("@timestamp")
	rangeQuery = rangeQuery.Gte(tsBegin)
	rangeQuery = rangeQuery.Lte(tsEnd)
	rangeQuery = rangeQuery.Format("epoch_millis")
	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(rangeQuery)

	strTotalQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND NOT terminaltype:0")
	strTotalQuery = strTotalQuery.AnalyzeWildcard(true)
	totalQuery := elastic.NewBoolQuery()
	totalQuery = totalQuery.Must(strTotalQuery)
	totalQuery = totalQuery.Filter(boolQuery)
	totalResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(totalQuery). // return all results, but ...
		Pretty(true).      // pretty print request and response JSON
		Do()               // execute
	if err != nil {
		infoLog.Println("client.Search failed err=", err)
		return
	}
	totalCount := totalResult.TotalHits()

	strAggQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND NOT terminaltype:0 AND NOT wns_code:0")
	strAggQuery = strAggQuery.AnalyzeWildcard(true)
	aggQuery := elastic.NewBoolQuery()
	aggQuery = aggQuery.Must(strAggQuery)
	aggQuery = aggQuery.Filter(boolQuery)

	termsAgg := elastic.NewTermsAggregation()
	termsAgg = termsAgg.Field("geoip.country_code3.raw")
	termsAgg = termsAgg.Size(30)
	termsAgg = termsAgg.OrderByCountDesc()

	searchResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(aggQuery).            // return all results, but ...
		Aggregation("2", termsAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	failTotalCount := searchResult.TotalHits()
	infoLog.Printf("============================================================================================================")
	infoLog.Printf("%s Android 重练总次数: %d Android 重练失败总次数:%d 失败占比:%f%%", strTargetDay, totalCount, failTotalCount, float32(failTotalCount)/float32(totalCount)*100)
	infoLog.Printf("------------------------------------------------==============================------------------------------")
	bucketKeyIterm, ok := searchResult.Aggregations.Terms("2")
	if ok {
		for i, iterm := range bucketKeyIterm.Buckets {
			infoLog.Printf("index:%2v 国家:%s 失败次数:%5v 占比:%f%%", i, iterm.Key.(string), iterm.DocCount, float32(iterm.DocCount)/float32(totalCount)*100)
		}
	} else {
		infoLog.Printf("searchResult.Aggregations.Terms trans err")
	}
	infoLog.Printf("============================================================================================================")
}

func GetReconnectUserData(client *elastic.Client, year, month, day int) {
	strTargetDay := fmt.Sprintf("%4d.%02d.%02d", year, month, day)
	tsBegin := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Unix()
	tsBegin = tsBegin * 1000 // change to ms
	tsEnd := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.Local).Unix()
	tsEnd = tsEnd * 1000 // change to ms
	infoLog.Printf("tsBegin=%v tsEnd=%v", tsBegin, tsEnd)

	indexBegin := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day-1)

	indexEnd := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day)
	infoLog.Printf("indexBegegin=%s indexEnd=%s", indexBegin, indexEnd)
	rangeQuery := elastic.NewRangeQuery("@timestamp")
	rangeQuery = rangeQuery.Gte(tsBegin)
	rangeQuery = rangeQuery.Lte(tsEnd)
	rangeQuery = rangeQuery.Format("epoch_millis")
	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(rangeQuery)

	totalAgg := elastic.NewCardinalityAggregation()
	totalAgg = totalAgg.Field("userid.raw")

	strTotalQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\"")
	strTotalQuery = strTotalQuery.AnalyzeWildcard(true)
	totalQuery := elastic.NewBoolQuery()
	totalQuery = totalQuery.Must(strTotalQuery)
	totalQuery = totalQuery.Filter(boolQuery)
	totalResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(totalQuery).          // return all results, but ...
		Aggregation("1", totalAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	if err != nil {
		infoLog.Println("client.Search failed err=", err)
		return
	}
	aggValueMetric, ok := totalResult.Aggregations.ValueCount("1")
	if !ok {
		infoLog.Println("totalResult.Aggregations.ValueCount(\"1\") failed")
		return
	}
	totalCount := *(aggValueMetric.Value)

	strAggQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND NOT wns_code:0")
	strAggQuery = strAggQuery.AnalyzeWildcard(true)
	aggQuery := elastic.NewBoolQuery()
	aggQuery = aggQuery.Must(strAggQuery)
	aggQuery = aggQuery.Filter(boolQuery)

	termsAgg := elastic.NewTermsAggregation()
	termsAgg = termsAgg.Field("geoip.country_code3.raw")
	termsAgg = termsAgg.Size(30)
	termsAgg = termsAgg.SubAggregation("1", totalAgg)
	termsAgg = termsAgg.OrderByAggregation("1", false)

	searchResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(aggQuery).            // return all results, but ...
		Aggregation("2", termsAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	infoLog.Printf("================================================================================================")
	bucketKeyIterm, ok := searchResult.Aggregations.Terms("2")
	if ok {
		failTotalCount := bucketKeyIterm.SumOfOtherDocCount
		infoLog.Printf("%s 重练总人数: %f 重练失败总人数:%d 失败占比:%f%%", strTargetDay, totalCount, failTotalCount, float32(failTotalCount)/float32(totalCount)*100)
		infoLog.Printf("------------------------------------------------------------------------------------------------")
		for i, iterm := range bucketKeyIterm.Buckets {
			valueMetric, ok := iterm.Aggregations.ValueCount("1")
			if ok {
				infoLog.Printf("index:%2v 国家:%s 失败次数:%5v 占比:%f%%", i, iterm.Key.(string), *(valueMetric.Value), float32(*(valueMetric.Value))/float32(totalCount)*100)
			} else {
				infoLog.Println("iterm.Aggregations.ValueCount()i failed index=", i)
			}
		}
	} else {
		infoLog.Printf("searchResult.Aggregations.Terms trans err")
	}
	infoLog.Printf("================================================================================================")
}

func GetIosReconnectUserData(client *elastic.Client, year, month, day int) {
	strTargetDay := fmt.Sprintf("%4d.%02d.%02d", year, month, day)
	tsBegin := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Unix()
	tsBegin = tsBegin * 1000 // change to ms
	tsEnd := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.Local).Unix()
	tsEnd = tsEnd * 1000 // change to ms
	infoLog.Printf("tsBegin=%v tsEnd=%v", tsBegin, tsEnd)

	indexBegin := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day-1)

	indexEnd := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day)
	infoLog.Printf("indexBegegin=%s indexEnd=%s", indexBegin, indexEnd)
	rangeQuery := elastic.NewRangeQuery("@timestamp")
	rangeQuery = rangeQuery.Gte(tsBegin)
	rangeQuery = rangeQuery.Lte(tsEnd)
	rangeQuery = rangeQuery.Format("epoch_millis")
	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(rangeQuery)

	totalAgg := elastic.NewCardinalityAggregation()
	totalAgg = totalAgg.Field("userid.raw")

	strTotalQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND terminaltype:0")
	strTotalQuery = strTotalQuery.AnalyzeWildcard(true)
	totalQuery := elastic.NewBoolQuery()
	totalQuery = totalQuery.Must(strTotalQuery)
	totalQuery = totalQuery.Filter(boolQuery)
	totalResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(totalQuery).          // return all results, but ...
		Aggregation("1", totalAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	if err != nil {
		infoLog.Println("client.Search failed err=", err)
		return
	}
	aggValueMetric, ok := totalResult.Aggregations.ValueCount("1")
	if !ok {
		infoLog.Println("totalResult.Aggregations.ValueCount(\"1\") failed")
		return
	}
	totalCount := *(aggValueMetric.Value)

	strAggQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND terminaltype:0 AND NOT wns_code:0")
	strAggQuery = strAggQuery.AnalyzeWildcard(true)
	aggQuery := elastic.NewBoolQuery()
	aggQuery = aggQuery.Must(strAggQuery)
	aggQuery = aggQuery.Filter(boolQuery)

	termsAgg := elastic.NewTermsAggregation()
	termsAgg = termsAgg.Field("geoip.country_code3.raw")
	termsAgg = termsAgg.Size(30)
	termsAgg = termsAgg.SubAggregation("1", totalAgg)
	termsAgg = termsAgg.OrderByAggregation("1", false)

	searchResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(aggQuery).            // return all results, but ...
		Aggregation("2", termsAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	infoLog.Printf("================================================================================================")
	bucketKeyIterm, ok := searchResult.Aggregations.Terms("2")
	if ok {
		failTotalCount := bucketKeyIterm.SumOfOtherDocCount
		infoLog.Printf("%s iOS 重练总人数: %f iOS重练失败总人次数:%v 失败占比:%f%%", strTargetDay, totalCount, failTotalCount, float32(failTotalCount)/float32(totalCount)*100)
		infoLog.Printf("------------------------------------------------------------------------------------------------")
		for i, iterm := range bucketKeyIterm.Buckets {
			valueMetric, ok := iterm.Aggregations.ValueCount("1")
			if ok {
				infoLog.Printf("index:%2v 国家:%s 失败次数:%5v 占比:%f%%", i, iterm.Key.(string), *(valueMetric.Value), float32(*(valueMetric.Value))/float32(totalCount)*100)
			} else {
				infoLog.Println("iterm.Aggregations.ValueCount()i failed index=", i)
			}
		}
	} else {
		infoLog.Printf("searchResult.Aggregations.Terms trans err")
	}
	infoLog.Printf("================================================================================================")
}

func GetAndroidReconnectUserData(client *elastic.Client, year, month, day int) {
	strTargetDay := fmt.Sprintf("%4d.%02d.%02d", year, month, day)
	tsBegin := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Unix()
	tsBegin = tsBegin * 1000 // change to ms
	tsEnd := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.Local).Unix()
	tsEnd = tsEnd * 1000 // change to ms
	infoLog.Printf("tsBegin=%v tsEnd=%v", tsBegin, tsEnd)

	indexBegin := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day-1)

	indexEnd := fmt.Sprintf("logstash-client-log-%4d.%02d.%02d", year, month, day)
	infoLog.Printf("indexBegegin=%s indexEnd=%s", indexBegin, indexEnd)
	rangeQuery := elastic.NewRangeQuery("@timestamp")
	rangeQuery = rangeQuery.Gte(tsBegin)
	rangeQuery = rangeQuery.Lte(tsEnd)
	rangeQuery = rangeQuery.Format("epoch_millis")
	boolQuery := elastic.NewBoolQuery()
	boolQuery = boolQuery.Must(rangeQuery)

	totalAgg := elastic.NewCardinalityAggregation()
	totalAgg = totalAgg.Field("userid.raw")

	strTotalQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND NOT terminaltype:0")
	strTotalQuery = strTotalQuery.AnalyzeWildcard(true)
	totalQuery := elastic.NewBoolQuery()
	totalQuery = totalQuery.Must(strTotalQuery)
	totalQuery = totalQuery.Filter(boolQuery)
	totalResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(totalQuery).          // return all results, but ...
		Aggregation("1", totalAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	if err != nil {
		infoLog.Println("client.Search failed err=", err)
		return
	}
	aggValueMetric, ok := totalResult.Aggregations.ValueCount("1")
	if !ok {
		infoLog.Println("totalResult.Aggregations.ValueCount(\"1\") failed")
		return
	}
	totalCount := *(aggValueMetric.Value)

	strAggQuery := elastic.NewQueryStringQuery("cmd:\"0x1005\" AND NOT terminaltype:0 AND NOT wns_code:0")
	strAggQuery = strAggQuery.AnalyzeWildcard(true)
	aggQuery := elastic.NewBoolQuery()
	aggQuery = aggQuery.Must(strAggQuery)
	aggQuery = aggQuery.Filter(boolQuery)

	termsAgg := elastic.NewTermsAggregation()
	termsAgg = termsAgg.Field("geoip.country_code3.raw")
	termsAgg = termsAgg.Size(30)
	termsAgg = termsAgg.SubAggregation("1", totalAgg)
	termsAgg = termsAgg.OrderByAggregation("1", false)

	searchResult, err := client.Search().
		Index(indexBegin, indexEnd). // search in indexBegin and indexEnd
		SearchType("count").
		IgnoreUnavailable(true).
		Query(aggQuery).            // return all results, but ...
		Aggregation("2", termsAgg). // add our aggregation to the query
		Pretty(true).               // pretty print request and response JSON
		Do()                        // execute
	infoLog.Printf("================================================================================================")
	bucketKeyIterm, ok := searchResult.Aggregations.Terms("2")
	if ok {
		failTotalCount := bucketKeyIterm.SumOfOtherDocCount
		infoLog.Printf("%s Android 重练总人数: %f Android 重练失败总人次数:%v 失败占比:%f%%", strTargetDay, totalCount, failTotalCount, float32(failTotalCount)/float32(totalCount)*100)
		infoLog.Printf("------------------------------------------------------------------------------------------------")
		for i, iterm := range bucketKeyIterm.Buckets {
			valueMetric, ok := iterm.Aggregations.ValueCount("1")
			if ok {
				infoLog.Printf("index:%2v 国家:%s 失败次数:%5v 占比:%f%%", i, iterm.Key.(string), *(valueMetric.Value), float32(*(valueMetric.Value))/float32(totalCount)*100)
			} else {
				infoLog.Println("iterm.Aggregations.ValueCount()i failed index=", i)
			}
		}
	} else {
		infoLog.Printf("searchResult.Aggregations.Terms trans err")
	}
	infoLog.Printf("================================================================================================")
}

func main() {
	// 处理命令行参数
	if _, err := parser.Parse(); err != nil {
		log.Fatalln("parse cmd line failed!")
	}

	if options.ClientConf == "" {
		log.Fatalln("Must input config file name")
	}

	// log.Println("config name =", options.ClientConf)
	// 读取配置文件
	cfg, err := ini.Load([]byte(""), options.ClientConf)
	if err != nil {
		log.Printf("load config file=%s failed", options.ClientConf)
		return
	}
	// 配置文件只读 设置此标识提升性能
	cfg.BlockMode = false
	// 定义一个文件
	fileName := cfg.Section("LOG").Key("path").MustString("/home/ht/clinet.log")
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
		return
	}

	// 创建一个日志对象
	infoLog = log.New(logFile, "[Info]", log.LstdFlags)
	// 配置log的Flag参数
	infoLog.SetFlags(infoLog.Flags() | log.LstdFlags)

	// 读取ip+port
	serverIp := cfg.Section("OUTER_SERVER").Key("server_ip").MustString("127.0.0.3")
	serverPort := cfg.Section("OUTER_SERVER").Key("server_port").MustString("9200")

	infoLog.Printf("server_ip=%v server_port=%v\n", serverIp, serverPort)

	// Obtain a client. You can also provide your own HTTP client here.
	client, err := elastic.NewClient(elastic.SetURL("http://"+serverIp+":"+serverPort),
		elastic.SetMaxRetries(10),
		elastic.SetBasicAuth("hellotalk_admin", "hellotalk_admin123456"),
		//elastic.SetTraceLog(infoLog),
		//elastic.SetInfoLog(infoLog),
		elastic.SetErrorLog(infoLog))
	if err != nil {
		// Handle error
		panic(err)
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://127.0.0.1:9200").Do()
	if err != nil {
		// Handle error
		panic(err)
	}
	infoLog.Printf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)

	// Getting the ES version number is quite common, so there's a shortcut
	esversion, err := client.ElasticsearchVersion("http://127.0.0.1:9200")
	if err != nil {
		// Handle error
		panic(err)
	}
	infoLog.Printf("Elasticsearch version %s", esversion)

	// Create an aggregation for users and a sub-aggregation for a date histogram of tweets (per year).
	//timeline := elastic.NewTermsAggregation().Field("userid").Size(10).OrderByCountDesc()

	// Search with a term query
	configDate := cfg.Section("CONFIG_DATE").Key("date").MustInt(1)
	tsNow := time.Now()
	year, month, day := tsNow.Date()
	infoLog.Printf("year=%4d month=%2d day=%02d", year, month, day)

	var targetYear, targetMonth, targetDay int
	if day-configDate <= 0 {
		targetMonth = int(month) - 1
		if targetMonth <= 0 {
			targetYear = year - 1
			targetMonth = 12
		}
		dayOfLastMonth := GetDayOfMonth(targetYear, targetMonth)
		targetDay = dayOfLastMonth + day - configDate
	} else {
		targetYear = year
		targetMonth = int(month)
		targetDay = day - configDate
	}
	// 总的分布
	//	GetReconnectData(client, targetYear, targetMonth, targetDay)

	//	GetIosReconnectData(client, targetYear, targetMonth, targetDay)

	//	GetAndroidReconnectData(client, targetYear, targetMonth, targetDay)

	// 首次登录分布
	//GetFirstReconnectData(client, targetYear, targetMonth, targetDay)
	//GetIosFirstReconnectData(client, targetYear, targetMonth, targetDay)
	//GetAndroidFirstReconnectData(client, targetYear, targetMonth, targetDay)

	// 首次登录失败安装分开统计情况
	GetFirstReconnectDistribute(client, targetYear, targetMonth, targetDay)
	GetIosFirstReconnectDistribute(client, targetYear, targetMonth, targetDay)
	GetAndroidFirstReconnectDistribute(client, targetYear, targetMonth, targetDay)

	// 失败人数分布情况
	//	GetReconnectUserData(client, targetYear, targetMonth, targetDay)

	//	GetIosReconnectUserData(client, targetYear, targetMonth, targetDay)

	//	GetAndroidReconnectUserData(client, targetYear, targetMonth, targetDay)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
