package helper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/net/html"
)

// const version = "1.0.13"

// region encrypt

// 加密
// 密钥长度应为 128，192，256
func EncryptAES(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	src = PaddingText(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	blockMode.CryptBlocks(src, src)
	return src, nil
}

// AES-ECB
//
// Deprecated: ECB 是不安全的, 请不要使用
func EncryptAESECB(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	size := block.BlockSize()
	src = PaddingText(src, size)

	// 官方不支持ECB，所有网上抄了一段
	l := len(src)
	res := make([]byte, l)
	for s, e := 0, size; s < l; s, e = s+size, e+size {
		block.Encrypt(res[s:e], src[s:e])
	}

	return res, nil
}

// See: https://www.saoniuhuo.com/question/detail-2866706.html
//
// Deprecated: ECB 是不安全的, 请不要使用
func DecryptAESECB(src, key []byte) ([]byte, error) {
	cipher, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(src))
	size := cipher.BlockSize()

	for bs, be := 0, size; bs < len(src); bs, be = bs+size, be+size {
		cipher.Decrypt(decrypted[bs:be], src[bs:be])
	}
	decrypted = UnPaddingText(decrypted)
	return decrypted, nil
}

// 解密
func DecryptAES(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(src)%block.BlockSize() != 0 {
		return nil, errors.New("helper/encrypt: input not full blocks")
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])

	blockMode.CryptBlocks(src, src)
	src = UnPaddingText(src)
	return src, nil
}

// 从 url 中获取指定键的值
func GetQuerySpecifiedFromUrl(query, k string) (string, bool) {
	// 先试图查找问号，如果找到，其后视为查询部。否则，全部视为查询部
	if i := strings.Index(query, "?"); i >= 0 {
		query = query[i+1:]
	}

	// 以下逻辑参考了 url.parseQuery
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			continue
		}
		if key == k {
			return value, true
		}
	}
	return "", false
}

// 移除末尾空格
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// 填充字符串（末尾）
func PaddingText(str []byte, blockSize int) []byte {
	//需要填充的数据长度
	paddingCount := blockSize - len(str)%blockSize
	//填充数据为：paddingCount ,填充的值为：paddingCount
	paddingStr := bytes.Repeat([]byte{byte(paddingCount)}, paddingCount)
	newPaddingStr := append(str, paddingStr...)
	return newPaddingStr
}

// 去掉字符（末尾）
func UnPaddingText(str []byte) []byte {
	n := len(str)
	if n == 0 {
		return []byte{}
	}
	// 末尾的字节标识了未使用的字节数量
	count := int(str[n-1])
	if count > n {
		return str
	}
	newPaddingText := str[:n-count]
	return newPaddingText
}

// region hash

// 32 位 MD5 值
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// region time

// 获取今天零点的时间
// 考虑时区偏移的
//
//	t 为空时，使用当前时间
func ZeroClock(t *time.Time) time.Time {
	if t == nil {
		temp := time.Now()
		t = &temp
	}
	temp := t.Truncate(time.Hour)
	t = &temp
	return t.Add(-time.Duration(t.Hour()) * time.Hour)
}

// 获取今天零点的时间戳
// 考虑时区偏移的
func ZeroClockUnix() int64 {
	today := time.Now()
	return today.Unix() - int64(today.Hour())*3600 - int64(today.Minute())*60 - int64(today.Second())
}

// 获取本周的周一零点的时间
func GetFirstDayOfWeek(today time.Time) time.Time {
	// 获取今天是周几
	weekday := today.Weekday()
	// 周日是 0，周一是 1，以此类推
	// 以周一为一周的开始, 因此要额外减去 1
	weekday--
	if weekday == -1 { // 周末
		weekday = 6
	}
	theDay := today.AddDate(0, 0, -int(weekday))
	// 减去今天是周几，即可得到本周一的日期
	return ZeroClock(&theDay)
}

// 获取本月月初的时间
func GetFirstDayOfMonth(today time.Time) time.Time {
	return time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
}

