package botapp

type BotApp struct{}

func New() *BotApp {
	return &BotApp{}
}

func (s *BotApp) Start() error {
	return nil
}

func (s *BotApp) Stop() error {
	return nil
}
