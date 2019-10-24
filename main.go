package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func main() {
	bts := getByts()
	if len(bts) == 0 {
		fmt.Println("go.mod is empty.")
		os.Exit(0)
	}

	rootMod := &mod{Name: []byte(getProjectName())}
	maps := make(map[string]*mod)
	maps[getProjectName()] = rootMod

	btsls := bytes.Split(bts, []byte{'\n'})
	for _, lv := range btsls {
		mods := bytes.Split(bytes.TrimSpace(lv), []byte{' '})
		if len(mods) < 2 {
			continue
		}
		m1, m2 := mods[0], mods[1]
		m1str, m2str := string(m1), string(m2)
		if m := maps[m1str]; m == nil {
			maps[m1str] = &mod{Name: m1}
		}
		if m := maps[m2str]; m == nil {
			maps[m2str] = &mod{Name: m2}
		}

		maps[m1str].Mods = append(maps[m1str].Mods, maps[m2str])
	}
	rootMod.String("", 0)
}

func getByts() []byte {
	buf := bytes.NewBuffer(nil)
	cmd := exec.Command("go", "mod", "graph")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = bufio.NewWriter(buf)
	if err := cmd.Run(); err != nil {
		fmt.Println("get the project dependencies err:", err.Error())
	}

	return buf.Bytes()
}

func getProjectName() string {
	bts, err := ioutil.ReadFile("go.mod")
	if err != nil {
		fmt.Println("read file go.mod with err:", err.Error())
		return ""
	}
	bts = bytes.TrimSpace(bts)
	if i := bytes.IndexByte(bts, '\n'); i > 7 {
		bts = bts[7:i]
	}
	return string(bts)
}

type mod struct {
	Name   []byte
	Mods   []*mod
	n      int
	traved bool
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
	if len(m.Mods) == 0 || m.traved {
		return
	}

	m.traved = true
	for _, v := range m.Mods {
		v.String(str+ft, n+1)
	}
}
