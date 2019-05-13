package config

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
	"testing"

	"golang.org/x/net/context"
)

func TestConfiguration_Service_GetConfiguration(t *testing.T) {
	defaultJSON := viper.New()
	_ = ReadConfigFromJsonString(defaultJSON, default_daemon_configuration)

	currentDeamonConfig := viper.New()
	_ = ReadConfigFromJsonString(currentDeamonConfig, defaultConfigJson)
	//For all the keys in the current Daemon Configuration
	for _, key := range defaultJSON.AllKeys() { //vip.AllSettings() {
	fmt.Println(key)
	}
fmt.Println("================================================================================================================")
	for _, key := range currentDeamonConfig.AllKeys() { //vip.AllSettings() {
		fmt.Println(key)
	}
}

func TestConfiguration_Service_UpdateConfiguration(t *testing.T) {
	type fields struct {
		address string
	}
	type args struct {
		in0 context.Context
		in1 *UpdateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ConfigurationResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Configuration_Service{
				address: tt.fields.address,
			}
			got, err := service.UpdateConfiguration(tt.args.in0, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configuration_Service.UpdateConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Configuration_Service.UpdateConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfiguration_Service_StopProcessingRequests(t *testing.T) {
	type fields struct {
		address string
	}
	type args struct {
		in0 context.Context
		in1 *CommandRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Configuration_Service{
				address: tt.fields.address,
			}
			got, err := service.StopProcessingRequests(tt.args.in0, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configuration_Service.StopProcessingRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Configuration_Service.StopProcessingRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfiguration_Service_StartProcessingRequests(t *testing.T) {
	type fields struct {
		address string
	}
	type args struct {
		in0 context.Context
		in1 *CommandRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Configuration_Service{
				address: tt.fields.address,
			}
			got, err := service.StartProcessingRequests(tt.args.in0, tt.args.in1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Configuration_Service.StartProcessingRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Configuration_Service.StartProcessingRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConfigurationService(t *testing.T) {
	tests := []struct {
		name string
		want *Configuration_Service
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfigurationService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfigurationService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeJSON(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"",""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeJSON(); got != tt.want {
				t.Errorf("mergeJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
