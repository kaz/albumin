package dedup

import (
	"encoding/binary"
	"math/bits"

	"github.com/kaz/albumin/model"
)

func GroupByHash(photos []*model.Photo) [][]*model.Photo {
	groups := make(map[string][]*model.Photo, len(photos))
	for _, photo := range photos {
		sHash := string(photo.Hash)
		if _, ok := groups[sHash]; !ok {
			groups[sHash] = []*model.Photo{}
		}

		groups[sHash] = append(groups[sHash], photo)
	}

	result := [][]*model.Photo{}
	for _, group := range groups {
		if len(group) < 2 {
			continue
		}

		result = append(result, group)
	}

	return result
}

func GroupByPHash(photos []*model.Photo, tolerance int) [][]*model.Photo {
	groups := make([][]*model.Photo, len(photos))
	phash := make([]uint64, len(photos))
	for i := 0; i < len(photos); i++ {
		groups[i] = []*model.Photo{photos[i]}
		phash[i] = binary.LittleEndian.Uint64(photos[i].PHash)
	}

	for i := 0; i < len(photos); i++ {
		for j := i + 1; j < len(photos); j++ {
			if bits.OnesCount64(phash[i]^phash[j]) <= tolerance {
				groups[i] = append(groups[i], photos[j])
			}
		}
	}

	result := [][]*model.Photo{}
	for _, group := range groups {
		if len(group) < 2 {
			continue
		}

		result = append(result, group)
	}

	return result
}
