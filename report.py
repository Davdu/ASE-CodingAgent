"""Summarize a harbor job: reward, turns, tokens and cost per trial, plus totals.

Usage:
    python report.py                 # newest job in ./jobs
    python report.py jobs/<job-dir>  # a specific job
    python report.py ... --append results.md   # also append the report to results.md
"""
import argparse
import json
from datetime import datetime
from pathlib import Path


def load_trials(job_dir: Path) -> list[dict]:
    trials = []
    for result_file in sorted(job_dir.glob("*/result.json")):
        r = json.loads(result_file.read_text())
        agent = r.get("agent_result") or {}
        meta = agent.get("metadata") or {}
        rewards = (r.get("verifier_result") or {}).get("rewards") or {}
        exec_times = r.get("agent_execution") or {}
        duration = None
        if exec_times.get("started_at") and exec_times.get("finished_at"):
            start = datetime.fromisoformat(exec_times["started_at"].replace("Z", "+00:00"))
            end = datetime.fromisoformat(exec_times["finished_at"].replace("Z", "+00:00"))
            duration = (end - start).total_seconds()
        trials.append({
            "task": r["task_name"].split("/")[-1].removeprefix("polyglot_"),
            "reward": rewards.get("reward"),
            "turns": meta.get("turns"),
            "reviews": meta.get("reviews"),
            "input": agent.get("n_input_tokens"),
            "output": agent.get("n_output_tokens"),
            "cost": agent.get("cost_usd"),
            "seconds": duration,
            "error": r.get("exception_info") is not None,
        })
    return trials


def fmt(value, spec="") -> str:
    return "-" if value is None else format(value, spec)


def mean(values):
    values = [v for v in values if v is not None]
    return sum(values) / len(values) if values else None


def total(values):
    values = [v for v in values if v is not None]
    return sum(values) if values else None


def render(job_dir: Path, trials: list[dict]) -> str:
    lines = [
        f"## {job_dir.name}",
        "",
        "| Task | Reward | Turns | Reviews | Input tok | Output tok | Cost (USD) | Agent time (s) |",
        "|---|---|---|---|---|---|---|---|",
    ]
    for t in trials:
        task = t["task"] + (" (exception)" if t["error"] else "")
        lines.append(
            f"| {task} | {fmt(t['reward'], '.1f')} | {fmt(t['turns'])} | {fmt(t['reviews'])} "
            f"| {fmt(t['input'], ',')} | {fmt(t['output'], ',')} | {fmt(t['cost'], '.4f')} "
            f"| {fmt(t['seconds'], '.0f')} |"
        )

    solved = sum(1 for t in trials if t["reward"] == 1.0)
    with_stats = sum(1 for t in trials if t["cost"] is not None)
    lines += [
        "",
        "| Metric | Value |",
        "|---|---|",
        f"| Trials | {len(trials)} (with token stats: {with_stats}) |",
        f"| Solved | {solved} |",
        f"| Mean reward | {fmt(mean([t['reward'] for t in trials]), '.3f')} |",
        f"| Total cost (USD) | {fmt(total([t['cost'] for t in trials]), '.4f')} |",
        f"| Mean cost per trial (USD) | {fmt(mean([t['cost'] for t in trials]), '.4f')} |",
        f"| Cost per solved task (USD) | {fmt(total([t['cost'] for t in trials]) / solved if solved and with_stats else None, '.4f')} |",
        f"| Mean turns | {fmt(mean([t['turns'] for t in trials]), '.1f')} |",
        f"| Total turns | {fmt(total([t['turns'] for t in trials]))} |",
        f"| Total input tokens | {fmt(total([t['input'] for t in trials]), ',')} |",
        f"| Total output tokens | {fmt(total([t['output'] for t in trials]), ',')} |",
        "",
    ]
    return "\n".join(lines)


def main():
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("job_dir", nargs="?", help="harbor job directory (default: newest in ./jobs)")
    parser.add_argument("--append", metavar="FILE", help="append the report to this markdown file")
    args = parser.parse_args()

    if args.job_dir:
        job_dir = Path(args.job_dir)
    else:
        jobs = sorted(p for p in Path("jobs").iterdir() if p.is_dir())
        if not jobs:
            raise SystemExit("no jobs found in ./jobs")
        job_dir = jobs[-1]

    report = render(job_dir, load_trials(job_dir))
    print(report)
    if args.append:
        with open(args.append, "a") as f:
            f.write("\n" + report)


if __name__ == "__main__":
    main()
