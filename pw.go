package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type PasswordsEntry struct {
	Key         string
	Password    string
	Description string
}

type Passwords struct {
	Entries []*PasswordsEntry
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func readPwFile() *Passwords {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(usr.HomeDir, ".pw.json")

	if !fileExists(path) {
		return &Passwords{}
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading file %v: %v", path, err)
	}
	passwords := &Passwords{}
	err = json.Unmarshal(file, passwords)
	if err != nil {
		log.Fatal(err)
	}
	return passwords
}

func writePwFile(passwords *Passwords) {
	jsonContent, err := json.MarshalIndent(passwords, "", "  ")
	if err != nil {
		log.Fatal(err)
		return
	}
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(usr.HomeDir, ".pw.json")
	err = ioutil.WriteFile(path, jsonContent, 0700)
	if err != nil {
		log.Fatal(err)
	}
}

func setPassword(key, password, description string, passwords *Passwords) {

	var entry *PasswordsEntry

	// entry exists, overwrite existing password and description
	if passwords != nil && passwords.Entries != nil {
		for i := range passwords.Entries {
			if strings.Contains(passwords.Entries[i].Key, key) {
				entry = passwords.Entries[i]
			}
		}
	}

	if entry != nil {

		// entry exists, overwrite existing password and description
		entry.Password = password
		entry.Description = description

	} else {

		// entry does not exist, create it
		entry = &PasswordsEntry{Key: key, Password: password, Description: description}
		if passwords.Entries == nil {
			passwords.Entries = make([]*PasswordsEntry, 0, 0)
		}
		passwords.Entries = append(passwords.Entries, entry)

	}
}

func printPassword(key string, passwords *Passwords) int {
	if passwords != nil && passwords.Entries != nil && len(passwords.Entries) > 0 {
		for i := range passwords.Entries {
			if strings.Contains(passwords.Entries[i].Key, key) {
				fmt.Printf(passwords.Entries[i].Password)
				return 0
			}
		}
		return 1
	}
	return 1
}

func keys(passwords *Passwords) []string {
	keys := make([]string, len(passwords.Entries))
	for i, e := range passwords.Entries {
		keys[i] = e.Key
	}
	sort.Strings(keys)
	return keys
}

func asMap(passwords *Passwords) map[string]*PasswordsEntry {
	m := make(map[string]*PasswordsEntry, len(passwords.Entries))
	for _, e := range passwords.Entries {
		m[(*e).Key] = e
	}
	return m
}

func printKeyList(passwords *Passwords) {
	if passwords != nil && passwords.Entries != nil && len(passwords.Entries) > 0 {
		maxLen := 0
		for _, entry := range passwords.Entries {
			if len(entry.Key) > maxLen {
				maxLen = len(entry.Key)
			}
		}
		keys := keys(passwords)
		pMap := asMap(passwords)
		fmt.Printf("%v\n\n", keys)
		fmt.Printf("%v\n\n", pMap)
		for _, key := range keys {
			fmt.Printf("%-"+strconv.Itoa(maxLen)+"v", key)
			if "" != (*pMap[key]).Description {
				fmt.Printf(" # %v", (*pMap[key]).Description)
			}
			fmt.Printf("\n")
		}
	}
}

type args struct {
	Key  string
	Pwd  string
	Desc string
}

func main() {

	if len(os.Args) == 1 || len(os.Args) > 4 {
		fmt.Printf("Usage: pw KEY     # prints the password for KEY")
		fmt.Printf("Usage: pw KEY PWD # sets the password for KEY")
		fmt.Printf("       pw --list  # prints available keys")
		return
	}

	passwords := readPwFile()
	retCode := 0

	if len(os.Args) == 2 {
		if os.Args[1] == "--list" {
			retCode = 0
			printKeyList(passwords)
		} else {
			retCode = printPassword(os.Args[1], passwords)
		}
	} else if len(os.Args) == 3 {
		setPassword(os.Args[1], os.Args[2], "", passwords)
		retCode = printPassword(os.Args[1], passwords)
		writePwFile(passwords)
	} else if len(os.Args) == 4 {
		setPassword(os.Args[1], os.Args[2], os.Args[3], passwords)
		retCode = printPassword(os.Args[1], passwords)
		writePwFile(passwords)
	}

	os.Exit(retCode)
}
