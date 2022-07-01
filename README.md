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

## deployments

```json
{
  "addressManager": "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0",
  "l1CrossLayerWitness": "0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9",
  "l2CrossLayerWitness": "0x0165878A594ca255338adfa4d48449f69242Eb8F",
  "feeToken": "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853",
  "rollupStateChain": "0x8A791620dd6260079BF849Dc5567aDC3F2FdC318",
  "challenge": "0x610178dA211FEF7D417bC0e6FeD39F05609AD788",
  "challengeBeacon": "0xB7f8BC63BbcaD18155201308C8f3540b07f84F5e",
  "challengeFactory": "0x0DCd1Bf9A1b36cE34237eEaFef220932846BCD82",
  "stakingManager": "0x9A9f2CCfdE556A7E9Ff0848998Aa4a0CFD8863AE",
  "rollupInputChain": "0x3Aa5ebB10DC797CAC828524e59A333d0A371443c",
  "stateStorageContainer": "0x59b670e9fA9D0A427751Af201D676719a970857b",
  "inputStorageContainer": "0x4ed7c70F96B99c776995fB64377f0d4aB3B0e1C1",
  "machineState": "0x322813Fd9A801c5507c9de605d63CEA4f2CE6c44",
  "stateTransition": "0x4ed7c70F96B99c776995fB64377f0d4aB3B0e1C1",
  "l1StandardBridge": "0x09635F643e140090A9A8Dcd712eD6285858ceBef"
}
```
