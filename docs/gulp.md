# Gulp qor如何打包后台的css 和js 

现有的gulpfile.js 根据 
// https://github.com/qor/qor/blob/master/gulpfile.js
修改来的， 
从qor/qor 下copy了此文件和.eslintignore, .eslintrc, package.json, yarn.lock

由于当前pacakge.json 中的Node-sass 不支持当前的Mac环境， 升级了该插件，其余的没动.

```
yarn remove node-sass
yarn add -D node-sass
```

/* layout.tmpl 里

stylesheet_tag 和 javascript_tag 是在qor/admin 中的funcmap

关于js引用的是  
vendors.js  
qor_admin_default.js = qor.js + app.js  
qor_demo.js  
前两个是admin qor_demo.js 是自定义逻辑  

关于css引用的是  
qor_admin_default.css  
qor_demo.css  
和其他插件的， 比如notification.css  
*/


```
gulp realease_js 
gulp realease_css
gulp 
``
