package main

import (
	"fmt"
	"os"

	"github.com/jwilder/k3a/pkg/spinner"
	"github.com/jwilder/k3a/pool"
	"github.com/spf13/cobra"
)

var poolCmd = &cobra.Command{
	Use:   "pool",
	Short: "Manage VMSS pools (list, create, delete, scale)",
}

var listPoolsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Virtual Machine Scale Sets (VMSS) in the specified resource group.",
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID, _ := cmd.Root().Flags().GetString("subscription")
		if subscriptionID == "" {
			return fmt.Errorf("--subscription flag is required (or set K3A_SUBSCRIPTION)")
		}
		cluster, _ := cmd.Flags().GetString("cluster")
		if cluster == "" {
			return fmt.Errorf("--cluster flag is required (or set K3A_CLUSTER)")
		}

		return pool.List(pool.ListPoolArgs{
			SubscriptionID: subscriptionID,
			Cluster:        cluster,
		})
	},
}

var createPoolCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new VMSS pool.",
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID, _ := cmd.Root().Flags().GetString("subscription")
		if subscriptionID == "" {
			return fmt.Errorf("--subscription flag is required (or set K3A_SUBSCRIPTION)")
		}
		cluster, _ := cmd.Flags().GetString("cluster")
		if cluster == "" {
			return fmt.Errorf("--cluster flag is required (or set K3A_CLUSTER)")
		}
		region, _ := cmd.Flags().GetString("region")
		role, _ := cmd.Flags().GetString("role")
		name, _ := cmd.Flags().GetString("name")
		sshKeyPath, _ := cmd.Flags().GetString("ssh-key")
		instanceCount, _ := cmd.Flags().GetInt("instance-count")
		k8sVersion, _ := cmd.Flags().GetString("k8s-version")
		sku, _ := cmd.Flags().GetString("sku")
		osDiskSize, _ := cmd.Flags().GetInt("os-disk-size")
		etcdAddr, _ := cmd.Flags().GetString("etcd-addr")
		hyperVGeneration, _ := cmd.Flags().GetString("hyperv-generation")

		// Accept one or more MSI resource IDs
		msiIDs, _ := cmd.Flags().GetStringArray("msi")

		// Add spinner for pool creation
		stopSpinner := spinner.Spinner("Creating VMSS pool...")
		defer stopSpinner()

		return pool.Create(pool.CreatePoolArgs{
			SubscriptionID:   subscriptionID,
			Cluster:          cluster,
			Region:           region,
			EtcdAddr:         etcdAddr,
			Role:             role,
			Name:             name,
			SSHKeyPath:       sshKeyPath,
			InstanceCount:    instanceCount,
			K8sVersion:       k8sVersion,
			SKU:              sku,
			OSDiskSizeGB:     osDiskSize,
			MSIIDs:           msiIDs,
			HyperVGeneration: hyperVGeneration,
		})
	},
}

var deletePoolCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a VMSS pool.",
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID, _ := cmd.Root().Flags().GetString("subscription")
		if subscriptionID == "" {
			return fmt.Errorf("--subscription flag is required (or set K3A_SUBSCRIPTION)")
		}
		cluster, _ := cmd.Flags().GetString("cluster")
		if cluster == "" {
			return fmt.Errorf("--cluster flag is required (or set K3A_CLUSTER)")
		}

		name, _ := cmd.Flags().GetString("name")

		// Add spinner for pool deletion
		stopSpinner := spinner.Spinner("Deleting VMSS pool...")
		defer stopSpinner()

		return pool.Delete(pool.DeletePoolArgs{
			SubscriptionID: subscriptionID,
			Cluster:        cluster,
			Name:           name,
		})
	},
}

var scalePoolCmd = &cobra.Command{
	Use:   "scale",
	Short: "Scale a VMSS pool to the desired number of instances.",
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID, _ := cmd.Root().Flags().GetString("subscription")
		if subscriptionID == "" {
			return fmt.Errorf("--subscription flag is required (or set K3A_SUBSCRIPTION)")
		}
		cluster, _ := cmd.Flags().GetString("cluster")
		if cluster == "" {
			return fmt.Errorf("--cluster flag is required (or set K3A_CLUSTER)")
		}
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("--name flag is required")
		}
		instanceCount, _ := cmd.Flags().GetInt("instance-count")
		if instanceCount < 1 {
			return fmt.Errorf("--instance-count must be greater than 0")
		}

		// Add spinner for pool scaling
		stopSpinner := spinner.Spinner("Scaling VMSS pool...")
		defer stopSpinner()

		return pool.Scale(pool.ScalePoolArgs{
			SubscriptionID: subscriptionID,
			Cluster:        cluster,
			Name:           name,
			InstanceCount:  instanceCount,
		})
	},
}

