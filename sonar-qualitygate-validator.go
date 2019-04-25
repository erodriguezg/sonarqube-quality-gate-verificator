package main

import "fmt"
import "strings"
import "os"
import "errors"
import b64 "encoding/base64"

// SONAR URL CONSTANTS

const API_URL = "api/project_branches/list"
const PROJECT_URL_PARAM = "project"

// APP PARAMS CONSTANTS

const APP_PARAM_SONAR_URL = "sonarUrl"
const APP_PARAM_TOKEN = "token"
const APP_PARAM_PROJECT_KEY = "projectKey"
const APP_PARAM_BRANCH_NAME = "branchName"

// Params Object

type Params struct {
	sonarUrl string
	token string
	projectKey string 
	branchName string
}

func (p *Params) init(args []string) error {
	if len(args) < 9 {
		return errors.New("the required parameters do not match")
	}

	fmt.Println("1")

	for i, arg := range args {
		if( i == 0) {
			continue
		}

		cleanArg := strings.Replace(arg, "-", "", 1)

		switch cleanArg {
		case APP_PARAM_SONAR_URL:
			p.sonarUrl = args[i+1]
		case APP_PARAM_TOKEN:
			fmt.Println("2")
			p.token = args[i+1]
		case APP_PARAM_PROJECT_KEY:
			p.projectKey = args[i+1]
		case APP_PARAM_BRANCH_NAME:
			p.branchName = args[i+1]
		}
	}

	fmt.Println("3")
	return nil;
}

func (p *Params) getAuthHeaderVal() string {
	plainText := fmt.Sprintf("%s:", p.token)
	return b64.StdEncoding.EncodeToString([]byte(plainText))
}

// End Params Object

func modoDeUso() string {
	return fmt.Sprintf("How to use: -%s $%s -%s $%s -%s $%s -%s $%s ", 
		APP_PARAM_SONAR_URL, APP_PARAM_SONAR_URL,
		APP_PARAM_TOKEN, APP_PARAM_TOKEN, 
		APP_PARAM_PROJECT_KEY, APP_PARAM_PROJECT_KEY, 
		APP_PARAM_BRANCH_NAME, APP_PARAM_BRANCH_NAME)
}


func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		fmt.Fprintf(os.Stderr, modoDeUso())
        os.Exit(1)
	}
	os.Exit(0)
}

func run(args []string) error {

	params := new(Params)
	errorInitParam := params.init(args)

	if(errorInitParam != nil) {
		return errorInitParam
	}





	name := ""
	if len(args) > 1 {
		name = args[1]
	} else {
		return errors.New("No viene el parametro")
	}

	fmt.Println("hello world", name)
	fmt.Println("token", params.token, "b64", params.getAuthHeaderVal())
	return nil
}