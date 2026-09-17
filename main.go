package main

import "fmt"
func twoSum(nums []int, target int) []int {
    m := make(map[int]int)
    for i := 0; i < len(nums); i++ {
        var dif int = target - nums[i];
        if j, ok := m[dif]; ok {
            return []int{j, i}
        }
        m[nums[i]] = i
    }
    return []int{}
}

func main() {
    nums := []int{1, 2, 3, 5}

	result := twoSum(nums, 5)

	fmt.Println(result)
}