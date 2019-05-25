"use strict";
function updateCart(t) {
  $.ajax({
    type: "PUT",
    data: t,
    dataType: "json",
    url: "/cart",
    beforeSend: function() {
      $(".cart__qty--btn").attr("disabled", !0);
    },
    success: function() {
      window.location.reload();
    },
    error: function(t) {
      window.alert(t.status + ": " + t.statusText), window.location.reload();
    }
  });
}
$(function() {}),
  $(function() {
    $(".cart__list--remove").on("click", function(t) {
      t.preventDefault(),
        updateCart({
          quantity: 0,
          size_variation_id: $(this)
            .closest("tr")
            .find(".cart__qty--num")
            .data("size-variation-id")
        });
    }),
      $(".cart__qty--minus").on("click", function(t) {
        var a = $(this)
            .attr("disabled", !0)
            .closest(".cart__qty")
            .find(".cart__qty--num"),
          e = a.data("size-variation-id"),
          i = parseInt(a.val()),
          r = i - 1;
        return (
          r < 1 && (r = 0),
          a.val(r),
          updateCart({ quantity: r, size_variation_id: e }),
          !1
        );
      }),
      $(".cart__qty--plus").on("click", function(t) {
        var a = $(this)
            .attr("disabled", !0)
            .closest(".cart__qty")
            .find(".cart__qty--num"),
          e = a.data("size-variation-id"),
          i = parseInt(a.val()),
          r = i + 1;
        return a.val(r), updateCart({ quantity: r, size_variation_id: e }), !1;
      }),
      $(".cart__qty--num").on("blur", function(t) {
        var a = $(this).attr("disabled", !0),
          e = a.data("size-variation-id"),
          i = a.val();
        return /^[0-9]*$/.test(i)
          ? (updateCart({ quantity: parseInt(i), size_variation_id: e }), !1)
          : void window.alert("please enter the right number!");
      });
  }),
  $(function() {
    $(".footer__change-language select").on("change", function() {
      var t = $(this).val();
      return t && (window.location = t), !1;
    });
  }),
  $(function() {
    $(".products__meta--size li").on("click", function(t) {
      console.log("select size");
      var a = $(t.target),
        e = $(".products__meta--size li"),
        i = $('[name="size_variation_id"]');
      e.removeClass("current"),
        a.addClass("current"),
        $(".products__meta--size li")
          .not(".current")
          .removeClass("selected"),
        a.toggleClass("selected"),
        a.hasClass("selected") ? i.val(a.attr("value")) : i.val(0);
    }),
      $(".products__meta--color span").on("click", function(t) {
        console.log("select color");
        var a = $(t.target),
          e = $(".products__meta--color span"),
          i = $('[name="color_variation_id"]');
        e.removeClass("current"),
          a.addClass("current"),
          $(".products__meta--color span")
            .not(".current")
            .removeClass("selected"),
          a.toggleClass("selected"),
          a.hasClass("selected") ? i.val(a.attr("value")) : i.val(0);
      }),
      $("#products__addtocart").on("submit", function(t) {
        if ((t.preventDefault(), "0" == $('[name="size_variation_id"]').val()))
          return void alert("please select size!");
        $.ajax({
          type: "PUT",
          url: "/cart",
          dataType: "json",
          data: $(t.target).serialize(),
          error: function(t) {
            alert(t.status + ": " + t.statusText);
          },
          success: function(t) {
            window.location.replace("/cart");
          }
        });
      }),
      $(".products__gallery--thumbs").length &&
        $(".products__gallery--thumbs").flexslider({
          animation: "slide",
          controlNav: !1,
          animationLoop: !1,
          slideshow: !1,
          itemWidth: 80,
          itemMargin: 16,
          asNavFor: ".products__gallery--top"
        }),
      $(".products__gallery--top").length &&
        $(".products__gallery--top").flexslider({
          animation: "slide",
          controlNav: !1,
          directionNav: !1,
          animationLoop: !1,
          slideshow: !1,
          sync: ".products__gallery--thumbs"
        });
    $(".products__featured--slider").width(),
      window.matchMedia("only screen and (max-width: 760px)").matches;
    $(".products__featured--slider").length &&
      $(".products__featured--slider").flexslider({
        animation: "slide",
        animationLoop: !1,
        controlNav: !1,
        itemWidth: 200,
        itemMargin: 16
      });
  });
