package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/magefile/mage/mg"
	"github.com/shirou/gopsutil/v3/process"

	"github.com/micro-service-lab/recs-seem-mono-container/internal/config"
	"github.com/micro-service-lab/recs-seem-mono-container/magefiles/utils"
)

// Kill kills all processes related to development server.
func Kill(ctx context.Context) {
	mg.CtxDeps(ctx, kill)
}

func kill(ctx context.Context) error {
	if err := killMageTasks(ctx); err != nil {
		return fmt.Errorf("kill mage tasks: %w", err)
	}

	if err := killArelo(ctx); err != nil {
		return fmt.Errorf("kill arelo processes: %w", err)
	}

	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("no .env file found: %w", err)
	}
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	if err := killProcessesListeningOnPort(ctx, uint16(cfg.Port)); err != nil {
		return fmt.Errorf("kill process listening on port %d: %w", cfg.Port, err)
	}

	dbgPort, err := debuggingPort()
	if err != nil {
		return fmt.Errorf("get debugging port: %w", err)
	}

	if err := killProcessesListeningOnPort(ctx, dbgPort); err != nil {
		return fmt.Errorf("kill process listening on port %d: %w", dbgPort, err)
	}

	fmt.Println("done.")

	return nil
}

func killMageTasks(ctx context.Context) error {
	processes, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return fmt.Errorf("list all processes: %w", err)
	}

	for _, v := range processes {
		if v.Pid == 0 {
			continue
		}

		commandline, err := v.CmdlineSliceWithContext(ctx)
		if err != nil {
			continue
		}

		//nolint:gomnd
		if len(commandline) < 2 {
			continue
		}

		command := commandline[0]
		if !isMageCommand(command) {
			continue
		}

		subcommand := commandline[1]
		if subcommand != "serve" && subcommand != "dev" {
			continue
		}

		fmt.Printf("kill `mage %s` process, pid: %d\n", subcommand, v.Pid)
		if err := killMageProcess(v.Pid); err != nil {
			return fmt.Errorf("kill mage process: %w", err)
		}
	}

	return nil
}

func isMageCommand(command string) bool {
	if utils.IsWindows() {
		return strings.HasSuffix(command, "\\mage.exe")
	}
	return command == "mage"
}

func killMageProcess(pid int32) error {
	if utils.IsWindows() {
		if err := killWindowsProcess(pid); err != nil {
			return fmt.Errorf("kill windows process: %w", err)
		}
		return nil
	}

	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return fmt.Errorf("find process: %w", err)
	}

	if err := proc.Signal(syscall.SIGINT); err != nil {
		return fmt.Errorf("send sigint: %w", err)
	}

	return nil
}

func killArelo(ctx context.Context) error {
	processes, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return fmt.Errorf("list all processes: %w", err)
	}

	for _, v := range processes {
		if v.Pid == 0 {
			continue
		}

		commandline, err := v.CmdlineSliceWithContext(ctx)
		if err != nil {
			continue
		}

		//nolint:gomnd
		if len(commandline) < 2 {
			continue
		}

		command := commandline[0]
		if !isAreloCommand(command) {
			continue
		}

		pkg := commandline[len(commandline)-1]
		if pkg != "." {
			continue
		}

		fmt.Printf("kill `arelo` process, pid: %d\n", v.Pid)
		if err := killAreloProcess(v.Pid); err != nil {
			return fmt.Errorf("kill arelo process: %w", err)
		}
	}

	return nil
}

func isAreloCommand(command string) bool {
	if utils.IsWindows() {
		return strings.HasSuffix(command, "\\arelo.exe")
	}
	return command == "arelo"
}

func killAreloProcess(pid int32) error {
	if utils.IsWindows() {
		if err := killWindowsProcess(pid); err != nil {
			return fmt.Errorf("kill windows process: %w", err)
		}
		return nil
	}

	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return fmt.Errorf("find process: %w", err)
	}

	if err := proc.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("send sigint: %w", err)
	}

	return nil
}

func killProcessesListeningOnPort(ctx context.Context, port uint16) error {
	processes, err := process.ProcessesWithContext(ctx)
	if err != nil {
		return fmt.Errorf("list all processes: %w", err)
	}

	for _, v := range processes {
		if v.Pid == 0 {
			continue
		}

		conns, err := v.ConnectionsWithContext(ctx)
		if err != nil {
			continue
		}

		for _, conn := range conns {
			if conn.Laddr.Port == uint32(port) {
				fmt.Printf("kill process listening on port %d, pid: %d\n", port, v.Pid)
				if utils.IsWindows() {
					if err := killWindowsProcess(v.Pid); err != nil {
						return fmt.Errorf("kill windows process: %w", err)
					}
				} else {
					proc, err := os.FindProcess(int(v.Pid))
					if err != nil {
						return fmt.Errorf("find process: %w", err)
					}

					if err := proc.Kill(); err != nil {
						return fmt.Errorf("kill: %w", err)
					}

					return nil
				}
				break
			}
		}
	}

	return nil
}

func killWindowsProcess(pid int32) error {
	//nolint:gosec
	if err := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(int(pid))).Run(); err != nil {
		return fmt.Errorf("execute TASKKILL: %w", err)
	}
	return nil
}
