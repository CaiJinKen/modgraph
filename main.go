package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var modName = ""

const _ModFile = "go.mod"

func main() {
	dependencies := getDependenciesByte()
	if len(dependencies) == 0 {
		fmt.Println("go.mod is empty.")
		os.Exit(0)
	}

	rootMod := &mod{Name: []byte(getProjectName())}
	nameModMaps := make(map[string]*mod)
	nameModMaps[getProjectName()] = rootMod

	lines := bytes.Split(dependencies, []byte{'\n'})
	for _, line := range lines {
		deps := bytes.Split(bytes.TrimSpace(line), []byte{' '})
		if len(deps) < 2 {
			continue
		}
		dep1, dep2 := deps[0], deps[1]
		dep1Str, dep2Str := string(dep1), string(dep2)
		if m := nameModMaps[dep1Str]; m == nil {
			nameModMaps[dep1Str] = &mod{Name: dep1}
		}
		if m := nameModMaps[dep2Str]; m == nil {
			nameModMaps[dep2Str] = &mod{Name: dep2}
		}

		nameModMaps[dep1Str].Mods = append(nameModMaps[dep1Str].Mods, nameModMaps[dep2Str])
	}
	rootMod.String("", 0)
}

func getDependenciesByte() []byte {
	buf := bytes.NewBuffer(nil)
	cmd := exec.Command("go", "mod", "graph")
	cmd.Stderr = os.Stderr
	cmd.Stdout = bufio.NewWriter(buf)
	if err := cmd.Run(); err != nil {
		fmt.Println("get the project dependencies err:", err.Error())
	}

	return buf.Bytes()
}

func getProjectName() string {
	if modName != "" {
		return modName
	}

	lines := getModFileContent()
	if len(lines) == 0 || len(lines[0]) <= 7 {
		return ""
	}
	modName = string(lines[0][7:])
	return modName
}

func getModFileContent() (content [][]byte) {
	modFileContent, err := ioutil.ReadFile(_ModFile)
	if err != nil {
		fmt.Println("read file go.mod with err:", err.Error())
		return
	}

	lines := bytes.Split(modFileContent, []byte{'\n'})
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) < 2 || bytes.HasPrefix(line, []byte("go ")) || bytes.HasSuffix(line, []byte(")")) || bytes.HasSuffix(line, []byte("require ")) {
			continue
		}
		content = append(content, line)
	}
	return
}

type mod struct {
	Name     []byte
	Mods     []*mod
	traveled bool
}

func (m *mod) String(str string, n int) {
	if m == nil {
		return
	}

	ft := "  |"
	if n == 0 {
		fmt.Println(string(m.Name))
	} else {
		fmt.Printf(str+"--%s\n", string(m.Name))
	}
	if len(m.Mods) == 0 || m.traveled {
		return
	}

	m.traveled = true
	for _, v := range m.Mods {
		v.String(str+ft, n+1)
	}
}
