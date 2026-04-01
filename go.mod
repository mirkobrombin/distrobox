module github.com/89luca89/distrobox

go 1.25.3

require (
	github.com/joho/godotenv v1.5.1
	github.com/mirkobrombin/go-foundation v0.4.0
	github.com/stretchr/testify v1.11.1
	github.com/urfave/cli/v3 v3.5.0
	gopkg.in/ini.v1 v1.67.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/mirkobrombin/go-cli-builder/v2 => /var/home/mirko/Projects/personal/go-tools/go-cli-builder
	github.com/mirkobrombin/go-conf-builder/v2 => /var/home/mirko/Projects/personal/go-tools/go-conf-builder
	github.com/mirkobrombin/go-foundation => /var/home/mirko/Projects/personal/go-tools/go-foundation
)
