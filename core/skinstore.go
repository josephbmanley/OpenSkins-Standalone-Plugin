package core

import (
	"bufio"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/josephbmanley/OpenSkins-Common/datatypes"
	"os"
)

type standaloneConfig struct {
	SkinDirectory    string `yaml:"skin_directory" env:"SKIN_DIR" env-default:"skins"`
	WebserverDomain  string `yaml:"domain" env:"DOMAIN" env-default:"localhost"`
	WebserverSubpath string `yaml:"subpath" env:"SUBPATH" env-default:""`
}

const configFile = "standalone_config.yaml"

var config standaloneConfig

// SkinstoreStandalone is root of the skinstore plugin
type SkinstoreStandalone struct{}

// Initialize intializes the skinstore module
func (s *SkinstoreStandalone) Initialize() error {

	for _, w := range warnings {
		fmt.Printf("WARNING: %v\n", w)
	}

	// Check if config file exists
	if _, err := os.Stat(configFile); err == nil {

		// Read config with config file & environment variables
		if err := cleanenv.ReadConfig(configFile, &config); err != nil {
			return err
		}
	} else if os.IsNotExist(err) {

		// Read config with only environment variables
		if err := cleanenv.ReadEnv(&config); err != nil {
			return err
		}
	} else {
		return err
	}

	if _, err := os.Stat(config.SkinDirectory); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		os.Mkdir(config.SkinDirectory, 0755)
	}

	return nil
}

// GetSkin returns a skin object
func (s *SkinstoreStandalone) GetSkin(skinID string) (*datatypes.Skin, error) {

	skinPath := fmt.Sprintf("%v/%v", config.SkinDirectory, skinID)

	// Check if skin exists
	if _, err := os.Stat(skinPath); err != nil {

		// Catch not exist errors
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	skin := datatypes.Skin{
		UID:      skinID,
		Name:     skinID,
		Location: fmt.Sprintf("%v/%v%v", config.WebserverDomain, config.WebserverSubpath, skinID),
		Metadata: map[string]string{},
	}

	return &skin, nil
}

// AddSkin returns a skin creates a new skin object
func (s *SkinstoreStandalone) AddSkin(skinID string, fileData []byte) error {

	skinPath := fmt.Sprintf("%v/%v", config.SkinDirectory, skinID)

	if _, err := os.Stat(skinPath); err != nil {

		// Check if error was not 404 error
		if !os.IsNotExist(err) {
			return err
		}

		// Create file if it didn't exist
		if _, err := os.Create(skinPath); err != nil {
			return err
		}

	}

	// Open file with rw permissions
	file, err := os.OpenFile(skinPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	// Write data to file
	writer := bufio.NewWriter(file)
	if _, err = writer.Write(fileData); err != nil {
		return err
	}

	return nil
}

// DeleteSkin deletes a skin object
// func (s *SkinstoreStandalone) DeleteSkin(skinID string) error {
// 	skinPath := fmt.Sprintf("%v/%v", config.SkinDirectory, skinID)

// 	// Check if skin exists
// 	if _, err := os.Stat(skinPath); err != nil {
// 		// Check if error was not 404 error
// 		if !os.IsNotExist(err) {
// 			return err
// 		}
// 		return nil
// 	}

// 	return os.Remove(skinPath)
// }
