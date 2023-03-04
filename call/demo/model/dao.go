package model 

var Xzm = struct {
	Count     [17]uint32
	Pos       [3 * 17]uint32
	NumStmt   [17]uint16
} {
	Pos: [3 * 17]uint32{
		5, 6, 0xf001d, // [0]
		9, 9, 0xe0002, // [1]
		6, 8, 0x3000f, // [2]
		9, 11, 0x3000e, // [3]
		11, 11, 0x150008, // [4]
		11, 13, 0x30015, // [5]
		13, 15, 0x30008, // [6]
		17, 18, 0xf001d, // [7]
		24, 24, 0x1a0002, // [8]
		28, 28, 0x120002, // [9]
		18, 19, 0xf000f, // [10]
		22, 22, 0x130003, // [11]
		19, 21, 0x4000f, // [12]
		24, 26, 0x3001a, // [13]
		31, 33, 0x2001b, // [14]
		35, 37, 0x2001b, // [15]
		38, 40, 0x2001b, // [16]
	},
	NumStmt: [17]uint16{
		1, // 0
		1, // 1
		1, // 2
		1, // 3
		1, // 4
		1, // 5
		1, // 6
		1, // 7
		1, // 8
		1, // 9
		1, // 10
		1, // 11
		1, // 12
		1, // 13
		1, // 14
		1, // 15
		1, // 16
	},
}
