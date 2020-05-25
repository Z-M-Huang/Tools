function logout() {
  postLink("/api/logout", (d) => {
    if (d != null && d != undefined && d) {
      location.reload();
    }
  })
}

function clearCookie(name) {
  document.cookie = name + "=; Max-age=0; path=/; domain=" + location.host;
}

function getCookieValue(name) {
  var value = "; " + document.cookie;
  var parts = value.split("; " + name + "=");
  if (parts.length == 2) return parts.pop().split(";").shift();
}

function getLocation(host, callback) {
  $.ajax({
    type: "GET",
    url: "http://ip-api.com/json/" + host
  })
  .done((d, status, jqXHR) => {
    if (d != null && d != undefined) {
      if (d.status == "success") {
        callback(d);
      } else {
        showAlertWarning("failed to get geo location for host: " + host, true, 3500);
      }
    }
  })
  .fail((xhr) => {
    if (xhr.status == 429 || xhr.status == 403) {
      showAlertDanger("too many requests. Please take a break.");
    } else {
      showAlertDanger("failed to get current geo location for host: " + host);
    }
  })
}

/*******************************************************
 *                    Like/Dislike Section
 *******************************************************/
function likeOnClick(name) {
  var ele = $("." + name);
  if (ele.hasClass("fas")) {
    //Unlike
    postLink("/api/app/" + name + "/dislike", (d) => {
      ele.parent().html("<i class=\"far fa-thumbs-up mr-1 hover-pointer hover-150 " + name + "\" onclick=\"likeOnClick('" + name +  "')\"></i>" + d)
    })
  } else {
    //like
    postLink("/api/app/" + name + "/like", (d) => {
      ele.parent().html("<i class=\"fas fa-thumbs-up mr-1 hover-pointer hover-150 " + name + "\" onclick=\"likeOnClick('" + name +  "')\"></i>" + d)
    })
  }
}


/*******************************************************
 *                    Ajax Section
 *******************************************************/
function bindForm(id, url, showSpinner, callback) {
  id = "#" + id;
  var form = $(id);
  var oldTexts = [];
  var submitButtons = $("button[type=submit]", form);
  if (showSpinner) {
    for (var i = 0; i < submitButtons.length; i++) {
      oldTexts.push(submitButtons[i].innerHTML);
    }
  }
  form.on("submit", (e) => {
    if (showSpinner) {
      submitButtons.attr('disabled','disabled');
      submitButtons.html('<span class="spinner-border spinner-border-sm mr-2" role="status" aria-hidden="true"></span>Loading...');
    }
    e.preventDefault();
    $.ajax({
      type: "POST",
      url: url,
      data: form.serialize()
    })
    .done((data) => {
      if (data != null && data != undefined) {
        if (data.Message != "" && data.Message != null && data.Message != undefined) {
          showAlertInfo(data.Message, true, 3500);
        }
        if (callback != null && callback != undefined) {
          callback(data.Data);
        }
      }
    })
    .fail((xhr) => {
      var data = xhr.responseJSON;
      if (data != "" && data != undefined && data != null) {
        if (xhr.status == 401) {
          showAlertDanger("Please login first");
        } else {
          showAlertCondition(xhr.status, data.Message);
        }
      } else {
        showAlertDanger("Failed to receive success response, please try again later.");
      }
    })
    .always(() => {
      if (showSpinner) {
        submitButtons.removeAttr('disabled');
        for (var i = 0; i < submitButtons.length; i++) {
          submitButtons[i].innerHTML = oldTexts[i];
        }
      }
    });
  });
}

