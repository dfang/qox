# REAMDME

## 1. 关于qor 对 go mod 的支持

由于qor是在go mod feature 出来之前就已经出了，后面因为没有继续开发（开发进度几乎停止），现在官方也没打算花大力气维护 。。。。。。

go 是从1.11开始引入go module，但是qor 对 go mod 的支持不太好，所以需要一些改动。

官方的github.com/qor/qor-example 和 github.com/qor/admin 上的demo 不做改动，是没法跑起来的

qor/admin 在运行 go run main.go 之后， 打开 http://localhost:9000/admin 就会报 runtime error: invalid memory address or nil pointer dereference

问题的原因在于 https://github.com/qor/admin/blob/master/admin.go#L102

最简单粗暴的做法是将 https://github.com/qor/admin/tree/master/views 里的所有复制到 main.go 同级的 app/views/qor/ 目录下

```
mkdir -p app/views/qor/
cp -R xxx app/views/qor/
go run main.go
````

或者  [modvendor](https://github.com/goware/modvendor)

```
go mod vendor
modvendor -copy="**/*.tmpl **/*.css **/*.js **/*woff **/*woff2 **/*ttf **/*jpg **/*png **/*ico" -v
```

注意:

用`go mod edit --replace` 跑起来并不能成功打开/admin页面, 需要好好的理解下go mod edit --replace的运行机制

go.mod

```
// replace github.com/qor/admin v0.0.0-20191021122103-f3db8244d2d2 => gitlab.com/qor2/admin v0.0.0-20190525133731-6329506ec305
```

删掉go.mod go.sum, 然后将main.go imports 里的 "github.com/qor/admin" 改为 "gitlab.com/qor2/admin"， 重跑 go mod init 和 go mod tidy， 能够成功打开localhost:9000/admin

但是这种方法对qor/qor-example 来说不是很好，得fork所有github.com/qor下的项目，代码里所有 qor/xxx 得改为 qor2/xxx

事实上这么干过，所有fork的子项目都在这里 https://gitlab.com/qor2


另外一个小测试:

go build之后，删掉 整个vendor 文件夹， 运行可执行文件 ./admin-demo会报错。 如果保留 admin/views 文件夹，其他的全部删掉是可以的。


所以对于qor-example 可以这么弄，

运行

```
go mod vendor

modvendor -copy="**/*.tmpl **/*.scss **/*.css **/*.js **/*.woff **/*.woff2 **/*.ttf **/*.jpg **/*.png **/*.ico **/*.yml **/*.yaml" -v

```

之后进入 vendor文件夹删掉除了views 之外的所有文件夹， 打包到线上的时候需要 一个可执行文件加上 vendor文件夹，这样Docker镜像会很小

用命令来操作就是

```
find ./vendor  -name "views" | xargs -I {} rsync -aR --delete {} ./tmp
rsync -avP --delete ./tmp/vendor .

cd vendor
fd test | xargs rm -rf
```

可以用此方法来进行测试：

将执行文件和vendor文件夹全部copy到一个临时文件夹如 /tmp/test/, 然后运行 ./qor-demo
当然config/locales 也需要一起copy过去，否则没有i18n

```
cp -r --parents config/locales /tmp/test/
```

最终的文件夹结构如下：

```
/tmp/test λ  tree -L 3
.
├── config
│   └── locales
│       ├── admin.zh-CN.yml
│       ├── auth.en-US.yml
│       ├── en-US.yml
│       └── zh-CN.yml
├── qor-demo
└── vendor
    └── github.com
        └── qor
````

整个可运行的文件夹也很小， 打包的镜像自然也会小很多，Dockefile 可以改进了
```
/tmp/test λ  du -d3 -h
24K	./config/locales
24K	./config
7.3M	./vendor/github.com/qor
7.3M	./vendor/github.com
7.3M	./vendor
37M	.
```


## 如果用vendor 文件夹的形式，如果进行自定以呢 覆盖vendor 里的css js呢

可以按如下方法测试

```
cd /tmp/test/
mkdir -p app/views/qor/assets/stylesheets/
cp app/views/qor/assets/stylesheets/qor_demo.css app/views/qor/assets/stylesheets/
```


public 文件夹的内容也可以复制过来


最终的文件夹结构如下：

```
/tmp/test λ  tree -L 1
.
├── app
├── config
├── public
├── qor-demo
└── vendor

4 directories, 1 file
```


一键生成可以测试的独立程序包:

```
go build
mkdir ~/tmp/test
cp .env ~/tmp/test
cp qor-demo ~/tmp/test
cp -r --parents config/locales ~/tmp/test/
cp -r --parents public ~/tmp/test/
cp -r --parents vendor ~/tmp/test/
cp -r --parents app/views ~/tmp/test/

cd ~/tmp/test/
./qor-demo
```
