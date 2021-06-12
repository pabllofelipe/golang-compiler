package main

import (
	"C"
	"bufio"
	"fmt"
	"io"
	"os"
)
func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

func createObjFile() {
	// check if file exists
	var _, err = os.Stat("programa.c")

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create("programa.c")
		if isError(err) {
			return
		}
		defer file.Close()
	}
}

func writeObjFile(str string, control int) {
	// Open file using READ & WRITE permission.
	var file, err = os.OpenFile("programa.c", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if isError(err) {
		return
	}
	defer file.Close()

	// Write some text line-by-line to file.
	_, err = file.WriteString(str)
	if isError(err) {
		return
	}

	// Save file changes.
	err = file.Sync()
	if isError(err) {
		return
	}
}

func readFile() {
// Open file for reading.
var file, err = os.OpenFile("programa.c", os.O_RDWR, 0644)
if isError(err) {
return
}
defer file.Close()

// Read file, line by line
var text = make([]byte, 1024)
for {
_, err = file.Read(text)

// Break if finally arrived at end of file
if err == io.EOF {
break
}

// Break if error occured
if err != nil && err != io.EOF {
isError(err)
break
}
}

fmt.Println("Reading from file.")
fmt.Println(string(text))
}

func deleteObjFile() {
	// delete file
	var err = os.Remove("programa.c")
	if isError(err) {
		return
	}

	fmt.Println("File Deleted")
}

func addBegFile(str []string){
	inicio := "#include<stdio.h>\ntypedef char literal[256];\nvoid main(void)\n{\n"
	variables := "/*----Variaveis temporarias----*/\n"
	for _, i:= range str{
		variables = variables + "int " + i + ";\n"
	}
	inicio = inicio + variables + "/*------------------------------*/\n"
	// make a temporary outfile
	outfile, err := os.Create("programa2.c")

	if err != nil {
		panic(err)
	}

	defer outfile.Close()
	// open the file to be appended to for read
	f, err := os.Open("programa.c")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	// append at the start
	_, err = outfile.WriteString(inicio)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)

	// read the file to be appended to and output all of it
	for scanner.Scan() {

		_, err = outfile.WriteString(scanner.Text())
		_, err = outfile.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	// ensure all lines are written
	outfile.Sync()
	// over write the old file with the new one
	err = os.Rename("programa2.c", "programa.c")
	if err != nil {
		panic(err)
	}
}