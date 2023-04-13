package main

import (
    "fmt"
    "io/fs"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"
    
    "github.com/qwenode/color"
    "github.com/qwenode/gogo/convert"
    "github.com/qwenode/gogo/file"
)

func main() {
    run()
    color.SuccessMessage("全部检查完毕,10分钟后关闭窗口,您也可以现在关闭")
    time.Sleep(time.Second * 600)
}
func run() {
    // catch all panic
    defer func() {
        err := recover()
        if err != nil {
            //debug.PrintStack()
            log.Println("System Error:", err)
        }
    }()
    sourceDir, _ := filepath.Abs("./")
    toDir := filepath.Join(sourceDir, "_dup")
    log.Println(sourceDir, toDir)
    if !file.Exist(sourceDir) || !file.IsDirectory(sourceDir) {
        log.Fatalln("source directory does not exist")
    }
    if !file.Exist(toDir) || !file.IsDirectory(toDir) {
        os.MkdirAll(toDir, os.ModePerm)
    }
    if sourceDir == toDir {
        log.Fatalln("source directory can't same as flat directory")
    }
    hashList := map[string]string{}
    moved := 0
    skipped := 0
    failed := 0
    _ = filepath.WalkDir(sourceDir, func(path string, info fs.DirEntry, err error) error {
        if strings.HasPrefix(path, toDir) {
            return nil
        }
        color.InfoMessage("检查:%s", path)
        if err != nil || info.IsDir() {
            color.InfoMessage("文件夹,跳过")
            return nil
        }
        //if info.Size() < 10 {
        //	skipped++
        //	log.Println("skipped for small size:", path)
        //	return nil
        //}
        //log.Println(path, info.IsDir(), info.Name(), err)
        sha, err := file.Sha1(path)
        if err != nil {
            return nil
        }
        if _, ok := hashList[sha]; !ok {
            hashList[sha] = path
            skipped++
            color.InfoMessage("文件未重复:%s", path)
            return nil
        }
        toFile := filepath.Join(toDir, info.Name())
        if file.Exist(toFile) {
            toFile = fmt.Sprintf("%s/%s_%s", toDir, convert.ToString(time.Now().Unix()), info.Name())
        }
        color.WarningMessage("重复文件=>%s", toFile)
        err = os.Rename(path, toFile)
        if err != nil {
            failed++
            color.ErrorMessage("重复文件移动失败:%s", err.Error())
            return nil
        }
        moved++
        color.SuccessMessage("文件移动成功:%s", path)
        return nil
    })
    color.SuccessMessage(fmt.Sprintf("数据统计=>跳过:%d,重复:%d,失败:%d", skipped, moved, failed))
    
}
