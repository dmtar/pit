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
    "album/edit/:objectId": "editAlbum",
    "album/list": "listAlbums",
    "album/public": "listPublicAlbums",
    "album/:objectId/pictures": "listAlbumPictures",
    "album/remove/:objectId": "removeAlbum",
    "picture/upload": "uploadPicture",
    "picture/view/:objectId": "viewPicture",
    "picture/edit/:objectId": "editPicture",
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
    if (app.CurrentUser.id) {
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
        },
        true
      );

      this.addSelctize('#pictureTags');
    }
  },

  editPicture: function(objectId) {
    var pictureModel = new app.PictureModel({id: objectId});
    var that = this;
    pictureModel.fetch({
      success: function(data){
        var picture = data.attributes;
        $('#main').html(new app.EditPictureView().render(picture).el);

        that.addSelctize('#pictureTags');
      }
    });
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

  activateMap: function(map, markersKey, selectors, setLocationName) {
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

        if(setLocationName) {
          that.setLocationName(latitude, longitude, selectors.location_name);
        }

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
    if (app.CurrentUser.id) {
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
        },
        true
      );
    }

    this.addSelctize('#albumTags');
  },

  editAlbum: function(objectId) {
    var albumModel = new app.AlbumModel({id: objectId});
    var that = this;
    albumModel.fetch({
      success: function(data){
        var album = data.attributes;
        $('#main').html(new app.EditAlbumView().render(album).el);
        var mapOptions = {
          zoom: 10,
          center: new google.maps.LatLng(album.location.lat, album.location.lng)
        };

        var editAlbumMap = new google.maps.Map(document.getElementById('editAlbumFormMap'), mapOptions);
        that.activateMap(
          editAlbumMap,
          'editAlbumMapMarkers', {
            lat: 'input[name="albumLocationLat"]',
            lng: 'input[name="albumLocationLng"]',
            location_name: 'input[name="albumLocationName"]'
          },
          false
        );

        that.addSelctize('#albumTags');

        var latLng = new google.maps.LatLng(album.location.lat, album.location.lng);
        new google.maps.Marker({
            position: latLng,
            map: editAlbumMap
        });
      }
    });
  },

  listAlbums: function() {
    if (app.CurrentUser.id) {
      var albumModel = new app.AlbumModel();
      albumModel.fetch({
        success: function(data) {
          var albums = data.attributes;
          $('#main').html(new app.ListAlbumsView().render(albums).el);
        }
      });

    }
  },

  listPublicAlbums: function() {
    if (app.CurrentUser.id) {
      $.ajax({
          url: "/albums/public",
          type: "GET",
          dataType: 'json',
          cache: false,
          success: function (data) {
            $('#main').html(new app.ListAlbumsView().render(data).el);
          },
          error: function(){ console.log(attributes); },
      });
    }
  },

  removeAlbum: function(objectId) {
    if (app.CurrentUser.id) {
      var albumModel = new app.AlbumModel({id: objectId});
      albumModel.destroy({
        success: function() {
          Backbone.trigger('flash', { message: 'Success!', type: 'success' });
          Backbone.history.navigate("#album/list", true);
        }
      });
    }
  },

  listAlbumPictures: function(objectId) {
    var albumModel = new app.AlbumModel();
    $.ajax({
      type: "GET",
      async: false,
      url: albumModel.urlRoot + objectId + "/pictures",
      dataType: "json",
      success: function(data) {
        $('#main').html(new app.ListAlbumPictures().render(data).el);
      },
      error: function(data) {
        Backbone.trigger('flash', { message: data.responseJSON.error, type: 'warning' });
      }
    });
  },

  addSelctize: function(selector) {
    $(selector).selectize({
        delimiter: ',',
        persist: false,
        create: function(input) {
            return {
                value: input,
                text: input
            }
        }
    });
  },

  logout: function() {
    app.CurrentUser.logout();
    Backbone.history.navigate("#", true);
    Backbone.trigger('flash', { message: 'Your are now logged out!', type: 'success' });
    this.home();
  }

});
