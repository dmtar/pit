var app = app || {};

$(function () {
  'use strict';

  app.loadTemplates = function(views, modalViews, callback){
    var deferreds = [];

    $.each(views, function(index, view) {
        if (app[view]) {
            deferreds.push($.get('tpl/' + view + '.html', function(data) {
                app[view].prototype.template = _.template(data);
            }, 'html'));
        } else {
            alert(view + " not found");
        }
    });

    $.each(modalViews, function(index, view) {
        if (app[view]) {
            deferreds.push($.get('tpl/modals/' + view + '.html', function(data) {
                app[view].prototype.body_template = _.template(data);
            }, 'html'));
        } else {
            alert(view + " not found");
        }
    });

    $.when.apply(null, deferreds).done(callback);
  };

});

function popNotification() {
  setTimeout(function () {
    $.ajax({
        url: "/notifications/pop",
        type: "GET",
        dataType: 'json',
        success: function (data) {
          Backbone.trigger('flash', { message: data.text, type: 'info', persist: true });
        },
        error: function(){},
        complete: popNotification
    });
  }, 5000);
}


$(document).on("ready", function () {
    app.CurrentUser = new app.UserModel();
    app.CurrentUser.getCurrentUser();

    app.loadTemplates(
      ["HeaderView",
       "HomeView",
        "PictureUploadView",
        "PictureViewView",
        "EditPictureView",
        "AddAlbumView",
        "EditAlbumView",
        "ListAlbumsView",
        "ListAlbumPictures"],
      ["LoginModal", "RegisterModal", "ProfileModal"],
      function () {
          app.router = new app.Router();
          Backbone.history.start();
          Backbone.Flash.initialize({el: "#flashes"});
          $("#loading").fadeOut(1000);
          popNotification();
      }
    );

});
