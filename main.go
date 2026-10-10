package main

import (
	"fmt"
	"github.com/vincentwc/dc/internal/docker"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法：dc <命令>")
		fmt.Println("  df           显示 docker 磁盘占用")
		fmt.Println("  stats        并发统计镜像/容器/数据卷数量")
		fmt.Println("  clean-cache  清理 build cache")
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "df":
		err = docker.ShowDF()

	case "clean-cache":
		// 危险操作，必须确认
		fmt.Print("⚠️  确认清理 Build Cache？此操作不可恢复 [y/N]: ")
		var ans string
		fmt.Scanln(&ans)
		if ans != "y" && ans != "Y" {
			fmt.Println("已取消")
			return
		}
		err = docker.CleanBuildCache()

	case "stats":
		err = docker.ParallelStats()

	default:
		fmt.Println("未知命令：", os.Args[1])
		os.Exit(1)
	}

	// main 是最上层，统一在这里处理错误
	if err != nil {
		fmt.Println("❌ 出错了：", err)
		os.Exit(1)
	}
}
