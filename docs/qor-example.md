# 如何快速把官方的qor-example 跑起来



1. 首先移除掉 enterprise相关的东西，否则go mod tidy 都无法正常跑完

```
rm -rf app/enterprise
rm -rf config/db/migrations/enterprise_migration.go
rm -rf config/db/seeds/data/enterprises.yml
rm -rf config/db/seeds/enterprise*
```

2.
go mod tidy

删掉 main.go enterprise相关的
删掉 config/bindatafs/bindatafs.go:118 （config.NoMetadata那一行)

```
go mod tidy
modvendor -copy="**/*.tmpl **/*.css **/*.js **/*.woff **/*.woff2 **/*.ttf **/*.jpg **/*.png **/*.ico **/*.yml **/*.yaml" -v

```

vendor里只保留views文件夹也是可以的
```
find ./vendor  -name "views" | xargs -I {} rsync -aR --delete {} ./tmp
rsync -avP --delete ./tmp/vendor .
```

3. go run main.go
