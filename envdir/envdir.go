package envdir

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Env struct {
	Value  string
	Delete bool
}

func getEnvFromFile(baseName string) (Env, error) {
	emptyEnv := Env{}
	file, err := os.Open(baseName)
	if err != nil {
		return emptyEnv, fmt.Errorf("failed to open file with path %s, got error %w", baseName, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return emptyEnv, fmt.Errorf("failed to get file stat: %w", err)
	}

	if !fileInfo.Mode().IsRegular() {
		return emptyEnv, errors.New("source file must be regular file")
	}

	var envValue string

	if fileInfo.Size() == 0 {
		return Env{Value: "", Delete: true}, nil
	}

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		envValue = strings.TrimRight(scanner.Text(), " \n\t")
	}

	return Env{Value: envValue}, nil
}

func getEnvs(rootPath string) map[string]Env {
	envs := make(map[string]Env)

	entries, err := os.ReadDir(rootPath)
	if err != nil {
		fmt.Printf("Could not list directory entries: %v\n", err)
		return envs
	}

	for _, entry := range entries {
		// Skip directories
		if entry.IsDir() {
			continue
		}

		entryName := entry.Name()

		// Skip file which contains '=' in the name
		if strings.Contains(entryName, "=") {
			continue
		}

		// Skip hidden files
		if entryName[0] == '.' {
			continue
		}

		env, err := getEnvFromFile(filepath.Join(rootPath, entryName))
		if err != nil {
			fmt.Printf("Could not open file: %v\n", err)
			return envs
		}

		envs[entryName] = env
	}

	return envs
}

func enrichCmdWithEnvs(cmd *exec.Cmd, envs map[string]Env) {
	for _, envStr := range os.Environ() {
		parts := strings.SplitN(envStr, "=", 2)

		// Remove env from command envs if it was marked for deletion
		if env, ok := envs[parts[0]]; ok && env.Delete {
			continue
		}
		cmd.Env = append(cmd.Env, envStr)
	}

	for envName, env := range envs {
		// Skip deleted envs
		if !env.Delete {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", envName, env.Value))
		}
	}
}

func Run(root string, args ...string) {
	cmd := exec.Command(args[0], args[1:]...) //#nosec G204

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	envs := getEnvs(root)
	enrichCmdWithEnvs(cmd, envs)

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
