# Benchmark the coding agent using harbor

# Compile the agent
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C code-editing-agent -o ../build/agent

[ -d venv ] || python3 -m venv venv
source ./venv/bin/activate
pip install -q harbor

# Run on 5 tasks first
PYTHONPATH=. harbor run \
  -d terminal-bench/terminal-bench \
  --agent go_agent:Clippy \
  -l 5

