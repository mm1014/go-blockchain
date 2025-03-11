package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyPairSignVerifySuccess(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()
	msg := []byte("hello world")
	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)
	assert.True(t, sig.Verify(pubKey, msg))
}

func TestKeyPairSignVerifyFail(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()
	msg := []byte("hello world")
	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)
	otherPrivKey := GeneratePrivateKey()
	otherPublicKey := otherPrivKey.PublicKey()
	assert.False(t, sig.Verify(otherPublicKey, msg))
	assert.False(t, sig.Verify(pubKey, []byte("xxxxx")))
}