// 获取 t 所在季度的初始时间
func GetFirstDayOfQuarter(t time.Time) time.Time {
	// 获取当前月份
	month := t.Month()
	if month >= 10 {
		month = 10
	} else if month >= 7 {
		month = 7
	} else if month >= 4 {
		month = 4
	} else {
		month = 1
	}
	// 返回季度的起始月份的开始时间
	return time.Date(t.Year(), month, 1, 0, 0, 0, 0, t.Location())
}

// 获取给定时间所在年的第一天零点
func GetFirstDayOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// 获取给定时间所在的季度（1-4）
func GetQuarterOfYear(t time.Time) int {
	return int((t.Month()-1)/3) + 1
}

// region rand

const numSequence = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// 获取指定长度的随机字符串，仅包含字母和数字
func RandBytes(len int) []byte {
	res := make([]byte, len)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Read(res)

	// 将生成的随机数固定在 str 的区间内，以生成符合要求的随机字符
	for i, x := range res {
		res[i] = numSequence[x%62]
	}
	return res
}

// 将十进制数转换成 36 进制的
// 仅包含数字和小写字母
func NumToAlphabet(num int64) []byte {
	if num == 0 {
		return []byte("0")
	}

	num_str := []byte{}
	ne := false
	if num < 0 {
		ne = true
	}
	for ; num != 0; num = num / 36 {
		offset := num % 36
		if ne {
			offset = -offset
		}
		num_str = append(num_str, numSequence[offset])
	}

	// 倒置 num_str
	if ne {
		num_str = append(num_str, '-')
	}
	for i, j := 0, len(num_str)-1; i < j; i, j = i+1, j-1 {
		num_str[i], num_str[j] = num_str[j], num_str[i]
	}

	return num_str
}

// 判断两个日期是否为同一天
func IsSameDay(x, y time.Time) bool {
	return x.Day() == y.Day() && x.Month() == y.Month() && x.Year() == y.Year()
}

// region array

// 判断数组里的元素是否唯一
func IsUnique[T comparable](arr []T) bool {
	occurred := map[T]bool{}
	for i := range arr {
		if occurred[arr[i]] {
			return false
		}

		occurred[arr[i]] = true
	}
	return true
}

// 判断数组中是否包含指定元素
//
//	arr: 数组
//	ele: 要判断的元素
//
// deprecated: 使用 slices.Contains
func InArray[T comparable](arr []T, ele T) bool {
	for i := range arr {
		if arr[i] == ele {
			return true
		}
	}
	return false
}

// Join 将一个数组连接成一个字符串，每个元素之间用指定的分隔符隔开。
// arr: 要连接的数组，元素类型可以是字符串或整数。
// sep: 用于分隔数组元素的字符串。
// 返回值: 连接后的字符串。
func Join[T string | uint | int | uint64](arr []T, sep string) string {
	// 创建一个字符串切片，用于存放数组中每个元素的字符串形式。
	rows := make([]string, len(arr))

	// 遍历数组，将每个元素转换为字符串并存入rows切片。
	for i := 0; i < len(arr); i++ {
		rows[i] = fmt.Sprintf("%v", arr[i])
	}

	// 使用strings.Join函数将rows切片中的元素连接成一个字符串，元素之间用sep分隔。
	return strings.Join(rows, sep)
}

// 使用 new 依次替换 s 中的 old 字符串，不会重复替换
//
//	old  字符串为空时，返回原字符串
func ReplaceInTurn(s string, old string, new []string) string {
	if old == "" {
		return s
	}
	l := len(old)
	offset := 0
	for i := range new {
		pos := strings.Index(s[offset:], old)
		if pos == -1 {
			break
		}
		pos += offset
		s = s[:pos] + new[i] + s[pos+l:]
		offset += len(new[i])
	}
	return s
}

