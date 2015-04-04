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
		c[i] = hex.EncodeToString(v)
	}

	fmt.Print(c)
	return cipherblocks
}

type Letter struct {
	Letter byte
	Count int
}

type LetterFreqList []Letter
type PositionFreq []LetterFreqList

func analyzeLetterFrequencies(cipherblocks []string) {
	kl := len(cipherblocks[0])
	posfreq := make(PositionFreq, kl)
	for i := 0; i < kl; i++ {
		mmap := make(map[string]int)
		for _, v := range cipherblocks {
			mmap[v[i]]++
		}
		l := make(DigraphPairList, len(mmap))
		for k, v := range mmap {
			l[k] = v

	}

	for i, v := range cipherblocks {
		posfreq[i] := 

func guessKey(ciphertext string, keylength int) {
	cipherblocks := breakCipherToKeyChunk(ciphertext, keylength)
	
	analyzeLetterFrequencies(cipherblocks)

}

func hex2ascii(ciphertext string) string {
	cipherbyte, err := hex.DecodeString(ciphertext)

	if err != nil {
		log.Println("Couldn't decode hex ciphertext")
	}

	return string(cipherbyte)
}

func main() {
	ciphertext := os.Args[1]
	ciphertext = hex2ascii(ciphertext)

	fmt.Println(ciphertext)
	keylength := guessKeyLength(ciphertext)
	guessKey(ciphertext, keylength)
}
