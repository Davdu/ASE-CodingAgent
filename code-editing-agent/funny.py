import random

ACTIONS = [
    "wrestle a rubber chicken",
    "teach a goldfish to moonwalk",
    "argue with a vending machine",
    "knit a sweater for a pickle",
    "challenge a pigeon to a staring contest",
]

JOKES = [
    "Why did the programmer quit his job? He didn't get arrays.",
    "I told my computer a joke about UDP. It didn't know if I'd get it.",
    "Why do cows wear bells? Because their horns don't work.",
    "What do you call a fake noodle? An impasta.",
    "I'm reading a book about anti-gravity. It's impossible to put down.",
]

def silly_fortune():
    name = input("What's your name? ").strip() or "Mysterious Stranger"
    print(f"\nWelcome, {name}! Your destiny today: you will {random.choice(ACTIONS)}.")
    print(f"Wise words: {random.choice(JOKES)}\n")

def main():
    print("=== The Ridiculous Machine ===")
    while True:
        choice = input("Type 'f' for a fortune, 'j' for a joke, or 'q' to quit: ").strip().lower()
        if choice == "f":
            silly_fortune()
        elif choice == "j":
            print(random.choice(JOKES) + "\n")
        elif choice == "q":
            print("Goodbye! May your socks always match.")
            break
        else:
            print("The machine does not understand. It is very confused. Beep boop.\n")

if __name__ == "__main__":
    main()
