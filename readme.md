# fundchain
**fundchain** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

## Get started

```
ignite chain serve
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

### Configure

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Web Frontend

Additionally, Ignite CLI offers both Vue and React options for frontend scaffolding:

For a Vue frontend, use: `ignite scaffold vue`
For a React frontend, use: `ignite scaffold react`
These commands can be run within your scaffolded blockchain project. 


For more information see the [monorepo for Ignite front-end development](https://github.com/ignite/web).

## Release
To release a new version of your blockchain, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

### Install
To install the latest version of your blockchain node's binary, execute the following command on your machine:

```
curl https://get.ignite.com/username/fundchain@latest! | sudo bash
```
`username/fundchain` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/allinbits/starport-installer).

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.gg/ignite)

## Demo scenario

This internal walkthrough helps me present FundChain quickly. Run the commands
from the project root in a terminal with [Ignite CLI](https://ignite.com/cli)
and Go installed.

### 1. Start the local chain

```bash
git clone https://github.com/username/FundChain.git
cd FundChain
ignite chain serve
```

The `serve` command builds, initializes, and starts a local node with two
preconfigured accounts: `alice` and `bob`.

### 2. Seed the community pool

Open a new terminal in the same directory and store Alice's and Bob's addresses:

```bash
fundchaind keys list --keyring-backend test
ALICE=$(fundchaind keys show alice -a --keyring-backend test)
BOB=$(fundchaind keys show bob -a --keyring-backend test)
```

Send tokens from Alice to the community pool:

```bash
fundchaind tx distribution fund-community-pool 100token \
  --from alice --fees 1stake --chain-id fundchain --keyring-backend test
```

Check the pool balance:

```bash
fundchaind q distribution community-pool
```

### 3. Bob proposes his project

Bob requests funding from the community pool for his project:

```bash
fundchaind tx gov submit-proposal community-pool-spend $BOB 50token \
  --title "Demo funding" \
  --description "Send 50token from the community pool to Bob" \
  --deposit 10000000stake \
  --from bob --fees 1stake --chain-id fundchain --keyring-backend test
```

The first proposal gets ID `1`:

```bash
PROPOSAL=1
fundchaind q gov proposal $PROPOSAL
```

### 4. Vote and finalize

Approve the proposal:

```bash
fundchaind tx gov vote $PROPOSAL yes --from alice --fees 1stake --chain-id fundchain --keyring-backend test
fundchaind tx gov vote $PROPOSAL yes --from bob --fees 1stake --chain-id fundchain --keyring-backend test
```

Wait for the voting period (about a minute) and confirm it passed:

```bash
fundchaind q gov proposal $PROPOSAL
```

Verify Bob received the funds:

```bash
fundchaind q bank balances $BOB
```

### 5. Stop the network

Return to the terminal running `ignite chain serve` and stop with `Ctrl+C`.
