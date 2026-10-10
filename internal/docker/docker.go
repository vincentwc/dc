package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// DiskUsage 对应docker system df 每一行json
type DiskUsage struct {
	Type        string `json:"Type"`
	TotalCount  string `json:"TotalCount"`
	Active      string `json:"Active"`
	Size        string `json:"Size"`
	Reclaimable string `json:"Reclaimable"`
}

// String 让 DiskUsage 自带"怎么打印自己"的能力
// 注意：这不是注解，不是 override，就是定义了一个叫 String 的方法
func (d DiskUsage) String() string {
	return fmt.Sprintf("%-15s %-10s %-10s %-12s %s",
		d.Type, d.TotalCount, d.Active, d.Size, d.Reclaimable)
}

// RunDocker 封装：执行 docker 命令，返回输出字符串和 error
// 这是一个通用的"跑 docker 子命令"的工具函数
func RunDocker(args ...string) (string, error) {
	cmd := exec.Command("docker", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("docker %s 执行失败：%w", strings.Join(args, " "), err)
	}
	return string(out), nil
}

// parseDiskUsage 把一行 JSON 解析成 DiskUsage
func parseDiskUsage(line string) (DiskUsage, error) {
	var d DiskUsage
	if err := json.Unmarshal([]byte(line), &d); err != nil {
		return DiskUsage{}, fmt.Errorf("解析行 [%s] 失败: %w", line, err)
	}
	return d, nil
}

// ShowDF 返回 error 而不是自己打印错误，把处理权交给 main
func ShowDF() error {
	start := time.Now()
	defer func() {
		fmt.Printf("{查询耗时 %v}\n", time.Since(start))
	}()

	out, err := RunDocker("system", "df", "--format", "{{json .}}")
	if err != nil {
		return err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	fmt.Printf("%-15s %-8s %-8s %-12s %s\n", "类型", "总数", "活跃", "占用", "可回收")
	fmt.Println(strings.Repeat("-", 55))

	for _, line := range lines {
		d, err := parseDiskUsage(line)
		if err != nil {
			return err
		}
		fmt.Println(d)
	}
	return nil
}

// CleanBuildCache 清理 build cache
func CleanBuildCache() error {
	cmd := exec.Command("docker", "builder", "prune", "-f")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("builder prune 失败: %w", err)
	}
	fmt.Println("✅ Build Cache 清理完成")
	return nil
}

// ParallelStats 并发查询4项 docker 统计
func ParallelStats() error {
	type query struct {
		name string
		args []string
	}

	queries := []query{
		{"镜像总数", []string{"images", "-q"}},
		{"运行中容器", []string{"ps", "-q"}},
		{"数据卷总数", []string{"volume", "ls", "-q"}},
	}

	// 定义 channel: 用来收 goroutine 的结果
	type result struct {
		name  string
		count int
		err   error
	}
	ch := make(chan result)

	// 启动3个 goroutine 每个跑一条 docker 命令
	for _, q := range queries {
		// q 是循环变量，在 goroutine 里用要先拷贝一份儿
		q := q
		go func() {
			out, err := RunDocker(q.args...)
			if err != nil {
				ch <- result{name: q.name, err: err}
				return
			}
			// 数行数 = 数量
			n := 0
			if s := strings.TrimSpace(out); s != "" {
				n = len(strings.Split(s, "\n"))
			}

			ch <- result{name: q.name, count: n}
		}()
	}

	// 等 3 个 goroutine 都返回结果
	fmt.Println("正在并发查询...")
	for i := 0; i < len(queries); i++ {
		r := <-ch
		if r.err != nil {
			return r.err
		}
		fmt.Printf("  %-12s %d\n", r.name, r.count)
	}
	return nil
}
