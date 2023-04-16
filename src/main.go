package main

import (
	"bufio"
	"flag"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

func get_shell() (string, string, string) {
	if runtime.GOOS == "windows" {
		out, _ := exec.Command("powershell", "(Get-PSReadlineOption).HistorySavePath").Output()
		target_path := strings.ReplaceAll(string(out), "\r", "")
		target_path = strings.ReplaceAll(target_path, "\n", "")
		return "pwsh", target_path, "`"
		// return "pwsh", strings.ReplaceAll(string(out), "\\", "\\\\"), "`"
	} else {
		shell_path := os.Getenv("SHELL")

		t := strings.Split(shell_path, "/")
		sh := t[len(t)-1]

		home_dir, _ := os.UserHomeDir()
		target_path := home_dir + "/." + sh + "_history"
		return sh, string(target_path), "\\"
	}
}

func read_and_replace(history_file string, sep string) []string {
	f, err := os.Open(history_file)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	re1 := regexp.MustCompile("^ +")
	re2 := regexp.MustCompile(" +$")
	re3 := regexp.MustCompile(" +")

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

	return sdata
}
func save_history_file(output_file string, data []string, sep string) {
	saved_file, err := os.Create(output_file)

	if err != nil {
		panic(err)
	}

	defer saved_file.Close()

	for _, v := range data {
		v = strings.ReplaceAll(v, sep, sep+"\n")

		saved_file.WriteString(v)
		saved_file.WriteString("\n")

	}
}

func main() {
	var (
		history_file = flag.String("f", "", "input history file path")
		output_file  = flag.String("o", "", "output history file path")
		shell        = flag.String("s", "", "select shell (bash, zsh, pwsh)")
	)
	flag.Parse()

	tmp_shell, tmp_history_file, separte_string := get_shell()
	if *shell == "" {
		*shell = tmp_shell
	}

	if *history_file == "" {
		*history_file = tmp_history_file
	}

	if *output_file == "" {
		*output_file = *history_file
	}

	sdata := read_and_replace(*history_file, separte_string)
	save_history_file(*output_file, sdata, separte_string)
}
