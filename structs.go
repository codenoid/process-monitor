package main

type config struct {
	InstanceName string `yaml:"instance_name`
	Notifier     struct {
		Telegram struct {
			Token   string  `yaml:"token`
			RoomIds []int64 `yaml:"rooms"`
		} `yaml:"telegram"`
	} `yaml:"notifier"`
	NotifConfig struct {
		RepeatEvery string `yaml:"repeat_every"`
	} `notif_config`
}
