package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/saintfish/chardet"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var charsetMap = make(map[string]encoding.Encoding)

var detector *chardet.Detector

func defaultRootPath() (root string, err error) {
	root, err = os.Getwd()
	return
}

var unifyCmd = &cobra.Command{
	Use:   "unify",
	Short: "unify 用于探测项目中的指定目标编码的要转换的文件",
	Long:  `unify 用于探测项目中的指定目标编码的要转换的文件`,
	Run:   unifyDir,
}

var (
	sourceCharset string
	targetCharset string
	rootPath      string
	fileSuffix    string
)

func init() {
	detector = chardet.NewTextDetector()
	charsetMap["GBK"] = simplifiedchinese.GBK
	charsetMap["GB-18030"] = simplifiedchinese.GB18030
	charsetMap["GB18030"] = simplifiedchinese.GB18030
	charsetMap["UTF8"] = unicode.UTF8
	charsetMap["UTF-8"] = unicode.UTF8
	unifyCmd.Flags().StringVarP(&sourceCharset, "source-charset", "s", "GB-18030",
		"source charset for matching file encoding when detecting")
	unifyCmd.Flags().StringVarP(&targetCharset, "target-charset", "t", "UTF8",
		"target charset for encoding file to")
	unifyCmd.Flags().StringVarP(&rootPath, "path", "p", "",
		"root path for unify file encoding")
	unifyCmd.Flags().StringVarP(&fileSuffix, "suffix", "x", "",
		"file suffix for matching file to detect")
	unifyCmd.MarkFlagRequired("suffix")
}

type transformContext struct {
	path   string
	source encoding.Encoding
	dest   encoding.Encoding
}

func compatUTF8(bs []byte) bool {
	for i := 0; i < len(bs); i++ {
		b := bs[i]
		// 单字节字符 0xxxxxx
		if b&0x80 == 0x0 {
			continue
		}
		// 双字节字符 110xxxxx 10xxxxxx
		if b&0xE0 == 0xC0 {
			if forwardByte(bs[i+1]) {
				i += 1
			} else {
				//fmt.Println(strconv.FormatInt(int64(bs[i+1]), 2))
				return false
			}
		}
		// 三字节字符 1110xxxx 10xxxxxx 10xxxxxx
		if b&0xF0 == 0xE0 {
			if forwardByte(bs[i+1]) && forwardByte(bs[i+2]) {
				i += 2
			} else {
				return false
			}
		}
		// 四字节字符 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
		if b&0xF8 == 0xF0 {
			if forwardByte(bs[i+1]) && forwardByte(bs[i+2]) && forwardByte(bs[i+3]) {
				i += 3
			} else {
				return false
			}
		}
	}
	return true
}

// 10xxxxxx
func forwardByte(b byte) bool {
	return b&0xC0 == 0x80
}

func unifyDir(cmd *cobra.Command, args []string) {
	root := rootPath
	if len(root) == 0 {
		r, err := defaultRootPath()
		if err != nil {
			panic(err)
		}
		root = r
	}

	log.Println("rootPath", root)
	log.Println("source charset", sourceCharset)
	log.Println("file suffix", fileSuffix)
	log.Println("target charset", targetCharset)

	err := filepath.Walk(root, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if matchExt(info, fileSuffix) {
			ctx, err := detectFile(p, sourceCharset, targetCharset)
			if err != nil {
				return err
			}
			if ctx != nil {
				transformFile(ctx)
			}
		}
		return err
	})

	if err != nil {
		panic(err)
	}
}

func detectFile(p string, charset string, destCharset string) (ctxt *transformContext, err error) {
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	dr, err := detector.DetectBest(b)
	if err != nil {
		return nil, err
	}
	if dr.Charset == charset {
		fmt.Printf("%s: %s\n", p, charset)
		return buildContext(p, sourceCharset, destCharset)
	} else {
		if !compatUTF8(b) && strings.HasPrefix(strings.ToUpper(destCharset), "UTF") {
			log.Println("unmappable character for UTF8 ", p)
			all, err := detector.DetectAll(b)
			if err != nil {
				return nil, err
			}
			for _, result := range all {
				if result.Charset == charset {
					return buildContext(p, sourceCharset, destCharset)
				}
			}
		}
	}
	return nil, nil
}

func matchExt(info fs.FileInfo, ext string) bool {
	return !info.IsDir() && strings.HasSuffix(info.Name(), ext)
}

func buildContext(p string, sourceCharset string, destCharset string) (ctxt *transformContext, err error) {
	ctx := &transformContext{
		path: p,
	}
	enc := charsetMap[sourceCharset]
	if enc == nil {
		return nil, errors.New("source charset not support")
	}
	ctx.source = enc
	destEnc := charsetMap[destCharset]
	if destEnc == nil {
		return nil, errors.New("target charset not support")
	}
	ctx.dest = destEnc
	return ctx, nil
}

func transformFile(ctx *transformContext) {
	file, err := os.OpenFile(ctx.path, os.O_RDONLY, os.ModePerm)
	destPath := ctx.path + ".tmp"
	destFile, err := os.Create(destPath)
	if err != nil {
		log.Fatal(err)
	}
	r := transform.NewReader(file, ctx.source.NewDecoder())
	scanner := bufio.NewScanner(r)
	w := transform.NewWriter(destFile, ctx.dest.NewEncoder())
	for scanner.Scan() {
		_, err := w.Write(scanner.Bytes())
		if err != nil {
			panic(err)
		}
		_, err = w.Write([]byte("\n"))
		if err != nil {
			panic(err)
		}
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}
	err = os.Remove(ctx.path)
	if err != nil {
		panic(err)
	}
	err = os.Rename(destPath, ctx.path)
	if err != nil {
		panic(err)
	}
	log.Println(ctx.path + " transformed")
}
