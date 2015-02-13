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

  markers: {},
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
    var mapOptions = {
      zoom: 10,
      center: new google.maps.LatLng(42.6833333, 23.3166667)
    };
    var pictureUploadMap = new google.maps.Map(document.getElementById('pictureUploadMap'), mapOptions);
    this.activateMap(
      pictureUploadMap,
      'pictureUploadMarkers', {
        lat: 'input[name="location[lat]"]',
        lng: 'input[name="location[lng]"]',
        location_name: 'input[name="location[name]"]'
      }
    );
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

  clearMarkers: function(markers) {
    for (var i = 0; i < markers.length; i++ ) {
      markers[i].setMap(null);
    }
    markers.length = 0;
  },

  activateMap: function(map, markersKey, selectors) {
      var that = this;
      google.maps.event.addListener(map, "click", function (event) {

        if(that.markers[markersKey] !== undefined) {
          that.clearMarkers(that.markers[markersKey]);
        }

        var latitude = event.latLng.lat();
        var longitude = event.latLng.lng();
        var marker = new google.maps.Marker({
            position: event.latLng,
            map: map
        });

        $(selectors.lat).val(latitude);
        $(selectors.lng).val(longitude);
        that.setLocationName(latitude, longitude, selectors.location_name);

        if(that.markers[markersKey] === undefined) {
          that.markers[markersKey] = [];
        }

        that.markers[markersKey].push(marker);
    });
  },

  setLocationName: function(lat, lng, selector) {
    var geocoder = new google.maps.Geocoder();
    var latlng = new google.maps.LatLng(lat, lng);
    geocoder.geocode({'latLng': latlng}, function(results, status) {
      if (status == google.maps.GeocoderStatus.OK) {
        $(selector).val(results[2].formatted_address);
      }
    });
  },

  addAlbum: function() {
     $('#main').html(new app.AddAlbumView().render().el);

      var mapOptions = {
        zoom: 10,
        center: new google.maps.LatLng(42.6833333, 23.3166667)
      };
      var addAlbumMap = new google.maps.Map(document.getElementById('albumFormMap'), mapOptions);
      this.activateMap(
        addAlbumMap,
        'albumMapMarkers', {
          lat: 'input[name="albumLocationLat"]',
          lng: 'input[name="albumLocationLng"]',
          location_name: 'input[name="albumLocationName"]'
        }
      );
  },

  logout: function() {
    app.CurrentUser.logout();
    Backbone.history.navigate("#");
    Backbone.trigger('flash', { message: 'Your are now logged out!', type: 'success' });
    this.home();
  }

});
