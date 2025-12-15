package server

import "github.com/gofiber/fiber/v2"

func FiberApplication() *fiber.App {
	return fiber.New(fiber.Config{
		AppName:            "Backend API Hub",
		Prefork:            false,
		EnablePrintRoutes:  false,
		StrictRouting:      false,
		ReduceMemoryUsage:  true,
		EnableIPValidation: true,
	})
}
