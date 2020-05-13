package main

import (
	"fmt"
	"os"
)

type point struct {
	i, j int
}

func (p point) add(step point) point {
	p.i += step.i
	p.j += step.j
	return p
}

func (p point) at(maze [][]int) (int, bool) {
	if p.i < 0 || p.i >= len(maze) {
		return 0, false
	}
	if p.j < 0 || p.j >= len(maze[p.i]) {
		return 0, false
	}
	return maze[p.i][p.j], true
}
// 广度优先算法走迷宫

func ReadMaze(path string) [][]int {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer func(){
		file.Close()
	}()
	var row, col int
	_, err = fmt.Fscanf(file, "%d %d", &row, &col)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	maze := make([][]int, row)
	for i := range maze {
		maze[i] = make([]int, col)
		for j := range maze[i] {
			_, err = fmt.Fscanf(file, "%d", &maze[i][j])
			if err != nil {
				fmt.Println(err)
				os.Exit(4)
			}
		}
	}
	return maze
}
// 每个位置都可以在上下左右方向探索
var direction = [4]point{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}

func walk(maze [][]int, start, end point) [][]int {
	// 首先 创建一张 跟迷宫一样的地图 记录走的步长和路径
	steps := make([][]int, len(maze))
	for i := range steps {
		steps[i] = make([]int, len(maze[i]))
	}

	Q := []point{start}

	for len(Q) > 0 {
		current_location := Q[0]
		Q = Q[1:]
		// 到达终点
		if current_location == end {
			break
		}
		for _, val := range direction {
			// 1. 当前位置在上下左右探索
			// 2. 如果碰到 1 或者出界了 就不往下走
			next := current_location.add(val)
			// 判断 next 是否走出边界或有效
			location, bl := next.at(maze)
			if !bl || location == 1 {
				continue
			}
			// 判断这一步是否走过了
			location, bl = next.at(steps)
			if !bl || location != 0 {
				continue
			}
			// 如果又回到了起点
			if next == start {
				continue
			}
			// 走到这说明这一步是有效的 获取当前位置的步数
			curSteps, _ := current_location.at(steps)
			// 获取下一步的位置 并让下一步数加1
			steps[next.i][next.j] = curSteps + 1
			// 下一步添加进队列里 用于下一次的循环探索
			Q = append(Q, next)
		}
	}
	return steps
}

func main(){
	dir, _ := os.Getwd()
	fmt.Println("当前路径 : ", dir)
	// 首先需要读取出来迷宫值
	maze := ReadMaze("src/go_dev/test/maze.in")
	steps := walk(maze, point{0,0}, point{len(maze)-1, len(maze[0])-1})
	// 开始 走 迷宫
	for i := range steps {
		for j := range steps[i] {
			fmt.Printf("%3d", steps[i][j])
		}
		fmt.Println()
	}
}
