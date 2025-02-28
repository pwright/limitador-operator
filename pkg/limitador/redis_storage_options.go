package limitador

import (
	"context"
	"errors"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	limitadorv1alpha1 "github.com/kuadrant/limitador-operator/api/v1alpha1"
)

func RedisDeploymentOptions(ctx context.Context, cl client.Client, defSecretNamespace string, redisObj limitadorv1alpha1.Redis) (DeploymentStorageOptions, error) {
	if redisObj.ConfigSecretRef == nil {
		return DeploymentStorageOptions{}, errors.New("there's no ConfigSecretRef set")
	}

	redisURL, err := getURLFromRedisSecret(ctx, cl, defSecretNamespace, *redisObj.ConfigSecretRef)
	if err != nil {
		return DeploymentStorageOptions{}, err
	}

	return DeploymentStorageOptions{
		Command: []string{"redis", redisURL},
	}, nil
}

func getURLFromRedisSecret(ctx context.Context, cl client.Client, defSecretNamespace string, secretRef v1.ObjectReference) (string, error) {
	secret := &v1.Secret{}
	if err := cl.Get(
		ctx,
		types.NamespacedName{
			Name: secretRef.Name,
			Namespace: func() string {
				if secretRef.Namespace != "" {
					return secretRef.Namespace
				}
				return defSecretNamespace
			}(),
		},
		secret,
	); err != nil {
		// Must exist, so if it does not, also return err
		return "", err
	}

	// nil map behaves as empty map when reading
	if url, ok := secret.Data["URL"]; ok {
		return string(url), nil
	}

	return "", errors.New("the storage config Secret doesn't have the `URL` field")
}
