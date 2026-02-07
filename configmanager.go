package main

type Config struct {
	firstName  string
	middleName string
	lastName   string
}

func LoadConfig() (Config, error) {
	var c Config
	// TODO

	// Load settings from the file, if there
	// configDir, dirErr := os.UserConfigDir()

	// var (
	// 	configPath string
	// 	origConfig []byte
	// )
	// if dirErr == nil {
	// 	configPath = filepath.Join(configDir, "ExampleUserConfigDir", "example.conf")
	// 	var err error
	// 	origConfig, err = os.ReadFile(configPath)
	// 	if err != nil && !os.IsNotExist(err) {
	// 		// The user has a config file but we couldn't read it.
	// 		// Report the error instead of ignoring their configuration.
	// 		log.Fatal(err)
	// 	}
	// }

	// // Use and perhaps make changes to the config.
	// config := bytes.Clone(origConfig)

	return c, nil
}

func SaveConfig(c Config) string {
	// TODO

	// configDir, dirErr := os.UserConfigDir()
	// if dirErr == nil {
	// 	configPath = filepath.Join(configDir, "ExampleUserConfigDir", "example.conf")
	// 	var err error
	// 	origConfig, err = os.ReadFile(configPath)
	// 	if err != nil && !os.IsNotExist(err) {
	// 		// The user has a config file but we couldn't read it.
	// 		// Report the error instead of ignoring their configuration.
	// 		log.Fatal(err)
	// 	}
	// }
	// if configPath == "" {
	// 	log.Printf("not saving config changes: %v", dirErr)
	// } else {
	// 	err := os.MkdirAll(filepath.Dir(configPath), 0700)
	// 	if err == nil {
	// 		err = os.WriteFile(configPath, config, 0600)
	// 	}
	// 	if err != nil {
	// 		log.Printf("error saving config changes: %v", err)
	// 	}
	// }

	return "Success" // TODO: BG localization
}
