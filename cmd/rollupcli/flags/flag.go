package flags

import "github.com/urfave/cli/v2"

var AmountFlag = &cli.Float64Flag{
	Name:     "amount",
	Usage:    "amount in decimal, support float value",
	Required: true,
}

var TargetFlag = &cli.StringFlag{
	Name:     "target",
	Usage:    "target contract address.",
	Required: true,
}

var MessageFlag = &cli.StringFlag{
	Name:     "message",
	Aliases:  []string{"m"},
	Usage:    "hex message to send to the target.",
	Required: true,
}

var ToFlag = &cli.StringFlag{
	Name:  "to",
	Usage: "l2 address to receive, default to sender",
}

var SubmitFlag = &cli.BoolFlag{
	Name:    "submit",
	Aliases: []string{"s"},
	Usage:   "submit transaction to remote chain",
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

var AccountFlag = &cli.StringFlag{
	Name:     "account",
	Usage:    "account address",
	Required: true,
}

var EnabledFlag = &cli.BoolFlag{
	Name:     "enabled",
	Usage:    "whether enable or disable",
	Required: true,
}
