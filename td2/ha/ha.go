package ha

import (
	"errors"
	"os/exec"
	"time"
)

var haMap map[string]*HaState = map[string]*HaState{}

func InitHaState(chainName string) *HaState {
	state := &HaState{
		ChainName: chainName,
		State:     "init",
		Status:    "offline",
		Jailed:    false,
	}

	haMap[chainName] = state

	return state
}

func ServiceAction(name string, action string) (string, error) {
	cmd := exec.Command("systemctl", action, name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "error", err
		// if exitErr, ok := err.(*exec.ExitError); ok {
		//   fmt.Printf("systemctl finished with non-zero: %v\n", exitErr)
		// } else {
		//   fmt.Printf("failed to run systemctl: %v", err)
		//   os.Exit(1)
		// }
	}
	return string(out), nil
}

func ProcessHaState(state *HaState, config *HaConfig, chainConfig *HaChainConfig) {
	if state.Jailed {
		if state.State == "active" || state.Status == "active" {
			_, err := ServiceAction(chainConfig.ServiceName, "stop")
			if err == nil {
				state.State = "standby"
				state.Status = "inactive"
			}
		}

		time.Sleep(time.Second)
		ProcessHaState(state, config, chainConfig)
		return
	}

	state.JailedNotified = false

	if state.State == "init" {
		status, err := ServiceAction(chainConfig.ServiceName, "check")

		if err != nil {
			state.Jailed = true
		} else {
			state.Status = status
			if status == "active" {
				state.State = "active"
			} else {
				state.State = "standby"
			}
		}

		state.Retry = 0
	} else if state.State == "active" {
		// Check state of another peer

		// Stop if another peer is running concurrently

		// If everything normal, ensure tmkms is running
		status, err := ServiceAction(chainConfig.ServiceName, "check")

		if err != nil {
			state.Jailed = true
		} else {
			state.Status = status
			if status != "active" {
				_, err := ServiceAction(chainConfig.ServiceName, "restart")

				if err != nil {
					state.Retry++
					if state.Retry >= 5 {
						state.Jailed = true
					}
				}
			} else {
				state.Retry = 0
			}
		}
	} else if state.State == "standby" {

	}

	ProcessHaState(state, config, chainConfig)
}

func Unjail(chainName string) error {
	state, ok := haMap[chainName]
	if !ok {
		return errors.New("chain not found")
	}

	if state.Jailed {
		state.State = "init"
		state.Status = "offline"
		state.Jailed = false
	}

	return nil
}

func OnAlert(chainName string) error {
	state, ok := haMap[chainName]
	if !ok {
		return errors.New("chain not found")
	}

	if state.Jailed {
		return nil
	}

	if state.State == "standby" {

	} else {
		state.Jailed = true
	}
}
