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

$(document).on("ready", function () {
    app.CurrentUser = new app.UserModel();
    app.CurrentUser.getCurrentUser();

    app.loadTemplates(
      ["HeaderView", "HomeView", "PictureUploadView", "PictureViewView", "AddAlbumView"],
      ["LoginModal", "RegisterModal", "ProfileModal"],
      function () {
          app.router = new app.Router();
          Backbone.history.start();
          Backbone.Flash.initialize({el: "#flashes"});
          $("#loading").fadeOut(1000);
      }
    );

    //Init google map
    //var map;

    //function initialize() {
      //var mapOptions = {
        //zoom: 10,
        //center: new google.maps.LatLng(42.6833333, 23.3166667)
      //};

      //map = new google.maps.Map(document.getElementById('map-canvas'), mapOptions);
    //}

    //google.maps.event.addDomListener(window, 'load', initialize);
});
