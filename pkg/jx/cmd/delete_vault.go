package cmd

import (
	"fmt"
	"io"

	"github.com/jenkins-x/jx/pkg/jx/cmd/templates"
	"github.com/jenkins-x/jx/pkg/kube"
	"github.com/jenkins-x/jx/pkg/log"
	"github.com/jenkins-x/jx/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

type DeleteVaultOptions struct {
	CommonOptions

	Namespace string
}

var (
	deleteVaultLong = templates.LongDesc(`
		Deletes a Vault
	`)

	deleteVaultExample = templates.Examples(`
		# Deletes a Vault from namespace my-namespace
		jx delete vault --namespace my-namespace my-vault
	`)
)

func NewCmdDeleteVault(f Factory, in terminal.FileReader, out terminal.FileWriter, errOut io.Writer) *cobra.Command {
	options := &DeleteVaultOptions{
		CommonOptions: CommonOptions{
			Factory: f,
			In:      in,
			Out:     out,
			Err:     errOut,
		},
	}

	cmd := &cobra.Command{
		Use:     "vault",
		Short:   "Deletes a Vault",
		Long:    deleteVaultLong,
		Example: deleteVaultExample,
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			CheckErr(err)
		},
	}

	cmd.Flags().StringVarP(&options.Namespace, "namespace", "n", "", "Namespace from where to delete the vault")
	return cmd
}

func (o *DeleteVaultOptions) Run() error {
	if len(o.Args) != 1 {
		return fmt.Errorf("Missing vault name")
	}
	vaultName := o.Args[0]

	client, ns, err := o.KubeClient()
	if err != nil {
		return errors.Wrap(err, "creating kubernetes client")
	}

	if o.Namespace == "" {
		o.Namespace = ns
	}

	vaultOperatorClient, err := o.VaultOperatorClient()
	if err != nil {
		return errors.Wrap(err, "creating vault operator client")
	}

	found := kube.FindVault(vaultOperatorClient, vaultName, o.Namespace)
	if !found {
		return fmt.Errorf("vault '%s' not found in namespace '%s'", vaultName, o.Namespace)
	}

	err = kube.DeleteVault(vaultOperatorClient, vaultName, o.Namespace)
	if err != nil {
		return errors.Wrap(err, "deleting the vault resource")
	}

	err = kube.DeleteIngress(client, o.Namespace, vaultName)
	if err != nil {
		return errors.Wrapf(err, "deleting the vault ingress '%s'", vaultName)
	}

	authServiceAccountName := kube.VaultAuthServiceAccountName(vaultName)
	err = kube.DeleteServiceAccount(client, o.Namespace, authServiceAccountName)
	if err != nil {
		return errors.Wrapf(err, "deleting the vault auth service account '%s'", authServiceAccountName)
	}

	log.Infof("Vault %s deleted\n", util.ColorInfo(vaultName))

	return nil
}
