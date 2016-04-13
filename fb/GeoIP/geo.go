/**
 * libgeo.go
 *
 * Copyright (c) 2010, Nikola Ranchev
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 * 	- Redistributions of source code must retain the above copyright
 * 	  notice, this list of conditions and the following disclaimer.
 * 	- Redistributions in binary form must reproduce the above copyright
 * 	  notice, this list of conditions and the following disclaimer in the
 * 	  documentation and/or other materials provided with the distribution.
 * 	- Neither the name of the <organization> nor the
 * 	  names of its contributors may be used to endorse or promote products
 * 	  derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package geoIP

// Dependencies
import (
	"errors"
	"os"
)

// Globals (const arrays that will be initialized inside init())
var (
	countryCode = []string{
		"--", "AP", "EU", "AD", "AE", "AF", "AG", "AI", "AL", "AM", "AN", "AO", "AQ", "AR",
		"AS", "AT", "AU", "AW", "AZ", "BA", "BB", "BD", "BE", "BF", "BG", "BH", "BI", "BJ",
		"BM", "BN", "BO", "BR", "BS", "BT", "BV", "BW", "BY", "BZ", "CA", "CC", "CD", "CF",
		"CG", "CH", "CI", "CK", "CL", "CM", "CN", "CO", "CR", "CU", "CV", "CX", "CY", "CZ",
		"DE", "DJ", "DK", "DM", "DO", "DZ", "EC", "EE", "EG", "EH", "ER", "ES", "ET", "FI",
		"FJ", "FK", "FM", "FO", "FR", "FX", "GA", "GB", "GD", "GE", "GF", "GH", "GI", "GL",
		"GM", "GN", "GP", "GQ", "GR", "GS", "GT", "GU", "GW", "GY", "HK", "HM", "HN", "HR",
		"HT", "HU", "ID", "IE", "IL", "IN", "IO", "IQ", "IR", "IS", "IT", "JM", "JO", "JP",
		"KE", "KG", "KH", "KI", "KM", "KN", "KP", "KR", "KW", "KY", "KZ", "LA", "LB", "LC",
		"LI", "LK", "LR", "LS", "LT", "LU", "LV", "LY", "MA", "MC", "MD", "MG", "MH", "MK",
		"ML", "MM", "MN", "MO", "MP", "MQ", "MR", "MS", "MT", "MU", "MV", "MW", "MX", "MY",
		"MZ", "NA", "NC", "NE", "NF", "NG", "NI", "NL", "NO", "NP", "NR", "NU", "NZ", "OM",
		"PA", "PE", "PF", "PG", "PH", "PK", "PL", "PM", "PN", "PR", "PS", "PT", "PW", "PY",
		"QA", "RE", "RO", "RU", "RW", "SA", "SB", "SC", "SD", "SE", "SG", "SH", "SI", "SJ",
		"SK", "SL", "SM", "SN", "SO", "SR", "ST", "SV", "SY", "SZ", "TC", "TD", "TF", "TG",
		"TH", "TJ", "TK", "TM", "TN", "TO", "TL", "TR", "TT", "TV", "TW", "TZ", "UA", "UG",
		"UM", "US", "UY", "UZ", "VA", "VC", "VE", "VG", "VI", "VN", "VU", "WF", "WS", "YE",
		"YT", "RS", "ZA", "ZM", "ME", "ZW", "A1", "A2", "O1", "AX", "GG", "IM", "JE", "BL",
		"MF", "BQ", "SS", "O1"}
	countryName = []string{
		"N/A", "Asia/Pacific Region", "Europe", "Andorra", "United Arab Emirates",
		"Afghanistan", "Antigua and Barbuda", "Anguilla", "Albania", "Armenia",
		"Netherlands Antilles", "Angola", "Antarctica", "Argentina", "American Samoa",
		"Austria", "Australia", "Aruba", "Azerbaijan", "Bosnia and Herzegovina",
		"Barbados", "Bangladesh", "Belgium", "Burkina Faso", "Bulgaria", "Bahrain",
		"Burundi", "Benin", "Bermuda", "Brunei Darussalam", "Bolivia", "Brazil", "Bahamas",
		"Bhutan", "Bouvet Island", "Botswana", "Belarus", "Belize", "Canada",
		"Cocos (Keeling) Islands", "Congo, The Democratic Republic of the",
		"Central African Republic", "Congo", "Switzerland", "Cote D'Ivoire",
		"Cook Islands", "Chile", "Cameroon", "China", "Colombia", "Costa Rica", "Cuba",
		"Cape Verde", "Christmas Island", "Cyprus", "Czech Republic", "Germany",
		"Djibouti", "Denmark", "Dominica", "Dominican Republic", "Algeria", "Ecuador",
		"Estonia", "Egypt", "Western Sahara", "Eritrea", "Spain", "Ethiopia", "Finland",
		"Fiji", "Falkland Islands (Malvinas)", "Micronesia, Federated States of",
		"Faroe Islands", "France", "France, Metropolitan", "Gabon", "United Kingdom",
		"Grenada", "Georgia", "French Guiana", "Ghana", "Gibraltar", "Greenland", "Gambia",
		"Guinea", "Guadeloupe", "Equatorial Guinea", "Greece",
		"South Georgia and the South Sandwich Islands", "Guatemala", "Guam",
		"Guinea-Bissau", "Guyana", "Hong Kong", "Heard Island and McDonald Islands",
		"Honduras", "Croatia", "Haiti", "Hungary", "Indonesia", "Ireland", "Israel", "India",
		"British Indian Ocean Territory", "Iraq", "Iran, Islamic Republic of",
		"Iceland", "Italy", "Jamaica", "Jordan", "Japan", "Kenya", "Kyrgyzstan", "Cambodia",
		"Kiribati", "Comoros", "Saint Kitts and Nevis",
		"Korea, Democratic People's Republic of", "Korea, Republic of", "Kuwait",
		"Cayman Islands", "Kazakhstan", "Lao People's Democratic Republic", "Lebanon",
		"Saint Lucia", "Liechtenstein", "Sri Lanka", "Liberia", "Lesotho", "Lithuania",
		"Luxembourg", "Latvia", "Libyan Arab Jamahiriya", "Morocco", "Monaco",
		"Moldova, Republic of", "Madagascar", "Marshall Islands",
		"Macedonia", "Mali", "Myanmar", "Mongolia",
		"Macau", "Northern Mariana Islands", "Martinique", "Mauritania", "Montserrat",
		"Malta", "Mauritius", "Maldives", "Malawi", "Mexico", "Malaysia", "Mozambique",
		"Namibia", "New Caledonia", "Niger", "Norfolk Island", "Nigeria", "Nicaragua",
		"Netherlands", "Norway", "Nepal", "Nauru", "Niue", "New Zealand", "Oman", "Panama",
		"Peru", "French Polynesia", "Papua New Guinea", "Philippines", "Pakistan",
		"Poland", "Saint Pierre and Miquelon", "Pitcairn Islands", "Puerto Rico",
		"Palestinian Territory", "Portugal", "Palau", "Paraguay", "Qatar",
		"Reunion", "Romania", "Russian Federation", "Rwanda", "Saudi Arabia",
		"Solomon Islands", "Seychelles", "Sudan", "Sweden", "Singapore", "Saint Helena",
		"Slovenia", "Svalbard and Jan Mayen", "Slovakia", "Sierra Leone", "San Marino",
		"Senegal", "Somalia", "Suriname", "Sao Tome and Principe", "El Salvador",
		"Syrian Arab Republic", "Swaziland", "Turks and Caicos Islands", "Chad",
		"French Southern Territories", "Togo", "Thailand", "Tajikistan", "Tokelau",
		"Turkmenistan", "Tunisia", "Tonga", "Timor-Leste", "Turkey", "Trinidad and Tobago",
		"Tuvalu", "Taiwan", "Tanzania, United Republic of", "Ukraine", "Uganda",
		"United States Minor Outlying Islands", "United States", "Uruguay", "Uzbekistan",
		"Holy See (Vatican City State)", "Saint Vincent and the Grenadines",
		"Venezuela", "Virgin Islands, British", "Virgin Islands, U.S.", "Vietnam",
		"Vanuatu", "Wallis and Futuna", "Samoa", "Yemen", "Mayotte", "Serbia",
		"South Africa", "Zambia", "Montenegro", "Zimbabwe", "Anonymous Proxy",
		"Satellite Provider", "Other", "Aland Islands", "Guernsey", "Isle of Man", "Jersey",
		"Saint Barthelemy", "Saint Martin", "Bonaire, Saint Eustatius and Saba",
		"South Sudan", "Other"}
	countryCnName = []string{
		"N/A", "亚洲/太平洋地区", "欧洲", "安道尔", "阿拉伯联合酋长国",
		"阿富汗", "安提瓜和巴布达", "安圭拉岛", "阿尔巴尼亚", "亚美尼亚",
		"荷属安的列斯群岛", "安哥拉", "南极洲", "阿根廷", "美属萨摩亚",
		"奥地利", "澳大利亚", "阿鲁巴岛", "阿塞拜疆", "波斯尼亚和黑塞哥维那",
		"巴巴多斯", "孟加拉国", "比利时", "布吉纳法索", "保加利亚", "巴林",
		"布隆迪", "贝宁", "百慕大", "文莱达鲁萨兰国", "玻利维亚", "巴西", "巴哈马群岛",
		"不丹", "布维岛", "博茨瓦纳", "白俄罗斯", "伯利兹城", "加拿大",
		"可可(Keeling)岛", "刚果民主共和国的",
		"中非共和国", "刚果", "瑞士", "象牙海岸",
		"库克群岛", "智利", "喀麦隆", "中国", "哥伦比亚", "哥斯达黎加", "古巴",
		"佛得角", "圣诞岛", "塞浦路斯", "捷克", "德国",
		"吉布提", "丹麦", "多米尼加", "多米尼加共和国", "阿尔及利亚", "厄瓜多尔",
		"爱沙尼亚", "埃及", "西撒哈拉", "厄立特里亚", "西班牙", "埃塞俄比亚", "芬兰",
		"斐济", "福克兰群岛(马尔维纳斯)", "密克罗尼西亚联邦",
		"法罗群岛", "法国", "法国,大都会", "加蓬", "联合王国",
		"格林纳达", "格鲁吉亚", "法属圭亚那", "加纳", "直布罗陀", "格陵兰岛", "冈比亚",
		"几内亚", "瓜德罗普岛", "赤道几内亚", "希腊",
		"南乔治亚岛和南桑威奇群岛", "危地马拉", "关岛",
		"几内亚比绍", "圭亚那", "香港", "听说岛和麦当劳群岛",
		"洪都拉斯", "克罗地亚", "海地", "匈牙利", "印尼", "爱尔兰", "以色列", "印度",
		"英属印度洋领地", "伊拉克", "伊朗伊斯兰共和国",
		"冰岛", "意大利", "牙买加", "乔丹", "日本", "肯尼亚", "吉尔吉斯斯坦", "柬埔寨",
		"基里巴斯", "科摩罗", "圣基茨和尼维斯",
		"朝鲜民主主义人民共和国", "朝鲜共和国", "科威特",
		"开曼群岛", "哈萨克斯坦", "老挝人民民主共和国", "黎巴嫩",
		"圣·露西亚", "列支敦士登", "斯里兰卡", "赖比瑞亚", "莱索托", "立陶宛",
		"卢森堡", "拉脱维亚", "阿拉伯利比亚民众国", "摩洛哥", "摩纳哥",
		"摩尔多瓦共和国", "马达加斯加", "马绍尔群岛",
		"马其顿", "马里", "缅甸", "蒙古",
		"澳门", "北马里亚纳群岛", "马提尼克岛", "毛里塔尼亚", "蒙特塞拉特岛",
		"马耳他", "毛里求斯", "马尔代夫", "马拉维", "墨西哥", "马来西亚", "莫桑比克",
		"纳米比亚", "新喀里多尼亚", "尼日尔", "诺福克岛", "尼日利亚", "尼加拉瓜",
		"荷兰", "挪威", "尼泊尔", "瑙鲁", "纽埃岛", "新西兰", "阿曼", "巴拿马",
		"秘鲁", "法属波利尼西亚", "巴布亚新几内亚", "菲律宾", "巴基斯坦",
		"波兰", "圣皮埃尔和密克隆群岛", "皮特凯恩群岛", "波多黎各",
		"巴勒斯坦领土", "葡萄牙", "帕劳", "巴拉圭", "卡塔尔",
		"团圆", "罗马尼亚", "俄罗斯联邦", "卢旺达", "沙特阿拉伯",
		"所罗门群岛", "塞舌尔群岛", "苏丹", "瑞典", "新加坡", "圣赫勒拿",
		"斯洛文尼亚", "斯瓦尔巴群岛和扬马延岛", "斯洛伐克", "狮子山", "圣马力诺",
		"塞内加尔", "索马里", "苏里南", "圣多美和普林西比", "萨尔瓦多",
		"阿拉伯叙利亚共和国", "斯威士兰", "特克斯和凯科斯群岛", "乍得",
		"法国南部地区", "多哥", "泰国", "塔", "托克劳",
		"土", "突尼斯", "汤加", "东帝汶", "土耳其", "特立尼达和多巴哥",
		"图瓦卢", "台湾", "坦桑尼亚联合共和国", "乌克兰", "乌干达",
		"美国小离岛", "美国", "乌拉圭", "乌兹别克斯坦",
		"教廷(梵蒂冈)", "圣文森特和格林纳丁斯",
		"委内瑞拉", "维尔京群岛", "英国", "维尔京群岛", "美国", "越南",
		"瓦努阿图", "瓦利斯群岛和富图纳群岛", "萨摩亚", "也门", "马约特岛", "塞尔维亚",
		"南非", "赞比亚", "黑山", "津巴布韦", "匿名代理",
		"卫星提供者", "其他", "阿兰群岛", "格恩西岛", "马恩岛", "泽西岛",
		"圣巴特尔米", "圣马丁", "博内尔岛, 圣Eustatius和萨巴",
		"南苏丹", "其他"}
)

// Constants
const (
	maxRecordLength = 4
	standardRecordLength = 3
	countryBegin = 16776960
	structureInfoMaxSize = 20
	fullRecordLength = 60
	segmentRecordLength = 3

	// DB Types
	dbCountryEdition = byte(1)
	dbCityEditionRev0 = byte(6)
	dbCityEditionRev1 = byte(2)
)

// These are some structs
type GeoIP struct {
	databaseSegment int    // No need to make an array of size 1
	recordLength    int    // Set to one of the constants above
	dbType          byte   // Store the database type
	data            []byte // All of the data from the DB file
}
type Location struct {
	CountryCode string // If country ed. only country info is filled
	CountryName string // If country ed. only country info is filled
	Region      string
	City        string
	PostalCode  string
	Latitude    float32
	Longitude   float32
}

// Load the database file in memory, detect the db format and setup the GeoIP struct
func Load(filename string) (gi *GeoIP, err error) {
	// Try to open the requested file
	dbInfo, err := os.Lstat(filename)
	if err != nil {
		return
	}
	dbFile, err := os.Open(filename)
	if err != nil {
		return
	}

	// Copy the db into memory
	gi = new(GeoIP)
	gi.data = make([]byte, dbInfo.Size())
	dbFile.Read(gi.data)
	dbFile.Close()

	// Check the database type
	gi.dbType = dbCountryEdition           // Default the database to country edition
	gi.databaseSegment = countryBegin      // Default to country DB
	gi.recordLength = standardRecordLength // Default to country DB

	// Search for the DB type headers
	delim := make([]byte, 3)
	for i := 0; i < structureInfoMaxSize; i++ {
		delim = gi.data[len(gi.data) - i - 3 - 1 : len(gi.data) - i - 1]
		if int8(delim[0]) == -1 && int8(delim[1]) == -1 && int8(delim[2]) == -1 {
			gi.dbType = gi.data[len(gi.data) - i - 1]
			// If we detect city edition set the correct segment offset
			if gi.dbType == dbCityEditionRev0 || gi.dbType == dbCityEditionRev1 {
				buf := make([]byte, segmentRecordLength)
				buf = gi.data[len(gi.data) - i - 1 + 1 : len(gi.data) - i - 1 + 4]
				gi.databaseSegment = 0
				for j := 0; j < segmentRecordLength; j++ {
					gi.databaseSegment += (int(buf[j]) << uint8(j * 8))
				}
			}
			break
		}
	}

	// Support older formats
	if gi.dbType >= 106 {
		gi.dbType -= 105
	}

	// Reject unsupported formats
	if gi.dbType != dbCountryEdition && gi.dbType != dbCityEditionRev0 && gi.dbType != dbCityEditionRev1 {
		err = errors.New("Unsupported database format")
		return
	}

	return
}

// Lookup by IP address and return location
func (gi *GeoIP) GetLocationByIP(ip string) *Location {
	return gi.GetLocationByIPNum(AddrToNum(ip))
}

// Lookup by IP number and return location
func (gi *GeoIP) GetLocationByIPNum(ipNum uint32) *Location {
	// Perform the lookup on the database to see if the record is found
	offset := gi.lookupByIPNum(ipNum)

	// Check if the country was found
	if gi.dbType == dbCountryEdition && offset - countryBegin == 0 ||
	gi.dbType != dbCountryEdition && gi.databaseSegment == offset {
		return nil
	}

	// Create a generic location structure
	location := new(Location)

	// If the database is country
	if gi.dbType == dbCountryEdition {
		location.CountryCode = countryCode[offset - countryBegin]
		location.CountryName = countryName[offset - countryBegin]

		return location
	}

	// Find the max record length
	recPointer := offset + (2 * gi.recordLength - 1) * gi.databaseSegment
	recordEnd := recPointer + fullRecordLength
	if len(gi.data) - recPointer < fullRecordLength {
		recordEnd = len(gi.data)
	}

	// Read the country code/name first
	location.CountryCode = countryCode[gi.data[recPointer]]
	location.CountryName = countryName[gi.data[recPointer]]
	readLen := 1
	recPointer += 1

	// Get the region
	for readLen = 0; gi.data[recPointer + readLen] != '\000' &&
	recPointer + readLen < recordEnd; readLen++ {
	}
	if readLen != 0 {
		location.Region = string(gi.data[recPointer : recPointer + readLen])
	}
	recPointer += readLen + 1

	// Get the city
	for readLen = 0; gi.data[recPointer + readLen] != '\000' &&
	recPointer + readLen < recordEnd; readLen++ {
	}
	if readLen != 0 {
		location.City = string(gi.data[recPointer : recPointer + readLen])
	}
	recPointer += readLen + 1

	// Get the postal code
	for readLen = 0; gi.data[recPointer + readLen] != '\000' &&
	recPointer + readLen < recordEnd; readLen++ {
	}
	if readLen != 0 {
		location.PostalCode = string(gi.data[recPointer : recPointer + readLen])
	}
	recPointer += readLen + 1

	// Get the latitude
	coordinate := float32(0)
	for j := 0; j < 3; j++ {
		coordinate += float32(int32(gi.data[recPointer + j]) << uint8(j * 8))
	}
	location.Latitude = float32(coordinate / 10000 - 180)
	recPointer += 3

	// Get the longitude
	coordinate = 0
	for j := 0; j < 3; j++ {
		coordinate += float32(int32(gi.data[recPointer + j]) << uint8(j * 8))
	}
	location.Longitude = float32(coordinate / 10000 - 180)

	return location
}

// Read the database and return record position
func (gi *GeoIP) lookupByIPNum(ip uint32) int {
	buf := make([]byte, 2 * maxRecordLength)
	x := make([]int, 2)
	offset := 0
	for depth := 31; depth >= 0; depth-- {
		for i := 0; i < 2 * maxRecordLength; i++ {
			buf[i] = gi.data[(2 * gi.recordLength * offset) + i]
		}
		for i := 0; i < 2; i++ {
			x[i] = 0
			for j := 0; j < gi.recordLength; j++ {
				var y int = int(buf[i * gi.recordLength + j])
				if y < 0 {
					y += 256
				}
				x[i] += (y << uint(j * 8))
			}
		}
		if (ip & (1 << uint(depth))) > 0 {
			if x[1] >= gi.databaseSegment {
				return x[1]
			}
			offset = x[1]
		} else {
			if x[0] >= gi.databaseSegment {
				return x[0]
			}
			offset = x[0]
		}
	}
	return 0
}

// Convert ip address to an int representation
func AddrToNum(ip string) uint32 {
	octet := uint32(0)
	ipnum := uint32(0)
	i := 3
	for j := 0; j < len(ip); j++ {
		c := byte(ip[j])
		if c == '.' {
			if octet > 255 {
				return 0
			}
			ipnum <<= 8
			ipnum += octet
			i--
			octet = 0
		} else {
			t := octet
			octet <<= 3
			octet += t
			octet += t
			c -= '0'
			if c > 9 {
				return 0
			}
			octet += uint32(c)
		}
	}
	if (octet > 255) || (i != 0) {
		return 0
	}
	ipnum <<= 8
	return uint32(ipnum + octet)
}

func CountryCnNameFromIndex(index int) string{

    return countryCnName[index]

}

func CountryCnNameFromCode(code string) string{
     for i,v:=range  countryCode{
	     if v==code{
		     return countryCnName[i]
	     }
     }
	return "N/A"
}

func CountryNameFromIndex(index int) string{
	return countryName[index]
}

func CountryNameFromCode(code string) string{
	for i,v:=range  countryCode{
		if v==code{
			return countryName[i]
		}
	}
	return "N/A"
}

func CountryCodes()[]string{
	return countryCode
}