function bindFormWithFile(id, url, showSpinner, callback) {
  id = "#" + id;
  var form = $(id);
  var oldTexts = [];
  var submitButtons = $("button[type=submit]", form);
  if (showSpinner) {
    for (var i = 0; i < submitButtons.length; i++) {
      oldTexts.push(submitButtons[i].innerHTML);
    }
  }
  form.on("submit", (e) => {
    if (showSpinner) {
      submitButtons.attr('disabled','disabled');
      submitButtons.html('<span class="spinner-border spinner-border-sm mr-2" role="status" aria-hidden="true"></span>Loading...');
    }
    e.preventDefault();
    var formData = new FormData(form[0]);
    $.ajax({
      type: "POST",
      enctype: 'multipart/form-data',
      url: url,
      data: formData,
      processData: false,
      contentType: false,
      cache: false
    })    
    .done((data) => {
      if (data != null && data != undefined) {
        if (data.Message != "" && data.Message != null && data.Message != undefined) {
          showAlertInfo(data.Message, true, 3500);
        }
        if (callback != null && callback != undefined) {
          callback(data.Data);
        }
      }
    })
    .fail((xhr) => {
      var data = xhr.responseJSON;
      if (data != "" && data != undefined && data != null) {
        if (xhr.status == 401) {
          showAlertDanger("Please login first");
        } else {
          showAlertCondition(xhr.status, data.Message);
        }
      } else {
        showAlertDanger("Failed to receive success response, please try again later.");
      }
    })
    .always(() => {
      if (showSpinner) {
        submitButtons.removeAttr('disabled');
        for (var i = 0; i < submitButtons.length; i++) {
          submitButtons[i].innerHTML = oldTexts[i];
        }
      }
    });
  });
}

function postJSONData(url, data, callback) {
  $.ajax({
    type: "POST",
    url: url,
    data: JSON.stringify(data),
    dataType: "json",
    contentType: "application/json"
  })
  .done((data) => {
    if (data != null && data != undefined) {
      if (data.Message != "" && data.Message != null && data.Message != undefined) {
        showAlertInfo(data.Message, true, 3500);
      }
      if (callback != null && callback != undefined) {
        callback(data.Data);
      }
    }
  })
  .fail((xhr) => {
    var data = xhr.responseJSON;
    if (data != "" && data != undefined && data != null) {
      if (xhr.status == 401) {
        showAlertDanger("Please login first");
      } else {
        showAlertCondition(xhr.status, data.Message);
      }
    } else {
      showAlertDanger("Failed to receive success response, please try again later.");
    }
  });
}

function postLink(url, callback) {
  $.ajax({
    type: "POST",
    url: url,
    dataType: "json",
    beforeSend: (xhr) => {
      var sessionToken = getCookieValue("session_token");
      if (
        sessionToken != "" &&
        sessionToken != null &&
        sessionToken != undefined
      ) {
        xhr.setRequestHeader("Authorization", "Bearer " + sessionToken);
      }
    }
  })
  .done((data) => {
    if (data != null && data != undefined) {
      if (data.Message != "" && data.Message != null && data.Message != undefined) {
        showAlertInfo(data.Message, true, 3500);
      }
      if (callback != null && callback != undefined) {
        callback(data.Data);
      }
    }
  })
  .fail((xhr) => {
    var data = xhr.responseJSON;
    if (data != "" && data != undefined && data != null) {
      if (xhr.status == 401) {
        showAlertDanger("Please login first");
      } else {
        showAlertCondition(xhr.status, data.Message);
      }
    } else {
      showAlertDanger("Failed to receive success response, please try again later.");
    }
  });
}

/*******************************************************
 *                    Dynamic Modal Section
 *******************************************************/
function getModalHTML(id, title, content, primaryButtonOnClick, primaryButtonText) {
  modal = '<div class="modal fade" id="' + id + '"><div class="modal-dialog" role="document"><div class="modal-content"><div class="modal-header"><h5 class="modal-title">' + title + '</h5><button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button></div><div class="modal-body">' + content + '</div><div class="modal-footer">';
  if (!(primaryButtonOnClick == null || primaryButtonOnClick == undefined)) {
    modal += '<button type="button" class="btn btn-primary" onclick="' + primaryButtonOnClick + '">' + primaryButtonText + '</button>';
  }
  modal += '<button type="button" class="btn btn-secondary" data-dismiss="modal">Close</button></div></div></div></div>';
  return modal;
}

/*******************************************************
 *                    Alert Section
 *******************************************************/
function getToastHTML(now, color, title, message, autohide, delay) {
  var html = '<div class="toast hover-pointer" role="alert" aria-live="assertive" aria-atomic="true" data-animation="true"';
  if (autohide) {
    html += ' data-autohide="true" data-delay="' + delay + '"';
  } else {
    html += ' data-autohide="false"';
  }
  html += ' onclick="toastOnClickDispose(this)" id="' + now + '"><div class="toast-header"><svg width="25", height="25" class="rounded-25 mr-2"><rect width="100%" height="100%" style="fill:' + color + '" /></svg><strong class="mr-auto">' + title + '</strong><small class="text-muted">' + new Date(now).toLocaleString('en-US') + '</small></div><div class="toast-body">' + message + '</div></div>';
  return html;
}

