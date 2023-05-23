package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Host            string `yaml:"host"`
	User            string `yaml:"user"`
	Pw              string `yaml:"pw"`
	BaseDN          string `yaml:"basedn"`
	Filter          string `yaml:"filter"`
	TargetAttribute string `yaml:"target_attribute"`
}

func ldapFind(ldapURL string,
	user string,
	pw string,
	baseDN string,
	filter string,
	searchString string,
	targetAttribute string) (result string) {

	client, err := ldap.DialURL(ldapURL)
	if err != nil {
		log.Fatal(err)
		return result
	}

	err = client.Bind(user, pw)
	if err != nil {
		log.Fatal(err)
		return result
	}
	result, err = list(client, baseDN, filter, searchString, targetAttribute)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
	return result
}

func list(conn *ldap.Conn,
	baseDN string,
	filter string,
	searchString string,
	targetAttribute string) (string, error) {
	result, err := conn.Search(ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		createFilter(filter, searchString),
		[]string{targetAttribute},
		nil,
	))

	if err != nil {
		return "", fmt.Errorf("Failed to search users. %s", err)
	}

	res := ""
	for _, entry := range result.Entries {
		res = entry.GetAttributeValue(targetAttribute)
	}
	return res, nil
}

func createFilter(filter string, needle string) string {
	res := strings.Replace(
		filter,
		"{placeholder}",
		needle,
		-1,
	)

	return res
}

func main() {
	searchString := os.Args[1]
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		os.Exit(1)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		os.Exit(1)
	}

	result := ldapFind(config.Host,
		config.User,
		config.Pw,
		config.BaseDN,
		config.Filter,
		searchString,
		config.TargetAttribute)

	if result == "" {
		os.Exit(0)
	}

    /*
	user := strings.Split(searchString, "@")

	// Prevent recursion loop
	if result == user[0] {
		os.Exit(0)
	}
    */

	fmt.Printf("%s", result)
}
