package main

import ("fmt"
		"strings"
		"io/ioutil"
		"regexp"
		"os"
		)

type mystring []interface{}

var parsed mystring

func objparser(jsonData string) string{
	fmt.Println("Inside array parser")
	var result string
	var key string
	m := make(map[string]interface{})
	if jsonData[0] == '{'{
		jsonData = jsonData[1:]
		for len(jsonData) > 0{
			result,jsonData = stringParser(jsonData)
			if result != "" {
				 key = result
			}
			result,jsonData = colonParser(jsonData)
			result,jsonData = elementParser(jsonData)
			if result != "" {
				m[key] = result
				result,jsonData = commaParser(jsonData)
			}
			if jsonData[0] == '[' {
				jsonData = arrayParser(jsonData)
			}
			if jsonData[0] == '}' {
				parsed = append(parsed,m)
				return  jsonData[1:]
			}
		}
	}
	return jsonData
}



func arrayParser(jsonData string) string{
	fmt.Println("Inside array parser")
	var result string
	if jsonData[0] == '['{
		jsonData = jsonData[1:]
		for len(jsonData) > 0{
			result,jsonData = elementParser(jsonData)
			if result != "" {
				if result == "nil" {
					parsed = append(parsed,nil)
				}else {
					parsed = append(parsed,result)
					result,jsonData = commaParser(jsonData)
					parsed = append(parsed,result)
				}
			}
			if jsonData[0] == '{' {
				jsonData = objparser(jsonData)
			}
			if jsonData[0] == ']' {
					return  jsonData[1:]
			}
		}
	}
	return jsonData
}

func stringParser(jsonData string) (string,string){
	fmt.Println("Inside String parser")
	var result string
	var index int
	jsonData = strings.TrimSpace(jsonData)
	if jsonData[0] =='"'{
		fmt.Println("string found ")
		jsonData = jsonData[1:]
		for jsonData[index] != '"'{
			index++
		}
		result = jsonData[:index]
		return result,jsonData[index +1:]
	
	fmt.Println("outside stringParser")
	}	
	return result,jsonData
}

func booleanParser(jsonData string) (string,string) {
	fmt.Println("Inside boolean parser")
	var result string
	jsonData = strings.Trim(jsonData, " ")
	if len(jsonData) > 4{
		if jsonData[0:4] == "true"{
			result = "true"
			return result,jsonData[4:]
		}else if jsonData[0:5] == "false" {
			result = "false"
			return result,jsonData[5:]
		}
	}
	return result,jsonData
		
}


func colonParser(jsonData string) (string,string) {
	fmt.Println("Inside colon parser")
	var result string
	jsonData = strings.Trim(jsonData," ")
	
	if jsonData[0] == ':'{
		 fmt.Println("colon found")
		 result = string(jsonData[0])
		return result,jsonData[1:]
	}
	return result,jsonData
}

func commaParser(jsonData string) (string,string){
	var result string
	fmt.Println("inside comma parser")
	jsonData = strings.Trim(jsonData, " ")
	if jsonData[0] ==','{
		result = ","
		fmt.Println("Comma found")
		return result,jsonData[1:]
	}
	return result,jsonData

}

func nullParser(jsonData string) (string,string) {
	fmt.Println("Inside null parser")
	var result string
	jsonData = strings.Trim(jsonData, " ")
	if len(jsonData) > 4{
		if jsonData[0:4] == "null"{
			fmt.Println("null found")
			result = "nil"
			return result,jsonData[4:]
		}
	}
	return result,jsonData
}
func numberParser(jsonData string) (string,string){
	jsonData = strings.Trim(jsonData, " ")
	fmt.Println("inside number parser")
	var result [][]int
	re:= regexp.MustCompile("^[-+]?[0-9]*.?[0-9]+([eE][-+]?[0-9]+)?")
	result1 := re.FindAllString(jsonData, -1)
	if len(result1) > 0{
		result = re.FindAllStringIndex(jsonData, -1)
		fmt.Println(result)
		index := result[0][1]
		return result1[0],jsonData[index:]
	}else{
		return  "",jsonData
	}
	
}
func elementParser(jsonData string) (string,string) {
	var result string
	fmt.Println("INside Element Parser")
	result,jsonData = stringParser(jsonData)
	if result != ""{
		return result,jsonData
	}
	result,jsonData = numberParser(jsonData)
	if result != ""{
		return result,jsonData
	}
	result,jsonData = booleanParser(jsonData)
	if result != ""{
		return result,jsonData
	}
	result,jsonData = nullParser(jsonData)
	if result != ""{
		return result,jsonData
	}
	result,jsonData = commaParser(jsonData)
	if result != ""{
		return result,jsonData
	}
	jsonData = arrayParser(jsonData)
	jsonData = objparser(jsonData)
	return result,jsonData

}

func main() {
	buf,err := ioutil.ReadFile("test.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	data:= string(buf)
	fmt.Println(data)
	arrayParser(data)
	fmt.Println(parsed)


}