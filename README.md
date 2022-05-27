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
yarn clean # clean built contracts
yarn abigen # build go binding
```
