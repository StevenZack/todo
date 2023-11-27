function encodeInputsUrlEncoded(inputs) {
    var body = '';
    for (var j = 0; j < inputs.length; j++) {
        if (j > 0) {
            body += '&';
        }
        body += inputs[j].getAttribute('name') + '=' + encodeURIComponent(inputs[j].value);
    }
    return body;
}
var vs = document.querySelectorAll('button[href-delete]');
for (var i = 0; i < vs.length; i++) {
    var v = vs[i];
    v.onclick = function (e) {
        var self = e.currentTarget;
        if (confirm('Are you sure to delete?')) {
            self.disabled = true;
            var before = self.innerHTML;
            self.innerHTML = '<progress style="width: 40px;"/>'
            var xhr = new XMLHttpRequest();
            xhr.onreadystatechange = function () {
                if (xhr.readyState !== 4) {
                    return
                }
                if (xhr.status !== 200) {
                    alert(xhr.responseText)
                    self.disabled = false;
                    self.innerHTML = before;
                    return
                }
                location.href = location.href;
            }

            xhr.open('DELETE', self.getAttribute('href-delete'))
            xhr.send();
        }
    }
}

var vs = document.querySelectorAll('form[method]');
for (var i = 0; i < vs.length; i++) {
    var v = vs[i];
    var method = v.getAttribute('method').toLowerCase();
    if (!method || method === 'get') continue;
    v.onsubmit = function (e) {
        e.preventDefault();
        var self = e.currentTarget;
        var before = self.innerHTML;

        var xhr = new XMLHttpRequest();
        xhr.onreadystatechange = function () {
            if (xhr.readyState !== 4) {
                return
            }
            if (xhr.status !== 200) {
                alert(xhr.responseText)
                self.disabled = false;
                self.innerHTML = before;
                return
            }
            location.href = location.href;
        }

        var method = self.getAttribute('method');
        if (!method) method = 'get';
        method = method.toUpperCase();
        var enctype = self.getAttribute('enctype');
        if (!enctype) enctype = 'application/x-www-form-urlencoded';
        var body = null;
        var url = self.getAttribute('action');

        var inputs = self.querySelectorAll('input[name]');
        switch (enctype) {
            case 'application/x-www-form-urlencoded':
                body = encodeInputsUrlEncoded(inputs);
                break;
            case 'multipart/form-data':
                body = new FormData();
                for (var j = 0; j < inputs.length; j++) {
                    var v = inputs[j]
                    var value = v.value;
                    if (v.files && v.files.length > 0)
                        value = v.files;
                    body.append(v.getAttribute('name'), value);
                }
                break;
        }

        self.disabled = true;
        self.innerHTML = '<progress style="width: 40px;"/>'
        xhr.open(method, url);
        xhr.setRequestHeader('Content-Type', enctype)

        xhr.send(body);
    }
}
