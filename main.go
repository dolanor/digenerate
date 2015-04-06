package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

var (
	letterFrequency = map[string]float32{"e": 13}
)

type DigraphPair struct {
	Digraph string
	Count   int
}

type DigraphPairList []DigraphPair

func (dp DigraphPairList) Swap(i, j int)      { dp[i], dp[j] = dp[j], dp[i] }
func (dp DigraphPairList) Len() int           { return len(dp) }
func (dp DigraphPairList) Less(i, j int) bool { return dp[i].Count < dp[j].Count }

func sortDigraphs(digraphFreq map[string]int) (DigraphPairList, DigraphPairList) {
	d := make(DigraphPairList, len(digraphFreq))

	i := 0
	for k, v := range digraphFreq {
		d[i] = DigraphPair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(d))
	s := make(DigraphPairList, len(digraphFreq))
	i = 0
	for _, v := range d {
		s[i] = DigraphPair{hex.EncodeToString([]byte(v.Digraph)), v.Count}
		i++
	}
	fmt.Println(s)
	return d, s
}

func findMostUsedDigraph(ciphertext string) string {
	digraphFreq := make(map[string]int)
	for j := 0; j < len(ciphertext)-1; j++ {
		currdigraph := ciphertext[j : j+2]
		for i := 0; i < len(ciphertext)-1; i++ {
			if ciphertext[i:i+2] == currdigraph {
				fmt.Printf("%d + 1 %x\n", digraphFreq[currdigraph], currdigraph)
				digraphFreq[currdigraph] += 1
			}
		}
	}

	fmt.Println(digraphFreq)
	d, _ := sortDigraphs(digraphFreq)
	return d[0].Digraph
	//fmt.Println(digraphFreq)
}

func guessKeyLength(ciphertext string) int {
	findMostUsedDigraph(ciphertext)

	for k := 1; k <= 13; k++ {
	}
	return 7 // found with Shani
}

func breakCipherToKeyChunk(ciphertext string, keylength int) []string {
	cipherblocks := make([]string, int(math.Ceil(float64(len(ciphertext))/float64(keylength))))
	fmt.Println("cipherblocks size: ", len(cipherblocks), " ", len(ciphertext))
	for i, idx := 0, 0; (idx + keylength) < len(ciphertext); idx, i = idx+keylength, i+1 {
		fmt.Println(i)
		cipherblocks[i] = ciphertext[i*keylength : i*keylength+keylength]
	}
	c := make([]string, len(cipherblocks))
	for i, v := range cipherblocks {
		c[i] = hex.EncodeToString([]byte(v))
	}

	fmt.Print(c)
	return cipherblocks
}

type Letter struct {
	Letter byte
	Count  int
}

type LetterFreqList []Letter

func (lfl LetterFreqList) Swap(i, j int)      { lfl[i], lfl[j] = lfl[j], lfl[i] }
func (lfl LetterFreqList) Len() int           { return len(lfl) }
func (lfl LetterFreqList) Less(i, j int) bool { return lfl[i].Count < lfl[j].Count }

type PositionFreq []LetterFreqList

func analyzeCipherblocks(cipherblocks []string) PositionFreq {
	kl := len(cipherblocks[0])

	analysis := make([]map[uint8]int, kl)
	for i := 0; i < kl; i++ {
		analysis[i] = make(map[uint8]int)
		for j := 0; j < len(cipherblocks)-1; j++ {
			analysis[i][cipherblocks[j][i]] += 1
		}
	}
	orderedAnalysis := make(PositionFreq, kl)
	for i, v := range analysis {
		orderedAnalysis[i] = make(LetterFreqList, len(v))
		j := 0
		for k, w := range v {
			orderedAnalysis[i][j] = Letter{k, w}
			j++
		}
		sort.Sort(sort.Reverse(orderedAnalysis[i]))
		orderedAnalysis[i] = orderedAnalysis[i][:3]
	}
	fmt.Println(orderedAnalysis)
	return orderedAnalysis
}

func guessKey(cipherblocks []string, freqAnalysis PositionFreq, lutsearch []int) []byte {
	kl := len(cipherblocks[0])
	guessedKey := make([]byte, kl)
	for i, v := range freqAnalysis {
		guessedKey[i] = v[lutsearch[i]].Letter ^ byte('e')
	}
	fmt.Printf("%x\n", guessedKey)
	return guessedKey
}

func hex2ascii(ciphertext string) string {
	cipherbyte, err := hex.DecodeString(ciphertext)

	if err != nil {
		log.Println("Couldn't decode hex ciphertext")
	}

	return string(cipherbyte)
}

func decrypt(cipherblocks []string, key []byte) {
	var plaintext string

	for _, v := range cipherblocks {
		for j, w := range v {
			decchar := byte(w) ^ key[j]
			if decchar < 32 || decchar > 127 {
				return
			}
			plaintext += string(decchar)
			//fmt.Printf("%c", byte(w)^key[j])
		}
	}
	fmt.Printf("%x\n", key)
	fmt.Printf("%s\n", plaintext)
}

func nextIndex(lutsearch []int, choices int) {
	for i := len(lutsearch) - 1; i >= 0; i-- {
		lutsearch[i]++
		if i == 0 || lutsearch[i] < choices {
			return
		}
		lutsearch[i] = 0
	}
}

func main() {
	ciphertext := os.Args[1]
	ciphertext = hex2ascii(ciphertext)

	fmt.Println(ciphertext)
	keylength := guessKeyLength(ciphertext)
	cipherblocks := breakCipherToKeyChunk(ciphertext, keylength)

	//	freqAnalysis := analyzeCipherblocks(cipherblocks)

	//	for i := 0; i < 3; i++ {
	for i := 0; i < 255; i++ {
		for j := 0; j < 255; j++ {
			key := []byte{byte(i), 0x1F, 0x95, byte(j), 0x53, 0x88, 0x3e}

			decrypt(cipherblocks, key)
		}
	}

	/*lutsearch := make([]int, keylength)
	for ; lutsearch[0] < 3; nextIndex(lutsearch, 3) {
		key := guessKey(cipherblocks, freqAnalysis, lutsearch)
		decrypt(cipherblocks, key)
	}*/

	//	}
}
