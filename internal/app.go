package internal

import "Cloud/logger"

// App представляет собой структуру приложения, содержащую необходимые зависимости
type App struct {
	RequestLogger *logger.RequestLogger // Логгер запросов
}

//internal представляет собой компонент вашего приложения и организует его зависимости.
//internal может содержать логику приложения и общие компоненты.
