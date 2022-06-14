# Optimistic Rollup Contracts

## Develepment

### Install Tools

forge: "https://github.com/gakonst/foundry"

nodejs: "https://nodejs.org/en/"

go: "https://go.dev/dl/"

yarn:
```bash
npm install --global yarn
```

go-web3:
```bash
go install github.com/laizy/web3/abigen/cmd@v0.1.9
mv $(go env GOBIN)/cmd $(go env GOBIN)/go-web3
```

### Build

```shell
yarn # install dependences 
yarn format # format contracts code
yarn build # build contracts
yarn test # run contracts testcase
yarn gasnap # run contracts testcase and generate gas snapshot file
yarn clean # clean built contracts
yarn abigen # build go binding
yarn go # build go cmd
```

### Rollup Cli Usage

```shell
cd build
cp ../config/*-config.json .

# deploy l2 contracts, addresses will save to addressl2.json
./rollup deploy l2 -submit

# copy  addressl2.json into rollup-config.json

# deploy l1 contracts, addresses will save to addressl1.json
./rollup deploy l1 -submit

# copy addressl1.json into rollup-config.json

# finish l2 bridge initialization
./rollup deploy l2init -submit

# deposit eth to l2 bridge
./rollup gateway depositEth -amount 1.234 -submit

```
