package main 

import ("fmt"
		"io/ioutil"
		"strings"
		"regexp"
		"os"

		)
func jsonarrayParser(jsonData string) []string{
	fmt.Println("Inside array parser")
	//fmt.Println(jsonData)
	parsedarray := make([]string,0)
	if jsonData[0] == '['{
		jsonData = jsonData[1:]
		for len(jsonData) > 0 {

			result := parserCombinator(jsonData, stringParser, numberParser, boolParser, nullParser, jsonarrayParser)
			//fmt.Println(result)
			if len(result) > 0{
				parsedarray = append(parsedarray,result[0])	
				jsonData = result[1]
				//fmt.Println(parsedarray)
			 	result = commaParser(jsonData)
				if result != nil{
					parsedarray = append(parsedarray,result[0])
					jsonData = jsonData[1:]
				}
			//fmt.Println(jsonData)
			if jsonData[0] == ']'{
				returnV := appendall2(parsedarray, jsonData[1:])
				return returnV
			}
		}
	}
}
	return nil
} 

func appendall2(result []string,data string) []string{
	returnV := make([]string,0)
	returnV = append(returnV, result...)
	returnV = append(returnV,data)

	return returnV
}

func stringParser(jsonData string)  []string {
	fmt.Println("Inside string parser")
	jsonData = strings.TrimSpace(jsonData)
	var index int
	returnV := make([]string,0)
	if string(jsonData[0]) ==string('"'){
		fmt.Println("string found ")
		jsonData = jsonData[1:]
		for jsonData[index] != '"'{
			index++
		}
		//fmt.Println(index)
		returnV = appendall(jsonData[:index],jsonData[index+1:])
		return returnV
	}
	return returnV	
}
func numberParser(jsonData string) []string{
	fmt.Println("Inside number parser")
	jsonData = strings.TrimSpace(jsonData)
	var result [][]int
	returnV := make([]string,0)
	re:= regexp.MustCompile("^[-+]?[0-9]*.?[0-9]+([eE][-+]?[0-9]+)?")
	result1 := re.FindAllString(jsonData, -1)
	if len(result1) > 0{
		//fmt.Println("inside num if ")
		result = re.FindAllStringIndex(jsonData, -1)
		//fmt.Println(result)
		index := result[0][1]
		//fmt.Println(index)
		//fmt.Println("Flag 1")
		returnV = append(returnV,result1[0])
		returnV = append(returnV, jsonData[index:])
		return returnV
	}else{
		//fmt.Println("Flag 2")
		return returnV
	}
}

func appendall(result string,data string) []string{
	returnV := make([]string,0)
	returnV = append(returnV, result)
	returnV = append(returnV,data)

	return returnV
}
func boolParser(jsonData string) []string{
	fmt.Println("Inside bool parser")
	jsonData = strings.TrimSpace(jsonData)
	if len(jsonData) > 4{
		if jsonData[0:4] == "true"{
			//fmt.Println("bool found")
			returnV := appendall("true",jsonData[4:])
			return returnV
		}else if jsonData[0:5] == "false" {
			//fmt.Println("bool found")
			returnV := appendall("false",jsonData[5:])
			return returnV
		}
	}
	return nil
}
func nullParser(jsonData string) []string {
	fmt.Println("Inside null parser")
	jsonData = strings.TrimSpace(jsonData)
	if len(jsonData) > 4{
		if jsonData[0:4] == "null"{
			fmt.Println("null found")
			returnV := appendall("nil",jsonData[4:])
			return returnV
		}
	}
	return nil
}
func commaParser(jsonData string) []string {
	fmt.Println("Inside comma parser")
	jsonData = strings.TrimSpace(jsonData)
	if jsonData[0] ==','{
		fmt.Println("Comma found")
		//fmt.Println(jsonData[1:])
		returnV := appendall(",",jsonData)
		return returnV
	}
	return nil
}
func coloParser(jsonData string) []string {
	fmt.Println("Inside colon parser")
	jsonData = strings.Trim(jsonData," ")
	result := string(jsonData[0])
	if jsonData[0] == ':'{
		 fmt.Println("colon found")
		 //fmt.Println(jsonData[1:])
		 returnV := appendall(result,jsonData[1:])
		return returnV
	}
	return nil
}
func parserCombinator(jsonData string,args ...func(string) []string) []string{
	
		for _,each := range args {
		result := each(jsonData)
		//fmt.Println(result)
		if len(result) > 0{
			return result
		}else {
			continue
		}
	}
	return nil
}
func main(){
	//file, err := os.Open(s)
	buf,err := ioutil.ReadFile("test.txt")// whole file at a time (Readall) 
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	data:= string(buf)
	//fmt.Println(data)
	fmt.Println(parserCombinator(data, jsonarrayParser))
}