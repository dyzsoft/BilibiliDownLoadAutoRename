package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var baseDir = ""

//var  baseDir = `G:\go\26604071`

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
	// 获取baseDir，并且切换到对应目录
	if baseDir == "" {
		if len(os.Args) >= 2 {
			baseDir = os.Args[1]
		} else {
			baseDir, _ = os.Getwd()
		}
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

		curdirname := filepath.Join(baseDir , dirNames[i])

		if realname, err := GetRealDirName(curdirname); err == nil {

			fmt.Println(filepath.Join(baseDir,realname))

			if err := ChangeDirName(curdirname, filepath.Join(baseDir,realname)); err != nil {
				fmt.Printf("重命名文件夹: %s 失败！\n",curdirname)
				fmt.Println(err)
				fmt.Println()
			}

		} else {
			fmt.Printf("获取真实文件夹 %s 的名称失败！\n", curdirname)
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
	fileContent, err = ioutil.ReadFile(filepath.Join(name,jsonConfigFile))

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

	ChangeExtBlvToMp4(old)  // 重命名 .blv  to .mp4

	return os.Rename(old, new)
}

func ChangeExtBlvToMp4(dirname string) error{
	fileInfos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}
	for i := 0; i < len(fileInfos); i++ {
		fileinfo := fileInfos[i]
		if  fileinfo.IsDir() {
			realdir := filepath.Join(dirname,fileinfo.Name())
			finfos, err := ioutil.ReadDir(realdir)
			if err != nil {
				continue
			}
			for i:=0 ;i<len(finfos);i++{
				filename := finfos[i].Name()
				//fmt.Println(filepath.Ext(filename))
				if filepath.Ext(filename) == ".blv" {
					err := os.Rename(filepath.Join(realdir,filename), filepath.Join(realdir,strings.Replace(filename, "blv", "mp4", -1)))
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}
	return nil
}
