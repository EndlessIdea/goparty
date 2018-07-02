package main

import (
	"bufio"
	"fmt"
	"errors"
	"log"
	"os"
	"strings"
	"os/exec"
	"flag"
)

var (
	sourceFile = "db.sql"
	targetFile = "dao.go"
	daoSuffix  = "Dao"
)

const (
	BOOL      = "BOOL"
	BOOLEAN   = "BOOLEAN"
	BIT       = "BIT"
	TINYINT   = "TINYINT"
	SMALLINT  = "SMALLINT"
	MEDIUMINT = "MEDIUMINT"
	INT       = "INT"
	INTEGER   = "INTEGER"
	BIGINT    = "BIGINT"
	DECIMAL   = "DECIMAL"
	DEC       = "DEC"
	FLOAT     = "FLOAT"
	DOUBLE    = "DOUBLE"
	DATE      = "DATE"
	DATETIME  = "DATETIME"
	TIMESTAMP = "TIMESTAMP"
	TIME      = "TIME"
	YEAR      = "YEAR"
	CHAR      = "CHAR"
	VARCHAR   = "VARCHAR"
	TEXT      = "TEXT"
)

func goNameStyle(name string) (ret string) {
	tokens := strings.Split(name, "_")
	for _, token := range tokens {
		ret += strings.Title(token)
	}
	return ret
}

func getDaoName(tableName string) string {
	daoName := goNameStyle(tableName)
	return daoName + daoSuffix
}

func getFieldType(typeInfo string) (fieldType string, err error) {
	schemaType := strings.ToUpper(strings.TrimSpace(typeInfo))
	typeEnd := strings.Index(typeInfo, "(")
	if typeEnd != -1 {
		schemaType = schemaType[0:typeEnd]
	}
	switch schemaType {
	case BOOL, BOOLEAN, BIT ,TINYINT ,SMALLINT:
		fieldType = "int16"
	case MEDIUMINT, INT, INTEGER:
		fieldType = "int32"
	case BIGINT:
		fieldType = "int64"
	case DECIMAL, DEC, FLOAT, DOUBLE:
		fieldType = "float64"
	case DATE, DATETIME, TIMESTAMP, TIME, YEAR, CHAR, VARCHAR, TEXT:
		fieldType = "string"
	default:
		err = errors.New(fmt.Sprintf("unknown field type: %s\n", typeInfo))
		return
	}
	return
}

func getDaoField(lineInfo []string) (daoField string, err error) {
	var fieldName, fieldType, fieldTag, fieldComment string
	if len(lineInfo) < 6 { //at least 6 sections, like: description text not null comment "描述"
		err = errors.New(fmt.Sprintf("scheme format error: %s", strings.Join(lineInfo, " ")))
		return
	}
	schemaName := strings.Trim(lineInfo[0], "`")
	fieldName = strings.Title(goNameStyle(schemaName))
	fieldType, err = getFieldType(lineInfo[1])
	if err != nil {
		return
	}
	fieldTag = fmt.Sprintf("`json:\"%s\"`", schemaName)
	fieldComment = strings.Trim(lineInfo[len(lineInfo)-1], "`")
	daoField = fmt.Sprintf("%s %s %s//%s\n", fieldName, fieldType, fieldTag, fieldComment)
	return
}

func main() {
	args := flag.String("s", sourceFile, "the origin schema file")
	argo := flag.String("o", targetFile, "the generated dao file")
	flag.Parse()
	sourceFile = *args
	targetFile = *argo
	fmt.Printf("generate %s depend on schema file %s\n", targetFile, sourceFile)

	sf, err := os.Open(sourceFile)
	if err != nil {
		log.Fatalf("open source file error: %v\n", err)
	}
	defer sf.Close()

	tf, err := os.Create(targetFile)
	if err != nil {
		log.Fatalf("create target file error: %v\n", err)
	}
	defer tf.Close()

	_, err = tf.WriteString("package dao\n\n")
	if err != nil {
		log.Fatalf("write target file error: %v\n", err)
	}
	scanner := bufio.NewScanner(sf)
	for scanner.Scan() {
		lineInfo := strings.Fields(strings.TrimSpace(scanner.Text()))
		if len(lineInfo) == 0 {
			continue
		}
		token := lineInfo[0]
		switch token {
		case "CREATE":
			tableName := strings.Trim(lineInfo[2], "`")
			daoName := getDaoName(tableName)
			if _, err = tf.WriteString(fmt.Sprintf("type %s struct {\n", daoName)); err != nil {
				log.Fatalf("write daoName error: %v", err)
			}
		case ")":
			if _, err := tf.WriteString("}\n\n"); err != nil {
				log.Fatalf("write error:%v", err)
			}
		default:
			if token[0:1] == "`" {
				daoField, err := getDaoField(lineInfo)
				if err != nil {
					log.Fatalf("get dao field error: %v", err)
				}
				if _, err := tf.WriteString(daoField); err != nil {
					log.Fatalf("write dao field error: %v", err)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("scan source file error: %v", err)
	}

	cmd := exec.Command("go", "fmt", targetFile)
	if err := cmd.Start(); err != nil {
		log.Fatalf("go fmt target file error: %v", err)
	}
	fmt.Println("generate success")

	type Person struct {
		Name string
		Friends map[string]string
	}
	Lilei := Person{"Lilei", map[string]string{"HanMei": "classmate"}}
	Mike := Lilei
	Mike.Friends["Michael"] = "superstar"
	fmt.Println(Lilei.Friends, Mike.Friends)
	var c chan struct{}
}
