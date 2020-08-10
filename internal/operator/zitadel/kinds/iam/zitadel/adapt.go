package zitadel

import (
	"github.com/caos/orbos/internal/operator/orbiter/kinds/clusters/kubernetes"
	"github.com/caos/orbos/internal/operator/orbiter/kinds/clusters/kubernetes/resources/namespace"
	"github.com/caos/orbos/internal/operator/zitadel"
	"github.com/caos/orbos/internal/operator/zitadel/kinds/databases"
	"github.com/caos/orbos/internal/operator/zitadel/kinds/iam/zitadel/ambassador"
	"github.com/caos/orbos/internal/operator/zitadel/kinds/iam/zitadel/configuration"
	"github.com/caos/orbos/internal/operator/zitadel/kinds/iam/zitadel/deployment"
	"github.com/caos/orbos/internal/operator/zitadel/kinds/iam/zitadel/imagepullsecret"
	"github.com/caos/orbos/internal/operator/zitadel/kinds/iam/zitadel/migration"
	"github.com/caos/orbos/internal/operator/zitadel/kinds/iam/zitadel/services"
	"github.com/caos/orbos/internal/operator/zitadel/kinds/networking"
	"github.com/caos/orbos/internal/tree"
	"github.com/caos/orbos/mntr"
	"github.com/pkg/errors"
	"sort"
	"strconv"
	"strings"
)

func AdaptFunc(timestamp string, features []string) zitadel.AdaptFunc {
	return func(
		monitor mntr.Monitor,
		desired *tree.Tree,
		current *tree.Tree,
	) (
		zitadel.QueryFunc,
		zitadel.DestroyFunc,
		error,
	) {
		internalMonitor := monitor.WithField("kind", "iam")

		desiredKind, err := parseDesiredV0(desired)
		if err != nil {
			return nil, nil, errors.Wrap(err, "parsing desired state failed")
		}
		desired.Parsed = desiredKind

		if !monitor.IsVerbose() && desiredKind.Spec.Verbose {
			internalMonitor.Verbose()
		}

		namespaceStr := "caos-zitadel"
		labels := map[string]string{
			"app.kubernetes.io/managed-by": "zitadel.caos.ch",
			"app.kubernetes.io/part-of":    "zitadel",
		}
		internalLabels := map[string]string{}
		for k, v := range labels {
			internalLabels[k] = v
		}
		internalLabels["app.kubernetes.io/component"] = "iam"

		// shared elements
		cmName := "zitadel-vars"
		secretName := "zitadel-secret"
		consoleCMName := "console-config"
		secretVarsName := "zitadel-secrets-vars"
		secretPasswordName := "zitadel-passwords"
		imagePullSecretName := "public-github-packages"
		//paths which are used in the configuration and also are used for mounting the used files
		certPath := "/home/zitadel/dbsecrets-zitadel"
		secretPath := "/secret"
		//services which are kubernetes resources and are used in the ambassador elements
		grpcServiceName := "grpc-v1"
		grpcPort := 80
		httpServiceName := "http-v1"
		httpPort := 80
		uiServiceName := "ui-v1"
		uiPort := 80

		users, migrationUser := getUsers(desiredKind)

		allZitadelUsers := make([]string, 0)
		for k := range users {
			if k != migrationUser {
				allZitadelUsers = append(allZitadelUsers, k)
			}
		}
		sort.Slice(allZitadelUsers, func(i, j int) bool {
			return allZitadelUsers[i] < allZitadelUsers[j]
		})

		allUsers := make([]string, 0)
		for k := range users {
			allUsers = append(allUsers, k)
		}
		sort.Slice(allUsers, func(i, j int) bool {
			return allUsers[i] < allUsers[j]
		})

		databaseCurrent := &tree.Tree{}
		queryDB, destroyDB, err := databases.GetQueryAndDestroyFuncs(
			internalMonitor,
			desiredKind.Database,
			databaseCurrent,
			namespaceStr,
			allUsers,
			labels,
			timestamp,
			secretPasswordName,
			migrationUser,
			features,
		)
		if err != nil {
			return nil, nil, err
		}

		queryNS, err := namespace.AdaptFuncToEnsure(namespaceStr)
		if err != nil {
			return nil, nil, err
		}
		destroyNS, err := namespace.AdaptFuncToDestroy(namespaceStr)
		if err != nil {
			return nil, nil, err
		}

		queryC, destroyC, configurationDone, err := configuration.AdaptFunc(
			internalMonitor,
			namespaceStr,
			labels,
			desiredKind.Spec.Configuration,
			cmName,
			certPath,
			secretName,
			secretPath,
			consoleCMName,
			secretVarsName,
			secretPasswordName,
			users,
		)
		if err != nil {
			return nil, nil, err
		}

		queryS, destroyS, err := services.AdaptFunc(internalMonitor, namespaceStr, internalLabels, grpcServiceName, grpcPort, httpServiceName, httpPort, uiServiceName, uiPort)
		if err != nil {
			return nil, nil, err
		}

		queryIPS, destroyIPS, err := imagepullsecret.AdaptFunc(internalMonitor, namespaceStr, imagePullSecretName, internalLabels)
		if err != nil {
			return nil, nil, err
		}

		queryM, destroyM, migrationDone, _, err := migration.AdaptFunc(internalMonitor, namespaceStr, "init", internalLabels, secretPasswordName, migrationUser, allZitadelUsers)
		if err != nil {
			return nil, nil, err
		}

		queryD, destroyD, deploymentDone, scaleDeployment, err := deployment.AdaptFunc(
			internalMonitor,
			namespaceStr,
			internalLabels,
			desiredKind.Spec.ReplicaCount,
			desiredKind.Spec.Version,
			imagePullSecretName,
			cmName,
			certPath,
			secretName,
			secretPath,
			consoleCMName,
			secretVarsName,
			secretPasswordName,
			allZitadelUsers,
			desiredKind.Spec.NodeSelector,
			migrationDone,
			configurationDone,
		)
		if err != nil {
			return nil, nil, err
		}

		networkingCurrent := &tree.Tree{}
		queryNW, destroyNW, err := networking.GetQueryAndDestroyFuncs(internalMonitor, desiredKind.Networking, networkingCurrent, namespaceStr, labels)
		if err != nil {
			return nil, nil, err
		}

		queryAmbassador, destroyAmbassador, err := ambassador.AdaptFunc(
			internalMonitor,
			namespaceStr,
			labels,
			grpcServiceName+"."+namespaceStr+":"+strconv.Itoa(grpcPort),
			"http://"+httpServiceName+"."+namespaceStr+":"+strconv.Itoa(httpPort),
			"http://"+uiServiceName+"."+namespaceStr,
		)
		if err != nil {
			return nil, nil, err
		}

		progLabels := map[string]string{
			"zitadel.caos.ch/restore-in-progress": "true",
			"zitadel.caos.ch/last-restore":        strings.ReplaceAll(timestamp, ":", "-"),
		}
		finishedLabels := map[string]string{
			"zitadel.caos.ch/restore-in-progress": "false",
		}

		queriers := make([]zitadel.QueryFunc, 0)
		for _, feature := range features {
			switch feature {
			case "networking":
				//networking
				queriers = append(queriers,
					zitadel.ResourceQueryToZitadelQuery(queryNS),
					queryNW,
					queryAmbassador,
				)
			case "zitadel":
				queriers = append(queriers, //namespace
					zitadel.ResourceQueryToZitadelQuery(queryNS),
					//database
					queryDB,
					//configuration
					queryC,
					//migration
					queryM,
					//services
					queryS,
					zitadel.ResourceQueryToZitadelQuery(queryIPS),
					queryD,
					zitadel.EnsureFuncToQueryFunc(deploymentDone),
				)
			case "restore":
				queriers = append(queriers,
					zitadel.EnsureFuncToQueryFunc(scaleDeployment(0, progLabels)),
					queryDB,
					zitadel.EnsureFuncToQueryFunc(scaleDeployment(desiredKind.Spec.ReplicaCount, finishedLabels)),
				)
			case "instantbackup":
				queriers = append(queriers,
					queryDB,
				)
			}
		}

		destroyers := make([]zitadel.DestroyFunc, 0)
		for _, feature := range features {
			switch feature {
			case "networking":
				destroyers = append(destroyers,
					destroyNW,
					destroyAmbassador,
					zitadel.ResourceDestroyToZitadelDestroy(destroyNS),
				)
			case "zitadel":
				destroyers = append(destroyers, //namespace
					destroyS,
					destroyM,
					zitadel.ResourceDestroyToZitadelDestroy(destroyIPS),
					destroyD,
					destroyDB,
					destroyC,
					zitadel.ResourceDestroyToZitadelDestroy(destroyNS),
				)
			case "restore":
				destroyers = append(destroyers,
					destroyDB,
				)
			}
		}

		return func(k8sClient *kubernetes.Client, queried map[string]interface{}) (zitadel.EnsureFunc, error) {
				return zitadel.QueriersToEnsureFunc(internalMonitor, true, queriers, k8sClient, queried)
			},
			zitadel.DestroyersToDestroyFunc(monitor, destroyers),
			nil
	}
}

