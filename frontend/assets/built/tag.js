$(function(){let r=null,l;function d(t,e,o,n){t.stop().fadeIn(300,"linear"),e.stop().fadeIn(300,"linear"),o.stop().fadeOut(300,"linear",function(){o.children().empty()}),n.stop(),n.destroy(),r=null}$(".header-tag").on("click",function(t){t.preventDefault();let e=$(this).data("slug");l=e,r&&r.abort();t=$(this).children(".gh-tag-nav-indicator");t.attr("id")?(t.removeAttr("id"),e=""):t.attr("id","active"),$(".header-tag").not(this).children(".gh-tag-nav-indicator").removeAttr("id");let o=$(".gh-post-feed"),n=$(".gh-post-feed-footer"),a=$(".gh-post-loader #loader"),i=(o.stop().fadeOut(300,"linear"),n.stop().fadeOut(300,"linear"),a.stop().fadeIn(300,"linear"),bodymovin.loadAnimation({container:a[0],renderer:"svg",loop:!0,autoplay:!0,path:"/assets/lottie/loading.json"}));r=$.ajax({url:"/blog/"+e,type:"GET",dataType:"html",success:function(t){l!==e?console.log(`Not updating after success. ${l} :: `+e):(t=$(t),$(".gh-post-feed").html(t.find(".gh-post-feed").html()),d(o,n,a,i),setupInfiniteScroll(t.find(".total-page-count").text(),e+"/"))},error:function(){l!==e?console.log(`Not updating after error. ${l} :: `+e):d(o,n,a,i)}})})});
//# sourceMappingURL=tag.js.map