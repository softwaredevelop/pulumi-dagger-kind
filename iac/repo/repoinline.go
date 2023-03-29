//revive:disable:package-comments,exported
package main

import (
	"context"
	"log"
	"os"

	"github.com/pulumi/pulumi-github/sdk/v5/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/debug"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/workspace"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {

	destroy := false
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "destroy" {
			destroy = true
		}
	}

	ctx := context.Background()

	projectName := "thesis"
	stackNameA := "repoinlineA"
	stackNameB := "repoinlineB"
	desc := "A inline source Go Pulumi program for repository creation in thesis project."
	ws, err := auto.NewLocalWorkspace(ctx, auto.Project(workspace.Project{
		Name:        tokens.PackageName(projectName),
		Runtime:     workspace.NewProjectRuntimeInfo("go", nil),
		Description: &desc,
	}))
	if err != nil {
		panic(err)
	}

	prj, err := ws.ProjectSettings(ctx)
	if err != nil {
		panic(err)
	}

	stackA, err := auto.NewStackInlineSource(ctx, stackNameA, prj.Name.String(), func(ctx *pulumi.Context) error {
		return nil
	})
	if err != nil && auto.IsCreateStack409Error(err) {
		log.Println("stack " + stackNameA + " already exists")
		stackA, err = auto.UpsertStackInlineSource(ctx, stackNameA, prj.Name.String(), func(ctx *pulumi.Context) error {
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
	if err != nil && !auto.IsCreateStack409Error(err) {
		panic(err)
	}

	stackB, err := auto.NewStackInlineSource(ctx, stackNameB, prj.Name.String(), func(ctx *pulumi.Context) error {
		return nil
	})
	if err != nil && auto.IsCreateStack409Error(err) {
		log.Println("stack " + stackNameB + " already exists")
		stackB, err = auto.UpsertStackInlineSource(ctx, stackNameB, prj.Name.String(), func(ctx *pulumi.Context) error {
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
	if err != nil && !auto.IsCreateStack409Error(err) {
		panic(err)
	}

	pat := os.Getenv("PULUMI_ACCESS_TOKEN")
	err = stackA.Workspace().SetEnvVars(map[string]string{
		"PULUMI_SKIP_UPDATE_CHECK": "true",
		"PULUMI_ACCESS_TOKEN":      pat,
	})
	if err != nil {
		panic(err)
	}

	err = stackB.Workspace().SetEnvVars(map[string]string{
		"PULUMI_SKIP_UPDATE_CHECK": "true",
		"PULUMI_ACCESS_TOKEN":      pat,
	})
	if err != nil {
		panic(err)
	}

	ght := os.Getenv("GITHUB_TOKEN")
	gho := os.Getenv("GITHUB_OWNER")
	err = stackA.SetAllConfig(ctx, auto.ConfigMap{
		"github:token": auto.ConfigValue{
			Value:  ght,
			Secret: true,
		},
		"github:owner": auto.ConfigValue{
			Value:  gho,
			Secret: true,
		},
	})
	if err != nil {
		panic(err)
	}

	err = stackB.SetAllConfig(ctx, auto.ConfigMap{
		"github:token": auto.ConfigValue{
			Value:  ght,
			Secret: true,
		},
		"github:owner": auto.ConfigValue{
			Value:  gho,
			Secret: true,
		},
	})
	if err != nil {
		panic(err)
	}

	if destroy {

		drstB, err := stackB.Destroy(ctx, optdestroy.Message("Successfully destroyed stack : "+stackNameB))
		if err != nil {
			panic(err)
		}
		log.Println(drstB.Summary.Kind + " " + drstB.Summary.Message)

		drstA, err := stackA.Destroy(ctx, optdestroy.Message("Successfully destroyed stack : "+stackNameA))
		if err != nil {
			panic(err)
		}
		log.Println(drstA.Summary.Kind + " " + drstA.Summary.Message)

		return
	}

	stackA.Workspace().SetProgram(func(pCtx *pulumi.Context) error {

		repositoryName := "pulumi-dagger-kind"
		repository, err := github.NewRepository(pCtx, "newRepository", &github.RepositoryArgs{
			DeleteBranchOnMerge: pulumi.Bool(true),
			Description:         pulumi.String("This is a repository for Pulumi repository creation with Dagger CI/CD for kind"),
			HasIssues:           pulumi.Bool(true),
			HasProjects:         pulumi.Bool(true),
			Name:                pulumi.String(repositoryName),
			Topics: pulumi.StringArray{
				pulumi.String("pulumi"),
				pulumi.String("dagger"),
				pulumi.String("github"),
				pulumi.String("kind"),
				pulumi.String("kind-cluster"),
			},
			Visibility: pulumi.String("public"),
		})
		if err != nil {
			return err
		}

		_, err = github.NewBranchProtection(pCtx, "branchProtection", &github.BranchProtectionArgs{
			RepositoryId:          repository.NodeId,
			Pattern:               pulumi.String("main"),
			RequiredLinearHistory: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		_, err = github.NewIssueLabel(pCtx, "newIssueLabelGhActions", &github.IssueLabelArgs{
			Color:       pulumi.String("E66E01"),
			Description: pulumi.String("This issue is related to github-actions dependencies"),
			Name:        pulumi.String("github-actions dependencies"),
			Repository:  repository.Name,
		})
		if err != nil {
			return err
		}

		_, err = github.NewIssueLabel(pCtx, "newIssueLabelGoModules", &github.IssueLabelArgs{
			Color:       pulumi.String("9BE688"),
			Description: pulumi.String("This issue is related to Go modules"),
			Name:        pulumi.String("go-modules dependencies"),
			Repository:  repository.Name,
		})
		if err != nil {
			return err
		}

		pCtx.Export("repository", repository.Name)
		pCtx.Export("repositoryUrl", repository.HtmlUrl)

		return nil

	})

	prev, err := stackA.Preview(ctx, optpreview.Message("Preview stack "+stackNameA), optpreview.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	if err != nil {
		panic(err)
	}
	log.Println(prev.StdOut)

	// refr, err := stackA.Refresh(ctx, optrefresh.Message("Refresh stack "+stackNameA), optrefresh.ProgressStreams(os.Stdout, os.Stderr))
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(refr.StdOut)

	// up, err := stackA.Up(ctx, optup.Message("Update stack "+stackNameA), optup.DebugLogging(debug.LoggingOptions{
	// 	Debug: true,
	// }))
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(up.StdOut)

	stackB.Workspace().SetProgram(func(pCtx *pulumi.Context) error {

		_, err = github.GetActionsPublicKey(pCtx, &github.GetActionsPublicKeyArgs{
			Repository: "pulumi-dagger-kind",
		}, nil)
		if err != nil {
			return err
		}

		_, err = github.NewActionsSecret(pCtx, "newActionsSecretPAT", &github.ActionsSecretArgs{
			Repository: pulumi.String("pulumi-dagger-kind"),
			SecretName: pulumi.String("PULUMI_ACCESS_TOKEN"),
		})
		if err != nil {
			return err
		}

		_, err = github.NewActionsSecret(pCtx, "newActionsSecretPON", &github.ActionsSecretArgs{
			Repository: pulumi.String("pulumi-dagger-kind"),
			SecretName: pulumi.String("PULUMI_ORG_NAME"),
		})
		if err != nil {
			return err
		}

		return nil
	})

	prev, err = stackB.Preview(ctx, optpreview.Message("Preview stack "+stackNameB), optpreview.DebugLogging(debug.LoggingOptions{
		Debug: true,
	}))
	if err != nil {
		panic(err)
	}
	log.Println(prev.StdOut)

	// refr, err = stackB.Refresh(ctx, optrefresh.Message("Refresh stack "+stackNameB), optrefresh.ProgressStreams(os.Stdout, os.Stderr))
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(refr.StdOut)

	// up, err = stackB.Up(ctx, optup.Message("Update stack "+stackNameB), optup.DebugLogging(debug.LoggingOptions{
	// 	Debug: true,
	// }))
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(up.StdOut)

}
