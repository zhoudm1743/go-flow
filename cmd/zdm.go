package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	// 版本信息
	version = "1.0.0"
)

var rootCmd = &cobra.Command{
	Use:   "zdm",
	Short: "Go-Flow命令行工具",
	Long:  `Go-Flow框架命令行工具，用于生成和管理项目结构`,
}

// versionCmd 版本信息命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看版本信息",
	Long:  "显示当前版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Go-Flow CLI v" + version)
		fmt.Println("作者: zhoudm1743")
	},
}

// cleanCmd 清理命令
var cleanCmd = &cobra.Command{
	Use:   "clean [module]",
	Short: "清理模块中自动生成的文件",
	Long:  "删除指定模块中的自动生成文件，保留基础结构",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleName := args[0]
		cleanModule(moduleName)
	},
}

// 清理模块中自动生成的文件
func cleanModule(moduleName string) {
	moduleDir := filepath.Join("internal", moduleName)

	// 检查模块是否存在
	if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
		fmt.Printf("错误: 模块 %s 不存在\n", moduleName)
		return
	}

	// 要清理的文件模式
	cleanPatterns := []string{
		filepath.Join(moduleDir, "model", "*.go"),
		filepath.Join(moduleDir, "repository", "*.go"),
		filepath.Join(moduleDir, "service", "*.go"),
		filepath.Join(moduleDir, "controller", "*.go"),
	}

	// 保留的基础文件
	preserveFiles := []string{
		filepath.Join(moduleDir, "module.go"),
		filepath.Join(moduleDir, "schemas", "req", "base.go"),
		filepath.Join(moduleDir, "schemas", "resp", "base.go"),
	}

	// 清理匹配的文件
	for _, pattern := range cleanPatterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Printf("查找文件失败: %s, 错误: %v\n", pattern, err)
			continue
		}

		for _, match := range matches {
			// 检查是否是要保留的文件
			isPreserved := false
			for _, preserve := range preserveFiles {
				if match == preserve {
					isPreserved = true
					break
				}
			}

			if !isPreserved {
				err := os.Remove(match)
				if err != nil {
					fmt.Printf("删除文件失败: %s, 错误: %v\n", match, err)
				} else {
					fmt.Printf("已删除: %s\n", match)
				}
			}
		}
	}

	fmt.Printf("模块 %s 已清理完成\n", moduleName)
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// 添加命令
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(cleanCmd)
}
