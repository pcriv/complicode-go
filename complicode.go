package complicode

import (
	"crypto/rc4"
	"encoding/hex"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/chtison/baseconverter"
	"github.com/osamingo/checkdigit"
)

var verhoeff = checkdigit.NewVerhoeff()

type Invoice struct {
	Nit    int
	Number int
	Amount float64
	Date   time.Time
}

type asciiSums struct {
	total    int
	partials []int
}

const base64 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz+/"

// Generate a code for the given authCode, key and invoice
func Generate(authCode string, key string, inv Invoice) string {
	seed1 := strconv.Itoa(inv.Number)
	seed2 := strconv.Itoa(inv.Nit)
	seed3 := inv.Date.Format("20060102")
	seed4 := strconv.FormatFloat(math.Round(inv.Amount), 'f', -1, 64)

	seeds := []string{seed1, seed2, seed3, seed4}
	seeds = appendVerificationsDigits(seeds, 2)
	digits := generateVerificationDigits(seeds)
	partialKeys := generatePartialKeys(key, digits)

	seeds = append([]string{authCode}, seeds...)
	seeds = appendPartialKeys(seeds, partialKeys)

	encryptionKey := key + digits
	encryptedData := encrypt(strings.Join(seeds, ""), encryptionKey)
	asciiSums := generateASCIISums(encryptedData, len(partialKeys))

	sum := 0

	for i, partialSum := range asciiSums.partials {
		sum += (asciiSums.total * partialSum / len(partialKeys[i]))
	}

	data := changeBase(sum)
	code := encrypt(data, encryptionKey)

	return format(code)
}

func appendVerificationsDigits(seeds []string, count int) []string {
	for index, seed := range seeds {
		seeds[index] = appendVerificationDigits(seed, count)
	}

	return seeds
}

func appendVerificationDigits(seed string, count int) string {
	for i := 1; i <= count; i++ {
		digit, _ := verhoeff.Generate(seed)
		seed += strconv.Itoa(digit)
	}

	return seed
}

func generateVerificationDigits(seeds []string) string {
	sum := sumSeeds(seeds)
	digits := appendVerificationDigits(strconv.Itoa(sum), 5)

	return digits[len(digits)-5:]
}

func sumSeeds(seeds []string) int {
	sum := 0

	for _, seed := range seeds {
		number, _ := strconv.Atoi(seed)
		sum += number
	}

	return sum
}

func generatePartialKeys(key string, verificationDigits string) []string {
	var partialKeySizes []int

	for _, digit := range strings.Split(verificationDigits, "") {
		size, _ := strconv.Atoi(digit)
		partialKeySizes = append(partialKeySizes, size+1)
	}

	var partialKeys []string

	start := 0
	for _, size := range partialKeySizes {
		end := start + size
		pk := key[start:end]
		partialKeys = append(partialKeys, pk)
		start += size
	}

	return partialKeys
}

func appendPartialKeys(seeds []string, partialKeys []string) []string {
	for index, seed := range seeds {
		seeds[index] = seed + partialKeys[index]
	}

	return seeds
}

func encrypt(data string, key string) string {
	cipher, _ := rc4.NewCipher([]byte(key))
	encrypted := make([]byte, len(data))
	cipher.XORKeyStream(encrypted, []byte(data))

	return strings.ToUpper(hex.EncodeToString(encrypted))
}

func generateASCIISums(data string, partialsCount int) asciiSums {
	sums := asciiSums{total: 0, partials: make([]int, partialsCount)}

	for i, b := range []byte(data) {
		sums.total += int(b)
		sums.partials[i%partialsCount] += int(b)
	}

	return sums
}

func changeBase(number int) string {
	result, _ := baseconverter.UInt64ToBase(uint64(number), base64)

	return result
}

func format(code string) string {
	formatted := ""

	for i, char := range code {
		if i != 0 && i%2 == 0 {
			formatted += "-"
		}
		formatted += string(char)
	}

	return formatted
}
