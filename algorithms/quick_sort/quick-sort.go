package main

import "fmt"

func main(){
	arr := [...]int{5,2,7,2,4,8,1,3}
	quickSort(arr[:], 0, len(arr)-1)
	fmt.Println(arr)
}

//思想:
// 1. 随机取一个数 找到这个数在列表中的位置(左边的都比它小，右边的都比它大)
// 2. 然后把左边、右边的重复 1 步骤
func quickSort(sc []int, begin, end int){
	if begin >= end {
		return
	}

	// 取第一个元素为比较元素
	val := sc[begin]
	k := begin

	for i := begin + 1; i <= end; i ++ {
		// 如果比选取的数小
		if sc[i] > val {
			// 这个小的数 放到 比较数的位置
			sc[k] = sc[i]
			// 把比较数 前面的数 放到空出来的位置里
			sc[i] = sc[k+1]
			// 比较数的位置往前移动一个(上一步已经把前面的数移走了)
			k++
		}
	}
	// 把比较数放入位置里
	sc[k] = val
	// 把左边都小的数再进行比较
	quickSort(sc, begin, k-1)
	// 把右边都大的数据再进行比较
	quickSort(sc, k+1, end)
}