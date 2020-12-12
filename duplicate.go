package main

import (
	"flag"
	"fmt"
	"github.com/qwenode/gogo/convert"
	"github.com/qwenode/gogo/file"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	source := flag.String("s", "", "set source directory")
	to := flag.String("t", "", "set backup directory")
	flag.Parse()
	sourceDir, _ := filepath.Abs(*source)
	toDir, _ := filepath.Abs(*to)
	log.Println(sourceDir, toDir)
	if !file.Exist(sourceDir) || !file.IsDirectory(sourceDir) {
		log.Fatalln("source directory does not exist")
	}
	if !file.Exist(toDir) || !file.IsDirectory(toDir) {
		log.Fatalln("backup directory does not exist")
	}
	if sourceDir == toDir {
		log.Fatalln("source directory can't same as flat directory")
	}
	hashList := map[string]string{}
	moved := 0
	skipped := 0
	failed := 0
	_ = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if info.Size() < 10 {
			skipped++
			log.Println("skipped for small size:", path)
			return nil
		}
		//log.Println(path, info.IsDir(), info.Name(), err)
		sha, err := file.Sha1(path)
		if err != nil {
			return nil
		}
		if _, ok := hashList[sha]; !ok {
			hashList[sha] = path
			skipped++
			log.Println("Skipped:", path)
			return nil
		}
		toFile := toDir + "/" + info.Name()
		if file.Exist(toFile) {
			toFile = fmt.Sprintf("%s/%s_%s", toDir, convert.ToString(time.Now().Unix()), info.Name())
		}
		log.Println("Moving File:", path, " To:", toFile)
		err = os.Rename(path, toFile)
		if err != nil {
			failed++
			log.Println("Move failed.")
			return nil
		}
		moved++
		log.Println("Success Moved.")
		return nil
	})
	log.Println(fmt.Sprintf("Skipped:%d,Moved:%d,Failed:%d", skipped, moved, failed))
}
