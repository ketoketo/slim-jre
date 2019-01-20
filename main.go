package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	fmt.Println(createDepModulesWithComma("/home/tmatsuzaki/Downloads/BOOT-INF/lib", []string{"logback-classic", "lombok"}))

}

func createDepModulesWithComma(searchPath string, excludeJarNames []string) string {
	fis, err := ioutil.ReadDir(searchPath)

	if err != nil {
		panic(err)
	}

	modulesSet := make(map[string]struct{})
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if isExcludeJar(fi.Name(), excludeJarNames) {
			continue
		}

		filepath := filepath.Join(searchPath, fi.Name())
		pos := strings.LastIndex(filepath, ".")
		if filepath[pos:] == ".jar" {
			fmt.Println(string(filepath))
			out, err := executeJdeps(filepath)
			if err != nil {
				panic(err)
			}
			jdepsResult := strings.Split(string(out), ",")
			if len(jdepsResult) == 1 {
				continue
			}
			createMoludesSet(jdepsResult, modulesSet)
		}

	}
	keys := make([]string, 0, len(modulesSet))
	for k := range modulesSet {
		keys = append(keys, k)
	}
	return strings.Join(keys, ",")
}

func isExcludeJar(e string, s []string) bool {
	for _, v := range s {
		if strings.Contains(e, v) {
			return true
		}
	}
	return false
}

func executeJdeps(filepath string) ([]byte, error) {
	out, err := exec.Command("jdeps", "--print-module-deps", "-q", filepath).CombinedOutput()
	if err != nil {
		out, err = exec.Command("jdeps", "--print-module-deps", "-q", "--multi-release", "11", filepath).CombinedOutput()
	}
	return out, err
}

func createMoludesSet(jdepsResult []string, modulesSet map[string]struct{}) {
	for _, module := range jdepsResult {
		newLineRepModule := strings.Replace(module, "\r\n", "\n", -1)
		newLineRepModule = strings.Replace(newLineRepModule, "\n", "", -1)
		modulesSet[newLineRepModule] = struct{}{}
	}
}
