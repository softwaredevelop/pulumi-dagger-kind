//revive:disable:package-comments,exported
package main

import (
	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		hello, err := local.NewCommand(ctx, "hello", &local.CommandArgs{
			Create: pulumi.String("echo \"Hello Pulumi\""),
		})
		if err != nil {
			return err
		}

		ctx.Export("hello", hello.Stdout)

		return nil
	})
}
