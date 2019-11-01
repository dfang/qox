"use strict";

let gulp = require("gulp"),
  babel = require("gulp-babel"),
  eslint = require("gulp-eslint"),
  plumber = require("gulp-plumber"),
  cleanCSS = require("gulp-clean-css"),
  concat = require("gulp-concat"),
  sass = require("gulp-sass"),
  uglify = require("gulp-uglify"),
  autoprefixer = require("gulp-autoprefixer"),
  del = require('del'),
  fs = require("fs"),
  path = require("path"),
  es = require("event-stream"),
  rename = require("gulp-rename");


let pathto = function (file) {
  return "app/views/qor/assets/" + file;
};

let vendorPathTo = function (file) {
  return "vendor/github.com/qor/admin/views/assets/" + file;
};

let scripts = {
  src: vendorPathTo("javascripts/app/*.js"),
  dest: pathto("javascripts"),
  qor: vendorPathTo("javascripts/qor/*.js"),
  qorInit: vendorPathTo("javascripts/qor/qor-config.js"),
  qorCommon: vendorPathTo("javascripts/qor/qor-common.js"),
  qorAdmin: [pathto("javascripts/qor.js"), pathto("javascripts/app.js")],
  all: ["gulpfile.js", pathto("javascripts/qor/*.js")],
};

let styles = {
  src: vendorPathTo("stylesheets/scss/{app,qor}.scss"),
  dest: pathto("stylesheets"),
  vendors: vendorPathTo("stylesheets/vendors"),
  main: pathto("stylesheets/{qor,app}.css"),
  qorAdmin: [
    pathto("stylesheets/vendors.css"),
    pathto("stylesheets/qor.css"),
    pathto("stylesheets/app.css"),
  ],
  scss: pathto("stylesheets/scss/**/*.scss"),
};

// 把 javascripts/qor 下所有的js 打包为 qor.js
gulp.task("qor", function () {
  return gulp
    .src([scripts.qorInit, scripts.qorCommon, scripts.qor])
    // .pipe(plumber())
    .pipe(concat("qor.js"))
    // .pipe(uglify())
    .pipe(gulp.dest(scripts.dest));
});

// 把 javascripts/app 下所有的js 打包为 app.js
gulp.task(
  "app",
  gulp.series("qor", function () {
    return gulp
      .src(scripts.src)
      // .pipe(plumber())
      .pipe(
        eslint({
          configFile: ".eslintrc",
        })
      )
      .pipe(concat("app.js"))
      // .pipe(uglify())
      .pipe(gulp.dest(scripts.dest));
  })
);

// Task for compress js and css vendor assets
gulp.task("combineJavaScriptVendor", function () {
  return gulp
    .src([
      "vendor/github.com/qor/admin/views/assets/javascripts/vendors/jquery.min.js",
      "vendor/github.com/qor/admin/views/assets/javascripts/vendors/*.js",
      "app/views/qor/assets/javascripts/vendors/*.js",
    ])
    .pipe(concat("vendors.js"))
    .pipe(gulp.dest("app/views/qor/assets/javascripts"));
});

gulp.task("combineDatetimePicker", function () {
  return gulp
    .src([
      "app/views/qor/assets/javascripts/qor/qor-config.js",
      "app/views/qor/assets/javascripts/qor/qor-material.js",
      "app/views/qor/assets/javascripts/qor/qor-modal.js",
      "app/views/qor/assets/javascripts/qor/datepicker.js",
      "app/views/qor/assets/javascripts/qor/qor-datepicker.js",
      "app/views/qor/assets/javascripts/qor/qor-timepicker.js",
    ])
    .pipe(plumber())
    .pipe(
      eslint({
        configFile: ".eslintrc",
      })
    )
    .pipe(
      babel({
        presets: ["@babel/env"],
      })
    )

    .pipe(concat("datetimepicker.js"))
    .pipe(uglify())
    .pipe(gulp.dest("app/views/qor/assets/javascripts"));
});

gulp.task("compressCSSVendor", function () {
  return gulp
    .src([
      "vendor/github.com/qor/admin/views/assets/stylesheets/vendors/*.css",
      "app/views/qor/assets/stylesheets/vendors/*.css"
    ])
    .pipe(concat("vendors.css"))
    .pipe(gulp.dest("app/views/qor/assets/stylesheets"));
});

gulp.task("qor+", function () {
  return gulp
    .src([scripts.qorInit, scripts.qorCommon, scripts.qor])
    .pipe(plumber())
    .pipe(
      eslint({
        configFile: ".eslintrc",
      })
    )
    .pipe(
      babel({
        presets: ["@babel/env"],
      })
    )
    .pipe(eslint.format())
    .pipe(concat("qor.js"))
    .pipe(uglify())
    .pipe(gulp.dest(scripts.dest));
});

gulp.task("app+", function () {
  return gulp
    .src(scripts.src)
    .pipe(plumber())
    .pipe(
      babel({
        presets: ["@babel/env"],
      })
    )
    .pipe(eslint.format())
    .pipe(concat("app.js"))
    .pipe(uglify())
    .pipe(gulp.dest(scripts.dest));
});

// gulp.task("js", gulp.series(["app", "combineDatetimePicker", "combineJavaScriptVendor"]));

// 把 qor.js 和 app.js 合并为 qor_admin_default.js
gulp.task("release_js", gulp.series("combineJavaScriptVendor", "qor+", "app+", function () {
  return gulp
    .src(scripts.qorAdmin)
    .pipe(concat("qor_admin_default.js"))
    .pipe(gulp.dest(scripts.dest));
}, function () {
  return del([
    "app/views/qor/assets/javascripts/qor.js",
    "app/views/qor/assets/javascripts/app.js",
  ])
}));

// 生成qor.css 和 app.css
gulp.task("sass", function () {
  return gulp
    .src(styles.src)
    .pipe(plumber())
    .pipe(sass().on("error", sass.logError))
    .pipe(gulp.dest(styles.dest));
});

gulp.task(
  "css",
  gulp.series("sass", function () {
    return gulp
      .src(styles.main)
      .pipe(plumber())
      .pipe(autoprefixer())
      .pipe(cleanCSS())
      .pipe(gulp.dest(styles.dest));
  })
);

gulp.task("release_css", gulp.series("compressCSSVendor", "css", function () {
  return gulp
    .src(styles.qorAdmin)
    .pipe(concat("qor_admin_default.css"))
    .pipe(gulp.dest(styles.dest));
}, function () {
  return del([
    "app/views/qor/assets/stylesheets/qor.css",
    "app/views/qor/assets/stylesheets/app.css",
  ])
}));


gulp.task("default", gulp.series(["release_js", "release_css"]))

gulp.task('clean', function () {
  return del([
    "app/views/qor/assets/stylesheets/vendors.css",
    "app/views/qor/assets/stylesheets/qor.css",
    "app/views/qor/assets/stylesheets/app.css",
    "app/views/qor/assets/stylesheets/qor_admin_default.css",
    "app/views/qor/assets/javascripts/vendors.js",
    "app/views/qor/assets/javascripts/qor_admin_default.js",
    "app/views/qor/assets/javascripts/qor.js",
    "app/views/qor/assets/javascripts/app.js",
    // we don't want to clean this file though so we negate the pattern
    // '!dist/mobile/deploy.json'
  ]);
});
