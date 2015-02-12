'use strict';
app.Router = Backbone.Router.extend({
  routes: {
    "": "home",
    "login": "openLoginModal",
    "register": "openRegisterModal",
    "profile": "openProfileModal",
    "me": "openProfileModal",
    "logout": "logout",
    "album/add": "addAlbum",
    "picture/upload": "uploadPicture",
    "picture/view/:objectId": "viewPicture",
  },
  markers: [],
  map: {},

  initialize: function () {
    $('header').html(new app.HeaderView().render().el);
    $('#main').html(new app.HomeView().render().el);
  },

  home: function() {
    $('#main').html(new app.HomeView().render().el);
  },

  uploadPicture: function() {
    $('#main').html(new app.PictureUploadView().render().el);
  },

  viewPicture: function(objectId) {
    $('#main').html(new app.PictureViewView(objectId).render().el);
  },

  openLoginModal: function() {
    new app.LoginModal().render();
  },

  openProfileModal: function() {
    new app.ProfileModal().render();
  },

  openRegisterModal: function() {
    new app.RegisterModal().render();
  },

  clearMarkers: function() {
    for (var i = 0; i < this.markers.length; i++ ) {
      this.markers[i].setMap(null);
    }
    this.markers.length = 0;
  },

  activate: function() {
      var mapOptions = {
        zoom: 10,
        center: new google.maps.LatLng(42.6833333, 23.3166667)
      };

      this.map = new google.maps.Map(document.getElementById('albumFormMap'), mapOptions);
      var that = this;
      google.maps.event.addListener(this.map, "click", function (event) {
        that.clearMarkers();
        var latitude = event.latLng.lat();
        var longitude = event.latLng.lng();
        var marker = new google.maps.Marker({
            position: event.latLng
        });

        marker.setMap(that.map);
        $('#albumLocationLat').val(latitude);
        $('#albumLocationLng').val(longitude);
        that.markers.push(marker);
    });
  },

  addAlbum: function() {
     $('#main').html(new app.AddAlbumView().render().el);
     this.activate();
  },

  logout: function() {
    app.CurrentUser.logout();
    Backbone.history.navigate("#");
    Backbone.trigger('flash', { message: 'Your are now logged out!', type: 'success' });
  }

});
