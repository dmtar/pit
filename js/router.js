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

  },

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

  openLoginModal: function() {
    new app.LoginModal().render();
  },

  openProfileModal: function() {
    new app.ProfileModal().render();
  },

  openRegisterModal: function() {
    new app.RegisterModal().render();
  },

  addAlbum: function() {
    console.log("Add Album");
  },

  logout: function() {
    app.CurrentUser.logout();
    Backbone.history.navigate("#");
    $(".message-success").text("You are now logged out!");
    $(".message-success").fadeIn(1000).fadeOut(1000);
  }

});
