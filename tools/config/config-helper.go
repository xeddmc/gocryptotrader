package main

import (
	"flag"
	"log"

	"github.com/thrasher-/gocryptotrader/common"
	"github.com/thrasher-/gocryptotrader/config"
)

func EncryptOrDecrypt(encrypt bool) string {
	if encrypt {
		return "encrypted"
	}
	return "decrypted"
}

func main() {
	var inFile, outFile, key string
	var encrypt bool
	var err error
	flag.StringVar(&inFile, "infile", "config.dat", "The config input file to process.")
	flag.StringVar(&outFile, "outfile", "config.dat.out", "The config output file.")
	flag.BoolVar(&encrypt, "encrypt", true, "Wether to encrypt or decrypt.")
	flag.StringVar(&key, "key", "", "The key to use for AES encryption.")
	flag.Parse()

	log.Println("GoCryptoTrader: config-helper tool.")

	if key == "" {
		result, err := config.PromptForConfigKey()
		if err != nil {
			log.Fatal("Unable to obtain encryption/decryption key.")
		}
		key = string(result)
	}

	file, err := common.ReadFile(inFile)
	if err != nil {
		log.Fatalf("Unable to read input file %s. Error: %s.", inFile, err)
	}

	if config.ConfirmECS(file) && encrypt {
		log.Println("File is already encrypted. Decrypting..")
		encrypt = false
	}

	if !config.ConfirmECS(file) && !encrypt {
		var result interface{}
		err := config.ConfirmConfigJSON(file, result)
		if err != nil {
			log.Fatal("File isn't in JSON format")
		}
		log.Println("File is already decrypted. Encrypting..")
		encrypt = true
	}

	var data []byte
	if encrypt {
		data, err = config.EncryptConfigFile(file, []byte(key))
		if err != nil {
			log.Fatalf("Unable to encrypt config data. Error: %s.", err)
		}
	} else {
		data, err = config.DecryptConfigFile(file, []byte(key))
		if err != nil {
			log.Fatalf("Unable to decrypt config data. Error: %s.", err)
		}
	}

	err = common.WriteFile(outFile, data)
	if err != nil {
		log.Fatalf("Unable to write output file %s. Error: %s", outFile, err)
	}
	log.Printf("Successfully %s input file %s and wrote output to %s.\n", EncryptOrDecrypt(encrypt), inFile, outFile)
}