// 差集
// a 中存在但 b 中不存在的值
// 1.0.5 泛型实现
func Difference[T uint | string](a, b []T) []T {
	res := make([]T, 0)
	for i := range a {
		equal := false
		for j := range b {
			if a[i] == b[j] {
				equal = true
				break
			}
		}
		if !equal {
			res = append(res, a[i])
		}
	}
	return res
}

// HasDifferentElement 检查两个切片代表的集合是否至少有一个不同的元素
// 函数会对元素去重，因此不考虑元素出现的数量，不关心元素出现的顺序
func IsDifferenceIgnoreOrder[T comparable](a, b []T) bool {
	set1 := toSet(a)
	set2 := toSet(b)

	// 如果大小不同，肯定有不同元素
	if len(set1) != len(set2) {
		return true
	}

	// 检查所有元素是否相同
	for k := range set1 {
		if _, exists := set2[k]; !exists {
			return true
		}
	}

	return false
}

// toSet 将切片转换为集合（去重）
func toSet[T comparable](slice []T) map[T]struct{} {
	set := make(map[T]struct{}, len(slice))
	for _, v := range slice {
		set[v] = struct{}{}
	}
	return set
}

// removeElement removes the element at index i from the slice without preserving order
func RemoveElementIgnoreOrder[T any](slice []T, i int) []T {
	if i < 0 || i >= len(slice) {
		return slice
	}

	slice[i] = slice[len(slice)-1] // Replace the element to be removed with the last element
	return slice[:len(slice)-1]    // Truncate the slice
}

// 判断给定的字符串是否为数值型
// 函数通过判断每一位是否为数字来实现
func IsNotNumber(str string) bool {
	for i, l := 0, len(str); i < l; i++ {
		if str[i] < 48 || str[i] > 57 {
			return true
		}
	}
	return false
}

// 判断给定的字符串是手机号(类似的)
// 要求 1 开头的 11 位数字即可
func IsMobile(str string) bool {
	if len(str) != 11 {
		return false
	}

	// 要求第一位为 1
	if str[0] != '1' {
		return false
	}

	// 后边必须都为数值
	for i := 1; i < 11; i++ {
		if str[i] < 48 || str[i] > 57 {
			return false
		}
	}

	return true
}

// region 身份证号

type idcard struct {
	Gender uint8 // 1nan 2nv
	Birth  time.Time
}

// 解析身份证号各部分
// 函数不检查身份证号合法性
func ParseIdCard(card string) (res *idcard, err error) {
	if len(card) != 18 {
		return nil, errors.New("invalid idcard")
	}

	res = new(idcard)
	b := card[16]
	if b%2 == 0 {
		res.Gender = 2
	} else {
		res.Gender = 1
	}

	res.Birth, err = time.ParseInLocation("20060102", card[6:14], time.Local)
	if err != nil {
		return nil, err
	}

	return
}

// 检查身份证号是否正确
func IsIDCard(idcard string) bool {
	if len(idcard) != 18 {
		return false
	}

	sum := 0
	for i, o := range idcard {
		n := 0
		if o >= '0' && o <= '9' {
			n = int(o) - 48
		} else if o == 'X' || o == 'x' {
			n = 10
		} else {
			return false
		}

		switch i {
		case 0:
			sum += n * 7
		case 1:
			sum += n * 9
		case 2:
			sum += n * 10
		case 3:
			sum += n * 5
		case 4:
			sum += n * 8
		case 5:
			sum += n * 4
		case 6:
			sum += n * 2
		case 7:
			sum += n * 1
		case 8:
			sum += n * 6
		case 9:
			sum += n * 3
		case 10:
			sum += n * 7
		case 11:
			sum += n * 9
		case 12:
			sum += n * 10
		case 13:
			sum += n * 5
		case 14:
			sum += n * 8
		case 15:
			sum += n * 4
		case 16:
			sum += n * 2
		case 17:
			// 校验, 余数对应 1－0－X－9－8－7－6－5－4－3－2
			re := sum % 11
			switch re {
			case 0:
				return n == 1
			case 1:
				return n == 0
			case 2:
				return n == 10
			case 3:
				return n == 9
			case 4:
				return n == 8
			case 5:
				return n == 7
			case 6:
				return n == 6
			case 7:
				return n == 5
			case 8:
				return n == 4
			case 9:
				return n == 3
			case 10:
				return n == 2
			}
		default:
			return false
		}
	}
	return false
}

