// Runs test.js and prints its output (and exit code) even on failure.
const { spawnSync } = require('node:child_process');
const r = spawnSync(process.execPath, [require('node:path').join(__dirname, 'test.js')], { encoding: 'utf8' });
const lines = (r.stdout + r.stderr).split('\n').filter(l => /^(✔|✖|ℹ (tests|pass|fail))/.test(l));
console.log(lines.join('\n'));
console.log('exit code:', r.status);
