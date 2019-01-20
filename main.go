package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// out, err := exec.Command("jdeps", "--list-deps", "C:\\Users\\tmatsuzaki\\Downloads\\loader.jar").CombinedOutput()
	// out, err := exec.Command("jar", "-xvf", "C:\\Users\\tmatsuzaki\\Downloads\\loader.jar").CombinedOutput()

	listFiles("/home/tmatsuzaki/Downloads/BOOT-INF/lib")

}

func listFiles(searchPath string) {
	fis, err := ioutil.ReadDir(searchPath)

	if err != nil {
		panic(err)
	}

	removes := []string{"logback-classic", "lombok"}

	modules := make(map[string]struct{})
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if isRemoveJar(fi.Name(), removes) {
			continue
		}

		filepath := filepath.Join(searchPath, fi.Name())
		pos := strings.LastIndex(filepath, ".")
		if filepath[pos:] == ".jar" {
			fmt.Println(string(filepath))
			out, err := exec.Command("jdeps", "--print-module-deps", "-q", filepath).CombinedOutput()
			if err != nil {
				out, err = exec.Command("jdeps", "--print-module-deps", "-q", "--multi-release", "11", filepath).CombinedOutput()
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
			outValue := strings.Split(string(out), ",")
			if len(outValue) == 1 {
				continue
			}
			for _, module := range outValue {
				newLineRepModule := strings.Replace(module, "\r\n", "\n", -1)
				newLineRepModule = strings.Replace(newLineRepModule, "\n", "", -1)
				modules[newLineRepModule] = struct{}{}
			}
		}

	}
	keys := make([]string, 0, len(modules))
	for k := range modules {
		keys = append(keys, k)
	}
	fmt.Println(strings.Join(keys, ","))

}

func isRemoveJar(e string, s []string) bool {
	for _, v := range s {
		if strings.Contains(e, v) {
			return true
		}
	}
	return false
}
