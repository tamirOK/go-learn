package envdir

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Env struct {
	Name  string
	Value *string
}

func (e Env) String() string {
	if e.Value == nil {
		return ""
	}

	return fmt.Sprintf("%s=%s", e.Name, *e.Value)
}

func getEnvFromFile(baseName string) *Env {
	file, err := os.Open(baseName)
	if err != nil {
		fmt.Printf("Failed to open file with path %s, got error %v\n", baseName, err)
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

	envName := fileInfo.Name()
	var envValue string

	if fileInfo.Size() == 0 {
		return &Env{envName, nil}
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		envValue = strings.TrimRight(scanner.Text(), " \n\t")
		break
	}

	return &Env{envName, &envValue}
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

		env := getEnvFromFile(filepath.Join(rootPath, entryName))

		if env != nil {
			envs[env.Name] = *env
		}
	}

	return envs
}

func enrichCmdWithEnvs(cmd *exec.Cmd, envs map[string]Env) {
	for _, envStr := range os.Environ() {
		parts := strings.SplitN(envStr, "=", 2)

		// Remove env from command envs if it was marked for deletion
		if env, ok := envs[parts[0]]; ok && env.Value == nil {
			continue
		}
		cmd.Env = append(cmd.Env, envStr)
	}

	for _, env := range envs {
		// Skip deleted envs
		if env.Value != nil {
			cmd.Env = append(cmd.Env, env.String())
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
