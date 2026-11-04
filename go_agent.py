# go_agent.py — Harbor wrapper for the Go coding agent in ./build/agent
import json
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

# Where the agent writes its turns/tokens/cost inside the container (see tools/usage.go).
REMOTE_STATS = "/tmp/agent_stats.json"


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
        env["AGENT_STATS_FILE"] = REMOTE_STATS
        result = await environment.exec(
            f"{REMOTE_PATH} {shlex.quote(instruction)}",
            env=env,
        )
        # Collect stats before checking the exit code, so failed runs still report their cost.
        await self._collect_stats(environment, context)
        # This Harbor version's AgentContext has no exit_code/error_message fields,
        # so surface failures as an exception (shows up in the trial's Exception tab).
        if result.return_code != 0:
            raise RuntimeError(
                f"agent exited with code {result.return_code}\n"
                f"--- stdout (tail) ---\n{(result.stdout or '')[-2000:]}\n"
                f"--- stderr (tail) ---\n{(result.stderr or '')[-2000:]}"
            )

    async def _collect_stats(
        self, environment: BaseEnvironment, context: AgentContext
    ) -> None:
        """Read the agent's stats file and put turns, tokens and cost into the trial result."""
        out = await environment.exec(f"cat {REMOTE_STATS}")
        if out.return_code != 0 or not out.stdout:
            self.logger.warning("agent stats file missing; no token/cost data for this trial")
            return
        try:
            stats = json.loads(out.stdout)
        except json.JSONDecodeError:
            self.logger.warning("agent stats file is not valid JSON")
            return

        total = stats["total"]
        # Harbor counts n_input_tokens including cached tokens.
        context.n_input_tokens = (
            total["input_tokens"] + total["cache_read_tokens"] + total["cache_write_tokens"]
        )
        context.n_cache_tokens = total["cache_read_tokens"]
        context.n_output_tokens = total["output_tokens"]
        context.cost_usd = total["cost_usd"]
        context.metadata = {
            "turns": total["turns"],
            "agent_turns": stats["agent"]["turns"],
            "reviewer_turns": stats["reviewer"]["turns"],
            "reviews": stats["reviews"],
            "stats": stats,
        }

        if self.logs_dir:
            Path(self.logs_dir).mkdir(parents=True, exist_ok=True)
            (Path(self.logs_dir) / "agent_stats.json").write_text(out.stdout)
