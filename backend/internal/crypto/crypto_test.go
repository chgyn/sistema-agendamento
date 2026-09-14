package crypto

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	secretMessage := "sk-ant-api03-secret-openai-or-gemini-key-123456"

	encrypted, err := Encrypt(secretMessage, key)
	if err != nil {
		t.Fatalf("falha ao criptografar: %v", err)
	}

	if encrypted == secretMessage {
		t.Fatalf("o texto criptografado não pode ser idêntico ao original")
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("falha ao descriptografar: %v", err)
	}

	if decrypted != secretMessage {
		t.Errorf("esperava '%s', mas obteve '%s'", secretMessage, decrypted)
	}
}

func TestDecryptWithInvalidKey(t *testing.T) {
	key1 := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	key2 := "fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210"
	secretMessage := "chave-secreta"

	encrypted, err := Encrypt(secretMessage, key1)
	if err != nil {
		t.Fatalf("falha ao criptografar: %v", err)
	}

	_, err = Decrypt(encrypted, key2)
	if err == nil {
		t.Fatalf("esperava erro de autenticação ao descriptografar com chave errada, mas obteve sucesso")
	}
}

func TestEmptyValues(t *testing.T) {
	key := "minha-chave-segura-123"

	enc, err := Encrypt("", key)
	if err != nil {
		t.Fatalf("erro inesperado com string vazia: %v", err)
	}
	if enc != "" {
		t.Errorf("esperava string vazia, obteve %s", enc)
	}

	dec, err := Decrypt("", key)
	if err != nil {
		t.Fatalf("erro inesperado com string vazia: %v", err)
	}
	if dec != "" {
		t.Errorf("esperava string vazia, obteve %s", dec)
	}
}
