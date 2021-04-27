// jsonHandler.go

package jsonHandler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
)

// JsonRead: datas from file to given interface / structure
// i.e: err := ReadJson(filename, &person)
// remember to put upper char at left of variables names to be saved.
func JsonRead(filename string, interf interface{}) (err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
		err = json.Unmarshal(textFileBytes, &interf)
	}
	return err
}

// JsonWrite: datas to file from given interface / structure
// i.e: err := WriteJson(filename, &person)
// remember to put upper char at left of variables names to be saved.
func JsonWrite(filename string, interf interface{}) (err error) {
	var out bytes.Buffer
	var jsonData []byte
	if jsonData, err = json.Marshal(&interf); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			if err = ioutil.WriteFile(filename, out.Bytes(), os.ModePerm); err == nil {
			}
		}
	}
	return err
}
