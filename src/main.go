package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

func get_shell() string {
	if runtime.GOOS == "windows" {
		return "pwsh"
	} else {
		cmd, _ := exec.Command("echo $HISTFILE").Output()
		fmt.Println(string(cmd))

		shell := strings.Split(os.Getenv("SHELL"), "/")
		sh := shell[len(shell)-1]
		return sh
		// if strings.Contains(shell, "zsh") {
		// 	return map[string]string{"shell": "zsh", "file": os.Getenv("HISTFILE")}
		// } else if strings.Contains(shell, "bash") {
		// } else {
		// 	return map[string]string{}
		// }
	}
}

func main() {
	var (
		history_file = flag.String("f", "default", "history file path")
	)
	flag.Parse()

	// path := ""

	// if strings.HasPrefix(*history_file, "~/") {
	// 	dirname, _ := os.UserHomeDir()
	// 	path = filepath.Join(dirname, (*history_file)[2:])
	// } else {
	// 	path, _ = filepath.Abs(*history_file)
	// }
	// fmt.Println(*history_file)
	// fmt.Println(path)

	fmt.Println("ファイル読み取り処理を開始します")

	f, err := os.Open(*history_file)

	if err != nil {
		fmt.Println("error")
	}
	// 関数が終了した際に確実に閉じるようにする
	defer f.Close()

	scanner := bufio.NewScanner(f)

	re1 := regexp.MustCompile("^ +")
	re2 := regexp.MustCompile(" +$")
	re3 := regexp.MustCompile(" +")

	sep := map[string]string{"pwsh": "`", "zsh": "\\", "bash": "\\"}[get_shell()]

	stock := ""

	data := ""

	for scanner.Scan() {
		txt := scanner.Text()
		txt = re1.ReplaceAllString(txt, "")
		txt = re2.ReplaceAllString(txt, "")
		txt = re3.ReplaceAllString(txt, " ")

		stock += txt
		if !strings.HasSuffix(stock, sep) {
			if !strings.Contains(data, stock) {
				data += stock + "\n"
			}
			stock = ""
		}
	}

	sdata := strings.Split(data, "\n")

	sort.Slice(sdata, func(i, j int) bool {
		return sdata[i] < sdata[j]
	})

	file2, err := os.Create("./output.txt")

	if err != nil {
		fmt.Println(err)
	}

	defer file2.Close()

	// re4 := regexp.MustCompile("\\")

	for _, v := range sdata {
		v = strings.ReplaceAll(v, "\\", "\\\n")

		file2.WriteString(v)
		file2.WriteString("\n")

	}
}
