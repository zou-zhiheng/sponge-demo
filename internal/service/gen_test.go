package service

import (
	"fmt"
	"gorm.io/gen"
	"log"
	"testing"
	"user/internal/dao"
	"user/internal/model"
)

func TestGormGen(t *testing.T) {

	g := gen.NewGenerator(gen.Config{
		OutPath: "../dao/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(model.GetDB())

	g.ApplyBasic(model.User{})

	g.ApplyInterface(func(dao.Querier) {}, model.User{})

	g.Execute()
}

func TestQueryParams(t *testing.T) {
	log.Println("123")

}

func TestDSA(t *testing.T) {
	fmt.Println(demo([]int{1, 3, 5, 4, 1}))
}

func demo(nums []int) int {

	if len(nums) == 0 {
		return 0
	} else if len(nums) == 1 {
		return nums[0]
	} else if len(nums) == 2 {
		return max(nums[0], nums[1])
	}

	var helper func(start, end int) int
	helper = func(start, end int) int {

		if start == end {
			return nums[start]
		}

		dp := make([]int, len(nums))
		dp[start] = nums[start]
		dp[start+1] = max(nums[start], nums[start+1])
		for i := start + 2; i <= end; i++ {
			dp[i] = max(dp[i-1], dp[i-2]+nums[i])
		}

		return dp[end]
	}

	res1 := helper(0, len(nums)-2)
	res2 := helper(1, len(nums)-1)

	return max(res1, res2)
}

func abs(a int) int {
	if a >= 0 {
		return a
	}

	return -a
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
