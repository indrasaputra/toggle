package feature_test

import "github.com/cucumber/godog"

func theToggleIsEmpty() error {
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^the toggle is empty$`, theToggleIsEmpty)
}
