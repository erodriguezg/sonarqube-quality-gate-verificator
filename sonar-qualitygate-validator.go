package main

import "fmt"
import "strings"
import "os"
import "errors"
import "encoding/base64"
import "net/http"
import "encoding/json"

// SONAR URL CONSTANTS

const SONAR_RELATIVE_API_URL = "api/project_branches/list"
const PROJECT_URL_PARAM = "project"

// Sonar Response Struts 

type SonarBranchesResponse struct {
	Branches []SonarBranch `json:"branches"`
}

type SonarBranch struct {
	Name string `json:"name"`
	IsMain bool `json:"isMain"`
	TypeVal string `json:"type"`
	MergeBranch string `json:"mergeBranch"`
	Status SonarBranchStatus `json:"status"`
	AnalysisDate string `json:"analysisDate"`
}

type SonarBranchStatus struct {
	QualityGateStatus string `json:"qualityGateStatus"`
	Bugs int `json:"bugs"`
	Vulnerabilities int `json:"vulnerabilities"`
	CodeSmells int `json:"codeSmells"`
}

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
	for i, arg := range args {
		if( i == 0) {
			continue
		}

		cleanArg := strings.Replace(arg, "-", "", 1)

		switch cleanArg {
		case APP_PARAM_SONAR_URL:
			p.sonarUrl = addEndSlashUrl(args[i+1])
		case APP_PARAM_TOKEN:
			p.token = args[i+1]
		case APP_PARAM_PROJECT_KEY:
			p.projectKey = args[i+1]
		case APP_PARAM_BRANCH_NAME:
			p.branchName = args[i+1]
		}
	}
	return nil
}

func addEndSlashUrl(url string) string {
	endChar := url[len(url)-1:]
	if endChar != "/" {
		return url + "/"
	} else {
		return url
	}
}


func (p *Params) getAuthHeaderVal() string {
	plainText := fmt.Sprintf("%s:", p.token)
	return base64.StdEncoding.EncodeToString([]byte(plainText))
}

// End Params Object

func showUsageHelp() string {
	return fmt.Sprintf("How to use: sonar-qualitygate-validator -%s <%s> -%s <%s> -%s <%s> -%s <%s> ", 
		APP_PARAM_SONAR_URL, APP_PARAM_SONAR_URL,
		APP_PARAM_TOKEN, APP_PARAM_TOKEN, 
		APP_PARAM_PROJECT_KEY, APP_PARAM_PROJECT_KEY, 
		APP_PARAM_BRANCH_NAME, APP_PARAM_BRANCH_NAME)
}


func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
        os.Exit(1)
	}
	os.Exit(0)
}

func run(args []string) error {

	params := new(Params)
	errorInitParam := params.init(args)

	if(errorInitParam != nil) {
		fmt.Fprintf(os.Stderr, "%s\n", showUsageHelp())
		return errorInitParam
	}

	sonarResp := new(SonarBranchesResponse)
	errorSonar := querySonarQube(params, sonarResp);
	if errorSonar != nil {
		return errorSonar
	}
	
	fmt.Println(sonarResp)

	if len(sonarResp.Branches) == 0 {
		return errors.New("0 branches found")
	}

	branchFound := false;
	for _, branch := range sonarResp.Branches {

		if params.branchName != branch.Name {
			continue
		} 

		branchFound = true;
		if "OK" == branch.Status.QualityGateStatus {
			fmt.Println("Quality Gates PASSED!")
			return nil
		}

		break
	}

	if branchFound {
		return errors.New("Quality Gate Failed!")
	} else {
		return errors.New(fmt.Sprintf("Branch '%s' not found", params.branchName))
	}

}

func querySonarQube(params *Params, target interface{}) error {

	req, errorSonar := createSonarRequest(params)

	if errorSonar != nil {
		return errorSonar
	}

	// Send req using http Client

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

func createSonarRequest(params *Params) (*http.Request, error) {
	url := fmt.Sprintf("%s%s?%s=%s", 
		params.sonarUrl, 
		SONAR_RELATIVE_API_URL, 
		PROJECT_URL_PARAM,
		params.projectKey)

	authHeader := fmt.Sprintf("Basic %s", params.getAuthHeaderVal())

    // Create a new request using http
    req, err := http.NewRequest("GET", url, nil)

    // add authorization header to the req
	req.Header.Add("Authorization", authHeader)
	
	return req, err
}
