package mock

//go:generate mockgen -package=mock -destination=mock_gen.go github.com/davidborzek/hetzner-ddns-updater/pkg/hetzner Client
