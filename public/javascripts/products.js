"use strict";

$(function() {
  $(".products__meta--size li").on("click", function(e) {
    let $target = $(e.target),
      $li = $(".products__meta--size li"),
      $size = $('[name="size_variation_id"]');

    $li.removeClass("current");
    $target.addClass("current");
    $(".products__meta--size li")
      .not(".current")
      .removeClass("selected");
    $target.toggleClass("selected");

    if ($target.hasClass("selected")) {
      $size.val($target.attr("value"));
    } else {
      $size.val(0);
    }
  });

  $(".products__meta--color li").on("click", function(e) {
    let $target = $(e.target),
      $li = $(".products__meta--color li"),
      $color = $('[name="color_variation_id"]');

    $li.removeClass("current");
    $target.addClass("current");
    $(".products__meta--color li")
      .not(".current")
      .removeClass("selected");
    $target.toggleClass("selected");

    if ($target.hasClass("selected")) {
      $color.val($target.attr("value"));
    } else {
      $color.val(0);
    }
  });

  // $(".products__meta--color li").on("click", function(t) {
  //   console.log("select color");
  //   var a = $(t.target),
  //     e = $(".products__meta--color li"),
  //     i = $('[name="color_variation_id"]');
  //   e.removeClass("current"),
  //     a.addClass("current"),
  //     $(".products__meta--color li")
  //       .not(".current")
  //       .removeClass("selected"),
  //     a.toggleClass("selected"),
  //     a.hasClass("selected") ? i.val(a.attr("value")) : i.val(0);
  // });

  // $(".products__meta--color span").on("click", function(t) {
  //   console.log("select color");
  //   var a = $(t.target),
  //     e = $(".products__meta--color span"),
  //     i = $('[name="color_variation_id"]');
  //   e.removeClass("current"),
  //     a.addClass("current"),
  //     $(".products__meta--color span")
  //       .not(".current")
  //       .removeClass("selected"),
  //     a.toggleClass("selected"),
  //     a.hasClass("selected") ? i.val(a.attr("value")) : i.val(0);
  // }),

  $("#products__addtocart").on("submit", function(event) {
    event.preventDefault();
    if ($('[name="size_color_id"]').val() == "0") {
      alert("please select color!");
      return;
    }
    if ($('[name="size_variation_id"]').val() == "0") {
      alert("please select size!");
      return;
    }
    $.ajax({
      type: "PUT",
      url: "/cart",
      dataType: "json",
      data: $(event.target).serialize(),
      error: function(xhr) {
        alert(xhr.status + ": " + xhr.statusText);
      },
      success: function(response) {
        window.location.replace("/cart");
      }
    });
  });

  $(".products__gallery--thumbs").length &&
    $(".products__gallery--thumbs").flexslider({
      animation: "slide",
      controlNav: false,
      animationLoop: false,
      slideshow: false,
      itemWidth: 80,
      itemMargin: 16,
      asNavFor: ".products__gallery--top"
    });

  $(".products__gallery--top").length &&
    $(".products__gallery--top").flexslider({
      animation: "slide",
      controlNav: false,
      directionNav: false,
      animationLoop: false,
      slideshow: false,
      sync: ".products__gallery--thumbs"
    });

  let productsFeaturedSliderH = $(".products__featured--slider").width(),
    isMobile = window.matchMedia("only screen and (max-width: 760px)").matches,
    columnNuber = isMobile ? 2 : 4;

  $(".products__featured--slider").length &&
    $(".products__featured--slider").flexslider({
      animation: "slide",
      animationLoop: false,
      controlNav: false,
      itemWidth: 200,
      itemMargin: 16
    });
});
