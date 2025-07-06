# 🗭 Pokédex CLI

A command-line Pokédex built in Go! Explore, catch, and inspect Pokémon using the [PokeAPI](https://pokeapi.co/). Fully interactive, supports caching, pagination, and stat-based catching logic.

---

## 🚀 Features

* 🔍 Explore Pokémon by location area (`map`, `mapb`)
* 🎯 Catch Pokémon with success based on base experience
* 🗂️ Inspect caught Pokémon (name, height, weight, stats, types)
* 📦 In-memory caching to reduce API calls
* 🧠 Built with Go and modular internal packages
* 💬 Simple REPL interface (`Pokedex >` prompt)

---

## 📦 Installation

1. **Clone the repository**

```bash
git clone https://github.com/yourusername/pokedex.git
cd pokedex
```

2. **Run the application**

```bash
go run main.go
```

---

## 🛠️ Usage

```txt
Pokedex > help
```

### Available Commands:

| Command              | Description                               |
| -------------------- | ----------------------------------------- |
| `help`               | List all available commands               |
| `exit`               | Exit the Pokédex                          |
| `map`                | Show next 20 location areas               |
| `mapb`               | Go back to the previous 20 location areas |
| `explore <location>` | Show wild Pokémon in a location area      |
| `catch <name>`       | Try catching a Pokémon by name            |
| `inspect <name>`     | View info about a caught Pokémon          |

---

## 📂 Project Structure

```
.
├── main.go                 # Entry point with CLI logic
├── internal/
│   └── pokecache/          # Caching logic using mutex and TTL
│       └── cache.go
└── go.mod / go.sum         # Module definitions
```

---

## 💪 Example

```txt
Pokedex > map
location-area-1
location-area-2
...

Pokedex > explore viridian-forest
Exploring viridian-forest...
Found Pokemon:
- caterpie
- weedle
- pikachu

Pokedex > catch pikachu
Throwing a Pokeball at pikachu...
Caught pikachu!

Pokedex > inspect pikachu
Name: pikachu
Height: 4
Weight: 60
Stats:
  -hp: 35
  -attack: 55
  -defense: 40
  -special-attack: 50
  -special-defense: 50
  -speed: 90
Types:
  - electric
```

---

## 🛠️ Built With

* [Go](https://golang.org/)
* [PokeAPI](https://pokeapi.co/)
* [JSON](https://pkg.go.dev/encoding/json)
* [Standard Library](https://pkg.go.dev/std)

---

## 🙌 Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you'd like to change.

---

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

---

## 📬 Acknowledgments

* Inspired by real-world CLI tooling
* Thanks to [PokeAPI](https://pokeapi.co/) for their awesome data
