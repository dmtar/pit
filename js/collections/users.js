"use strict";
app.UsersCollection = Backbone.Collection.extend({
  urlRoot: 'users/',
  model: User,

  auth: function() {
    $.ajax({
      type: "POST",
      contentType:"application/json; charset=utf-8",
      url: this.urlRoot + "auth",
      data: this.toJSON(),
      success: this.loggedIn,
      dataType: "json"
    });
  },

  getCurrentUser: function(callback) {
    $.ajax({
      type: "GET",
      url: this.urlRoot + "me",
      success: callback,
      dataType: "json"
    });
  },

});
