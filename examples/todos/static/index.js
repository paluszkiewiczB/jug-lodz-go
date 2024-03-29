function taskToggled(id) {
    var xhr = new XMLHttpRequest();
    xhr.open("PUT", "/api/todos/" + id + "/toggle", true);
    xhr.send();
}