# Realfin

Layer 1 blockchain for real-world asset (RWA) tokenization, built with Cosmos SDK.

**Tech stack:** Go 1.25 | Cosmos SDK v0.53.5 | CometBFT v0.38.21 | IBC-Go v10.5.0

## Overview

Realfin is a purpose-built Layer 1 blockchain designed to empower small and medium enterprises (SMEs) through the tokenization of real-world assets (RWAs). Built with decentralization at its core, Realfin streamlines the process of tokenization, enabling businesses to unlock liquidity, access global markets, and participate in the future of finance.

The platform extends the Cosmos SDK with six custom modules that provide on-chain data infrastructure for price feeds, credit ratings, real estate valuations, real-world asset tokenization, and embedded insurance policies. These modules serve as the foundation for asset classification, risk assessment, valuation, tokenization, and insurance coverage within the Realfin ecosystem.

As a sovereign L1 chain, Realfin operates independently using CometBFT (Tendermint) Proof-of-Stake consensus, providing fast transaction finality, high throughput, and energy efficiency. Its modular design grants full sovereignty over protocol rules and customization capabilities tailored to RWA compliance requirements.

Key characteristics:

- **Sovereign L1 chain** with CometBFT (Tendermint) Proof-of-Stake consensus
- **IBC-enabled** for cross-chain interoperability — supports interchain transfers and interchain accounts, connecting Realfin to the broader Cosmos ecosystem
- **Six custom modules**: `oracle` (price data), `creditscore` (credit ratings), `realestate` (real estate valuations), `tokenization` (real-world asset tokenization), `insurance` (embedded insurance policies), `realfin` (base module with governance parameters)
- **Binary**: `realfind` — the node daemon and CLI client
- **Default denomination**: `urlf`
- **Address prefix**: `cosmos` (coin type 118)

Realfin is designed for a wide range of participants:

- **SMEs** looking to tokenize assets, unlock liquidity, or access new funding opportunities
- **Investors** interested in trading or owning fractionalized real-world assets
- **Developers** building applications and tools on a secure, scalable L1 blockchain
- **Validators** securing the network through staking and consensus participation

## Architecture

### Project Structure

```
realfin-core/
├── app/                    # Application wiring
│   ├── app.go              # App struct, keepers, depinject setup
│   ├── app_config.go       # Module registration, execution order
│   └── ibc.go              # IBC keeper setup (manual, non-depinject)
├── cmd/realfind/           # Binary entry point
│   ├── main.go
│   └── cmd/
│       ├── root.go         # Root command with depinject + AutoCLI
│       └── commands.go     # Standard Cosmos commands registration
├── proto/realfin/          # Protobuf definitions (source of truth for API)
│   ├── oracle/v1/          # Price feed proto definitions
│   ├── creditscore/v1/     # Credit rating proto definitions
│   ├── realestate/v1/      # Real estate rating proto definitions
│   ├── tokenization/v1/    # Asset tokenization proto definitions
│   ├── insurance/v1/       # Insurance policy proto definitions
│   └── realfin/v1/         # Base module proto definitions
├── x/                      # Custom modules implementation
│   ├── oracle/             # Price data module
│   ├── creditscore/        # Credit rating module
│   ├── realestate/         # Real estate rating module
│   ├── tokenization/       # Asset tokenization module
│   ├── insurance/          # Insurance policy module
│   └── realfin/            # Base module (params only)
├── docs/static/            # OpenAPI/Swagger specification
├── config.yml              # Development chain configuration
└── Makefile                # Build, test, and development commands
```

### Module Layout

Each custom module in `x/<module>/` follows the standard Cosmos SDK module structure. This modular approach allows each component to be developed, tested, and maintained independently:

```
x/<module>/
├── keeper/         # State management, message handlers (msg_server_*.go),
│                   # query handlers (query_*.go), genesis, and tests
├── module/         # AppModule definition, depinject wiring (depinject.go),
│                   # AutoCLI configuration (autocli.go), simulation setup
├── types/          # Protobuf-generated types (*_pb.go), codec registration,
│                   # error definitions, store keys, parameters, genesis types
└── simulation/     # Simulation test helpers for fuzz testing
```

The modules `oracle`, `creditscore`, `realestate`, `tokenization`, and `insurance` are structurally identical — each manages a single `collections.Map` collection keyed by `symbol` or `policy_id` (string). The `realfin` module is the base module and contains only governance-controlled `Params` (no data map).

### Standard Cosmos Modules

In addition to the six custom modules, Realfin includes the full suite of standard Cosmos SDK modules. These provide the foundational blockchain functionality — account management, token transfers, staking, governance, and more. All modules are registered via depinject in `app/app_config.go`:

`auth`, `authz`, `bank`, `consensus`, `distribution`, `epochs`, `evidence`, `feegrant`, `genutil`, `gov`, `group`, `mint`, `nft`, `params`, `slashing`, `staking`, `upgrade`, `circuit`, `vesting`

### Execution Order

The Cosmos SDK processes module logic at specific points in the block lifecycle. The execution order is configured in `app/app_config.go` and determines when each module's hooks are called:

