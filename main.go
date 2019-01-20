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

	if err := os.Mkdir("slim-jre-tmp", 0777); err != nil {
		fmt.Println(err)
	}

	copy("loader.jar", filepath.Join("slim-jre-tmp", "loader.jar"))

	// fmt.Println(createDepModulesWithComma("/home/tmatsuzaki/Downloads/BOOT-INF/lib", []string{"logback-classic", "lombok"}))

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

		fullpath := filepath.Join(searchPath, fi.Name())
		if filepath.Ext(fullpath) == ".jar" {
			fmt.Println(string(fullpath))
			out, err := executeJdeps(fullpath)
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
