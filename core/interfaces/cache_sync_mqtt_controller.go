package interfaces

type ICacheSyncMQTTController interface {
	Process(forceSync bool)
}
