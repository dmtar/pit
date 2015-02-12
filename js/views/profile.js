app.ProfileModal = Backbone.ModalView.extend({
  title: "Profile",

  buttons: [{
    className: "btn-default cancel",
    label: "Cancel",
    close: true
  },{
    className: "btn-primary ok profile-form-ok",
    label: "Save"
  }],

  events: {
    "click .modal-footer a.ok": "onOk",
    "click .modal-footer a.cancel": "onCancel",
    "submit #profile-form": "profile",
    "hidden.bs.modal": "onHidden"
  },

  onOk: function(e) {
    e.preventDefault();
    $(".profile-form-ok").attr("disabled", "disabled");
    $("#profile-form").find("#profile-form-submit").click();
    $(".profile-form-ok").removeAttr("disabled");
  },

  profile: function(e) {
    var that = this;
    $(".profile-error").hide();
    $(".profile-form-ok").attr("disabled", "disabled");
    e.preventDefault();
    app.CurrentUser.set({display_name: $("#profileDisplayName").val()});
    app.CurrentUser.update({
      success: function(data) {
        that.success(data)
      },
      error: that.displayError
    });
    $(".profile-form-ok").removeAttr("disabled");
  },

  displayError: function(data) {
    $error = $(".profile-error");
    $error.text(data.responseJSON.error);
    $error.show();
  },

  success: function(data) {
    $(".profile-form-ok").attr("disabled", "disabled");

    this.close();
    Backbone.trigger('flash', { message: 'Your profile is updated!', type: 'success' });
  },

  onCancel: function(e) {
  },

  onHidden: function(e) {
    $(".profile-error").hide();
    Backbone.history.history.back();
  },

  postRender: function() {
    $(".profile-error").hide();
  }

});
