app.HeaderView = Backbone.View.extend({

    events: {
      "click .login-menu": "openLoginModal",
      "click .register-menu": "openRegisterModal",
      "click .profile-menu": "openProfileModal",
      "click .logout-menu": "logout",
    },

    initialize: function() {
      app.CurrentUser.on('change', this.render, this)
    },

    render: function () {
        $(this.el).html(this.template({user: app.CurrentUser}));
        return this;
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
