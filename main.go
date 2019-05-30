package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

var (
	configmsg = `	gup config <host> <username> <password>
	- Initial setup required to make the program work.
	  The details are saved in a configuration file
	  in your home directory.`
	uploadmsg = `	gup <file>
	- Uploads the supplied file.`
	deletemsg = `	gup delete <url>
	- Deletes the file at the supplied url.`
	helpmsg = `	gup help | -h | --help
	- Shows this message.`
	usagemsg = fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s\n",
		configmsg, uploadmsg, deletemsg, helpmsg)
)

func main() {
	printMsg(usagemsg, 2)

	switch os.Args[1] {
	case "delete":
		printMsg(deletemsg, 3)
		deleteFile(os.Args[2])
	case "config":
		printMsg(configmsg, 5)
		writeConfig(os.Args[2], os.Args[3], os.Args[4])
	case "help", "-h", "--help":
		fmt.Println(usagemsg)
	default:
		uploadFile(os.Args[1])
	}
}

func printMsg(msg string, argLen int) {
	if len(os.Args) < argLen {
		fmt.Printf("Usage:\n%s\n", msg)
		os.Exit(0)
	}
}

func uploadFile(file string) {
	f, err := os.Open(file)
	check(err)
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	check(err)

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	ff, err := writer.CreateFormFile("file", f.Name())
	check(err)
	ff.Write(content)
	err = writer.Close()
	check(err)

	config := parseConfig()
	req, err := http.NewRequest("POST", config[0], &b)
	check(err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config[1]))

	client := &http.Client{}
	res, err := client.Do(req)
	check(err)

	resContent, err := ioutil.ReadAll(res.Body)
	check(err)
	defer res.Body.Close()
	fmt.Println(string(resContent))
}

func deleteFile(url string) {
	client := &http.Client{}
	config := parseConfig()
	req, err := http.NewRequest("DELETE", url, nil)
	check(err)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config[1]))

	res, err := client.Do(req)
	check(err)
	resContent, err := ioutil.ReadAll(res.Body)
	check(err)
	defer res.Body.Close()
	fmt.Println(string(resContent))
}

func writeConfig(host, user, pass string) {
	res, err := http.PostForm(host+"/login",
		url.Values{"user": {user}, "pass": {pass}})
	check(err)

	resContent, err := ioutil.ReadAll(res.Body)
	check(err)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", string(resContent))
		os.Exit(1)
	}

	home, err := homedir.Dir()
	check(err)

	path := filepath.Join(home, ".gupcfg")
	content := fmt.Sprintf("%s %s", host, resContent)

	err = ioutil.WriteFile(path, []byte(content), 0644)
	check(err)
	fmt.Println("Successfully saved configuration!")
}

func parseConfig() []string {
	home, err := homedir.Dir()
	check(err)

	content, err := ioutil.ReadFile(filepath.Join(home, ".gupcfg"))
	check(err)

	return strings.Split(string(content), " ")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
