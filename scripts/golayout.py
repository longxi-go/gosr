"""Delegate repository formatting and headers to the pinned golayout module."""

import json
import os
from pathlib import Path
import subprocess
import sys

from go_module import read_module
from tools import environment, workspace_environment

OPERATIONS = {
    "setup": ("format.sh", "setup"),
    "write": ("format.sh", "write"),
    "check": ("format.sh", "check"),
    "test": ("format_test.sh", None),
    "headers-write": ("license_headers.sh", "write"),
    "headers-check": ("license_headers.sh", "check"),
}


def run(root: Path, operation: str) -> int:
    """Resolve shared scripts and keep Go flags out of their subprocesses."""
    go = os.environ.get("GO", "go")
    with environment(root, published=False) as env:
        env = workspace_environment(root, env)
        download = subprocess.run(
            [go, "-C", str(root / "tools"), "mod", "download", "-json", "github.com/longxi-go/golayout"],
            env=env,
            capture_output=True,
            text=True,
            check=False,
        )
        if download.returncode:
            raise ValueError("Cannot resolve pinned golayout module:\n" + download.stderr + download.stdout)
        try:
            module = json.loads(download.stdout)
            directory = module.get("Dir") if isinstance(module, dict) else None
            if not isinstance(directory, str) or not directory or not Path(directory).is_absolute():
                raise ValueError("expected an absolute module directory")
        except ValueError as error:
            raise ValueError(f"Invalid golayout download metadata: {error}") from error
        script_name, mode = OPERATIONS[operation]
        directory = Path(directory)
        script = directory / "scripts" / script_name
        if not script.is_file():
            raise ValueError(f"Pinned golayout script is missing: {script}")
        language = read_module(root / "go.mod").get("Go")
        if not language:
            raise ValueError("A go directive is required in the root go.mod")
        command = ["bash", str(script)]
        if mode is not None:
            command.append(mode)
        command.extend([str(root), str(root / "tools"), "go" + language])
        result = subprocess.run(command, env=env, check=False)
    return result.returncode


def main() -> None:
    if len(sys.argv) != 2 or sys.argv[1] not in OPERATIONS:
        raise SystemExit("expected a golayout operation: " + ", ".join(OPERATIONS))
    try:
        result = run(Path(__file__).resolve().parents[1], sys.argv[1])
    except (ValueError, OSError, subprocess.SubprocessError) as error:
        raise SystemExit(str(error)) from error
    raise SystemExit(result)


if __name__ == "__main__":
    main()
