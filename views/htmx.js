var vs = document.querySelectorAll('button[href-delete]');
for (var i = 0; i < vs.length; i++) {
    var v = vs[i];
    v.onclick = function (e) {
        if (confirm('Are you sure to delete?')) {
            var xhr = new XMLHttpRequest();
            xhr.onreadystatechange = function () {
                if (xhr.readyState !== 4) {
                    return
                }
                if (xhr.status !== 200) {
                    alert(xhr.responseText)
                    return
                }
                location.href = location.href;
            }

            xhr.open('DELETE', e.currentTarget.getAttribute('href-delete'))
            xhr.send();
        }
    }
}

var vs = document.querySelectorAll('form[method]');
for (var i = 0; i < vs.length; i++) {
    var v = vs[i];
    if (v.getAttribute('method') === 'patch') {
        v.onsubmit = function (e) {
            e.preventDefault();
            var xhr = new XMLHttpRequest();
            xhr.onreadystatechange = function () {
                if (xhr.readyState !== 4) {
                    return
                }
                if (xhr.status !== 200) {
                    alert(xhr.responseText)
                    return
                }
                location.href = location.href;
            }
            xhr.open('PATCH', e.currentTarget.getAttribute('action'));
            xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded')
            var inputs = e.currentTarget.querySelectorAll('input[name]');
            var body = '';
            for (var j = 0; j < inputs.length; j++) {
                if (j > 0) {
                    body += '&';
                }
                body += inputs[j].getAttribute('name') + '=' + encodeURIComponent(inputs[j].value);
            }
            xhr.send(body);
        }
    } else if (v.getAttribute('method') === 'put') {
        v.onsubmit = function (e) {
            e.preventDefault();
            var xhr = new XMLHttpRequest();
            xhr.onreadystatechange = function () {
                if (xhr.readyState !== 4) {
                    return
                }
                if (xhr.status !== 200) {
                    alert(xhr.responseText)
                    return
                }
                location.href = location.href;
            }
            xhr.open('PUT', e.currentTarget.getAttribute('action'));
            xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded')
            var inputs = e.currentTarget.querySelectorAll('input[name]');
            var body = '';
            for (var j = 0; j < inputs.length; j++) {
                if (j > 0) {
                    body += '&';
                }
                body += inputs[j].getAttribute('name') + '=' + encodeURIComponent(inputs[j].value);
            }
            xhr.send(body);
        }
    }
}