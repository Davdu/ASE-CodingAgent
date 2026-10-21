#!/usr/bin/env node
'use strict';

/**
 * Lazily yields the moves that solve Tower of Hanoi for n disks,
 * moving from `from` to `to` with `via` as the spare peg.
 */
function* hanoiMoves(n, from = 'A', to = 'C', via = 'B') {
  if (n === 0) return;
  yield* hanoiMoves(n - 1, from, via, to);
  yield [from, to];
  yield* hanoiMoves(n - 1, via, to, from);
}

/** Number of moves needed for n disks: 2^n - 1 (BigInt, so any n works). */
function moveCount(n) {
  return (1n << BigInt(n)) - 1n;
}

/** Replays moves on a fresh tower and throws if any move is illegal. */
function verify(n, moves) {
  const pegs = { A: [], B: [], C: [] };
  for (let d = n; d >= 1; d--) pegs.A.push(d);
  for (const [f, t] of moves) {
    if (pegs[f].length === 0) throw new Error(`Empty peg ${f}`);
    const disk = pegs[f].pop();
    const top = pegs[t][pegs[t].length - 1];
    if (top !== undefined && top < disk) throw new Error(`Illegal move: disk ${disk} onto ${top}`);
    pegs[t].push(disk);
  }
  if (pegs.C.length !== n) throw new Error('Puzzle not solved');
  return true;
}

function main(argv) {
  const n = Number(argv[2] ?? 3);
  if (!Number.isInteger(n) || n < 1) {
    console.error('Usage: node hanoi_v2.js <disks>=1..');
    process.exit(1);
  }
  if (n > 20) {
    console.log(`n=${n} needs ${moveCount(n)} moves; printing all of them is impractical.`);
    return;
  }
  const moves = [...hanoiMoves(n)];
  verify(n, moves);
  moves.forEach(([f, t], i) => console.log(`${String(i + 1).padStart(3)}. ${f} -> ${t}`));
  console.log(`Solved ${n} disks in ${moves.length} moves (verified).`);
}

if (require.main === module) main(process.argv);

module.exports = { hanoiMoves, moveCount, verify };
