"use strict";
app.UserModel = Backbone.Model.extend({
  urlRoot: 'users/',

  auth: function(callbacks) {
    var that = this;
    $.ajax({
      type: "POST",
      async: false,
      url: that.urlRoot + "auth",
      contentType:"application/json",
      dataType: "json",
      data: JSON.stringify(that.toJSON()),
      success: function(data) { that.parse(data); callbacks.success(data); },
      error: function(data) { callbacks.error(data); }
    });
  },

  logout: function() {
    var that = this;
    $.ajax({
      type: "GET",
      url: this.urlRoot + "logout",
      dataType: "json",
      success: function(data) {that.clear();},
      error: function(data) {that.clear();}
    });
  },

  register: function(callbacks) {
    var that = this;
    $.ajax({
      type: "POST",
      async: false,
      url: that.urlRoot,
      contentType:"application/json",
      dataType: "json",
      data: JSON.stringify(that.toJSON()),
      success: function(data) { that.parse(data); callbacks.success(data); },
      error: function(data) { callbacks.error(data); }
    });
  },

  getCurrentUser: function() {
    var that = this;
    $.ajax({
      type: "GET",
      async: false,
      url: that.urlRoot + "me",
      success: function(data) { that.parse(data); },
      dataType: "json"
    });
  },

  parse: function(data) {
    this.unset("password");
    this.set(data);
    return data;
  }

});
