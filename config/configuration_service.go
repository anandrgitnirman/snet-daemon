package config

import (
	"encoding/json"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"strings"
)

type Configuration_Service struct {
	address string
}

//TODO
const default_daemon_configuration = `
{
  "registry_address_key": {
    "mandatory": false,
    "value": "",
    "description": "domain name for which ....)",
    "type": "string",
    "editable": true,
    "restart_daemon": false,
    "section": "blockchain"
  }
,
  "ethereum_json_rpc_endpoint": {
    "mandatory": false,
    "value": "https://ropsten.infura.io",
    "description": "domain name for which ....)",
    "type": "string",
    "editable": true,
    "restart_daemon": false,
    "section": "blockchain"
  },

  "blockchain_network_selected": {
    "mandatory": false,
    "value": "ropsten",
    "description": "domain name for which ....)",
    "type": "string",
    "editable": true,
    "restart_daemon": false,
    "section": "general"
  },
  "passthrough_enabled":{
    "mandatory": false,
    "value": true,
    "description": "domain name for which ....)",
    "type": "string",
    "editable": true,
    "restart_daemon": false,
    "section": "general"
  },
  "passthrough_endpoint": {
    "mandatory": false,
    "value": "http://127.0.0.1:7003",
    "description": "domain name for which ....)",
    "type": "string",
    "editable": true,
    "restart_daemon": false,
    "section": "general"
  },
  "daemon_end_point": {
    "mandatory": false,
    "value":"localhost:8086",
    "description": "domain name for which ....)",
    "type": "string",
    "editable": true,
    "restart_daemon": false,
    "section": "general"
  },
  "ipfs_end_point":{
    "mandatory": false,
    "value":"http://ipfs.singularitynet.io:80",
    "description": "domain name for which ....)",
    "type": "string",
    "editable": true,
    "restart_daemon": false,
    "section": "blockchain"
  },
  "organization_id": {
    "mandatory": false,
    "value": "test-snet",
    "description": "domain name for which ....)",
    "type": "string",
    "editable": true,
    "restart_daemon": false,
    "section": "blockchain"
  },
  "service_id": {
    "mandatory": false,
    "value":"test-example-service-7",
    "description": "domain name for which ....)",
    "type": true,
    "editable": true,
    "restart_daemon": false,
    "section": "blockchain"
  },
  "blockchain_enabled": {
    "mandatory": false,
    "value": "",
    "description": "domain name for which ....)",
    "type": "string",
    "editable": true,
    "restart_daemon": false,
    "section": "blockchain"
  },
  "burst_size":{
    "mandatory": false,
    "value": 100,
    "description": "domain name for which ....)",
    "type": "string",
    "editable": true,
    "restart_daemon": false,
    "section": "rateLimit"
  },
  "rate_limit_per_minute": {
    "mandatory": false,
    "value": 100,
    "description": "domain name for which ....)",
    "type": "string",
    "editable": true,
    "restart_daemon": false,
    "section": "rateLimit"
  },

  "log": {
    "level": {
      "mandatory": false,
      "value":"debugTESTING",
      "description": "domain name for which ....)",
      "type": "string",
      "editable": true,
      "restart_daemon": false,
      "section": "certificate"
    },
    "output": {
      "current_link":  {
        "mandatory": false,
        "value": "./daemonTEST.log",
        "description": "domain name for which ....)",
        "type": "string",
        "editable": true,
        "restart_daemon": false,
        "section": "certificate"
      },
      "file_pattern": {
        "mandatory": false,
        "value": "./daemon.%Y%m%d.log",
        "description": "domain name for which ....)",
        "type": "string",
        "editable": true,
        "restart_daemon": false,
        "section": "certificate"
      },
      "rotation_count": {
        "mandatory": false,
        "value": 0,
        "description": "domain name for which ....)",
        "type": "string",
        "editable": true,
        "restart_daemon": false,
        "section": "certificate"
      },
      "rotation_time_in_sec": {
        "mandatory": false,
        "value": "",
        "description": "domain name for which ....)",
        "type": "string",
        "editable": true,
        "restart_daemon": false,
        "section": "certificate"
      },
      "type":{
        "mandatory": false,
        "value": "file",
        "description": "domain name for which ....)",
        "type": "string",
        "editable": true,
        "restart_daemon": false,
        "section": "certificate"
      }
    }
  }
}`

func (service *Configuration_Service) GetConfiguration(context.Context, *ReadRequest) (*ConfigurationResponse, error) {
	response := &ConfigurationResponse{}
	defaultConfig := viper.New()
	ReadConfigFromJsonString(defaultConfig,default_daemon_configuration)
	response.JsonData = mergeJSON(defaultConfig,vip)
	return response, nil
}

func (service *Configuration_Service) UpdateConfiguration(context.Context, *UpdateRequest) (*ConfigurationResponse, error) {
	panic("implement me")
}

func (service *Configuration_Service) StopProcessingRequests(context.Context, *CommandRequest) (*Response, error) {
	panic("implement me")
}

func (service *Configuration_Service) StartProcessingRequests(context.Context, *CommandRequest) (*Response, error) {
	panic("implement me")
}

func NewConfigurationService() *Configuration_Service {
	return &Configuration_Service{}
}


//Used to map the attribute values to a struct
type Configuration_Details struct {
	Mandatory      string `json:"mandatory"`
	Value          string `json:"value"`
	Description    string `json:"description"`
	Type           string `json:"type"`
	Editable       string `json:"editable"`
	Restart_daemon string `json:"restart_daemon"`
	Section        string `json:"section"`
}

func getLeafNodeKey(key string) (bool,string) {
	if strings.Contains(key,".restart_daemon") {
		newKey := strings.Replace(key, ".restart_daemon", "",-1)
		return true, newKey
	}
	return false,""
}


func mergeJSON(defaultJSON *viper.Viper,currentDaemonConfig *viper.Viper) string {
	//For all the keys in the default Daemon Configuration
	for _, key := range defaultJSON.AllKeys() {
		//Find out if the given key is the key of a Leaf or not .
		if isLeaf,leafKey := getLeafNodeKey(key); isLeaf{
			//Get the value of this key from the current Daemon configuration if available
			if  currentDaemonConfig.Get(leafKey) != nil  {
				//Replace the value attribute from default to what has been set up in the configuration file.
				configurationDetailsJSON, _ := ConvertStructToJSON(defaultJSON.Get(leafKey))
				configDetails := &Configuration_Details{}
				_ = json.Unmarshal(configurationDetailsJSON, configDetails)
				configDetails.Value = currentDaemonConfig.GetString(leafKey)
				defaultJSON.Set(leafKey, configDetails)
			}
		}
	}
	data, _ := json.Marshal(defaultJSON.AllSettings())
	return string(data)
}

//convert the given struct to its corresponding json.
func ConvertStructToJSON(payLoad interface{}) ([]byte, error) {

	b, err := json.Marshal(&payLoad)
	if err != nil {

		return nil, err
	}
	return b, nil
}
