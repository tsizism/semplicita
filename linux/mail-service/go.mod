module mail

go 1.23.3

replace github.com/tsizism/semplicita/linux/shared => ./../shared

require (
	github.com/rs/cors v1.11.1
	github.com/tsizism/semplicita/linux/shared v0.0.0-00010101000000-000000000000
	github.com/vanng822/go-premailer v1.22.0
	github.com/xhit/go-simple-mail/v2 v2.16.0
	go.mongodb.org/mongo-driver v1.17.1
	golang.org/x/text v0.18.0
)

require (
	github.com/PuerkitoBio/goquery v1.9.2 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/go-test/deep v1.1.1 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/toorop/go-dkim v0.0.0-20201103131630-e1cd1a0a5208 // indirect
	github.com/vanng822/css v1.0.1 // indirect
	golang.org/x/net v0.29.0 // indirect
)
