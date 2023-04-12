package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"testing"
)

type StructInfo struct {
	Name      string `json:"name"`
	Comment   string `json:"comment"`
	ModelPath string `json:"modelPath"`
	ModelName string `json:"modelName"`
}

type FuncInfo struct {
	Name      string         `json:"name"`
	Comment   string         `json:"comment"`
	ModelPath string         `json:"modelPath"`
	ModelName string         `json:"modelName"`
	Params    []*FuncVarInfo `json:"params"`
	Return    *FuncVarInfo   `json:"return"`
	HasError  bool           `json:"hasError"`
	Func      interface{}    `json:"-"`
}

type FuncVarInfo struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Type    string `json:"type"`
}

func TestGen(t *testing.T) {
	var srcDir = `C:\Program Files\Go\src`
	fmt.Println("srcDir:", srcDir)
	genGoSrc(srcDir)
}

func genGoSrc(srcDirPath string) {

	var funcInfoList []*FuncInfo
	var structInfoList []*StructInfo
	loadGoSrc(srcDirPath, "", &funcInfoList, &structInfoList)
	fmt.Println("all funcInfoList:", len(funcInfoList))
	fmt.Println("all structInfoList:", len(structInfoList))
	_ = os.MkdirAll("./packages", os.ModePerm)

	var imports []string
	var modelPathName = map[string]string{}
	var importNames []string
	addImport := func(modelPath string) (name string) {
		if StringIndexOf(imports, modelPath) >= 0 {
			name = modelPathName[modelPath]
			return
		}
		imports = append(imports, modelPath)
		name = modelPath
		if strings.Contains(name, "/") {
			name = modelPath[strings.Index(name, "/")+1:]
		}
		var name_ = name
		var index int
		for StringIndexOf(importNames, name) >= 0 {
			index++
			name = fmt.Sprintf("%s%d", name_, index)
		}
		importNames = append(importNames, name)
		modelPathName[modelPath] = name
		return
	}
	for _, one := range funcInfoList {
		name := addImport(one.ModelPath)
		one.ModelName = name
	}
	for _, one := range structInfoList {
		name := addImport(one.ModelPath)
		one.ModelName = name
	}

	genContent := `package packages

import (
`
	sort.Strings(imports)
	for _, modelPath := range imports {
		name := modelPathName[modelPath]
		if strings.HasSuffix(modelPath, "/"+name) {
			//genContent += "\t" + `"` + modelPath + `"` + "\n"
			//genContent += "\t" + `_ "` + modelPath + `"` + "\n"
		} else {
			//genContent += "\t" + name + ` "` + modelPath + `"` + "\n"
			//genContent += "\t" + `_ "` + modelPath + `"` + "\n"
		}
	}
	genContent += `)` + "\n"

	genContent += `func init(){

`
	//for _, one := range structInfoList {
	//	if one.Comment != "" {
	//		genContent += "\t" + `// ` + one.Name + ` ` + one.Comment + "\n"
	//	}
	//	genContent += "\t" + `` + one.ModelName + `.` + one.Name + "\n"
	//}
	//for _, one := range funcInfoList {
	//	if one.Comment != "" {
	//		genContent += "\t" + `// ` + one.Name + ` ` + one.Comment + "\n"
	//	}
	//	genContent += "\t" + `` + one.ModelName + `.` + one.Name + "\n"
	//}
	genContent += `}` + "\n"

	f, err := os.Create("./packages/golang.go")
	if err != nil {
		panic("os.Create error:" + err.Error())
	}
	defer func() { _ = f.Close() }()
	_, _ = f.WriteString(genContent)
}

