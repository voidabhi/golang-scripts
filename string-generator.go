package main

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numerals = "0123456789"
	Ascii    = Alphabet + Numerals + "~!@#$%^&*()-_+={}[]\\|<,>.?/\"';:`"
)

type GeneratorExprRanges [][]byte

func seedAndReturnRandom(n int) int {
	return rand.Intn(n)
}

func alphabetSlice(from, to byte) (string, error) {
	leftPos := strings.Index(Ascii, string(from))
	rightPos := strings.LastIndex(Ascii, string(to))
	if leftPos > rightPos {
		return "", fmt.Errorf("Invalid range specified: %s-%s", string(from), string(to))
	}
	return Ascii[leftPos:rightPos], nil
}

func replaceWithGenerated(s *string, expresion string, ranges [][]byte, length int) error {
	var alphabet string
	for _, r := range ranges {
		switch string(r[0]) + string(r[1]) {
		case `\w`:
			alphabet += Ascii
		case `\d`:
			alphabet += Numerals
		default:
			if slice, err := alphabetSlice(r[0], r[1]); err != nil {
				return err
			} else {
				alphabet += slice
			}
		}
	}
	if len(alphabet) == 0 {
		return fmt.Errorf("Empty range in expresion: %s", expresion)
	}
	result := make([]byte, length, length)
	for i := 0; i <= length-1; i++ {
		result[i] = alphabet[seedAndReturnRandom(len(alphabet))]
	}
	*s = strings.Replace(*s, expresion, string(result), 1)
	return nil
}

func findExpresionPos(s string) GeneratorExprRanges {
	rangeExp, _ := regexp.Compile(`([\\]?[a-zA-Z0-9]\-?[a-zA-Z0-9]?)`)
	matches := rangeExp.FindAllStringIndex(s, -1)
	result := make(GeneratorExprRanges, len(matches), len(matches))
	for i, r := range matches {
		result[i] = []byte{s[r[0]], s[r[1]-1]}
	}
	return result
}

func rangesAndLength(s string) (string, int, error) {
	expr := s[0:strings.LastIndex(s, "{")]
	length, err := parseLength(s)
	return expr, length, err
}

func parseLength(s string) (int, error) {
	lengthStr := string(s[strings.LastIndex(s, "{")+1 : len(s)-1])
	if l, err := strconv.Atoi(lengthStr); err != nil {
		return 0, fmt.Errorf("Unable to parse length from %v", s)
	} else {
		return l, nil
	}
}

func Generate(template string) (string, error) {
	result := template
	generatorsExp, _ := regexp.Compile(`\[([a-zA-Z0-9\-\\]+)\](\{([0-9]+)\})`)
	matches := generatorsExp.FindAllStringIndex(template, -1)
	for _, r := range matches {
		ranges, length, err := rangesAndLength(template[r[0]:r[1]])
		if err != nil {
			return "", err
		}
		positions := findExpresionPos(ranges)
		if err := replaceWithGenerated(&result, template[r[0]:r[1]], positions, length); err != nil {
			return "", err
		}
	}
	return result, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Generate 8 character long random string that contains *all* Ascii
	// characters (See Ascii const)
	result, _ := Generate(`foo[\w]{8}bar`)
	log.Printf("result0=%v", result)

	// Generate 3 character long random suffix using just lowercase alpha+numberic
	// characters
	result, _ = Generate(`admin[a-z0-9]{3}`)
	log.Printf("result1=%v", result)

	result, _ = Generate(`admin[a-z0-9]{3}something[\w]{3}`)
	log.Printf("result1=%v", result)

	// Generate 12 character long random password combining all alphabetic characters
	// and numeric characters
	result, _ = Generate(`pass[a-zA-Z0-9]{12}`)
	log.Printf("result2=%v", result)

	// Generate 8 character long password by using alphabetic characters
	// lowercase + uppercase
	result, _ = Generate(`pass[a-Z]{8}`)
	log.Printf("result3=%v", result)

	// Invalid range error
	result, err := Generate(`pass[z-a]{8}`)
	log.Printf("result4=%v; err=%v", result, err)

	// Generate 8 random numeric characters
	result, _ = Generate(`foo[\d]{8}bar`)
	log.Printf("result5=%v", result)
}
