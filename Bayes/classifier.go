package bayes

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

var filtered_words []string = []string{
	"there",
	"is",
	"are",
	"a",
	"my",
	"mine",
	"what",
	"someone",
	"you",
	"on",
	"if",
	"at",
	"ever",
	"never",
	"just",
	"but",
	"to",
	"then",
	"whenever",
	"also",
	"ok",
}

func GetFiles(folder string) []string {
	out, err := exec.Command("ls", folder).CombinedOutput()
	if err != nil {
		panic(err)
	}
	var files []string = make([]string, 0)

	i := 0
	for {
		if i >= len(out) {
			break
		}
		if out[i] != '\n' || out[i] != ' ' {
			var file []byte = make([]byte, 0)
			for {
				if out[i] == '\n' || out[i] == ' ' {
					break
				}
				file = append(file, out[i])
				i++
			}
			files = append(files, string(file))
		}
		i++
	}

	return files
}

func GetContents(folder string) []string {
	out_files := GetFiles(folder)

	var texts []string = make([]string, 0)

	for i := 0; i < len(out_files); i++ {
		out, err := os.ReadFile(folder + "/" + out_files[i])
		if err != nil {
			panic(err)
		}
		texts = append(texts, string(out))
	}

	return texts
}

func GetContent(file string) string {
	var text []byte = make([]byte, 0)

	out, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	for _, el := range out {
		text = append(text, el)
	}
	return string(text)
}
func TokenizeWords(folder string) map[string]float64 {
	tokenmap := make(map[string]float64)

	var file_contents []string = GetContents(folder)

	for _, el := range file_contents {
		j := 0
		token := make([]byte, 0)
		for {
			if j >= len(el) {
				break
			}
			str1 := el
			char1 := el[j]
			if len(str1) > 0 && (!unicode.IsLetter(rune(char1))) {
				if !unicode.IsLetter(rune(char1)) {
					if len(token) >= 12 {
						token = nil
					} else {
						tokenmap[strings.ToLower(string(token))]++
						token = nil
					}
				}
			}

			if unicode.IsLetter(rune(char1)) {
				token = append(token, byte(char1))
			}
			j++
		}
	}

	return tokenmap
}

func TokenizeWords_File(file string) map[string]float64 {
	j := 0
	var content string = GetContent(file)
	tokenmap := make(map[string]float64)
	token := make([]byte, 0)

	for {
		if j >= len(content) {
			break
		}
		str1 := content
		char1 := content[j]
		if len(str1) > 0 && (!unicode.IsLetter(rune(char1))) {
			if !unicode.IsLetter(rune(char1)) {
				if len(token) >= 12 {
					token = nil
				} else {
					tokenmap[strings.ToLower(string(token))]++
					token = nil
				}
			}
		}
		if unicode.IsLetter(rune(char1)) {
			token = append(token, byte(char1))
		}
		j++
	}
	return tokenmap
}
func ParseFile(file string) []string {
	j := 0
	var content string = GetContent(file)
	tokenlist := make([]string, 0)
	token := make([]byte, 0)

	for {
		if j >= len(content) {
			break
		}
		str1 := content
		char1 := content[j]
		if len(str1) > 0 && (!unicode.IsLetter(rune(char1))) {
			if !unicode.IsLetter(rune(char1)) {
				if len(token) >= 12 {
					token = nil
				} else {
					found := false
					for _, el := range filtered_words {
						if strings.Compare(strings.ToLower(el), strings.ToLower(string(token))) == 0 {
							found = true
							break
						}
					}
					if !found {
						tokenlist = append(tokenlist, strings.ToLower(string(token)))
					}
					token = nil
				}
			}
		}
		if unicode.IsLetter(rune(char1)) {
			token = append(token, byte(char1))
		}
		j++
	}
	return tokenlist
}
func CheckFile(file string, label string, crime map[string]float64, entertainment map[string]float64, politics map[string]float64, science map[string]float64) bool {
	file_parsed := ParseFile(file)
	return Classify_Text(file_parsed, label, crime, entertainment, politics, science)
}

func P_Word_Class(map1 map[string]float64) map[string]float64 {
	var return_map map[string]float64 = make(map[string]float64)
	var wordcount float64 = 0.0
	for _, value := range map1 {
		wordcount += value
	}
	for key, value := range map1 {
		return_map[key] = value / wordcount
	}
	return return_map
}

func Display(map1 map[string]float64) {
	for key, value := range map1 {
		fmt.Printf("\n%s : %.8f", key, value)
	}
}

func GetTotalWordSize(map1 map[string]float64) float64 {
	var sum float64 = 0.0
	for _, value := range map1 {
		sum += value
	}
	return sum
}

func MergeMaps(map1 map[string]float64, map2 map[string]float64) map[string]float64 {
	var mergedmap map[string]float64 = make(map[string]float64)

	for key, value := range map1 {
		mergedmap[key] = value
	}
	for key, value := range map2 {
		mergedmap[key] += value
	}
	return mergedmap
}

func P_Word(map1 map[string]float64, key string) float64 {
	var wordcount float64 = 0.0
	for _, val := range map1 {
		wordcount += val
	}

	return map1[key] / wordcount

}

func Predict(X map[string]float64, W map[string]float64) float64 {
	var sum float64 = 0.0
	for key, value := range X {
		sum += value * W[key]
	}
	return sum
}

func P_Class_Word(p_word_class map[string]float64, p_class float64, p_word map[string]float64) map[string]float64 {
	var p_c_w map[string]float64 = make(map[string]float64)

	for key, value := range p_word_class {
		p_c_w[key] = p_class * value / p_word[key]
	}

	return p_c_w
}

func log_sum(list []float64) float64 {
	var sum float64 = 0.0

	for _, el := range list {
		sum += math.Log(el)
	}
	return sum
}

func Classify_Text(text []string, label string, crime map[string]float64, entertainment map[string]float64, politics map[string]float64, science map[string]float64) bool {
	var crime_mul []float64 = make([]float64, 0)
	var entertainment_mul []float64 = make([]float64, 0)
	var politics_mul []float64 = make([]float64, 0)
	var science_mul []float64 = make([]float64, 0)
	for _, val := range text {
		crime_mul = append(crime_mul, crime[val])
		entertainment_mul = append(entertainment_mul, entertainment[val])
		politics_mul = append(politics_mul, politics[val])
		science_mul = append(science_mul, science[val])
	}
	var crime_mul_ float64 = log_sum(crime_mul)
	var entertainment_mul_ float64 = log_sum(entertainment_mul)
	var politics_mul_ float64 = log_sum(politics_mul)
	var science_mul_ float64 = log_sum(science_mul)

	var max_ float64 = math.Max(math.Max(crime_mul_, entertainment_mul_), math.Max(politics_mul_, science_mul_))

	if max_ == crime_mul_ {
		if label == "c" {
			return true
		}
	} else if max_ == entertainment_mul_ {
		if label == "e" {
			return true
		}
	} else if max_ == politics_mul_ {
		if label == "p" {
			return true
		}
	} else if max_ == science_mul_ {
		if label == "s" {
			return true
		}
	}
	return false
}
