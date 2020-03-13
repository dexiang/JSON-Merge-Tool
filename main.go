package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func main()  {

	var input string
	var output string
	var outputCSV bool

	flag.StringVar(&input, "input", "", "input file or folder")
	flag.StringVar(&output, "output", "", "output file")
	flag.BoolVar(&outputCSV, "csv", false, "output csv type file")
	flag.Parse()

	fi, err := os.Stat(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	var filePaths []string

	switch mode := fi.Mode(); {
	case mode.IsDir():
		filePaths, _ = GetOnlyJSONInDir(input)
	case mode.IsRegular():
		filePaths = append(filePaths, input)
	default:
		fmt.Println(err)
		return
	}

	var mergeResult []interface{}

	for _, filePath := range filePaths {
		fileContent, _ := ioutil.ReadFile(filePath)
		var fileContentInJSONFormat interface{}
		_ = json.Unmarshal(fileContent, &fileContentInJSONFormat)

		switch value := fileContentInJSONFormat.(type) {
		case string:
			fmt.Println("String")
			return
		case float64:
			fmt.Println("Int")
			return
		case []interface{}:
			for _, object := range value {
				mergeResult = append(mergeResult, object)
			}
		case map[string]interface{}:
			fmt.Println("Object")
			return
		default:
			fmt.Println("(unknown)")
			fmt.Println(value)
			return
		}
	}

	if outputCSV {
		f, err := os.Create(output)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()

		w := csv.NewWriter(f)

		// get Object Key
		var objectKeys []string
		for _, object := range mergeResult {
			switch object := object.(type) {
			case map[string]interface{}:
				for key, _ := range object {
					if ! StringInSlice(key, objectKeys) {
						objectKeys = append(objectKeys, key)
					}
				}
			default:
				fmt.Println("It is not Object.")
				return
			}
		}

		_ = w.Write(objectKeys)

		for _, object := range mergeResult {

			var record []string

			switch object := object.(type) {
			case map[string]interface{}:

				for _, objectKey := range objectKeys {
					switch value := object[objectKey].(type) {
					case string:
						record = append(record, value)
					case float64:
						record = append(record, strconv.FormatFloat(value, 'f', -1,64))
					case bool:
						record = append(record, strconv.FormatBool(value))
					default:
						fmt.Println("It is not String.")
						record = append(record, "")
					}
				}
			default:
				fmt.Println("It is not Object.")
				break
			}

			_ = w.Write(record)
		}
		w.Flush()

	} else {
		jsonString, _ := json.Marshal(mergeResult)
		_ = ioutil.WriteFile(output, jsonString, 0644)
	}
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetOnlyJSONInDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if filepath.Ext(info.Name()) == ".json" {
				files = append(files, path)
			}
		}
		return nil
	})
	return files, err
}