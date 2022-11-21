import {ajaxPostHandler} from './ajax.js';

function likeControl() {
    document.addEventListener('click', function (evt){
        let target = evt.target;

        if(target.classList.contains('play-btn')){
            let songTarget = target.closest('.content-wrapper');
            console.log(songTarget.getAttribute('data-song-id'));
            updateLike(songTarget.getAttribute('data-song-id'), target.parentElement);

        }
    });

    function updateLike(id, parentElement){
        console.log('Attempting to update total plays for song...');
        let data = {
            songID : id,
        }
        ajaxPostHandler('/updatePlays', data).then((data) => {
            console.log('Total plays being updated for song...');
            console.log(data);
            updateSongCount(parentElement, data.Total_plays);
        }).catch((error) => {
            console.log('Error updating total plays for song...');
            console.log(error);
        });
    }

    function updateSongCount(parentElement, songPlayCount){
        let songPlays = parentElement.querySelector('.song-play-count');
        if(songPlays)
            songPlays.innerHTML = songPlayCount;
    }
}

window.addEventListener('DOMContentLoaded', function (evt){
    likeControl();
});