function toastOnClickDispose(ele) {
  $(ele).toast('hide')
}

function showAlertDanger(message, autohide, delay) {
  var d = new Date();
  var now = d.getTime();
  var toastHTML = getToastHTML(now, window.alertColors.danger, "Error!", message, autohide, delay);
  $("#toasts").append(toastHTML);
  var toast = $("#" + now);
  toast.toast('show');
  toast.on('hidden.bs.toast', function() {
    $(this).remove();
  })
}

function showAlertWarning(message, autohide, delay) {
  var d = new Date();
  var now = d.getTime();
  var toastHTML = getToastHTML(now, window.alertColors.warning, "Warning!", message, autohide, delay);
  $("#toasts").append(toastHTML);
  var toast = $("#" + now);
  toast.toast('show');
  toast.on('hidden.bs.toast', function() {
    $(this).remove();
  })
}

function showAlertSuccess(message, autohide, delay) {
  var d = new Date();
  var now = d.getTime();
  var toastHTML = getToastHTML(now, window.alertColors.success, "Well done!", message, autohide, delay);
  $("#toasts").append(toastHTML);
  var toast = $("#" + now);
  toast.toast('show');
  toast.on('hidden.bs.toast', function() {
    $(this).remove();
  })
}

function showAlertInfo(message, autohide, delay) {
  var d = new Date();
  var now = d.getTime();
  var toastHTML = getToastHTML(now, window.alertColors.info, "Heads up!", message, autohide, delay);
  $("#toasts").append(toastHTML);
  var toast = $("#" + now);
  toast.toast('show');
  toast.on('hidden.bs.toast', function() {
    $(this).remove();
  })
}

function showAlertCondition(statusCode, message) {
  if (message != "" && message != null && message != undefined) {
    if (statusCode == 200) {
      showAlertInfo(message, true, 3500)
    } else if (statusCode == 400) {
      showAlertWarning(message, true, 3500)
    } else {
      showAlertDanger(message, false, 0);
    }
  }
}

/*******************************************************
 *                    onClick/onChange functions
 *******************************************************/
function styleChangeOnClick(styleName) {
  var d = new Date();
  //100 years should be more than enough right?
  d.setTime(d.getTime() + 3153600000000);
  document.cookie = "page_style=" + styleName + "; expires=" + d.toUTCString() + ";path=/";
  location.reload();
}

function copyValueOnClick(ele) {
  var e = $(ele);
  if (e.val().length > 0) {
    navigator.permissions.query({ name: "clipboard-write" }).then((result) => {
      if (result.state == "granted" || result.state == "prompt") {
        navigator.clipboard.writeText(e.val()).then(
          function () {
            showAlertSuccess("Text Copied!.", true, 3000);
          },
          function () {
            showAlertWarning("Failed to copy text.", true, 3000);
          }
        );
      }
    });
  }
}

function onClickRedirect(url) {
  window.location.href = url
}

function onChangeFileInput(ele) {
  var fileName = $(ele).val();
  $(ele).next('.custom-file-label').html(fileName);
}

/*******************************************************
 *                    Colors
 *******************************************************/
window.chartColors = {
  black: "rgb(0,0,0)",
  navy: "rgb(0,0,128)",
  blue: "rgb(0,0,255)",
  green: "rgb(0,128,0)",
  teal: "rgb(0,128,128)",
  lime: "rgb(0,255,0)",
  aqua: "rgb(0,255,255)",
  maroon: "rgb(128,0,0)",
  purple: "rgb(128,0,128)",
  olive: "rgb(128,128,0)",
  gray: "rgb(128,128,128)",
  silver: "rgb(192,192,192)",
  red: "rgb(255,0,0)",
  fuchsia: "rgb(255,0,255)",
  yellow: "rgb(255,255,0)",
  white: "rgb(255,255,255)"
};

window.alertColors = {
  danger: "#dc3545",
  warning: "#ffc107",
  success: "#28a745",
  info: "#17a2b8"
}