func loadGoSrc(srcDirPath string, modelPath string, funcInfoList *[]*FuncInfo, structInfoList *[]*StructInfo) {
	ds, err := os.ReadDir(srcDirPath)
	if err != nil {
		panic("ReadDir error:" + err.Error())
	}

	for _, d := range ds {
		if d.IsDir() {
			if d.Name() == "test" {
				continue
			}
			if d.Name() == "testdata" {
				continue
			}
			if d.Name() == "internal" {
				continue
			}
			if d.Name() == "vendor" {
				continue
			}

			subModelPath := d.Name()
			if modelPath != "" {
				subModelPath = modelPath + "/" + subModelPath
			}
			loadGoSrc(fmt.Sprintf("%s/%s", srcDirPath, d.Name()), subModelPath, funcInfoList, structInfoList)
		} else {
			if !strings.HasSuffix(d.Name(), ".go") {
				continue
			}
			if strings.HasSuffix(d.Name(), "_test.go") {
				continue
			}
			funcList, structList := loadGoSrcFile(fmt.Sprintf("%s/%s", srcDirPath, d.Name()), modelPath)
			*funcInfoList = append(*funcInfoList, funcList...)
			*structInfoList = append(*structInfoList, structList...)
		}
	}

}
func loadGoSrcFile(filePath string, modelPath string) (funcInfoList []*FuncInfo, structInfoList []*StructInfo) {
	var lines []string
	lines, err := ReadLine(filePath)
	if err != nil {
		panic("genGoFile " + filePath + " ReadLine error:" + err.Error())
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "//go:build") {
			return
		}
		if strings.HasPrefix(line, "package main") {
			return
		}

	}
	for row, line := range lines {
		if strings.HasPrefix(line, "func ") && strings.Contains(line, "(") {
			if row == 0 {
				continue
			}
			funcName := line[len("func "):strings.Index(line, "(")]
			if funcName == "" {
				continue
			}
			var commandLines []string
			var lastComment string
			var i = row - 1
			for {
				if !strings.HasPrefix(lines[i], "//") {
					break
				}
				lastComment = lines[i]
				commandLines = append(commandLines, lastComment)
				i--
			}
			vv := []rune(funcName)
			if vv[0] >= 97 && vv[0] <= 122 {
				continue
			}
			var fS = "// " + funcName + " "
			comment := ""
			if strings.HasPrefix(lastComment, fS) {
				comment = lastComment[len(fS):]
			}
			funcInfo := &FuncInfo{
				Name:      funcName,
				Comment:   comment,
				ModelPath: modelPath,
			}
			//for i = len(commandLines) - 1; i >= 0; i-- {
			//fmt.Println(commandLines[i])
			//}
			//fmt.Println("funcName", funcName)
			funcInfoList = append(funcInfoList, funcInfo)
		} else if strings.HasPrefix(line, "type ") && strings.Contains(line, "{") {
			if row == 0 {
				continue
			}
			structName := line[len("type "):strings.Index(line, "{")]
			if structName == "" {
				continue
			}
			var commandLines []string
			var lastComment string
			var i = row - 1
			for {
				if !strings.HasPrefix(lines[i], "//") {
					break
				}
				lastComment = lines[i]
				commandLines = append(commandLines, lastComment)
				i--
			}
			vv := []rune(structName)
			if vv[0] >= 97 && vv[0] <= 122 {
				continue
			}
			var fS = "// " + structName + " "
			comment := ""
			if strings.HasPrefix(lastComment, fS) {
				comment = lastComment[len(fS):]
			}
			structInfo := &StructInfo{
				Name:      structName,
				Comment:   comment,
				ModelPath: modelPath,
			}
			//for i = len(commandLines) - 1; i >= 0; i-- {
			//fmt.Println(commandLines[i])
			//}
			//fmt.Println("funcName", funcName)
			structInfoList = append(structInfoList, structInfo)
		}

	}

	return
}

// StringIndexOf 返回 某个值 在数组中的索引位置，未找到返回 -1
func StringIndexOf(array []string, v string) (index int) {
	index = -1
	size := len(array)
	for i := 0; i < size; i++ {
		if array[i] == v {
			index = i
			return
		}
	}
	return
}

// ReadLine 逐行读取文件
func ReadLine(filename string) (lines []string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	buf := bufio.NewReader(f)
	var line string
	for {
		line, err = buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF { //读取结束，会报EOF
				err = nil
				return
			}
			return nil, err
		}
		lines = append(lines, line)
	}
	return
}
