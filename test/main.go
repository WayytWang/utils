package main

import "log"

type files struct {
	list []*file
	fileMap map[string]int
}

type file struct {
	src *files
	name string
}

func main() {
	files := &files{
		list:    make([]*file,0),
		fileMap: make(map[string]int),
	}
	file1 := &file{
		src:  files,
		name: "file1",
	}

	file2 := &file{
		src:  files,
		name: "file2",
	}
	files.list = append(files.list, file1)
	files.list = append(files.list, file2)
	files.fileMap[file1.name] = 1
	files.fileMap[file2.name] = 2

	log.Printf("%+v",files.fileMap)
	log.Printf("%+v",file1.src.fileMap)
	log.Printf("%+v",file2.src.fileMap)
}