package main

import (
	"fmt"
    "bufio"
    "os"
    "strings"
)

func ReadSiritori(data_name string)(map[string][]string){
    fp, _ := os.Open(data_name)
    defer fp.Close()
    scanner := bufio.NewScanner(fp)
    var ans_map map[string][]string = make(map[string][]string)
    for scanner.Scan(){
        line := scanner.Text()
        first_word = strings.Split(line, "")[0]
        if 
    }
    return ans_map
}

func main(){
    ReadSiritori("siritori_data.txt")
    fmt.Println("aaa")
}