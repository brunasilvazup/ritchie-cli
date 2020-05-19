package formula

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/google/uuid"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

var ErrNotEnableDocker = errors.New("this formula is not enabled to run in a container")

type DockerPreRunner struct {
	sDefault Setuper
}

func NewDockerPreRunner(setuper Setuper) DockerPreRunner {
	return DockerPreRunner{sDefault: setuper}
}

func (d DockerPreRunner) PreRun(def Definition) (Setup, error) {
	setup, err := d.sDefault.Setup(def)
	if err != nil {
		return Setup{}, err
	}

	if err := validate(setup.tmpBinDir); err != nil {
		return Setup{}, err
	}

	containerId, err := uuid.NewRandom()
	if err != nil {
		return Setup{}, err
	}

	setup.containerId = containerId.String()
	if err := buildImg(setup.containerId); err != nil {
		return Setup{}, err
	}

	return setup, nil
}

func validate(tmpBinDir string) error {
	dockerFile := fmt.Sprintf("%s/Dockerfile", tmpBinDir)
	if !fileutil.Exists(dockerFile) {
		return ErrNotEnableDocker
	}

	return nil
}

func buildImg(containerId string) error {
	fmt.Println("Building docker image...")
	args := []string{dockerBuildCmd, "-t", containerId, "."}
	cmd := exec.Command(docker, args...)
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	fmt.Println("Docker image was built :)")
	return nil
}
