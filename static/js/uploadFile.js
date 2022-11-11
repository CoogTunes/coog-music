class uploadManager {
    constructor(form, routed) {
        this.form = form;
        this.route = routed;
    }

    inputElements() {
        let inputs = this.form.querySelectorAll("input, select");
        return inputs;
    }

    inputLog(elems) {
        elems.forEach((elem) => console.log(elem.name, elem.type, elem.value));
    }

    gatherData() {
        let data = new FormData();
        let elems = this.inputElements();
        elems.forEach((elem) => {
            if (elem.type === "file")
                data.append(elem.getAttribute("name"), elem.files[0]);
            else data.append(elem.getAttribute("name"), elem.value);
        });
        return data;
    }

    send() {
        let data = this.gatherData();
        fetch(this.route, {
            method: "POST",
            body: data,
            // headers: {
            //     'Content-Type': 'application/json'
            //     // 'Content-Type': 'application/x-www-form-urlencoded',
            // },
        })
            .then((response) => response.json())
            .then((success) => {
                console.log(success);
            })
            .catch((error) => {
                console.log(error);
            });
    }
}

function uploadForm() {
    let submitBtn = document.querySelector(".upload-btn");
    let formUpload = document.querySelector("#upload-form");
    let uploadCenter = new uploadManager(formUpload, "/upload");

    submitBtn.addEventListener("click", function (evt) {
        evt.preventDefault();
        console.log(uploadCenter.route);
        uploadCenter.inputLog(uploadCenter.inputElements());
        console.log([...uploadCenter.gatherData().entries()]);
        uploadCenter.send();
    });
}

window.addEventListener("DOMContentLoaded", function () {
    uploadForm();
});
