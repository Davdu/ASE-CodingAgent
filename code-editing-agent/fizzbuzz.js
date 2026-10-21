// Clyppi: your 2026 FizzBuzz assistant. It looks like you're counting to 15!

function clyppiSay(message) {
  const bubble = "+" + "-".repeat(message.length + 2) + "+";
  console.log(bubble);
  console.log("| " + message + " |");
  console.log(bubble);
  console.log("  \\  (o_o)  Clyppi");
}

function intro() {
  console.log("Hello, friend! I'm Clyppi, your coding companion.");
  console.log("Let's count together and have some fun before we dive in.");
  console.log("Ready? Here we go!\n");
}

function fizzbuzz(n) {
  for (let i = 1; i <= n; i++) {
    if (i % 15 === 0) {
      clyppiSay("FizzBuzz! Great job, you hit a multiple of both 3 and 5!");
    } else if (i % 3 === 0) {
      clyppiSay("Fizz! Looks like a multiple of 3.");
    } else if (i % 5 === 0) {
      clyppiSay("Buzz! Looks like a multiple of 5.");
    } else {
      console.log(`${i}`);
    }
  }
  console.log("Clyppi: Done! Want me to optimize this? (Yes / Yes)");
}

intro();
fizzbuzz(15);
