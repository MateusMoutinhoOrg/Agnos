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
