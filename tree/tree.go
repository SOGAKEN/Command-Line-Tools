package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"sort"
)

type Counter struct {
	dirs  int
	files int
}

func (counter *Counter) index(path string) {
	stat, _ := os.Stat(path)
	if stat.IsDir() {
		counter.dirs += 1
	} else {
		counter.files += 1
	}
}

func (counter *Counter) output() string {
	return fmt.Sprintf("\n%d directories, %d files", counter.dirs, counter.files)
}

func dirnamesFrom(base string) []string {
	file, err := os.Open(base)
	if err != nil {
		fmt.Println(err)
	}
	names, _ := file.Readdirnames(0)
	file.Close()

	sort.Strings(names)
	return names
}

// tree関数に現在の深さ(currentDepth)と最大深さ(maxDepth)を追加
func tree(counter *Counter, base string, prefix string, currentDepth int, maxDepth int) {
	names := dirnamesFrom(base)

	for index, name := range names {
		// 隠しファイル・ディレクトリはスキップ
		if name[0] == '.' {
			continue
		}
		subpath := path.Join(base, name)
		counter.index(subpath)

		// ツリーの記号で出力
		if index == len(names)-1 {
			fmt.Println(prefix + "└── " + name)
		} else {
			fmt.Println(prefix + "├── " + name)
		}

		// maxDepthが-1の場合は無制限、または現在の深さが maxDepth-1 より浅い場合のみ再帰する
		if maxDepth < 0 || currentDepth < maxDepth-1 {
			var newPrefix string
			if index == len(names)-1 {
				newPrefix = prefix + "    "
			} else {
				newPrefix = prefix + "│   "
			}
			tree(counter, subpath, newPrefix, currentDepth+1, maxDepth)
		}
	}
}

func main() {
	// -L オプションで最大階層を指定（デフォルトは -1: 無制限）
	var maxDepth int
	flag.IntVar(&maxDepth, "L", -1, "maximum depth level (-1 for unlimited)")
	flag.Parse()

	// 残りの引数から対象ディレクトリを取得（なければカレントディレクトリ）
	var directory string
	if flag.NArg() > 0 {
		directory = flag.Arg(0)
	} else {
		directory = "."
	}

	counter := new(Counter)
	fmt.Println(directory)

	// ルートを深さ0として tree を呼び出す
	tree(counter, directory, "", 0, maxDepth)
	fmt.Println(counter.output())
}
