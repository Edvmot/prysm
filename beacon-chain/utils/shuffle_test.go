package utils

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prysmaticlabs/prysm/shared/params"
)

func TestFaultyShuffleIndices(t *testing.T) {
	var list []uint32

	upperBound := 1<<(params.BeaconConfig().RandBytes*8) - 1

	for i := 0; i < upperBound+1; i++ {
		list = append(list, uint32(i))
	}

	if _, err := ShuffleIndices(common.Hash{'a'}, list); err == nil {
		t.Error("Shuffle should have failed when validator count exceeds ModuloBias")
	}
}

func TestShuffleIndices(t *testing.T) {
	hash1 := common.BytesToHash([]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'a', 'b', 'c', 'd', 'e', 'f', 'g'})
	hash2 := common.BytesToHash([]byte{'1', '2', '3', '4', '5', '6', '7', '1', '2', '3', '4', '5', '6', '7', '1', '2', '3', '4', '5', '6', '7', '1', '2', '3', '4', '5', '6', '7', '1', '2', '3', '4', '5', '6', '7'})
	var list1 []uint32

	for i := 0; i < 10; i++ {
		list1 = append(list1, uint32(i))
	}

	list2 := make([]uint32, len(list1))
	copy(list2, list1)

	list1, err := ShuffleIndices(hash1, list1)
	if err != nil {
		t.Errorf("Shuffle failed with: %v", err)
	}

	list2, err = ShuffleIndices(hash2, list2)
	if err != nil {
		t.Errorf("Shuffle failed with: %v", err)
	}

	if reflect.DeepEqual(list1, list2) {
		t.Errorf("2 shuffled lists shouldn't be equal")
	}
	if !reflect.DeepEqual(list1, []uint32{5, 4, 9, 6, 7, 3, 0, 1, 8, 2}) {
		t.Errorf("list 1 was incorrectly shuffled")
	}
	if !reflect.DeepEqual(list2, []uint32{9, 0, 1, 5, 3, 2, 4, 7, 8, 6}) {
		t.Errorf("list 2 was incorrectly shuffled")
	}
}

func TestSplitIndices(t *testing.T) {
	var l []uint32
	validators := 64000
	for i := 0; i < validators; i++ {
		l = append(l, uint32(i))
	}
	split := SplitIndices(l, params.BeaconConfig().EpochLength)
	if len(split) != int(params.BeaconConfig().EpochLength) {
		t.Errorf("Split list failed due to incorrect length, wanted:%v, got:%v", params.BeaconConfig().EpochLength, len(split))
	}

	for _, s := range split {
		if len(s) != validators/int(params.BeaconConfig().EpochLength) {
			t.Errorf("Split list failed due to incorrect length, wanted:%v, got:%v", validators/int(params.BeaconConfig().EpochLength), len(s))
		}
	}
}
