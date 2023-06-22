package challenge

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"github.com/Hermes-Bird/faraway-test-task.git/config"
	"github.com/Hermes-Bird/faraway-test-task.git/internal/hash"
)

func GetChallenge() []byte {
	randCrypto, _ := rand.Prime(rand.Reader, 128)
	return randCrypto.Bytes()
}

func CheckChallenge(challenge, nonce []byte) bool {
	result := append(challenge, nonce...)
	sum := hash.HashBytes(result)
	return bytes.HasPrefix(sum[:], config.ChallengeCondition)
}

func SolveChallenge(bs []byte) []byte {
	var nonce uint32 = 1
	for !bytes.HasPrefix(hash.HashBytes(binary.BigEndian.AppendUint32(bs, nonce)), config.ChallengeCondition) {
		nonce += 1
	}

	res := make([]byte, 4)
	binary.BigEndian.PutUint32(res, nonce)

	return res
}
