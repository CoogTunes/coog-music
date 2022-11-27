import {ajaxPostHandler} from './ajax.js';

function likeControl() {
    document.addEventListener('click', function (evt){
        let target = evt.target;

        if(target.classList.contains('bi-heart') || target.classList.contains('bi-heart-fill')){
            let songTarget = target.closest('.content-wrapper') ?? target.closest('.table-song-item');
            console.log(songTarget.getAttribute('data-song-id'));
            updateLike(songTarget.getAttribute('data-song-id'), target.parentElement);
        }
    });

    function updateLike(id, parentElement){
        console.log('Attempting to add like...');
        let data = {
            songID : id,
            check: "true",
        }
        ajaxPostHandler('/like', data).then((data) => {
            console.log('Like added to song...');
            updateSongCount(parentElement, data.Likes);
        }).catch((error) => {
            console.log('Error adding like to song...');
            console.log(error);
        });
    }

    function updateSongCount(parentElement, songLikes){
        let currentLike = parentElement.querySelector('.current-likes');
        currentLike.innerHTML = songLikes;
        if (parentElement.querySelector('.bi-heart')) {
            parentElement.querySelector('.bi-heart').classList.add('bi-heart-fill')
            parentElement.querySelector('.bi-heart').classList.remove('bi-heart')
        } else if (parentElement.querySelector('.bi-heart-fill')) {
            parentElement.querySelector('.bi-heart-fill').classList.add('bi-heart')
            parentElement.querySelector('.bi-heart-fill').classList.remove('bi-heart-fill')
        }
    }

}

window.addEventListener('DOMContentLoaded', function (evt){
    likeControl();
});