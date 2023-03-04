package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", upload)
	http.HandleFunc("/get/children/", func(w http.ResponseWriter, r *http.Request) {
		// Form：存储了post、put和get参数，在使用之前需要调用ParseForm方法。
		// PostForm：存储了post、put参数，在使用之前需要调用ParseForm方法。
		// MultipartForm：存储了包含了文件上传的表单的post参数，在使用前需要调用ParseMultipartForm方法
		fmt.Println(r.ParseForm())
		fmt.Println("req", r.PostForm.Get("file"))
		w.Write(getNodes("./static/upload/" + r.PostForm.Get("file")))
		// w.Write([]byte(`[{"id":"0","parent":"#","state":{"disabled":false,"opened":true,"selected":false},"text":"树形结构"},{"id":"69","parent":"0","text":"工程"},{"id":"5","parent":"0","text":"行政"},{"id":"71","parent":"0","text":"迷"},{"id":"1","parent":"0","text":"技术"}]`))
	})

	http.HandleFunc("/save/children/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("req")
		fmt.Println(*r)

		err := r.ParseForm()
		if err != nil {
			log.Fatal("parse form error ", err)
		}
		// 初始化请求变量结构
		formData := make(map[string]interface{})
		// 调用json包的解析，解析请求body
		json.NewDecoder(r.Body).Decode(&formData)
		for key, value := range formData {
			log.Println("key:", key, " => value :", value)
		}

		w.Write([]byte(`[{"id":"0","parent":"#","state":{"disabled":false,"opened":true,"selected":false},"text":"树形结构"},{"id":"69","parent":"0","text":"保存后的数据","a_attr":{"style":"color:red"}},{"id":"5","parent":"0","text":"保存后的数据刷新"},{"id":"71","parent":"0","text":"展示用：可自定义"},{"id":"1","parent":"0","text":"自定义返回数据"}]`))
	})
	// Continue to process new requests until an error occurs
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func http_resp(code int, msg string, w http.ResponseWriter) {
	var Result map[string]interface{}
	Result = make(map[string]interface{})

	Result["code"] = code
	Result["msg"] = msg

	data, err := json.Marshal(Result)
	if err != nil {
		log.Printf("%v\n", err)
	}

	w.Write([]byte(string(data)))
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./static/upload.html")
		t.Execute(w, nil)
		return
	}

	contentType := r.Header.Get("content-type")
	contentLen := r.ContentLength

	if !strings.Contains(contentType, "multipart/form-data") {
		http_resp(-1001, "The content-type must be multipart/form-data", w)
		return
	}
	//限制最大文件大小
	if contentLen >= 50*1024*1024 {
		http_resp(-1002, "Failed to large,limit 50MB", w)
		return
	}

	err := r.ParseMultipartForm(50 * 1024 * 1024)
	if err != nil {
		http_resp(-1003, "Failed to ParseMultipartForm", w)
		return
	}

	if len(r.MultipartForm.File) == 0 {
		http_resp(-1004, "File is NULL", w)
		return
	}

	var Result map[string]interface{}
	Result = make(map[string]interface{})

	// Result["code"] = 0
	DownLoadUrl := "http://127.0.0.1:8080/static/upload/"

	FileCount := 0
	for _, headers := range r.MultipartForm.File {
		for _, header := range headers {
			log.Printf("recv file: %s\n", header.Filename)

			filePath := filepath.Join("./static/upload", header.Filename)
			dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				log.Printf("Create file %s error: %s\n", filePath, err)
				return
			}
			srcFile, err := header.Open()
			if err != nil {
				log.Printf("Open header failed: %s\n", err)
			}
			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				log.Printf("Write file %s error: %s\n", filePath, err)
			}
			_, _ = srcFile.Close(), dstFile.Close()
			FileCount++
			// name := fmt.Sprintf("file%d_url", FileCount)
			name := header.Filename
			Result[name] = (DownLoadUrl + header.Filename)
		}
	}
	// data, err := json.Marshal(Result)
	// if err != nil {
	// 	log.Printf("%v \n", err)
	// }
	// w.Write([]byte(string(data)))

	vv := []struct {
		Name string
		Url  string
	}{}
	for k, v := range Result {
		vv = append(vv, struct {
			Name string
			Url  string
		}{
			Name: k,
			Url:  v.(string),
		})
	}

	t, err := template.ParseFiles("./static/parse.html")
	fmt.Println(err, vv)
	t.Execute(os.Stdout, vv)
	t.Execute(w, vv)
}