| Phase | Purpose | Modules (in order) |
|---|---|---|
| **PreBlockers** | Critical operations before block processing | `upgrade`, `auth` |
| **BeginBlockers** | Start-of-block logic (inflation, slashing, etc.) | `mint` → `distribution` → `slashing` → `evidence` → `staking` → `authz` → `epochs` → `ibc` → `realfin` → `oracle` → `creditscore` → `realestate` → `tokenization` → `insurance` |
| **EndBlockers** | End-of-block logic (governance tallying, etc.) | `gov` → `staking` → `feegrant` → `group` → `realfin` → `oracle` → `creditscore` → `realestate` → `tokenization` → `insurance` |
| **InitGenesis** | One-time initialization from genesis state | `consensus` → `auth` → `bank` → `distribution` → `staking` → `slashing` → `gov` → `mint` → `genutil` → `evidence` → `authz` → `feegrant` → `vesting` → `nft` → `group` → `upgrade` → `circuit` → `epochs` → `ibc` → `transfer` → `interchainaccounts` → `realfin` → `oracle` → `creditscore` → `realestate` → `tokenization` → `insurance` |

The custom modules (`realfin`, `oracle`, `creditscore`, `realestate`, `tokenization`, `insurance`) are always positioned after all standard Cosmos modules, ensuring that the core blockchain infrastructure is fully initialized before custom logic executes.

### Keeper Model

The keeper is the central component of each module, responsible for reading from and writing to the module's state. Each keeper contains:

- `storeService` — provides access to the module's dedicated KV store partition
- `cdc` — codec for binary serialization/deserialization of state objects
- `addressCodec` — handles conversion between human-readable addresses and raw bytes
- `authority` — the address authorized to update module parameters (defaults to the `x/gov` module account)

State is managed via `cosmossdk.io/collections`, a type-safe abstraction over raw KV store operations:

- `Params` — `collections.Item` stores a single parameters object for each module
- `Price` / `Rate` / `Asset` / `Policy` — `collections.Map[string, Type]` stores data entries keyed by a string symbol or policy ID, supporting efficient lookup, iteration, and pagination

### depinject Wiring

Realfin uses the Cosmos SDK's dependency injection framework (`depinject`) to automatically wire modules together. Each module's `depinject.go` file calls `appconfig.Register()` in an `init()` function and defines `ModuleInputs` (dependencies the module needs) and `ModuleOutputs` (what the module provides to others). The `OnePerModuleType` marker ensures each module is instantiated exactly once.

The `App` struct in `app/app.go` is assembled via `depinject.Inject()`, which resolves all keeper dependencies and assigns them to the application. The `runtime.App` is then constructed with `appBuilder.Build()`. Optimistic execution is enabled via `baseapp.SetOptimisticExecution()`, allowing the node to begin executing the next block's transactions speculatively while the previous block is being committed.

### IBC Integration

Realfin is IBC-enabled, allowing it to communicate with other Cosmos SDK-based blockchains. This is critical for cross-chain asset transfers, interchain accounts, and future interoperability with the broader Cosmos ecosystem.

IBC modules are registered manually (not via depinject) in `app/ibc.go`, because the IBC module does not yet support the app wiring framework. The following IBC components are configured:

- **IBCKeeper** — core IBC protocol implementation, managing connections, channels, and packet routing
- **TransferKeeper** — handles IBC token transfers between chains. Supports both IBC v1 (classic) and IBC v2 routing, meaning tokens can be sent to and received from any IBC-connected chain
- **ICAHostKeeper** — allows external chains to execute transactions on Realfin via interchain accounts (host side)
- **ICAControllerKeeper** — allows Realfin to control accounts on external chains via interchain accounts (controller side)
- **Light clients**: Tendermint (`07-tendermint`) for standard Cosmos chains, and Solo Machine (`06-solomachine`) for single-signer scenarios

## Getting Started

### Prerequisites

