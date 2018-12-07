// Package utils defines utility functions for the beacon-chain.
package utils

import (
	"errors"
	"math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/shared/hashutil"
	"github.com/prysmaticlabs/prysm/shared/params"
)

// ShuffleIndices returns a list of pseudorandomly sampled
// indices. This is used to shuffle validators on ETH2.0 beacon chain.
func ShuffleIndices(seed common.Hash, indicesList []uint32) ([]uint32, error) {
	// Each entropy is consumed from the seed in randBytes chunks.
	randBytes := int(params.BeaconConfig().RandBytes)
	upperBound := uint64(math.Pow(2, float64(randBytes) * 8)) - 1

	// Since we are consuming randBytes of entropy at a time in the loop,
	// we have a bias at 2**24, this check defines our max list size and is used to remove the bias.
	// more info on modulo bias: https://stackoverflow.com/questions/10984974/why-do-people-say-there-is-modulo-bias-when-using-a-random-number-generator.
	if uint64(len(indicesList)) >= upperBound {
		return nil, errors.New("input list exceeded upper bound and reached modulo bias")
	}

	// Rehash the seed to obtain a new pattern of bytes.
	hashSeed := hashutil.Hash(seed[:])
	totalCount := len(indicesList)
	index := 0
	for index < totalCount-1 {
		// Iterate through the hashSeed bytes in chunks of size randBytes.
		for i := 0; i < 32-(32%randBytes); i += randBytes {
			// Determine the number of indices remaining and exit if last index reached.
			remaining := uint64(totalCount - index)
			if remaining == 1 {
				break
			}
			// Read randBytes of hashSeed as a 3 x randBytes big-endian integer.
			randChunk := hashSeed[i : i+randBytes]
			var randValue uint64
			randValue |= uint64(randChunk[0])
			randValue |= uint64(randChunk[1])
			randValue |= uint64(randChunk[2])

			// Sample values greater than or equal to sampleMax will cause
			// modulo bias when mapped into the remaining range.
			randMax := upperBound - upperBound%remaining

			// Perform swap if the consumed entropy will not cause modulo bias.
			if randValue < randMax {
				// Select replacement index from the current index.
				replacementIndex := (randValue % remaining) + uint64(index)
				indicesList[index], indicesList[replacementIndex] = indicesList[replacementIndex], indicesList[index]
				index++
			}
		}
	}
	return indicesList, nil
}

// SplitIndices splits a list into n pieces.
func SplitIndices(l []uint32, n uint64) [][]uint32 {
	var divided [][]uint32
	var lSize = uint64(len(l))
	for i := uint64(0); i < n; i++ {
		start := lSize * i / n
		end := lSize * (i + 1) / n
		divided = append(divided, l[start:end])
	}
	return divided
}
