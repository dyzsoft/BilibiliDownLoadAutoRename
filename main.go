package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var baseDir = ""
var jsonConfigFile = "entry.json"

type PageData struct {
	Part string `json:"part"`
}

type Entry struct {
	PageData `json:"page_data"`
}

/***
		说明： b站离线现在的文件  拷贝到pc之后，真实文件夹名称都在 entry.json 这个文件 的page_data字段 中的part字段
        程序直接读取当目录，寻找entry.json 并且进行批量重命名
	    如果程序第二个参数不为空，则第二个参数为指定的 baseDir

*/

func main() {
	//  获取baseDir，并且切换到对应目录
	//
	if len(os.Args) >= 2 {
		baseDir = os.Args[1] + "/"
	} else {
		baseDir, _ = os.Getwd()
		baseDir = baseDir + "/"
	}

	err := os.Chdir(baseDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有文件夹
	dirNames, err := GetAllDirName()

	if err != nil {
		fmt.Println("读取指定目录失败")
		fmt.Println(err)
		return
	}
	// 循环判断每个文件夹，如果有entry.json 文件，则通过GetRealDirName 获取json文件中的真实文件名，并且进行重命名

	for i := 0; i < len(dirNames); i++ {

		if realname, err := GetRealDirName(baseDir + dirNames[i]); err == nil {

			if err := ChangeDirName(baseDir+dirNames[i], baseDir+realname); err != nil {
				fmt.Printf("重命名文件夹: %s 失败！\n", baseDir+dirNames[i])
				fmt.Println(err)
				fmt.Println()
			}

		} else {
			fmt.Printf("获取真实文件夹 %s 的名称失败！\n", baseDir+dirNames[i])
			fmt.Println(err)
			fmt.Println()
			continue
		}
	}

}

func GetAllDirName() (names []string, err error) {

	fileDir, err := ioutil.ReadDir(baseDir)

	if err == nil {
		for i := 0; i < len(fileDir); i++ {
			if fileDir[i].IsDir() {
				names = append(names, fileDir[i].Name())
			}
		}
	}
	return
}

func GetRealDirName(name string) (realname string, err error) {

	var fileContent []byte
	fileContent, err = ioutil.ReadFile(name + "/" + jsonConfigFile)

	if err == nil {
		var entry Entry
		err = json.Unmarshal(fileContent, &entry)
		if err == nil {
			realname = entry.PageData.Part
		}
	}

	return realname, err
}
func ChangeDirName(old string, new string) error {
	return os.Rename(old, new)
}
