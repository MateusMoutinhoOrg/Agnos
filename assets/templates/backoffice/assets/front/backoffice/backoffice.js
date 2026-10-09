// The behaviour of the backoffice pages. It lives in this file, served from
// the server itself, because their Content-Security-Policy runs no script
// written inline: neither a <script> block nor an on* attribute.
(function () {
  // A form carrying data-confirm asks before it is sent.
  document.querySelectorAll('form[data-confirm]').forEach(function (form) {
    form.addEventListener('submit', function (event) {
      if (!window.confirm(form.dataset.confirm)) {
        event.preventDefault();
      }
    });
  });

  // The create token form: the date field shows only for a custom expiration.
  var expiration = document.getElementById('expiration');
  var group = document.getElementById('date-group');
  var date = document.getElementById('date');
  if (expiration && group && date) {
    var sync = function () {
      var custom = expiration.value === 'custom';
      group.hidden = !custom;
      date.required = custom;
    };
    expiration.addEventListener('change', sync);
    sync();
  }

  // The create token form: add the browser's own ip to the allowed ones.
  var mine = document.getElementById('use-my-ip');
  if (mine) {
    mine.addEventListener('click', function () {
      var field = document.getElementById('ips');
      var ip = mine.dataset.ip;
      var listed = field.value.split(',').map(function (entry) { return entry.trim(); }).filter(Boolean);
      if (listed.indexOf(ip) === -1) {
        listed.push(ip);
      }
      field.value = listed.join(', ');
    });
  }

  // The snapshot upload form: the server reads no multipart body, so the
  // file is sent as the whole body of the request, and the list is reloaded
  // with the outcome — or the reason it was refused is shown.
  var upload = document.getElementById('snapshot-upload');
  if (upload) {
    upload.addEventListener('submit', function (event) {
      event.preventDefault();
      var file = document.getElementById('snapshot-file').files[0];
      var button = document.getElementById('snapshot-upload-button');
      var error = document.getElementById('upload-error');
      if (!file) {
        return;
      }
      var fail = function (message) {
        error.textContent = message;
        error.hidden = false;
        button.disabled = false;
        button.textContent = 'Upload';
      };
      error.hidden = true;
      button.disabled = true;
      button.textContent = 'Uploading…';
      fetch(upload.action, {
        method: 'POST',
        headers: { 'Content-Type': 'application/zip' },
        body: file
      }).then(function (response) {
        if (response.ok) {
          location.href = '/admin/root/list-backups?notice=uploaded';
          return;
        }
        return response.json().then(function (body) {
          fail(body.error || 'The upload was refused.');
        }, function () {
          fail('The upload was refused (' + response.status + ').');
        });
      }, function () {
        fail('The upload could not reach the server.');
      });
    });
  }

  // The token just created: select it on focus, copy it, and spell the
  // example request against this server.
  var origin = document.getElementById('origin');
  if (origin) {
    origin.textContent = location.origin;
  }
  var token = document.getElementById('new-token');
  if (token) {
    token.addEventListener('focus', function () { token.select(); });
  }
  var copy = document.getElementById('copy-token');
  if (copy && token) {
    copy.addEventListener('click', function () {
      token.select();
      var done = function () { copy.textContent = 'Copied'; };
      if (navigator.clipboard) {
        navigator.clipboard.writeText(token.value).then(done, function () { document.execCommand('copy'); done(); });
      } else {
        document.execCommand('copy');
        done();
      }
    });
  }
})();
