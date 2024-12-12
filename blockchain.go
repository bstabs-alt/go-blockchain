package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func calcHash(b Block) string {
	r := fmt.Sprint(b.Index) + b.Timestamp + fmt.Sprint(b.Coins) + b.PrevHash
	h := sha256.New()
	h.Write([]byte(r))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func genBlock(prev Block, coins int) (Block, error) {
	t := time.Now().UTC()

	newBlock := Block{
		Index:     prev.Index + 1,
		Timestamp: t.String(),
		Coins:     coins,
		PrevHash:  prev.Hash,
		Hash:      calcHash(prev),
	}

	return newBlock, nil
}

func isValidBlock(prev Block, curr Block) error {
	switch {
	case prev.Index != curr.Index-1:
		return fmt.Errorf("Incorrect Index increment")
	case prev.Hash != curr.PrevHash:
		return fmt.Errorf("Previous Hash does not match")
	case calcHash(prev) != curr.Hash:
		return fmt.Errorf("Current Hash does not match")
	default:
		return nil
	}
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}
