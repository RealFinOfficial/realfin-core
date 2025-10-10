# realfin
**realfin** is a blockchain built using Cosmos SDK and Tendermint.

## Get started

```
realfin chain
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

### Configure

Your blockchain in development can be configured with `config.yml`.

### Web Frontend
These commands can be run within your scaffolded blockchain project. 

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
curl https://get.ignite.com/username/realfin@latest! | sudo bash
```
`username/realfin` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/allinbits/starport-installer).


### Price oracle

./realfind tx oracle create-price ETH 5001 0 --from alice
./realfind tx oracle update-price ETH 5001 0 --from alice
./realfind tx oracle delete-price ETH --from alice

./realfind q oracle show-price ETH
./realfind q oracle list-price


./realfind tx oracle create-price [symbol] [price] [timestamp]
./realfind tx oracle update-price [symbol] [price] [timestamp]
./realfind tx oracle delete-price [symbol]

./realfind q oracle show-price [symbol]
./realfind q oracle list-price

### Credit Score

./realfind tx creditscore create-rate credit1 1000 name description --from alice
./realfind tx creditscore update-rate credit1 1001 name description --from alice
./realfind tx creditscore delete-rate credit --from alice

./realfind q creditscore show-rate credit1
./realfind q creditscore list-rate

./realfind tx creditscore create-rate [symbol] [rate] [name] [description]
./realfind tx creditscore update-rate [symbol] [rate] [name] [description]
./realfind tx creditscore delete-rate [symbol]

./realfind q creditscore show-rate [symbol]
./realfind q creditscore list-rate

### Real estate

./realfind tx realestate create-rate realstate 1000 name description --from alice
./realfind tx realestate update-rate realstate 1001 name description --from alice
./realfind tx realestate delete-rate realstate --from alice

./realfind q realestate show-rate realstate
./realfind q realestate list-rate

./realfind tx realestate create-rate [symbol] [rate] [name] [description]
./realfind tx realestate update-rate [symbol] [rate] [name] [description]
./realfind tx realestate delete-rate [symbol]

./realfind q realestate show-rate [symbol]
./realfind q realestate list-rate