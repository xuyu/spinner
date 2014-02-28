package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"git.code4.in/spinner"
)

var (
	versionFile string
)

func init() {
	flag.StringVar(&versionFile, "version", "spinner-node.version", "spinner node version file")
}

func checkUpdate() (string, []string) {
	version := spinner.ReadVersion(versionFile)
	u := fmt.Sprintf("%s/spinner/central/checkupdate?version=%s&h=%s",
		centralSpinner, url.QueryEscape(version), url.QueryEscape(hostname))
	resp, err := http.Get(u)
	if err != nil {
		log.Println(err.Error())
		return "", nil
	}
	if resp.StatusCode != http.StatusOK {
		return "", nil
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return "", nil
	}
	defer resp.Body.Close()
	var result []string
	if err := json.Unmarshal(data, &result); err != nil {
		log.Println(err.Error())
		return "", nil
	}
	return result[0], result[1:]
}

func updateFile(filename string) error {
	s, err := spinner.FileMd5(filename)
	if err != nil && !os.IsNotExist(err) {
		log.Println(err.Error())
		return err
	}
	u := fmt.Sprintf("%s/spinner/central/update?file=%s&md5=%s&h=%s",
		centralSpinner, url.QueryEscape(filename), s, url.QueryEscape(hostname))
	resp, err := http.Get(u)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if resp.StatusCode == http.StatusNotModified {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer resp.Body.Close()
	perm, err := spinner.FilePerm(filename, 0644)
	if err != nil && !os.IsNotExist(err) {
		log.Println(err.Error())
		return err
	}
	if err := ioutil.WriteFile(filename, data, perm); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func update(newVersion string, files []string) {
	upgrade := true
	for _, filename := range files {
		if err := updateFile(filename); err != nil {
			upgrade = false
		}
	}
	if upgrade {
		perm, err := spinner.FilePerm(versionFile, 0644)
		if err != nil && !os.IsNotExist(err) {
			log.Println(err.Error())
			return
		}
		if err := ioutil.WriteFile(versionFile, []byte(newVersion), perm); err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func updateRunner() {
	for {
		newVersion, files := checkUpdate()
		if newVersion != "" && files != nil {
			update(newVersion, files)
		}
		time.Sleep(10 * time.Minute)
	}
}
