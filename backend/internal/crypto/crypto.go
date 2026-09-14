package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

var (
	ErrInvalidCiphertext = errors.New("texto cifrado inválido ou corrompido")
	ErrEmptyKey          = errors.New("chave de criptografia não pode ser vazia")
)

// deriveKey assegura uma chave de 32 bytes (AES-256) a partir da string fornecida
func deriveKey(keyStr string) ([]byte, error) {
	if keyStr == "" {
		return nil, ErrEmptyKey
	}

	// Se a chave for hexadecimal de 64 caracteres (32 bytes), decodifica diretamente
	if len(keyStr) == 64 {
		if keyBytes, err := hex.DecodeString(keyStr); err == nil {
			return keyBytes, nil
		}
	}

	// Se a chave tiver exatamente 32 caracteres ASCII
	if len(keyStr) == 32 {
		return []byte(keyStr), nil
	}

	// Caso contrário, deriva deterministicamente usando SHA-256
	hash := sha256.Sum256([]byte(keyStr))
	return hash[:], nil
}

// Encrypt cifra um texto plano utilizando AES-256-GCM com nonce aleatório e retorna base64
func Encrypt(plainText, keyStr string) (string, error) {
	if plainText == "" {
		return "", nil
	}

	key, err := deriveKey(keyStr)
	if err != nil {
		return "", fmt.Errorf("falha ao derivar chave de criptografia: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("falha ao instanciar cifra AES: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("falha ao inicializar GCM: %w", err)
	}

	// Nonce de 12 bytes padrão para GCM
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("falha ao gerar nonce criptográfico: %w", err)
	}

	// O nonce é prefixado no texto cifrado
	cipherBytes := gcm.Seal(nonce, nonce, []byte(plainText), nil)

	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

// Decrypt decifra um texto cifrado em base64 utilizando AES-256-GCM
func Decrypt(cipherTextB64, keyStr string) (string, error) {
	if cipherTextB64 == "" {
		return "", nil
	}

	key, err := deriveKey(keyStr)
	if err != nil {
		return "", fmt.Errorf("falha ao derivar chave de criptografia: %w", err)
	}

	cipherBytes, err := base64.StdEncoding.DecodeString(cipherTextB64)
	if err != nil {
		return "", fmt.Errorf("base64 inválido: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("falha ao instanciar cifra AES: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("falha ao inicializar GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(cipherBytes) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	nonce, ciphertext := cipherBytes[:nonceSize], cipherBytes[nonceSize:]
	plainBytes, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("falha na autenticação/decifragem dos dados: %w", err)
	}

	return string(plainBytes), nil
}
