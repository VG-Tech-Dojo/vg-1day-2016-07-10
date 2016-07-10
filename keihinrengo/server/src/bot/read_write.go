package read_write

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
        first_word := strings.Split(line, "")[0]
        _, ok := ans_map[first_word]
        if ok == true{
            ans_map[first_word] = append(ans_map[first_word], line)
        }else{
            ans_map[first_word] = []string{line}
        }
    }
    return ans_map
}