package envdir

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getEnvs(rootPath string) map[string]string {
	envs := make(map[string]string)

	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path == rootPath {
				return nil
			}

			// Do not enter into subdirectories
			return filepath.SkipDir
		}

		// Skip file which contains '=' in the name
		if strings.Contains(d.Name(), "=") {
			return nil
		}

		// Skip hidden files
		if d.Name()[0] == '.' {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Failed to open file with path %s, got error %v\n", path, err)
			return nil
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			fmt.Printf("failed to get file stat: %v\n", err)
			return nil
		}

		if !fileInfo.Mode().IsRegular() {
			fmt.Println("source file must be regular file")
			return nil
		}

		envName := d.Name()

		if fileInfo.Size() == 0 {
			delete(envs, envName)
		}

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			envValue := strings.TrimRight(scanner.Text(), " \n\t")
			envs[envName] = envValue
			break
		}

		return nil
	})
	if err != nil {
		fmt.Printf("Error during walking root: %v\n", err)
	}

	return envs
}

func enrichWithEnvs(cmd *exec.Cmd, envs map[string]string) {
	cmd.Env = append(cmd.Env, os.Environ()...)
	for k, v := range envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
}

func Run(root string, args ...string) {
	cmd := exec.Command(args[0], args[1:]...) //#nosec G204

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	envs := getEnvs(root)
	enrichWithEnvs(cmd, envs)

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error during running command: %v", err)
		os.Exit(cmd.ProcessState.ExitCode())
	}
}

// func main() { //nolint:deadcode
// 	flag.Parse()
// 	args := flag.Args()

// 	if len(args) < 2 {
// 		fmt.Println("Usage: go-envdir /path/to/env/dir command [arg]")
// 		os.Exit(1)
// 	}

// 	Run(args[0], args[1:]...)
// }