var kubeadmInstallCmd = &cobra.Command{
	Use:   "kubeadm-install",
	Short: "Install kubeadm on an existing VMSS pool.",
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID, _ := cmd.Root().Flags().GetString("subscription")
		if subscriptionID == "" {
			return fmt.Errorf("--subscription flag is required (or set K3A_SUBSCRIPTION)")
		}
		cluster, _ := cmd.Flags().GetString("cluster")
		if cluster == "" {
			return fmt.Errorf("--cluster flag is required (or set K3A_CLUSTER)")
		}
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			return fmt.Errorf("--name flag is required")
		}
		role, _ := cmd.Flags().GetString("role")
		if role == "" {
			return fmt.Errorf("--role flag is required")
		}
		k8sVersion, _ := cmd.Flags().GetString("k8s-version")
		sshKeyPath, _ := cmd.Flags().GetString("ssh-key")

		region, _ := cmd.Flags().GetString("region")
		etcdAddr, _ := cmd.Flags().GetString("etcd-addr")

		// Add spinner for kubeadm installation
		stopSpinner := spinner.Spinner("Installing kubeadm on VMSS pool...")
		defer stopSpinner()

		return pool.KubeadmInstall(pool.KubeadmInstallArgs{
			SubscriptionID: subscriptionID,
			Cluster:        cluster,
			Name:           name,
			Role:           role,
			K8sVersion:     k8sVersion,
			SSHKeyPath:     sshKeyPath,
			Region:         region,
			EtcdAddr:       etcdAddr,
		})
	},
}

func init() {
	clusterDefault := ""
	if v := os.Getenv("K3A_CLUSTER"); v != "" {
		clusterDefault = v
	}
	// Pool list flags
	listPoolsCmd.Flags().String("cluster", clusterDefault, "Cluster name (or set K3A_CLUSTER) (required)")

	// Pool create flags
	createPoolCmd.Flags().String("cluster", clusterDefault, "Cluster name (or set K3A_CLUSTER) (required)")
	createPoolCmd.Flags().String("name", "", "Name of the node pool (required)")
	createPoolCmd.Flags().String("role", "control-plane", "Role of the node pool (control-plane or worker)")
	createPoolCmd.Flags().String("region", "canadacentral", "Azure region for the pool")
	createPoolCmd.Flags().String("etcd-addr", "", "Etcd address for the first control-plane node (required for first control-plane pool)")
	createPoolCmd.Flags().Int("instance-count", 1, "Number of VMSS instances")
	createPoolCmd.Flags().String("ssh-key", os.ExpandEnv("$HOME/.ssh/id_rsa.pub"), "Path to the SSH public key file")
	createPoolCmd.Flags().String("k8s-version", "v1.33.1", "Kubernetes version (e.g. v1.33.1)")
	createPoolCmd.Flags().String("sku", "Standard_D2s_v3", "VM SKU type (default: Standard_D2s_v3)")
	createPoolCmd.Flags().String("hyperv-generation", "gen2", "Hypervisor generation: gen1 or gen2 (default: gen2)")
	createPoolCmd.Flags().Int("os-disk-size", 30, "OS disk size in GB (default: 30)")
	createPoolCmd.Flags().StringArray("msi", nil, "Additional user-assigned MSI resource IDs to add to the VMSS (can be specified multiple times)")

	_ = createPoolCmd.MarkFlagRequired("name")
	_ = createPoolCmd.MarkFlagRequired("role")

	// Pool delete flags
	deletePoolCmd.Flags().String("cluster", clusterDefault, "Cluster name (or set K3A_CLUSTER) (required)")
	deletePoolCmd.Flags().String("name", "", "Name of the node pool (required)")
	_ = deletePoolCmd.MarkFlagRequired("name")

	// Pool scale flags
	scalePoolCmd.Flags().String("cluster", clusterDefault, "Cluster name (or set K3A_CLUSTER) (required)")
	scalePoolCmd.Flags().String("name", "", "Name of the node pool (required)")
	scalePoolCmd.Flags().Int("instance-count", 1, "Number of VMSS instances (required)")
	_ = scalePoolCmd.MarkFlagRequired("name")
	_ = scalePoolCmd.MarkFlagRequired("instance-count")

	// Pool kubeadm install flags
	kubeadmInstallCmd.Flags().String("cluster", clusterDefault, "Cluster name (or set K3A_CLUSTER) (required)")
	kubeadmInstallCmd.Flags().String("name", "", "Name of the node pool (required)")
	kubeadmInstallCmd.Flags().String("role", "", "Role of the node pool (control-plane or worker) (required)")
	kubeadmInstallCmd.Flags().String("k8s-version", "v1.33.1", "Kubernetes version (e.g. v1.33.1)")
	kubeadmInstallCmd.Flags().String("ssh-key", "", "SSH public key path (default: ~/.ssh/id_rsa.pub)")
	kubeadmInstallCmd.Flags().String("etcd-addr", "", "Etcd address for the first control-plane node (required for first control-plane pool)")
	kubeadmInstallCmd.Flags().String("region", "canadacentral", "Azure region for the pool")
	_ = kubeadmInstallCmd.MarkFlagRequired("name")
	_ = kubeadmInstallCmd.MarkFlagRequired("role")

	poolCmd.AddCommand(instancesPoolCmd)
	poolCmd.AddCommand(listPoolsCmd, createPoolCmd, deletePoolCmd, scalePoolCmd, kubeadmInstallCmd)

	rootCmd.AddCommand(poolCmd)
}
