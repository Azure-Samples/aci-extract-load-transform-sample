package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	defaultFilePath   = "./cleaned_hm.csv"
	defaultResultPath = "./result.csv"
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	LIMIT_WORDS := [...]string{"happy", "day", "got", "went", "today", "made", "one", "two", "time", "last", "first", "going", "getting", "took", "found", "lot", "really", "saw", "see", "month", "week", "day", "yesterday", "year", "ago", "now", "still", "since", "something", "great", "good", "long", "thing", "toi", "without", "yesteri", "2s", "toand", "ing", "came", "able", "bought", "go"}

	STOP_WORDS := map[string]bool{"a": true, "about": true, "above": true, "after": true, "again": true, "against": true, "all": true, "am": true, "an": true, "and": true, "any": true, "are": true, "aren't": true, "as": true, "at": true, "be": true, "because": true, "been": true, "before": true, "being": true, "below": true, "between": true, "both": true, "but": true, "by": true, "can": true, "can't": true, "cannot": true, "com": true, "could": true, "couldn't": true, "did": true, "didn't": true, "do": true, "does": true, "doesn't": true, "doing": true, "don't": true, "down": true, "during": true, "each": true, "else": true, "ever": true, "few": true, "for": true, "from": true, "further": true, "get": true, "had": true, "hadn't": true, "has": true, "hasn't": true, "have": true, "haven't": true, "having": true, "he": true, "he'd": true, "he'll": true, "he's": true, "her": true, "here": true, "here's": true, "hers": true, "herself": true, "him": true, "himself": true, "his": true, "how": true, "how's": true, "http": true, "i": true, "i'd": true, "i'll": true, "i'm": true, "i've": true, "if": true, "in": true, "into": true, "is": true, "isn't": true, "it": true, "it's": true, "its": true, "itself": true, "just": true, "k": true, "let's": true, "like": true, "me": true, "more": true, "most": true, "mustn't": true, "my": true, "myself": true, "no": true, "nor": true, "not": true, "of": true, "off": true, "on": true, "once": true, "only": true, "or": true, "other": true, "ought": true, "our": true, "ours ": true, "ourselves": true, "out": true, "over": true, "own": true, "r": true, "same": true, "shall": true, "shan't": true, "she": true, "she'd": true, "she'll": true, "she's": true, "should": true, "shouldn't": true, "so": true, "some": true, "such": true, "than": true, "that": true, "that's": true, "the": true, "their": true, "theirs": true, "them": true, "themselves": true, "then": true, "there": true, "there's": true, "these": true, "they": true, "they'd": true, "they'll": true, "they're": true, "they've": true, "this": true, "those": true, "through": true, "to": true, "too": true, "under": true, "until": true, "up": true, "very": true, "was": true, "wasn't": true, "we": true, "we'd": true, "we'll": true, "we're": true, "we've	": true, "were": true, "weren't": true, "what": true, "what's": true, "when": true, "when's": true, "where": true, "where's": true, "which": true, "while": true, "who": true, "who's": true, "whom": true, "why": true, "why's": true, "with": true, "won't": true, "would": true, "wouldn't": true, "www": true, "you": true, "you'd": true, "you'll": true, "you're": true, "you've": true, "your": true, "yours": true, "yourself": true, "yourselves": true}

	filePath, ok := os.LookupEnv("FILE_PATH")
	if !ok {
		log.Println("FILE_PATH not set. Using default.")
		filePath = defaultFilePath
	}

	resultPath, ok := os.LookupEnv("RESULT_PATH")
	if !ok {
		log.Println("RESULT_PATH not set. Using default.")
		resultPath = defaultResultPath
	}

	for {
		exists, _ := pathExists(filePath)
		if exists {
			break
		}
		log.Println("Source file is not exist, retry in 1s")
		time.Sleep(1 * time.Second)
	}

	in, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(strings.NewReader(string(in)))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer

	for i := 1; i < len(records); i++ {
		buf.WriteString(records[i][4])
		buf.WriteString(" ")
	}

	text := strings.ToLower(buf.String())
	allSignals := regexp.MustCompile("[^a-z]+")
	for _, v := range LIMIT_WORDS {
		text = strings.Replace(text, " "+v, "", -1)
		text = strings.Replace(text, v+" ", "", -1)
		text = allSignals.ReplaceAllString(text, " ")
	}

	regex := regexp.MustCompile("[a-z]*")
	// allwords := regex.FindAllString(text, -1)
	allwords := strings.Fields(text)

	m := make(map[string]int)
	for _, n := range allwords {
		if n == "" || STOP_WORDS[n] || STOP_WORDS[regex.FindString(n)] {
			continue
		} else if count, ok := m[n]; ok == true {
			m[n] = count + 1
		} else {
			m[n] = 1
		}
	}

	out, err := os.Create(resultPath + ".tmp")
	if err != nil {
		log.Fatal(err)
	}
	out.WriteString("Word,Count\n")
	i := 0
	count := 0
	for k, v := range m {
		if v >= 500 && len(k) > 1 {
			count++
			out.WriteString(fmt.Sprintf("%s,%d\n", k, v))
		}
		i++
	}
	out.Close()
	os.Rename(resultPath+".tmp", resultPath)

	log.Println(count)
}
