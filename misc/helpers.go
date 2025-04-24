package misc

import "github.com/joho/godotenv"

func LoadEnv(path string) error {
	err := godotenv.Overload(path)
	return err
}
