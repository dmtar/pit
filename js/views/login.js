app.LoginModal = Backbone.ModalView.extend({
  title: "Login",

  buttons: [{
    className: "btn-default cancel",
    label: "Cancel",
    close: true
  },{ className: "btn-primary ok login-form-ok",
    label: "Login"
  }],

  events: {
    "click .modal-footer a.ok": "onOk",
    "click .modal-footer a.cancel": "onCancel",
    "submit #login-form": "login",
    "hidden.bs.modal": "onHidden"
  },

  onOk: function(e) {
    e.preventDefault();
    $(".login-form-ok").attr("disabled", "disabled");
    $("#login-form").find("#login-form-submit").click();
    $(".login-form-ok").removeAttr("disabled");
  },

  login: function(e) {
    var that = this;
    $(".login-error").hide();
    $(".login-form-ok").attr("disabled", "disabled");
    e.preventDefault();
    app.CurrentUser.set({email: $("#loginEmail").val()});
    app.CurrentUser.set({password: $("#loginPassword").val()});
    app.CurrentUser.auth({
      success: function(data) {
        that.success(data)
      },
      error: that.displayError
    });
    $(".login-form-ok").removeAttr("disabled");
  },

  displayError: function(data) {
    $error = $(".login-error");
    $error.text(data.responseJSON.error);
    $error.show();
  },

  success: function(data) {
    $(".login-form-ok").attr("disabled", "disabled");

    this.close();
    Backbone.trigger('flash', { message: 'Your are now logged in!', type: 'success' });
  },

  onCancel: function(e) {
  },

  onHidden: function(e) {
    $(".login-error").hide();
    Backbone.history.navigate("#");
  },

  postRender: function() {
    $(".login-error").hide();
  }

});
