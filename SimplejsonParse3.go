package main 

import ("fmt"
		"io/ioutil"
		"strings"
		"os"
		"regexp"
		"strconv"
		//"reflect"
		)


type json struct {
	 Bool bool
	 String string
	 Number int
	 Float float64 
	 Array []json
	 Object map[string]json
	 Parsed bool
	 Null interface{}
	 Type string
}



func stringParser(jsonData string) (json, string) {
	//This function parses string elements and returns the string and remaining jsondata
	var result json
	var index int                          // stores the index of end of string
	jsonData = strings.TrimSpace(jsonData) // Trim the data for spaces
	if jsonData[0] == '"' {  
		jsonData = jsonData[1:]
		for jsonData[index] != '"' { //looping over till the end of string
			index++
		}
		return json{String: jsonData[:index], Parsed: true, Type: "String"}, jsonData[index+1:]
	}
	return result, jsonData
}


func numberParser(jsonData string) (json, string) {
	//Parses numbers in jsonstring and returns the number and remaining jsonData
	jsonData = strings.Trim(jsonData, " ")
	var number json
	var result [][]int
	re := regexp.MustCompile("^[-+]?[0-9]*.?[0-9]+([eE][-+]?[0-9]+)?") //init regex
	result1 := re.FindAllString(jsonData, -1)                          //Library function to find occurrence of string and returns a slice of successive elements
	if len(result1) > 0 {
		result = re.FindAllStringIndex(jsonData, -1) //Gives start and end index of the string
		index := result[0][1] //index is the end  index of the number
		j,err := strconv.Atoi(result1[0])
		if err == nil{
			return json{Number: j, Parsed: true, Type: "Number"}, jsonData[index:]
		}
		i,err := strconv.ParseFloat(result1[0],64)
		if err == nil {
			return json{Float: i,Parsed: true, Type: "Float"}, jsonData[index:]
		}
	} else {
		return number, jsonData
	}
	return number,jsonData
}

func arrayParser(jsonData string) (json,string) {
	//Parses json arrays and stores them into the user defined type else returns the data
	var result json
	parsed := make([]json,0)
	if jsonData[0] == '[' {
		jsonData = jsonData[1:]
		for len(jsonData) > 0 {
			result, jsonData = elementParser(jsonData)
			if result.Parsed == true {
				parsed = append(parsed,result)
				jsonData = commaParser(jsonData)
			}
			if jsonData[0] == ']'{
				fmt.Println("End of array")
				return json{Array: parsed, Parsed: true, Type: "Array"},jsonData[1:]
			}
		}
	}
	return result,jsonData
}

func objParser(jsonData string) (json,string) {
	//Parses json arrays and stores them into the user defined type else returns the data
	var result json
	var key string
	parsed := make(map[string]json)
	if jsonData[0] == '{' {
		jsonData = jsonData[1:]
		for len(jsonData) > 0 {
			result,jsonData = stringParser(jsonData)
			if result.Parsed == true {
				key = result.String
			}
			jsonData = colonParser(jsonData)
			result, jsonData = elementParser(jsonData)
			if result.Parsed == true {
				parsed[key] = result
				jsonData = commaParser(jsonData)
			}
			if jsonData[0] == '}'{
				fmt.Println("End of Object")
				return json{Object: parsed, Parsed: true, Type: "Object"},jsonData[1:]
			}
		}
	}
	return result,jsonData
}

func colonParser(jsonData string) string {
	//Parses colon in json objects and returns the colon and the remaining jsonData
	jsonData = strings.Trim(jsonData, " ")
	if jsonData[0] == ':' {
		return jsonData[1:]
	}
	return jsonData
}

func commaParser(jsonData string) string {
	//Parses comma which saperate elements and returns comma and the remaining jsonData
	jsonData = strings.Trim(jsonData, " ")
	if jsonData[0] == ',' {
		return jsonData[1:]
	}
	return jsonData
}

func booleanParser(jsonData string) (json, string) {
	// Function parses boolean elements and returns the bool value and the remaining jsondata
	var result json
	jsonData = strings.Trim(jsonData, " ")
	if len(jsonData) > 4 {
		if jsonData[0:4] == "true" {
			return json{Bool: true, Parsed: true, Type: "Bool"}, jsonData[4:] //slicing the jsonData
		} else if jsonData[0:5] == "false" {
			return json{Bool: false, Parsed: true ,Type:"Bool"}, jsonData[5:]
		}
	}
	return result, jsonData
}

func nullParser(jsonData string)  (json ,string) {
	//Parses null values form jsonstring and returns the nil value and remaining jsonData
	jsonData = strings.Trim(jsonData, " ")
	var result json
	if len(jsonData) > 4 {
		if jsonData[0:4] == "null" {
			return json{Null:nil, Parsed: true, Type: "Null"}, jsonData[4:]
		}
	}
	return result, jsonData
}


func elementParser(jsonData string) (json, string) {
	//Function tries all the parsers one by one on each element and returns the result and the remaining jsonData
	var result json
	result, jsonData = stringParser(jsonData)
	if result.Parsed == true {
		return result, jsonData
	}
	result, jsonData = numberParser(jsonData)
	if result.Parsed == true {
		return result, jsonData
	}
	
	result, jsonData = booleanParser(jsonData)
	if result.Parsed == true {
		return result, jsonData
	}
	result, jsonData = nullParser(jsonData)
	if result.Parsed == true {
		return result, jsonData
	}
	result, jsonData = arrayParser(jsonData)
	if result.Parsed == true {
		return result,jsonData
	}
	result, jsonData = objParser(jsonData)
	if result.Parsed == true {
		return result,jsonData
	}
	return result, jsonData
}

func (m json) getElement() interface{} {
	Type := m.Type
	switch Type {
		case "String":
			return m.String
		case "Object":
			return m.Object
		case "Number":
			return m.Number
		case "Array":
			return m.Array
		case "Null":
			return m.Null
		case "Float":
			return m.Float
		case "Bool":
			return m.Bool
		default:
			fmt.Println("Don't know how to handle the type")
			os.Exit(1)
	}
	return ""
}

func main(){
	buf, err := ioutil.ReadFile("test.txt") //Reading the file completely
	if err != nil {                          //error check
		fmt.Println(err)
		os.Exit(1)
	}
	data := string(buf) //bytes to string
	if data[0] == '['{
		result,_ :=arrayParser(data)
		final := result.getElement()
		fmt.Println(final)          // we can extract the parsed data using getElement() if we know the structure
	
	}else if data[0] =='{' {
		result,_ :=objParser(data)
		fmt.Println(result)
		k := result.getElement()
		fmt.Println(k)
		/*m := k.(map[string]json) // if we know the structure of json we can extract using the getElement()
		for a,b := range m{
			fmt.Println("Key:",a,",Value:",b.getElement())

		}*/
		
	}
	os.Exit(0)
}
