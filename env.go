package main

import (
  "strings"
	"os"
	"log"
	"path/filepath"
	"bufio"
	"errors"
	"fmt"
)

var (
	ErrEnvFileLoadedAlready = errors.New(".env file was loaded already")
)

type Env struct {
	fileData map[string]string
}

func (env Env) Get(key string, defaultValue string) (val string, ok bool) {

	if val, ok = os.LookupEnv(key); !ok {
		val, ok = env.fileData[key]
		if !ok {
			val = defaultValue
		}
	}

	val = strings.Trim(val, " \n\t;.=")

	return val, ok
}

func (env * Env) LoadFile() error {
	if env.fileData != nil && len(env.fileData) != 0 {
		return ErrEnvFileLoadedAlready
	}

	envFilename, err := filepath.Abs(".env")
	if err != nil {
		message := fmt.Sprintf("filepath.Abs() error: %s", err.Error())
		return errors.New(message)
	}
	log.Println("Reading env file", envFilename)

	f, err := os.Open(envFilename)
	if err != nil {
		return err
	}

	defer f.Close()

	env.fileData = make(map[string]string, 0)

	scanner := bufio.NewScanner(f)
  for scanner.Scan() {
      line := scanner.Text()

			if strings.HasPrefix(strings.Trim(line, " \t"), "#") {
				continue
			}

			parts := strings.Split(line, "=")
			if len(parts) != 2 {
				continue
			}

			for i, _ := range parts {
				parts[i] = strings.Trim(parts[i], " ,\t;#\"")
			}

			env.fileData[parts[0]] = parts[1]
  }

  if err := scanner.Err(); err != nil {
      return err
  }

	return nil
}
