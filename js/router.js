'use strict';
app.Router = Backbone.Router.extend({
  routes: {
    "": "home",
    "login": "openLoginModal",
    "register": "openRegisterModal",
    "profile": "openProfileModal",
    "logout": "logout",
  },

  initialize: function () {
    var headerView = new app.HeaderView();
    $('header').html(headerView.render().el);
  },

  openLoginModal: function(e) {
    new app.LoginModal().render();
  },

  openProfileModal: function(e) {
    new app.ProfileModal().render();
  },

  openRegisterModal: function(e) {
    new app.RegisterModal().render();
  },

  logout: function(e) {
    app.CurrentUser.logout();
    $(".message-success").text("You are now logged out!");
    $(".message-success").fadeIn(1000).fadeOut(1000);
  }
});
