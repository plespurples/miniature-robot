package config

import "github.com/spf13/viper"

// Config represents a structure of app configuration
type Config struct {
	Purples struct {
		Year string
	}
	Net struct {
		Port string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Name     string
		Password string
	}
	Security struct {
		AuthorizationString string
	}
}

// Data contains the loaded configuration
var Data Config = Config{}

// Load loads the config file and stores it in the specified variable
func Load() error {
	// set up config file name
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")

	// read the file
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// unmarshal the config
	viper.Unmarshal(&Data)
	return nil
}
