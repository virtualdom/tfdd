package terraform

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

const (
	TF_CMD      = "terraform"
	TF_INIT_CMD = "init"
	TF_PLAN_CMD = "plan"
)

type TFClient struct {
	dir string
}

func read(r io.ReadCloser) (string, error) {
	b := make([]byte, 8)
	contents := ""

	for {
		n, err := r.Read(b)
		contents += string(b[:n])
		if err == io.EOF {
			break
		} else if err != nil {
			return contents, err
		}
	}

	return contents, nil
}

func getOutput(cmd *exec.Cmd) (string, string, error) {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", "", errors.Wrap(err, "failed to get stderr")
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", errors.Wrap(err, "failed to get stdout")
	}

	args := ""
	for _, a := range cmd.Args {
		args += a + " "
	}

	err = cmd.Start()
	if err != nil {
		return "", "", errors.Wrap(err, "failed to start command "+args)
	}

	output, err := read(stdout)
	if err != nil {
		return output, "", errors.Wrap(err, "failed to parse output of "+args)
	}

	errlog, err := read(stderr)
	if err != nil {
		return output, errlog, errors.Wrap(err, "failed to parse error log of "+args)
	}

	return output, errlog, cmd.Wait()
}

func GetServices(root string) ([]string, error) {
	services := make([]string, 0)

	toExplore := []string{root}

	for len(toExplore) > 0 {
		currentDir := toExplore[0]

		dir, err := os.Open(currentDir)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open "+currentDir)
		}

		dirContents, err := dir.Readdir(-1)
		dir.Close()
		if err != nil {
			return nil, errors.Wrap(err, "failed to read contents of "+currentDir)
		}

		isService := false

		for _, file := range dirContents {
			if file.IsDir() && strings.Index(file.Name(), ".") != 0 {
				toExplore = append(toExplore, fmt.Sprintf("%s/%s", currentDir, file.Name()))
			}

			if !file.IsDir() && strings.Index(file.Name(), ".tf") == len(file.Name())-3 {
				isService = true
			}
		}

		if isService {
			services = append(services, currentDir)
		}

		toExplore = toExplore[1:]
	}

	return services, nil
}

func New(dir string) *TFClient {
	return &TFClient{
		dir: dir,
	}
}

func (client *TFClient) Init() error {
	init := exec.Command(TF_CMD, TF_INIT_CMD)
	init.Dir = client.dir

	_, errlog, err := getOutput(init)
	return errors.Wrap(err, "failed to run `terraform init`:\n"+errlog)
}

func (client *TFClient) Plan() (string, string, error) {
	plan := exec.Command(TF_CMD, TF_PLAN_CMD)
	plan.Dir = client.dir

	return getOutput(plan)
}
