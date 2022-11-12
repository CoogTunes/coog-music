class uploadManager {
    constructor(form, route) {
        this.form = form;
        this.route = route;
    }

    getElements() {
        let inputs = [...this.form.querySelectorAll('input, select')].filter(elem => elem.getAttribute('data-state') != 'disabled');
        return inputs;
    }

    inputLog(elems) {
        elems.forEach(elem => console.log(elem.name, elem.value));
    }

    gatherData(){
        let data = new FormData();
        let elems = this.getElements();
        elems.forEach(elem => {
            if(elem.type === 'file'){
                if(elem.files.length > 1){
                    for(let i = 0; i < elem.files.length; i++){
                        data.append(elem.getAttribute('name'), elem.files[i]);
                    }
                }
                else{
                    data.append(elem.getAttribute('name'), elem.files[0]);
                }
            }
            else
                data.append(elem.getAttribute('name'), elem.value);
        });
        return data;
    }

    send() {
        let data = this.gatherData();
        fetch(this.route, {
            method: 'POST',
            body: data
        }).then(
            response => response.json()
        ).then(
            success => console.log(success)
        ).catch(
            error => console.log(error)
        );
    }
};

function uploadForm(){
    let submitBtn = document.querySelector('.upload-btn');
    let formUpload = document.querySelector('#upload-form');
    let uploader = new uploadManager(formUpload, '/upload');

    submitBtn.addEventListener('click', function (event) {
        event.preventDefault();
        console.log('Submitting Upload information');
        uploader.inputLog(uploader.getElements());
        console.log([...uploader.gatherData().entries()]);
        uploader.send();
    });
}

window.addEventListener('DOMContentLoaded', function () {
    uploadForm();
});

