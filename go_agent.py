# go_agent.py — Harbor wrapper for the Go coding agent in ./build/agent
import os
import shlex
from pathlib import Path

from harbor.agents.base import BaseAgent
from harbor.environments.base import BaseEnvironment
from harbor.models.agent.context import AgentContext

# Built by the benchmark script: go build -C code-editing-agent -o ../build/agent .
BINARY = Path(__file__).parent / "build" / "agent"
REMOTE_PATH = "/usr/local/bin/agent"

# API keys to forward from your shell into the container (only those that are set).
FORWARD_ENV = ["ANTHROPIC_API_KEY", "OPENAI_API_KEY", "OPENROUTER_API_KEY"]


class Clippy(BaseAgent):
    @staticmethod
    def name() -> str:
        return "code-editing-agent"

    def version(self) -> str | None:
        return "0.1"

    async def setup(self, environment: BaseEnvironment) -> None:
        if not BINARY.exists():
            raise FileNotFoundError(f"Agent binary not found at {
                                    BINARY} — build it first.")
        await environment.upload_file(BINARY, REMOTE_PATH)
        await environment.exec(f"chmod +x {REMOTE_PATH}")

    async def run(
        self,
        instruction: str,
        environment: BaseEnvironment,
        context: AgentContext,
    ) -> None:
        env = {k: os.environ[k] for k in FORWARD_ENV if k in os.environ}
        result = await environment.exec(
            f"{REMOTE_PATH} {shlex.quote(instruction)}",
            env=env,
        )
        # This Harbor version's AgentContext has no exit_code/error_message fields,
        # so surface failures as an exception (shows up in the trial's Exception tab).
        if result.return_code != 0:
            raise RuntimeError(
                f"agent exited with code {result.return_code}\n"
                f"--- stdout (tail) ---\n{(result.stdout or '')[-2000:]}\n"
                f"--- stderr (tail) ---\n{(result.stderr or '')[-2000:]}"
            )
