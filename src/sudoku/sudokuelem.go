package main

import (
	"fmt"
)

type SudokuElem struct {
	internal_ uint16
}

func (elem SudokuElem) String() string {
	return fmt.Sprint("%v", elem.GetValue())
}

func (elem *SudokuElem) SetValue(value int8) {
	var _Tp uint16
	_Tp = (uint16(value)) << 12
	elem.internal_ |= _Tp
}

func (elem SudokuElem) GetValue() int8 {
	return int8(elem.internal_ >> 12)
}

func (elem SudokuElem) CacheNum() int8 {
	var _Ret int8 = 0
	var i uint = 0
	for ; i < 9; i++ {
		_Ret += int8((elem.internal_ >> i) & 0x01)
	}
	return _Ret
}

func (elem *SudokuElem) PushCacheBack(num int8) {
	elem.internal_ |= (0x01 << uint(num-1))
}

func (elem *SudokuElem) RemoveFromCache(num int8) {
	var _Tp uint16 = (0x01 << uint(num-1))
	_Tp = ^_Tp
	elem.internal_ &= _Tp
}

func (elem *SudokuElem) PopCacheFront() int8 {
	var num uint8 = 1
	for ; num <= 9; num++ {
		if (uint16(0x01)<<(num-1))&elem.internal_ > 0 {
			elem.internal_ &= ^(0x01 << (num - 1))
			return int8(num)
		}
	}
	return int8(0)
}

func (elem *SudokuElem) RemoveAllCache() {
	var _Tp uint16 = 511
	elem.internal_ &= ^_Tp
}

func (elem *SudokuElem) PushAllToCache() {
	var _Tp uint16 = 511
	elem.internal_ |= _Tp
}
