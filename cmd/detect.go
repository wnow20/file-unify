package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "detect 检测文件字符集",
	Long:  `detect 检测文件字符集`,
	Run:   detectFileForCmd,
}

var filePath string

func init() {
	rootCmd.AddCommand(detectCmd)
	detectCmd.Flags().StringVarP(&filePath, "file", "f", "",
		"file to detect charset of file encoding")
	detectCmd.MarkFlagRequired("file")
}

func detectFileForCmd(cmd *cobra.Command, args []string) {
	if !path.IsAbs(filePath) {
		filePath = path.Join(defaultRootPath(), filePath)
	}
	fmt.Println("detecting ", filePath)
	_, err := os.Lstat(filePath)
	if err != nil {
		panic(err)
	}
	f, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	dr, err := detector.DetectAll(f)
	if err != nil {
		panic(err)
	}

	fmt.Println("Possible charset")
	for _, result := range dr {
		fmt.Println(result.Charset)
	}
}
