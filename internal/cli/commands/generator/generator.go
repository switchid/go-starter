package generator

import (
	"GoStarter/internal/pkg/config"
	"GoStarter/pkg/utils/crypts/encoding"
	"crypto/rand"
	"log"
)

func GenerateApplicationKey() {
	cfg, errCfg := config.Load()
	if errCfg != nil {
		log.Fatalf("Failed to load configuration: %v", errCfg)
	}
	key, errKey := generateKey()
	if errKey != nil {
		log.Fatalf("Failed to generate key: %v", errKey)
	}
	cfg.SetAppKey(key)
	errSaveCfg := cfg.SaveConfig()
	if errSaveCfg != nil {
		log.Fatalf("Failed to save configuration of application key: %v", errSaveCfg)
	}

	log.Println("Generated application key successfully")
	log.Printf("Application key : %s \n", key)
}

func generateKey() (string, error) {
	randBytes := make([]byte, 32)
	_, errRandStr := rand.Read(randBytes)
	if errRandStr != nil {
		return "", errRandStr
	}

	strEncode := encoding.NewEncodingBASE32(string(randBytes))

	return strEncode.EncodingText()
}
