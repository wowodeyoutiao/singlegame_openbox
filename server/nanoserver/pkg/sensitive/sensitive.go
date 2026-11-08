package sensitive

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"nanoserver/global"
	"regexp"
	"strings"
)

type SensitiveLexicon struct {
	Words         []string `json:"words"`
	FilteredWords []string
	Re            *regexp.Regexp
	CERegex       *regexp.Regexp
}

var lexicon SensitiveLexicon

func IsValidName(text string) bool {
	if !lexicon.CERegex.MatchString(text) {
		return false
	}

	return true
}

func ContainsSensitiveWord(text string) bool {
	for _, word := range lexicon.Words {
		if strings.Contains(text, word) {
			return true
		}
	}
	return false
}

func ReplaceSensitiveWords(text string) string {
	// 检查是否有敏感词
	if len(lexicon.Words) == 0 {
		return text
	}

	// 替换文本中的敏感词
	return lexicon.Re.ReplaceAllString(text, "*")
}

func loadSensitiveLexicon(filename string) (SensitiveLexicon, error) {
	var lexicon SensitiveLexicon

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return lexicon, err
	}

	err = json.Unmarshal(data, &lexicon)
	if err != nil {
		return lexicon, err
	}

	return lexicon, nil
}

// 去除字符串切片中的重复元素
func removeDuplicates(words []string) []string {
	// 使用map来存储已经遇到的元素
	encountered := make(map[string]bool)

	// 用于存储唯一元素的切片
	uniqueWords := []string{}

	// 遍历原始切片
	for _, word := range words {
		// 如果元素尚未被添加到map中，则添加到map和结果切片中
		if !encountered[word] {
			encountered[word] = true
			uniqueWords = append(uniqueWords, word)
		}
	}

	return uniqueWords
}

func Load() {
	var err error
	lexicon, err = loadSensitiveLexicon(global.BasePath + "/configs/gameconfig/SensitiveLexicon.json")
	if err != nil {
		log.Panic("Error loading sensitive lexicon:", err)
		return
	}

	lexicon.Words = removeDuplicates(lexicon.Words)

	// 避免正则表达式错误：如果敏感词列表中有空字符串，需要过滤掉
	lexicon.FilteredWords = make([]string, 0, len(lexicon.Words))
	for _, word := range lexicon.Words {
		if word != "" {
			// 使用regexp.QuoteMeta确保正则表达式中的特殊字符被正确处理
			lexicon.FilteredWords = append(lexicon.FilteredWords, regexp.QuoteMeta(word))
		}
	}
	// 编译正则表达式
	lexicon.Re, err = regexp.Compile(strings.Join(lexicon.FilteredWords, "|"))
	if err != nil {
		log.Panic("Error compiling regex:", err)
	}

	// 汉字+字母
	lexicon.CERegex = regexp.MustCompile(`^[\p{Han}a-zA-Z]+$`)
}
