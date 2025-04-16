package handlers_test

const ValidLoginJSON = `
{
	"username": "admin",
	"password": "secret123"
}`

const InvalidLoginJSON = `
{
	"username": "admin",
	"password": "secret123",`

const MissingPasswordJSON = `
{
	"username": "admin"
}`
