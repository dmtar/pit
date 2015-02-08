app.RegisterModal = Backbone.ModalView.extend({
    title: "Register",

    buttons: [{
      className: "btn-default cancel",
      label: "Cancel",
      close: true
    },{
      className: "btn-primary ok register-form-ok",
      label: "Register"
    }],

    events: {
      "click .modal-footer a.ok": "onOk",
      "click .modal-footer a.cancel": "onCancel",
      "submit #register-form": "register",
      "hidden.bs.modal": "onHidden"
    },

    onOk: function(e) {
      e.preventDefault();
      $(".register-form-ok").attr("disabled", "disabled");
      $("#register-form").find("#register-form-submit").click();
      $(".register-form-ok").removeAttr("disabled");
    },

    register: function(e) {
      var that = this;
      $(".register-error").hide();
      $(".register-form-ok").attr("disabled", "disabled");
      e.preventDefault();
      app.CurrentUser.set({email: $("#registerEmail").val()});
      app.CurrentUser.set({display_name: $("#registerDisplayName").val()});
      app.CurrentUser.set({username: $("#registerUsername").val()});
      app.CurrentUser.set({password: $("#registerPassword").val()});
      app.CurrentUser.register({
        success: function(data) {
          that.success(data)
        },
        error: that.displayError
      });
      $(".register-form-ok").removeAttr("disabled");
    },

    displayError: function(data) {
      $error = $(".register-error");
      $error.text(data.responseJSON.error);
      $error.show();
    },

    success: function(data) {
      $(".register-form-ok").attr("disabled", "disabled");

      this.close();
      $(".message-success").text("You are now logged in!");
      $(".message-success").fadeIn(1000).fadeOut(1000);
    },

    onCancel: function(e) {
    },

    onHidden: function(e) {
      $(".register-error").hide();
    },

    postRender: function() {
      $(".register-error").hide();
    }

});