// 港澳身份证号格式正则（匹配前已去除括号、空格与斜杠分隔符）
var (
	// 香港：1-2 位字母 + 6 位数字 + 1 位校验码（数字或字母 A）
	hkIDRegex = regexp.MustCompile(`^[A-Z]{1,2}[0-9]{6}[0-9A]$`)
	// 澳门：首位为 1、5 或 7（代表取证时代）+ 6 位数字 + 1 位查核用数码（数字或字母 A）
	macauIDRegex = regexp.MustCompile(`^[157][0-9]{6}[0-9A]$`)
)

// normalizeIDCard 统一港澳身份证号的书写形式：
// 转为大写，并去除首尾空格以及号码中的括号、空格与斜杠等分隔符
func normalizeIDCard(idcard string) string {
	s := strings.ToUpper(strings.TrimSpace(idcard))
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(', ')', ' ', '/':
			// 忽略括号、空格与斜杠等分隔符
		default:
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// hkCharValue 返回港澳身份证校验中单个字符对应的数值
//
//	字母 A-Z 依次取值 10-35，数字取值本身，空格取值 36
func hkCharValue(c byte) int {
	switch {
	case c >= 'A' && c <= 'Z':
		return int(c-'A') + 10
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c == ' ':
		return 36
	default:
		return 0
	}
}

// 检查香港身份证号是否正确
//
//	格式为 X123456(A) 或 XY123456(A)：1-2 位字母 + 6 位数字 + 括号内的校验码（0-9 或 A）
//	兼容不带括号、字母小写以及夹带空格的写法
//
// 校验码采用模 11 算法：字母 A-Z 依次取值 10-35，空格取值 36，数字取值本身；
// 从左到右权重依次为 9、8、7、6、5、4、3、2，仅一位字母时在最左侧补一个空格；
// 加权和的余数对应校验码 (11 - sum%11) % 11，其中结果 10 用字母 A 表示
func IsHKIDCard(idcard string) bool {
	s := normalizeIDCard(idcard)
	if !hkIDRegex.MatchString(s) {
		return false
	}

	// 分离主体（字母 + 6 位数字）与末尾校验码
	body := s[:len(s)-1]
	check := s[len(s)-1]

	// 仅一位字母时在最左侧补空格，使主体固定为 8 位
	if len(body) == 7 {
		body = " " + body
	}

	// 权重从左到右依次为 9..2
	sum := 0
	for i := 0; i < len(body); i++ {
		sum += hkCharValue(body[i]) * (9 - i)
	}

	// 余数对应校验码，其中 10 用字母 A 表示
	expect := (11 - sum%11) % 11
	if expect == 10 {
		return check == 'A'
	}
	return check == byte('0'+expect)
}

// 检查澳门身份证号是否正确
//
//	格式为 XNNNNNN(Y) 或旧式 X/NNNNNN/Y：首位为 1、5 或 7（代表取证时代）+ 6 位数字 + 括号内的查核用数码（0-9 或 A）
//	兼容不带括号、以及使用斜杠分隔的写法
//
// 注意：澳门身份证的查核用数码（校验码）算法并未公开标准化，
// 因此本函数仅校验号码格式，不校验查核用数码本身的正确性
func IsMacauIDCard(idcard string) bool {
	return macauIDRegex.MatchString(normalizeIDCard(idcard))
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func IsFileExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// 当给定的文件夹不存在时, 创建它
func MakeDir(dir string) error {
	if !IsFileExist(dir) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// 当给定的文件夹不存在时, 创建它
// 如果给定的路径末尾包含了文件后缀，那么忽略它，只创建之前的路径
func MakeDirTrimFileName(dir string) error {
	if dir == "" {
		return nil
	}

	i := strings.LastIndex(dir, "/")
	if i == -1 {
		i = 0
	}
	j := strings.LastIndex(dir[i:], ".") // 在最后一个目录中查找文件后缀。查到则忽略整个区间
	if j != -1 {
		dir = dir[:i]
	}
	if dir == "" {
		return nil
	}

	if !IsFileExist(dir) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// region string

// 截取前 end 个字符（rune）
//
//	如果 end 小于等于 0，返回空字符串
func SubStr(s string, end int) string {
	temp := []rune(s)
	if end <= 0 {
		return ""
	}
	if end >= len(temp) {
		return s
	}
	return string(temp[:end])
}

// 对用户姓名进行匿名操作, 只显示开头结尾的名字
// 如: 李晓明 -> 李*明
//
//	onlyFirstName: 是否只显示第一个名字. 默认仅隐藏中间字.
func Anonymous(name string, onlyFirstName bool) string {
	if name == "" {
		return name
	}

	uni := []rune(name)
	if len(uni) <= 2 {
		return string(uni[0]) + "*"
	}

	if onlyFirstName {
		return string(uni[0]) + "**"
	}

	return string(uni[0]) + "*" + string(uni[len(uni)-1])
}

// 通过指定保留的前缀长度，后缀长度，隐藏字符串中间的部分
// 可以用于诸如手机号，身份证号等的匿名化处理
func AnonymizeByOffset(name string, prefix int, suffix int) string {
	if name == "" {
		return name
	}
	uni := []rune(name)
	if len(uni) <= prefix+suffix {
		return name
	}
	return string(uni[:prefix]) + strings.Repeat("*", len(uni)-prefix-suffix) + string(uni[len(uni)-suffix:])
}

// 检查要求字符串只能包含uint和逗号，且数字必须合法
func IsNumberCombinedWithComma(s string) bool {
	ss := strings.Split(s, ",")
	for i := range ss {
		num, err := strconv.Atoi(ss[i])
		// err: maybe out of range
		if err != nil || num == 0 {
			return false
		}
	}
	return true
}

// 提取 html 中会被渲染出的文本
//
// 使用 html.Parse 解析并递归提取 TextNode，会跳过 script 和 style 标签中的文本
//
// 如果要解析的 html 相对简单，比如没有 script，不会出现针对性异常属性，则可以直接使用正则版本，效率会更高 ExtractHtmlTextWithRegex
func ExtractHtmlText(r *strings.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		// 忽略 <script> 和 <style> 的全部子节点
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			return
		}
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(doc)
	return sb.String(), nil
}

// 预编译正则，避免在每次调用时重复编译
var tagRegex = regexp.MustCompile(`<[^>]*>`)

// 提取 html 中的文本
//
// 正则版本：直接替换所有标签为空字符串
//
// 正则版本速度会更快，但是无法处理所有的情况，比如 <a title="a > b">link</a>，<script>console.log('world')</script> 内的文本会被展示，因此如果需要更准确的结果，建议使用 ExtractHtmlText
//
// 你也可以在部分数据库中直接使用 REGEXP_REPLACE(s.profile, '<[^>]*>', '') 来实现同样的功能，避免取回大量数据
func ExtractHtmlTextWithRegex(htmlStr string) string {
	return tagRegex.ReplaceAllString(htmlStr, "")
}

// region float

// 对浮点数进行四舍五入,保留prec位
func RoundPrec(num float64, prec int) float64 {
	var factor float64 = 1
	for i := 0; i < prec; i++ {
		factor *= 10
	}
	return math.Round(num*factor) / factor
}

/*
 * 中国正常GCJ02坐标---->百度地图BD09坐标
 * 腾讯地图用的也是GCJ02坐标
 * @param float64 x 经度
 * @param float64 y 纬度
 * @return lng:经度 lat: 纬度
 */
func ConvertGCJ02ToBD09(x, y float64) (lng, lat float64) {
	pi := math.Pi * 3000.0 / 180.0
	z := math.Sqrt(x*x+y*y) + 0.00002*math.Sin(y*pi)
	theta := math.Atan2(y, x) + 0.000003*math.Cos(x*pi)
	lng = z*math.Cos(theta) + 0.0065
	lat = z*math.Sin(theta) + 0.006
	return
}

/**
 * 百度地图BD09坐标---->中国正常GCJ02坐标
 * 腾讯地图用的也是GCJ02坐标
 * @param float64 x 经度
 * @param float64 y 纬度
 * @return lng:经度 lat: 纬度
 */
func ConvertBD09ToGCJ02(x, y float64) (lng, lat float64) {
	pi := math.Pi * 3000.0 / 180.0
	x = x - 0.0065
	y = y - 0.006
	z := math.Sqrt(x*x+y*y) + 0.00002*math.Sin(y*pi)
	theta := math.Atan2(y, x) + 0.000003*math.Cos(x*pi)
	lng = z * math.Cos(theta)
	lat = z * math.Sin(theta)
	return
}

// 获取随机数
func RandNumber(length int) int {
	num := 10
	for i := 1; i < length-1; i++ {
		num *= 10
	}
	number := rand.Intn(9*num) + num
	return number
}

// 判断字符串是否为汉字
func IsChineseCharacter(str string) bool {
	if str == "" {
		return false
	}
	for _, v := range str {
		if !unicode.Is(unicode.Han, v) {
			return false
		}
	}
	return true
}
func CheckMobile(mobile string) bool {
	matched, _ := regexp.MatchString(`^1[3456789]\d{9}$`, strings.TrimSpace(mobile))
	return matched
}

// des加密,返回加密后的16进制字符串
func EncryptByDESCBC(message, key, iv string) (string, error) {
	// 将key和iv转换为字节切片
	keyBytes := []byte(key)
	ivBytes := []byte(iv)

	// 创建DES密钥
	block, err := des.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// 初始化CBC加密模式
	blockMode := cipher.NewCBCEncrypter(block, ivBytes)

	// 对明文进行PKCS7填充
	data := []byte(message)
	blockSize := block.BlockSize()
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	paddedMessage := append(data, padText...)

	// 创建足够大的字节切片来存放加密后的数据
	ciphertext := make([]byte, len(paddedMessage))

	// 执行加密
	blockMode.CryptBlocks(ciphertext, paddedMessage)

	// 返回Base64编码的密文
	return strings.ToUpper(hex.EncodeToString(ciphertext)), nil
}

// 根据出生日期计算年龄
func CalAge(birthday, today time.Time) (age int) {
	age = today.Year() - birthday.Year()
	if today.Month() < birthday.Month() {
		age--
	} else if today.Month() == birthday.Month() {
		if today.Day() < birthday.Day() {
			age--
		}
	}

	return age
}

// 计算 a 的增长率. 即 (b-a)/a ,结果保留两位小数
// 返回单位 %
func CalculateGrowthPercent(a, b int) float64 {
	if a == b {
		return 0
	}
	if a == 0 { // 防止除以0的情况
		if b > 0 {
			return float64(100)
		} else if b < 0 {
			return float64(-100)
		}
		return 0
	}
	percent := float64(b-a) / float64(a) * 100
	return math.Round(percent*100) / 100 // 保留两位小数, 四舍五入
}

// 计算a/b，结果保留两位浮点数
func Divide(a, b int) float64 {
	if b == 0 {
		return 0
	}
	result := float64(a) / float64(b)
	res := math.Round(result*100) / 100 // 先乘以100四舍五入，再除以100得到保留两位的结果
	return res
}

// region print

// 将 data 转换为 json 字符串
func PrintJson(data any) (string, error) {
	a, err := json.Marshal(data)
	return string(a), err
}

// 将 data 转换为 json 字符串
// 总是不返回错误，便于某些场景下的行内打印
func PrintJsonWithoutErr(data any) string {
	a, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(a)
}