func getUsers(desired *DesiredV0) (map[string]string, string) {
	passwords := &configuration.Passwords{}
	if desired.Spec != nil && desired.Spec.Configuration != nil && desired.Spec.Configuration.Passwords != nil {
		passwords = desired.Spec.Configuration.Passwords
	}
	users := make(map[string]string, 0)

	migrationUser := "flyway"
	migrationPassword := migrationUser
	if passwords.Migration != nil {
		migrationPassword = passwords.Migration.Value
	}
	users[migrationUser] = migrationPassword

	mgmtUser := "management"
	mgmtPassword := mgmtUser
	if passwords != nil && passwords.Management != nil {
		mgmtPassword = passwords.Management.Value
	}
	users[mgmtUser] = mgmtPassword

	adminUser := "adminapi"
	adminPassword := adminUser
	if passwords != nil && passwords.Adminapi != nil {
		adminPassword = passwords.Adminapi.Value
	}
	users[adminUser] = adminPassword

	authUser := "auth"
	authPassword := authUser
	if passwords != nil && passwords.Auth != nil {
		authPassword = passwords.Auth.Value
	}
	users[authUser] = authPassword

	authzUser := "authz"
	authzPassword := authzUser
	if passwords != nil && passwords.Authz != nil {
		authzPassword = passwords.Authz.Value
	}
	users[authzUser] = authzPassword

	notUser := "notification"
	notPassword := notUser
	if passwords != nil && passwords.Notification != nil {
		notPassword = passwords.Notification.Value
	}
	users[notUser] = notPassword

	esUser := "eventstore"
	esPassword := esUser
	if passwords != nil && passwords.Eventstore != nil {
		esPassword = passwords.Notification.Value
	}
	users[esUser] = esPassword

	return users, migrationUser
}
