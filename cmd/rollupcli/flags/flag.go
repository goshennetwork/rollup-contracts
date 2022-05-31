package flags

import "github.com/urfave/cli/v2"

var AmountFlag = &cli.Float64Flag{
	Name:     "amount",
	Usage:    "amount in ether",
	Required: true,
}

var TargetFlag = &cli.StringFlag{
	Name:     "target",
	Usage:    "target contract address.",
	Required: true,
}

var SenderFlag = &cli.StringFlag{
	Name:     "sender",
	Usage:    "sender contract address.",
	Required: true,
}

var QueueIndexFlag = &cli.Uint64Flag{
	Name:     "queueIndex",
	Usage:    "queueIndex",
	Required: true,
}

var MessageFlag = &cli.StringFlag{
	Name:     "message",
	Usage:    "hex message to send to the target.",
	Required: true,
}

var GasLimitFlag = &cli.Uint64Flag{
	Name:     "gasLimit",
	Usage:    "gas limit for the provided message.",
	Required: true,
}

var ToFlag = &cli.StringFlag{
	Name:  "to",
	Usage: "l2 address to receive, default to sender",
}

var AddressFlag = &cli.StringFlag{
	Name:  "address",
	Usage: "address",
}

var SubmitFlag = &cli.BoolFlag{
	Name:  "submit",
	Usage: "submit transaction to remote chain",
}

var ConfigFlag = &cli.StringFlag{
	Name:  "cfg",
	Usage: "specify config file",
	Value: "rollup-config.json",
}

var L1TokenFlag = &cli.StringFlag{
	Name:  "l1Token",
	Usage: "l1Token address",
}

var L2TokenFlag = &cli.StringFlag{
	Name:  "l2Token",
	Usage: "l2Token address",
}

var DataHash = &cli.StringFlag{
	Name:  "datahash",
	Usage: "cross layer message hash",
}
