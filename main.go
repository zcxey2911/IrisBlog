package main

import (
	"IrisBlog/handler"
	"IrisBlog/mytool"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
)

func main() {

	app := newApp()

	app.HandleDir("/assets", iris.Dir("./assets"))
	app.Favicon("./favicon.ico")
	app.Listen(":5000")
}

func newApp() *iris.Application {

	app := iris.New()

	tmpl := iris.HTML("./views", ".html")

	tmpl.Delims("${", "}")

	tmpl.Reload(true)

	app.RegisterView(tmpl)

	app.Use(iris.Compression)

	app.Post("/captcha/", mytool.GetCaptchaId)
	app.Get("/captcha/*/", mytool.GetCaptchaImg)

	// app.Get("/captcha/", func(ctx iris.Context) {

	// 	id := captcha.New()

	// 	//res := captcha.VerifyString(id, "123456")

	// 	//fmt.Println(res)

	// 	captcha.WriteImage(ctx, id, captcha.StdWidth, captcha.StdHeight)
	// })

	app.Get("/signin/", handler.User_signin)
	app.Post("/signin/", handler.Signin)

	// app.Get("/", func(ctx iris.Context) {

	// 	ctx.ViewData("message", "你好，女神")

	// 	ctx.View("index.html")
	// })

	verifier := jwt.NewVerifier(jwt.HS256, mytool.SigKey)
	//verifier.WithDefaultBlocklist()

	verifyMiddleware := verifier.Verify(func() interface{} {

		return new(mytool.PlayLoad)
	})

	adminhandler := app.Party("/admin")
	// {
	// 	adminhandler.Use(iris.Compression)
	// 	adminhandler.Get("/userlist/", handler.Admin_userlist)
	// 	adminhandler.Delete("/user_action/", handler.Admin_userdel)
	// 	adminhandler.Put("/user_action/", handler.Admin_userupdate)
	// 	adminhandler.Post("/user_action/", handler.Admin_useradd)
	//adminhandler.Get("/user/", handler.Admin_user_page)

	// }

	adminhandler.Use(verifyMiddleware)

	adminhandler.Get("/userlist/", handler.Admin_userlist)
	adminhandler.Delete("/user_action/", handler.Admin_userdel)
	adminhandler.Put("/user_action/", handler.Admin_userupdate)
	adminhandler.Post("/user_action/", handler.Admin_useradd)
	adminhandler.Get("/user/", handler.Admin_user_page)

	adminhandler.Get("/logout", mytool.Logout)

	app.Get("/", mytool.TestToken())

	return app

}
