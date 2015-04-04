package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

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

func displayDecryptedMessage(ciphertext string, key []byte) {
	cipherblocks := breakCipherToKeyChunk(ciphertext, len(key))

	for i, v := range key {
		for j, b := range cipherblocks {
			if v == 0 {
				b[i] = '*'
			} else {
				b[i] = string(int(b[i]) ^ int(v))
			}
		}
	}
	fmt.Print(cipherblocks)
}

func main() {
	//ciphertext := os.Args[1]
	keylength := flag.Int("keylength", 0, "The estimated length of key")

	flag.Parse()

	reader := bufio.NewReader(os.Stdin)
	key := make([]byte, *keylength)

	for i := 0; i < *keylength; i++ {
		fmt.Println("Write your estimate for", i, "th position of the key (in Hex) :")
		val, _ := reader.ReadString('\n')
		est, err := hex.DecodeString(val[:2])
		fmt.Println("estimates is", len(est), "long")
		if err != nil {
			fmt.Fprint(os.Stderr, "Couldn't decode, is it hex ?\n")
			i--
		} else {
			fmt.Println("Taking", int(est[0]), "as", i, "th value for the key")
			key[i] = est[0]
		}
	}

	fmt.Println("Key is ", key)
	displayDecryptedMessage(ciphertext, key)

}
