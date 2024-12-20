package core

import mt "github.com/minetest-go/types"

// +/- 2000 mapblocks
const key_offset = 2000

func MapKey(mbpos *mt.Pos) int64 {
	return int64(mbpos[0]+key_offset) +
		int64(mbpos[1]+key_offset)<<16 +
		int64(mbpos[2]+key_offset)<<32
}

func ParseMapKey(i int64) *mt.Pos {
	x := int((i>>0)&0xFFFF) - key_offset
	y := int((i>>16)&0xFFFF) - key_offset
	z := int((i>>32)&0xFFFF) - key_offset
	return mt.NewPos(x, y, z)
}
