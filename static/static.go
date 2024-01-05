package static

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"unicode"

	"gitlab.com/daneofmanythings/wohrdle/utils"
)

// TODO: fix this. It is hella error prone
const (
	targetRawWordListAbsPath string = "/etc/dictionaries-common/words" // for linux
	localWordTxt             string = "./static/words.txt"             // relative to the root
	localWordJSON            string = "./static/words.json"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func MakeWordFile() {
	file, err := os.Open(targetRawWordListAbsPath)
	check(err)
	defer file.Close()

	var words string
	var intermediate string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		intermediate = ""
		intermediate += scanner.Text()
		for _, r := range intermediate {
			if unicode.IsUpper(r) || r == '\'' || !utils.RuneIsAlpha(r) {
				intermediate = ""
				break
			}
		}
		if intermediate != "" {
			intermediate += "\n"
		}
		words += intermediate
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	f, err2 := os.Create(localWordTxt)
	check(err2)
	defer f.Close()

	f.WriteString(words)
	f.Sync()
}

func CleanWordFile() {
	file, err := os.Open(targetRawWordListAbsPath)
	check(err)
	defer file.Close()

	var words string
	var intermediate string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		intermediate = ""
		intermediate += scanner.Text()
		for _, r := range intermediate {
			if unicode.IsUpper(r) {
				intermediate = ""
			} else if r == '\'' {
				intermediate = ""
			}
		}
		words += intermediate
		words += "\n"
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	f, err2 := os.Create(localWordTxt)
	check(err2)
	defer f.Close()

	f.WriteString(words)
	f.Sync()
}

func CreateJSONFile() {
	wordRepo := utils.WordRepository{
		Words: map[string][]string{},
	}

	file, err := os.Open(localWordTxt)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := scanner.Text()
		wordLen := len(word)
		wordLenStr := strconv.Itoa(wordLen)
		wordRepo.Words[wordLenStr] = append(wordRepo.Words[wordLenStr], word)
	}

	b, e := json.MarshalIndent(wordRepo, "", "\t")
	check(e)

	f, err2 := os.Create(localWordJSON)
	check(err2)

	f.Write(b)
}