- **Go 1.25+** — the minimum Go version required to build the project
- **make** — for running build and test commands
- **git** — for cloning the repository
- **Ignite CLI** (optional) — for quick-start development chain and protobuf generation. Install from [docs.ignite.com](https://docs.ignite.com)

### Build and Install

Clone the repository and build the `realfind` binary:

```bash
git clone <repository-url>
cd realfin-core
make install
```

The `make install` command first verifies that dependencies have not been modified (`go mod verify`), then compiles and installs the `realfind` binary to your `$GOPATH/bin` directory. Make sure `$GOPATH/bin` is in your system's `PATH`.

### Verify Installation

After building, confirm the binary is available:

```bash
realfind version
```

This prints the version string, which includes the git branch and commit hash.

### Development Chain Configuration

Realfin includes a `config.yml` file that defines a local development chain with pre-funded test accounts. This allows developers to quickly spin up a local node for testing without manual genesis configuration.

| Setting | Value |
|---|---|
| Default denomination | `urlf` — the smallest unit of the native token |
| Test accounts | `alice` (200,000,000 urlf + 20,000 token), `bob` (100,000,000 urlf + 10,000 token) |
| Validators | `alice` (bonded 100M urlf), `validator1` (bonded 200M urlf), `validator2` (bonded 100M urlf) |
| Faucet | Operated by `bob` — dispenses 5 token + 100,000 urlf per request |
| OpenAPI spec | Generated at `docs/static/openapi.json` |

The development chain runs in isolation on your local machine, providing a sandbox environment for testing transactions, module interactions, and smart contract logic without connecting to a live network.

### Running a Local Development Chain

There are several ways to start a local chain for development and testing.

#### Quick Start with Ignite (Recommended)

The fastest way to start a development chain is with Ignite CLI. A single command installs dependencies, builds the binary, initializes the chain using `config.yml`, creates the test accounts and validators defined there, and starts the node:

```bash
ignite chain serve
```

This reads `config.yml` and automatically:
- Compiles and installs the `realfind` binary
- Initializes the chain with the default denomination `urlf`
- Creates the pre-configured accounts (`alice`, `bob`) and funds them
- Sets up validators (`alice`, `validator1`, `validator2`) with their bonded stakes
- Starts the faucet (operated by `bob`)
- Generates the OpenAPI spec at `docs/static/openapi.json`
- Starts the blockchain node with hot-reload — when you modify Go source files, Ignite automatically rebuilds and restarts the chain

The chain resets on each restart. To preserve state between restarts, use `ignite chain serve --reset-once` (resets only on first run).

To stop the chain, press `Ctrl+C` in the terminal.

#### Single-Node Manual Setup

If you prefer not to use Ignite, you can manually initialize and start a single-node chain using the standard Cosmos SDK commands:

```bash
# 1. Initialize the node with a moniker (display name)
realfind init my-node --chain-id realfin-local-1

# 2. Create a key pair for your test account
realfind keys add alice

# 3. Add a genesis account with initial funds
realfind genesis add-genesis-account $(realfind keys show alice -a) 200000000urlf,20000token

# 4. Create a genesis validator transaction
realfind genesis gentx alice 100000000urlf --chain-id realfin-local-1

# 5. Collect genesis transactions into the genesis file
realfind genesis collect-gentxs

# 6. Validate the genesis file
realfind genesis validate

# 7. Start the node
realfind start
```

Once started, the node produces blocks and exposes the following default endpoints:

| Service | Endpoint |
|---|---|
| RPC (CometBFT) | `tcp://localhost:26657` |
| REST API (gRPC-gateway) | `http://localhost:1317` |
| gRPC | `localhost:9090` |

You can verify the node is running with:

```bash
realfind status
```

#### Multi-Node Local Testnet

For testing with multiple validators running as separate processes (e.g., via Docker Compose), use the `multi-node` testnet command. It generates configuration directories for each validator node:

```bash
realfind testnet multi-node \
  --v 4 \
  --output-dir ./.testnets \
  --validators-stake-amount 1000000,200000,300000,400000 \
  --list-ports 47222,50434,52851,44210 \
  --chain-id realfin-testnet-1
```

| Flag | Description |
|---|---|
| `--v` | Number of validator nodes to generate |
| `--output-dir` | Directory where node configurations will be written |
| `--validators-stake-amount` | Comma-separated stake amounts for each validator |
| `--list-ports` | Comma-separated custom ports for each node |
| `--starting-ip-address` | Base IP for peer addresses (default: localhost) |
| `--chain-id` | Chain identifier for the testnet |

This creates one subdirectory per validator inside `--output-dir`, each containing a full node configuration (private keys, genesis file, `config.toml`, `app.toml`). Start each validator as a separate process pointing to its directory.

#### In-Place Testnet (from Existing State)

For advanced scenarios such as testing upgrades against real chain state, the `in-place-testnet` command replaces the existing validator set with a single local validator:

```bash
realfind in-place-testnet [new-chain-id] [validator-operator-address] \
  --home $HOME/.realfind/validator1 \
  --validator-privkey=[base64-encoded-key] \
  --accounts-to-fund="cosmos1abc...,cosmos1def..."
```

This modifies both application and consensus stores to remove the old validator set and introduce a new one suitable for local testing. The `--accounts-to-fund` flag accepts a comma-separated list of addresses that will be funded with tokens for testing.

## Custom Modules

Realfin extends the Cosmos SDK with six custom modules that provide on-chain data infrastructure. Five of these modules (`oracle`, `creditscore`, `realestate`, `tokenization`, `insurance`) follow an identical CRUD pattern for managing data entries, while the sixth (`realfin`) serves as the base module with governance-controlled parameters.

Each data module stores its entries in a `collections.Map` keyed by a string symbol. All entries include a `creator` field that tracks the address which originally created the record, enabling ownership-based access control for updates and deletions.

### Oracle (`x/oracle`) — Price Data

The oracle module provides on-chain storage for price feed data. It allows authorized users to publish, update, and remove price entries identified by a unique symbol. This module serves as the data source for asset pricing within the Realfin ecosystem.

**Entity: Price**

Each price entry consists of the following fields:

| Field | Type | Description |
|---|---|---|
| `symbol` | `string` | Unique identifier for the price entry (e.g., `ETH`, `BTC`). Used as the map key — must be unique across all entries in this module. |
| `rate` | `uint64` | The price value. Stored as an unsigned 64-bit integer to avoid floating-point precision issues. Applications should define their own decimal scaling convention. |
| `name` | `string` | A human-readable name for the asset (e.g., `Ethereum`, `Bitcoin`). Informational only — not used for lookups. |
| `description` | `string` | A free-text description providing additional context about the price entry. |
| `creator` | `string` | The bech32-encoded address of the account that created this entry. This address is the owner — only the creator can update or delete the entry. |

**Transaction Commands:**

```bash
# Create a new price entry. The symbol must not already exist.
# All four positional arguments are required.
realfind tx oracle create-price [symbol] [rate] [name] [description] --from <key>

# Update an existing price entry. The symbol must exist, and the --from address
# must match the original creator. All fields are overwritten.
realfind tx oracle update-price [symbol] [rate] [name] [description] --from <key>

# Delete a price entry. The symbol must exist, and the --from address
# must match the original creator.
realfind tx oracle delete-price [symbol] --from <key>
```

**Query Commands:**

```bash
# Retrieve a single price entry by its symbol.
# Aliases: get-price, show-price
realfind q oracle get-price [symbol]

# List all price entries with pagination support.
# Supports standard Cosmos pagination flags: --limit, --offset, --count-total
realfind q oracle list-price

# Show the oracle module's current parameters.
realfind q oracle params
```

**Example usage:**

```bash
# Create a price entry for ETH
realfind tx oracle create-price ETH 5001 "Ethereum" "Ethereum spot price" --from alice

# Query the price
realfind q oracle get-price ETH

# Update the price
realfind tx oracle update-price ETH 5200 "Ethereum" "Updated ETH price" --from alice

# List all prices
realfind q oracle list-price

# Delete the entry
realfind tx oracle delete-price ETH --from alice
```

**Access control:** Only the original creator (the address that submitted the `create-price` transaction) can update or delete a price entry. Attempting to modify another user's entry returns an `ErrUnauthorized` error.

---

### Creditscore (`x/creditscore`) — Credit Ratings

The creditscore module provides on-chain storage for credit rating data. It enables the publication of credit scores for entities identified by a unique symbol. Within the Realfin ecosystem, credit ratings support risk assessment by assigning probability-of-default (PD) scores to assets, helping investors evaluate risk-return profiles.

The module is structurally identical to oracle but uses the `Rate` entity name instead of `Price`.

**Entity: Rate**

| Field | Type | Description |
|---|---|---|
| `symbol` | `string` | Unique identifier for the credit rating entry (e.g., `SME-001`, `BOND-XYZ`). Used as the map key. |
| `rate` | `uint64` | The credit rating value. Interpretation is application-defined — could represent a PD score, a numeric grade, or a custom metric. |
| `name` | `string` | A human-readable name for the rated entity. |
| `description` | `string` | Additional context about the credit rating — methodology, date, scope, etc. |
| `creator` | `string` | The bech32-encoded address of the rating publisher. Only this address can update or delete the entry. |

**Transaction Commands:**

```bash
# Create a new credit rating. The symbol must not already exist.
realfind tx creditscore create-rate [symbol] [rate] [name] [description] --from <key>

# Update an existing credit rating. Requires creator ownership.
realfind tx creditscore update-rate [symbol] [rate] [name] [description] --from <key>

# Delete a credit rating. Requires creator ownership.
realfind tx creditscore delete-rate [symbol] --from <key>
```

**Query Commands:**

```bash
# Retrieve a single credit rating by symbol.
# Aliases: get-rate, show-rate
realfind q creditscore get-rate [symbol]

# List all credit ratings with pagination.
realfind q creditscore list-rate

# Show the creditscore module's current parameters.
realfind q creditscore params
```

**Example usage:**

```bash
# Publish a credit rating
realfind tx creditscore create-rate SME-001 850 "Acme Corp" "Annual PD assessment" --from alice

# Query the rating
realfind q creditscore get-rate SME-001

# Update the rating after reassessment
realfind tx creditscore update-rate SME-001 870 "Acme Corp" "Updated Q2 assessment" --from alice

# List all ratings
realfind q creditscore list-rate
```

**Access control:** Only the original creator can update or delete a rate entry.

---

### Realestate (`x/realestate`) — Real Estate Ratings

The realestate module provides on-chain storage for real estate valuation and rating data. It enables the publication of property assessments identified by a unique symbol. Within the Realfin ecosystem, real estate ratings support the tokenization of property assets by providing transparent, on-chain valuation data.

The module is structurally identical to creditscore — same entity structure, same CRUD operations, same access control model.

**Entity: Rate**

| Field | Type | Description |
|---|---|---|
| `symbol` | `string` | Unique identifier for the real estate entry (e.g., `PROP-SF-101`, `LAND-BG-42`). Used as the map key. |
| `rate` | `uint64` | The real estate rating or valuation value. |
| `name` | `string` | A human-readable name for the property or asset. |
| `description` | `string` | Additional context — location, property type, valuation methodology, etc. |
| `creator` | `string` | The bech32-encoded address of the entity that published this rating. Only this address can modify or remove the entry. |

**Transaction Commands:**

```bash
# Create a new real estate rating. The symbol must not already exist.
realfind tx realestate create-rate [symbol] [rate] [name] [description] --from <key>

# Update an existing real estate rating. Requires creator ownership.
realfind tx realestate update-rate [symbol] [rate] [name] [description] --from <key>

# Delete a real estate rating. Requires creator ownership.
realfind tx realestate delete-rate [symbol] --from <key>
```

**Query Commands:**

```bash
# Retrieve a single real estate rating by symbol.
# Aliases: get-rate, show-rate
realfind q realestate get-rate [symbol]

# List all real estate ratings with pagination.
realfind q realestate list-rate

# Show the realestate module's current parameters.
realfind q realestate params
```

**Example usage:**

```bash
# Publish a property valuation
realfind tx realestate create-rate PROP-SF-101 2500000 "123 Main St" "Commercial property in SF" --from alice

# Query the valuation
realfind q realestate get-rate PROP-SF-101

# Update after reappraisal
realfind tx realestate update-rate PROP-SF-101 2650000 "123 Main St" "Q3 reappraisal" --from alice

# List all property ratings
realfind q realestate list-rate
```

**Access control:** Only the original creator can update or delete a rate entry.

---

### Tokenization (`x/tokenization`) — Asset Tokenization

The tokenization module provides on-chain storage for tokenized real-world asset (RWA) metadata. It enables authorized users to register, update, and remove asset entries identified by a unique symbol. Within the Realfin ecosystem, this module serves as the registry for tokenized assets — recording their classification, provenance, and descriptive metadata on-chain.

Unlike the oracle, creditscore, and realestate modules which use a `uint64` rate/price field, the tokenization module uses only string fields, making it suitable for rich metadata storage including JSON-encoded provenance and classification data.

**Entity: Asset**

| Field | Type | Description |
|---|---|---|
| `symbol` | `string` | Unique identifier for the tokenized asset (e.g., `RWA-SF-101`, `INV-2024-001`). Used as the map key — must be unique across all entries in this module. |
| `name` | `string` | A human-readable name for the asset (e.g., `Main Street Property`, `Acme Inventory Q4`). |
| `description` | `string` | A free-text description providing additional context about the tokenized asset. |
| `asset_type` | `string` | Classification of the asset. Recommended values: `real_estate`, `inventory`, `invoice`, `ip`, `receivable` — but the field is free-form and application-defined. |
| `metadata` | `string` | Embedded metadata for provenance, classification, and additional structured data. Typically a JSON string (e.g., `{"location":"Sofia","appraised_value":"500000"}`). |
| `creator` | `string` | The bech32-encoded address of the account that registered this asset. This address is the owner — only the creator can update or delete the entry. |

**Transaction Commands:**

```bash
# Register a new tokenized asset. The symbol must not already exist.
# All five positional arguments are required.
realfind tx tokenization create-asset [symbol] [name] [description] [asset_type] [metadata] --from <key>

# Update an existing tokenized asset. The symbol must exist, and the --from address
# must match the original creator. All fields are overwritten.
realfind tx tokenization update-asset [symbol] [name] [description] [asset_type] [metadata] --from <key>

# Delete a tokenized asset entry. The symbol must exist, and the --from address
# must match the original creator.
realfind tx tokenization delete-asset [symbol] --from <key>
```

**Query Commands:**

```bash
# Retrieve a single asset entry by its symbol.
# Aliases: get-asset, show-asset
realfind q tokenization get-asset [symbol]

# List all asset entries with pagination support.
# Supports standard Cosmos pagination flags: --limit, --offset, --count-total
realfind q tokenization list-asset

# Show the tokenization module's current parameters.
realfind q tokenization params
```

**Example usage:**

```bash
# Register a tokenized real estate asset
realfind tx tokenization create-asset RWA-SF-101 "123 Main St" "Commercial property in SF" real_estate '{"location":"San Francisco","sqft":5000}' --from alice

# Query the asset
realfind q tokenization get-asset RWA-SF-101

# Update the asset metadata after reappraisal
realfind tx tokenization update-asset RWA-SF-101 "123 Main St" "Commercial property in SF - reappraised" real_estate '{"location":"San Francisco","sqft":5000,"appraised_value":"2650000"}' --from alice

# List all tokenized assets
realfind q tokenization list-asset

# Remove the asset entry
realfind tx tokenization delete-asset RWA-SF-101 --from alice
```

**Access control:** Only the original creator (the address that submitted the `create-asset` transaction) can update or delete an asset entry. Attempting to modify another user's entry returns an `ErrUnauthorized` error.

---

### Insurance (`x/insurance`) — Insurance Policies

The insurance module provides on-chain storage for insurance policies linked to tokenized real-world assets. Insurance providers can register coverage details against any tokenized asset, creating a transparent and auditable record of which assets are insured, by whom, and to what extent. This module is a foundational building block for asset grading and risk assessment on the Realfin platform, enabling the creation of insured and non-insured asset tranches.

Like the tokenization module, insurance uses only string fields, making it suitable for flexible metadata about coverage terms and provider information.

**Entity: Policy**

| Field | Type | Description |
|---|---|---|
| `policy_id` | `string` | Unique identifier for the insurance policy (e.g., `POL-001`, `INS-RWA-101`). Used as the map key — must be unique across all entries in this module. |
| `asset_symbol` | `string` | The symbol of the tokenized asset being insured (e.g., `RWA-SF-101`). This is a reference to an entry in the tokenization module, though it is not enforced at the protocol level. |
| `provider` | `string` | Name of the insurance company or underwriter providing the coverage. |
| `coverage_type` | `string` | Classification of the coverage. Recommended values: `full`, `partial` — but the field is free-form and application-defined. |
| `coverage_percentage` | `string` | The percentage of asset value covered by this policy (e.g., `100`, `75`). Stored as a string to allow flexible formatting. |
| `creator` | `string` | The bech32-encoded address of the account that registered this policy. This address is the owner — only the creator can update or delete the entry. |

**Transaction Commands:**

```bash
# Create a new insurance policy. The policy_id must not already exist.
# All five positional arguments are required.
realfind tx insurance create-policy [policy_id] [asset_symbol] [provider] [coverage_type] [coverage_percentage] --from <key>

# Update an existing insurance policy. The policy_id must exist, and the --from address
# must match the original creator. All fields are overwritten.
realfind tx insurance update-policy [policy_id] [asset_symbol] [provider] [coverage_type] [coverage_percentage] --from <key>

# Delete an insurance policy. The policy_id must exist, and the --from address
# must match the original creator.
realfind tx insurance delete-policy [policy_id] --from <key>
```

**Query Commands:**

```bash
# Retrieve a single policy by its policy_id.
# Aliases: get-policy, show-policy
realfind q insurance get-policy [policy_id]

# List all policies with pagination support.
# Supports standard Cosmos pagination flags: --limit, --offset, --count-total
realfind q insurance list-policy

# Show the insurance module's current parameters.
realfind q insurance params
```

**Example usage:**

```bash
# Create an insurance policy for a tokenized property
realfind tx insurance create-policy POL-001 RWA-SF-101 "AIG Insurance" full 100 --from alice

# Query the policy
realfind q insurance get-policy POL-001

# Update the coverage terms
realfind tx insurance update-policy POL-001 RWA-SF-101 "AIG Global" full 95 --from alice

# List all policies
realfind q insurance list-policy

# Remove the policy
realfind tx insurance delete-policy POL-001 --from alice
```

**Access control:** Only the original creator (the address that submitted the `create-policy` transaction) can update or delete a policy entry. Attempting to modify another user's entry returns an `ErrUnauthorized` error.

---

### Realfin (`x/realfin`) — Base Module

The realfin module is the base module of the chain. Unlike the four data modules above, it does not manage any data map — it exists solely to hold governance-controlled module parameters.

```bash
# Query the base module's parameters
realfind q realfin params
```

The `UpdateParams` message is restricted to the `x/gov` module authority address. It is not exposed via CLI (marked with `Skip: true` in AutoCLI) — parameter updates can only be submitted through a governance proposal. This ensures that any changes to the base module's configuration require network-wide consensus.

---

### CRUD Operation Pattern

All five data modules (oracle, creditscore, realestate, tokenization, insurance) implement identical handler logic for their Create, Update, and Delete operations. This consistency simplifies development and ensures predictable behavior across the entire data layer:

| Operation | Handler Logic |
|---|---|
| **Create** | 1. Validate the sender's address format. 2. Check if the symbol already exists via `Map.Has()` — reject with `ErrInvalidRequest` if duplicate. 3. Construct the entity with all provided fields. 4. Write to state via `Map.Set()`. |
| **Update** | 1. Validate the sender's address format. 2. Load the existing entry via `Map.Get()` — reject with `ErrKeyNotFound` if missing. 3. Verify the sender matches the stored `creator` — reject with `ErrUnauthorized` if mismatch. 4. Overwrite all fields and save via `Map.Set()`. |
| **Delete** | 1. Validate the sender's address format. 2. Load the existing entry via `Map.Get()` — reject with `ErrKeyNotFound` if missing. 3. Verify the sender matches the stored `creator` — reject with `ErrUnauthorized` if mismatch. 4. Remove from state via `Map.Remove()`. |

**Query operations** are provided via a separate `queryServer` that embeds the keeper:

- **Get**: retrieves a single entry by symbol via `Map.Get()`
- **List**: returns all entries with automatic pagination via `query.CollectionPaginate()`

## CLI Reference

The `realfind` binary serves as both the node daemon and the CLI client. All interactions with the blockchain — submitting transactions, querying state, managing keys — go through this single binary.

### General Syntax

```bash
# Submit a transaction (state-changing operation)
realfind tx <module> <command> [args] --from <key> [flags]

# Query state (read-only operation)
realfind q <module> <command> [args] [flags]
```

**Common flags for transactions:**

| Flag | Description |
|---|---|
| `--from <key>` | Name or address of the signing key (required for all transactions) |
| `--chain-id <id>` | Chain identifier (required when not configured in client config) |
| `--gas auto` | Automatically estimate gas required |
| `--fees <amount>` | Transaction fee (e.g., `500urlf`) |
| `--broadcast-mode sync` | Wait for tx to be included in the mempool (default) |
| `--yes` | Skip confirmation prompt |

**Common flags for queries:**

| Flag | Description |
|---|---|
| `--output json` | Output results in JSON format |
| `--limit <n>` | Maximum number of results per page (for list queries) |
| `--offset <n>` | Number of results to skip (for list queries) |
| `--count-total` | Include total count in paginated responses |
| `--node <url>` | RPC endpoint to query (default: `tcp://localhost:26657`) |

### Module Commands Summary

| Module | Transaction Commands | Query Commands |
|---|---|---|
| `oracle` | `create-price`, `update-price`, `delete-price` | `get-price` (alias: `show-price`), `list-price`, `params` |
| `creditscore` | `create-rate`, `update-rate`, `delete-rate` | `get-rate` (alias: `show-rate`), `list-rate`, `params` |
| `realestate` | `create-rate`, `update-rate`, `delete-rate` | `get-rate` (alias: `show-rate`), `list-rate`, `params` |
| `tokenization` | `create-asset`, `update-asset`, `delete-asset` | `get-asset` (alias: `show-asset`), `list-asset`, `params` |
| `insurance` | `create-policy`, `update-policy`, `delete-policy` | `get-policy` (alias: `show-policy`), `list-policy`, `params` |
| `realfin` | — | `params` |

### Standard Node Commands

In addition to the module-specific commands, `realfind` provides all standard Cosmos SDK node management commands:

```bash
# Node initialization and operation
realfind init [moniker]         # Initialize a new node with a given moniker name
realfind start                  # Start the blockchain node
realfind status                 # Query the current status of the node

# Key management — create, import, export, and list signing keys
realfind keys add <name>        # Create a new key
realfind keys list              # List all keys
realfind keys show <name>       # Show key details and address
realfind keys delete <name>     # Delete a key

# Genesis file management
realfind genesis [subcommand]   # Genesis file utilities (add accounts, gentx, validate, etc.)

# Data management
realfind snapshot [subcommand]  # Create and manage state snapshots for fast sync
realfind pruning [subcommand]   # Configure state pruning strategies
realfind config [subcommand]    # Application configuration management (confix)

# Blockchain queries
realfind query block <height>   # Query a block by height
realfind query blocks           # Query blocks with filtering
realfind query tx <hash>        # Query a transaction by hash
realfind query txs --events ... # Query transactions by events

# Transaction utilities
realfind tx sign                # Sign a transaction offline
realfind tx broadcast           # Broadcast a signed transaction
realfind tx encode              # Encode a transaction to amino bytes
realfind tx decode              # Decode an amino-encoded transaction

# Testing and development
realfind testnet [subcommand]   # Set up single-node or multi-node test networks
```

## API Endpoints

Realfin exposes REST and gRPC APIs for programmatic access to all module functionality. REST endpoints are automatically generated from protobuf definitions via gRPC-gateway, meaning every gRPC service method has a corresponding HTTP endpoint.

### REST (gRPC-gateway)

All REST endpoints use the GET method and return JSON responses. The base URL depends on your node's API server configuration (default: `http://localhost:1317`).

**Oracle module:**

| Endpoint | Description |
|---|---|
| `/realfin/oracle/v1/params` | Returns the oracle module's current parameters. |
| `/realfin/oracle/v1/price/{symbol}` | Returns a single price entry by its symbol. The `{symbol}` path parameter is the unique identifier used when the price was created. |
| `/realfin/oracle/v1/price` | Returns all price entries with pagination. Accepts optional query parameters: `pagination.limit`, `pagination.offset`, `pagination.count_total`. |

**Creditscore module:**

| Endpoint | Description |
|---|---|
| `/realfin/creditscore/v1/params` | Returns the creditscore module's current parameters. |
| `/realfin/creditscore/v1/rate/{symbol}` | Returns a single credit rating by its symbol. |
| `/realfin/creditscore/v1/rate` | Returns all credit ratings with pagination support. |

**Realestate module:**

| Endpoint | Description |
|---|---|
| `/realfin/realestate/v1/params` | Returns the realestate module's current parameters. |
| `/realfin/realestate/v1/rate/{symbol}` | Returns a single real estate rating by its symbol. |
| `/realfin/realestate/v1/rate` | Returns all real estate ratings with pagination support. |

**Tokenization module:**

| Endpoint | Description |
|---|---|
| `/realfin/tokenization/v1/params` | Returns the tokenization module's current parameters. |
| `/realfin/tokenization/v1/asset/{symbol}` | Returns a single tokenized asset entry by its symbol. |
| `/realfin/tokenization/v1/asset` | Returns all tokenized asset entries with pagination support. |

**Insurance module:**

| Endpoint | Description |
|---|---|
| `/realfin/insurance/v1/params` | Returns the insurance module's current parameters. |
| `/realfin/insurance/v1/policy/{policy_id}` | Returns a single insurance policy by its policy ID. |
| `/realfin/insurance/v1/policy` | Returns all insurance policies with pagination support. |

**Realfin base module:**

| Endpoint | Description |
|---|---|
| `/realfin/realfin/v1/params` | Returns the base module's current parameters. |

### gRPC Services

Each module exposes both a `Query` and a `Msg` service via gRPC. The `Query` service handles read-only state queries, while the `Msg` service processes state-changing transactions:

| Module | Query Service | Msg Service |
|---|---|---|
| Oracle | `realfin.oracle.v1.Query` | `realfin.oracle.v1.Msg` |
| Creditscore | `realfin.creditscore.v1.Query` | `realfin.creditscore.v1.Msg` |
| Realestate | `realfin.realestate.v1.Query` | `realfin.realestate.v1.Msg` |
| Tokenization | `realfin.tokenization.v1.Query` | `realfin.tokenization.v1.Msg` |
| Insurance | `realfin.insurance.v1.Query` | `realfin.insurance.v1.Msg` |
| Realfin | `realfin.realfin.v1.Query` | `realfin.realfin.v1.Msg` |

The default gRPC port is `9090`. You can use any gRPC client (e.g., `grpcurl`) to interact with these services directly.

### OpenAPI / Swagger

The full OpenAPI specification is available at `docs/static/openapi.json`. When the node's API server is running with Swagger enabled, an interactive documentation interface is served at the API endpoint, allowing you to browse and test all available REST endpoints in a web browser.

## IBC Integration

Realfin's IBC integration enables cross-chain communication with any blockchain in the Cosmos ecosystem that supports the Inter-Blockchain Communication protocol. This is a foundational capability for RWA tokenization, as it allows tokenized assets, credit ratings, and price data to be referenced and utilized across multiple interconnected chains.

The IBC stack is configured in `app/ibc.go` and includes:

### Token Transfers

The **Transfer module** allows users to send and receive tokens between Realfin and any IBC-connected chain. Both IBC v1 (classic packet-based) and IBC v2 routing are supported, providing broad compatibility with existing Cosmos chains and forward compatibility with the evolving IBC protocol.

### Interchain Accounts (ICA)

Interchain Accounts enable remote execution of transactions across chains:

- **Host module**: Allows external chains to register and control accounts on Realfin. An account on another chain can execute Realfin transactions (e.g., creating a price entry) without the user needing a local Realfin account.
- **Controller module**: Allows Realfin accounts to register and control accounts on other IBC-connected chains. This enables Realfin-based applications to interact with DeFi protocols, governance systems, or asset registries on external chains.

### Light Clients

Two light client types are registered for verifying the state of connected chains:

- **Tendermint (`07-tendermint`)** — the standard light client for Cosmos SDK chains using CometBFT consensus. Used for most IBC connections.
- **Solo Machine (`06-solomachine`)** — a light client for single-signer scenarios, useful for connecting off-chain systems or custodial services to the blockchain.

## Governance

Realfin inherits the full Cosmos SDK governance framework, which allows token holders to propose and vote on protocol changes. Within the custom modules, governance plays a specific role in parameter management.

Module parameters for all six custom modules (`realfin`, `oracle`, `creditscore`, `realestate`, `tokenization`, `insurance`) are updated exclusively through governance proposals. The `UpdateParams` message in each module requires the sender to match the `x/gov` module authority address — any other sender is rejected with `ErrUnauthorized`.

The `UpdateParams` CLI command is intentionally hidden from AutoCLI (`Skip: true` in all six modules). This is a deliberate design choice: parameter changes affect the entire network and should go through the standard governance process (submit proposal → deposit period → voting period → execution), rather than being callable directly from the CLI.

The governance flow for updating module parameters:

1. A governance proposal is submitted containing the `MsgUpdateParams` message with the new parameter values
2. Token holders deposit `urlf` to meet the minimum deposit threshold
3. The proposal enters the voting period, during which staked token holders vote
4. If the proposal passes, the `UpdateParams` message is executed with the governance module's authority address as the sender
5. The module's parameters are updated on-chain

## Development

### Makefile Commands

The `Makefile` provides all common development workflows. Run `make <target>` from the repository root:

| Command | Description |
|---|---|
| `make install` | Verify dependencies and build + install the `realfind` binary to `$GOPATH/bin`. This is the primary build command. |
| `make test` | Run the full test suite: `go vet` for static analysis, `govulncheck` for known vulnerabilities, then all unit tests. Use this before committing. |
| `make test-unit` | Run only unit tests (`go test -mod=readonly -v -timeout 30m ./...`). Faster than `make test` when you've already run vet/vulncheck. |
| `make test-race` | Run unit tests with Go's race condition detector enabled. Use this to catch concurrent access bugs. Slower than regular tests. |
| `make test-cover` | Run unit tests with code coverage analysis. Generates an HTML coverage report. |
| `make bench` | Run benchmarks across all packages. Useful for measuring performance of keeper operations. |
| `make lint` | Run `golangci-lint` with a 15-minute timeout. Checks for style issues, potential bugs, and code quality problems. |
| `make lint-fix` | Same as `make lint`, but automatically fixes issues where possible. |
| `make govet` | Run `go vet` separately. Checks for suspicious constructs that the compiler doesn't catch. |
| `make proto-gen` | Regenerate all Go code from protobuf definitions. Requires Ignite CLI to be installed. Run this after modifying any `.proto` file. |

### Running Individual Tests

To run a specific test function within a module, use the Go test command with the `-run` flag:

```bash
# Run a specific test by name (supports regex patterns)
go test -mod=readonly -v -timeout 30m -run TestFunctionName ./x/oracle/keeper/...

# Run all tests in the oracle keeper package
go test -mod=readonly -v -timeout 30m ./x/oracle/keeper/...

# Run all tests across the entire project
go test -mod=readonly -v -timeout 30m ./...
```

The `-mod=readonly` flag prevents accidental modification of `go.mod` during test runs. The `-timeout 30m` flag sets a generous timeout for the full test suite.

### Test Infrastructure

Tests use the **fixture pattern** where each test constructs an isolated testing environment:

1. An in-memory store is created via `testutil.DefaultContextWithDB()`, providing a clean state for each test
2. A keeper is instantiated with a test codec from `moduletestutil.MakeTestEncodingConfig()`, wiring up all necessary type registrations
3. Assertions are written using the `testify` library (`require` for fatal checks, `assert` for non-fatal checks)

This approach ensures tests are fast (no disk I/O), isolated (no shared state between tests), and deterministic (same inputs always produce the same results).

### Simulation Testing

Each module includes simulation support in `module/simulation.go`, which defines weighted randomized operations for Create, Update, and Delete. The `SimulationManager` coordinates:

- **Random genesis generation** — creates valid initial state with random data
- **Weighted operations** — each CRUD operation has a configurable probability of being executed during simulation
- **State validation** — invariants are checked after each operation to ensure the module's state remains consistent

Simulation testing is particularly valuable for discovering edge cases and race conditions that unit tests might miss, as it exercises the modules with thousands of randomized operations.

### Code Quality

- **Linting**: `golangci-lint` with a comprehensive set of rules and a 15-minute timeout for large codebases
- **Static analysis**: `go vet` catches common issues like unreachable code, incorrect format strings, and suspicious constructs
- **Vulnerability scanning**: `govulncheck` checks all dependencies against the Go vulnerability database
- **Proto tooling**: `buf` for protobuf linting, formatting, and breaking change detection
