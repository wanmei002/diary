package main

import (
	"fmt"
	"math"
)

func main(){
	arr := [...]int{1,2,2,3,4,5,6,7,8,9}
	index := binarySearch(arr[:], 5, 0, len(arr)-1)

	fmt.Println("get index : ",index)

}

// 二分查找
// 原理 :
// 1. 针对的是一个有序的数据集合(必须是有序), 查找思想有点类似分治思想
// 2. 每次通过跟区间的中间元素对比, 将待查找的区间缩小为之前的一半, 直到找到要查找的元素, 或者区间被缩小为0

func binarySearch(nums []int, num, begin, end int) int {

	if begin > end {
		return -1
	}

	// 先取中间值
	fmt.Println("aaa")
	// 取中间值
	middle := int(math.Floor(float64(begin+end)/2))

	if num < nums[middle] {
		return binarySearch(nums, num, begin, middle-1)
	} else if num > nums[middle] {
		return binarySearch(nums, num, middle+1, end)
	} else {
		return middle
	}

}
