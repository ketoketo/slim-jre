package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const WORKDIR = "slim-jre-tmp"

func main() {
	// app := cli.NewApp()
	// app.Name = "Goshirase"
	// app.Usage = "test usage"
	// app.Version = "0.0.1"

	// // flags
	// configName := ".goshirase/config"
	// app.Flags = []cli.Flag{
	// 	cli.StringFlag{
	// 		Name:        "profile, p",
	// 		Value:       ".goshirase/config",
	// 		Usage:       "config file name",
	// 		Destination: &configName,
	// 	},
	// }

	// app.Commands = []cli.Command{
	// 	{
	// 		Name:    "configure",
	// 		Aliases: []string{"c"},
	// 		Usage:   "set config file",
	// 		Action: func(c *cli.Context) error {
	// 			err := registerConfig(configName)
	// 			return err
	// 		},
	// 	},
	// }
	// sort.Sort(cli.FlagsByName(app.Flags))
	// sort.Sort(cli.CommandsByName(app.Commands))

	// err := app.Run(os.Args)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// 対象のjarの依存関係を調べる
	targetJar := "loader.jar"
	modulesSet := make(map[string]struct{})
	jdepsResult := executeJdeps(targetJar)
	createMoludesSet(jdepsResult, modulesSet)

	// jar内に依存ライブラリがある場合ワークディレクトリに解答し、そこで依存関係を調べる
	mkWorkDir()
	executeUnzipJar(WORKDIR, targetJar)
	addInnerDepModules("C:\\Users\\tmatsuzaki\\Documents\\go\\src\\slim-jre\\slim-jre-tmp\\BOOT-INF\\lib", []string{"logback-classic", "lombok"}, modulesSet)
	fmt.Println(createModulesStringWithComma(modulesSet))
	delete(WORKDIR)
}

func mkWorkDir() {
	if _, err := os.Stat(WORKDIR); !os.IsNotExist(err) {
		delete(WORKDIR)
	}

	if err := os.Mkdir(WORKDIR, 0777); err != nil {
		panic(err)
	}
}

func delete(name string) {
	if err := os.RemoveAll(name); err != nil {
		panic(err)
	}
}

func copy(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func addInnerDepModules(searchPath string, excludeJarNames []string, modulesSet map[string]struct{}) {
	fis, err := ioutil.ReadDir(searchPath)

	if err != nil {
		panic(err)
	}

	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if isExcludeJar(fi.Name(), excludeJarNames) {
			continue
		}

		fullpath := filepath.Join(searchPath, fi.Name())
		if filepath.Ext(fullpath) == ".jar" {
			fmt.Println(string(fullpath))
			jdepsResult := executeJdeps(fullpath)
			if len(jdepsResult) == 1 && (jdepsResult[0] == "\r\n" || jdepsResult[0] == "\n") {
				// 依存ライブラリが何もない場合
				continue
			}
			createMoludesSet(jdepsResult, modulesSet)
			fmt.Println(createModulesStringWithComma(modulesSet))
		}
	}
}

func isExcludeJar(e string, s []string) bool {
	for _, v := range s {
		if strings.Contains(e, v) {
			return true
		}
	}
	return false
}

func executeUnzipJar(workPath string, jarName string) {
	pwd, _ := os.Getwd()
	jarPath := filepath.Join(pwd, jarName)
	os.Chdir(workPath)
	err := exec.Command("jar", "-xf", jarPath).Run()
	if err != nil {
		panic(err)
	}
	os.Chdir(pwd)
}

func executeJdeps(filepath string) []string {
	out, err := exec.Command("jdeps", "--print-module-deps", "-q", filepath).CombinedOutput()
	if err != nil {
		out, err = exec.Command("jdeps", "--print-module-deps", "-q", "--multi-release", "11", filepath).CombinedOutput()
	}
	if err != nil {
		panic(err)
	}
	jdepsResult := strings.Split(string(out), ",")
	return jdepsResult
}

func createMoludesSet(jdepsResult []string, modulesSet map[string]struct{}) {
	for _, module := range jdepsResult {
		newLineRepModule := strings.Replace(module, "\r\n", "\n", -1)
		newLineRepModule = strings.Replace(newLineRepModule, "\n", "", -1)
		modulesSet[newLineRepModule] = struct{}{}
	}
}

func createModulesStringWithComma(modulesSet map[string]struct{}) string {
	keys := make([]string, 0, len(modulesSet))
	for k := range modulesSet {
		keys = append(keys, k)
	}
	return strings.Join(keys, ",")
}
