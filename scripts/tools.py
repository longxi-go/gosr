"""Run Go module checks with explicit source or published dependency policy."""

from contextlib import contextmanager
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile

from go_module import ROOT_MODULE, read_module


@contextmanager
def environment(root: Path, published: bool):
    """Provide an isolated Go environment without inheriting module overrides."""
    if os.environ.get("GOFLAGS"):
        raise ValueError("GOFLAGS must be unset for isolated module checks")

    env = os.environ.copy()
    env["GOWORK"] = "off"
    env["GOFLAGS"] = ""
    go = os.environ.get("GO", "go")
    effective_flags = subprocess.run(
        [go, "env", "GOFLAGS"],
        env=env,
        capture_output=True,
        text=True,
        check=True,
    ).stdout.strip()
    if effective_flags:
        raise ValueError("go env GOFLAGS must be empty for isolated module checks")
    if published:
        replacements = read_module(root / "tools/go.mod").get("Replace") or []
        if any(not replacement["New"].get("Version") for replacement in replacements):
            raise ValueError("published checks forbid local module replacements")
        if any(replacement["Old"]["Path"] == ROOT_MODULE for replacement in replacements):
            raise ValueError("published checks forbid root module replacements")
        yield env
        return

    with tempfile.TemporaryDirectory(prefix="gosr-tools-mod-") as directory:
        modfile = Path(directory) / "gosr-tools.mod"
        sumfile = modfile.with_suffix(".sum")
        shutil.copy2(root / "tools/go.mod", modfile)
        source_sum = root / "tools/go.sum"
        if source_sum.exists():
            # Tidy omits checksums for the root module when it uses local source.
            lines = source_sum.read_text(encoding="utf-8").splitlines(keepends=True)
            sumfile.write_text(
                "".join(line for line in lines if not line.startswith(f"{ROOT_MODULE} ")),
                encoding="utf-8",
            )
        else:
            sumfile.touch()
        subprocess.run(
            [
                go,
                "mod",
                "edit",
                f"-modfile={modfile}",
                f"-replace={ROOT_MODULE}={root.resolve()}",
            ],
            check=True,
            env=env,
        )
        env["GOFLAGS"] = f"-modfile={modfile}"
        yield env


def command_environment(command, env):
    """Keep source-modfile flags on the Go invocation, away from tested programs."""
    command = list(command)
    env = env.copy()
    if command[0] == os.environ.get("GO", "go") and env.get("GOFLAGS"):
        position = 3 if command[1] == "mod" else 2
        command.insert(position, env["GOFLAGS"])
        env["GOFLAGS"] = ""
    return command, env


def workspace_environment(root, env):
    """Use a private tools workspace for linters that run their own Go commands."""
    env = env.copy()
    modfile = Path(env["GOFLAGS"].removeprefix("-modfile="))
    env["GOFLAGS"] = ""
    work = modfile.parent / "go.work"
    go = os.environ.get("GO", "go")
    subprocess.run(
        [go, "work", "init", str(root / "tools")],
        cwd=work.parent,
        env=env,
        check=True,
    )
    env["GOWORK"] = str(work)
    subprocess.run(
        [go, "work", "edit", f"-replace={ROOT_MODULE}={root.resolve()}"],
        cwd=work.parent,
        env=env,
        check=True,
    )
    return env


def main() -> None:
    if len(sys.argv) < 3 or sys.argv[1] not in {"source", "published"}:
        raise SystemExit("usage: tools.py source|published -- command [arguments...]")
    mode = sys.argv[1]
    command = sys.argv[2:]
    if command[0] == "--":
        command = command[1:]
    if not command:
        raise SystemExit("a command is required")

    root = Path(__file__).resolve().parents[1]
    with environment(root, published=mode == "published") as env:
        if mode == "source" and command[0] != os.environ.get("GO", "go"):
            env = workspace_environment(root, env)
        command, env = command_environment(command, env)
        result = subprocess.run(command, cwd=root / "tools", env=env, check=False)
    raise SystemExit(result.returncode)


if __name__ == "__main__":
    main()
