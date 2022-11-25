import { templateReplace } from "./razer.js";

class musicQueue {
  constructor(queueContainer) {
    this.songQueue = new Array();
    this.songElements = new Array();
    this.counter = 0;
    this.queueContainer = queueContainer;
  }

  setSongQueue(songList){
    this.songQueue = songList;
  }

  createSongElements(){
    this.songQueue.forEach(entry => {
      console.log(entry);
      const albumSongInfo = {
        "{{songID}}" : entry.songID,
        "{{audioPath}}" : entry.audioPath,
        "{{cover}}" : entry.cover,
        "{{artistName}}" : entry.artistName,
        "{{songTitle}}" : entry.songTitle,
      }

      let songTemplate = `
      <div class="playlist-item-flex content-wrapper" data-audio-path="{{audioPath}}" data-music-state="paused" data-song-id="{{songID}}">
        <div class="playlist-img-contain audio-cover">
          <img src="{{cover}}">
        </div>
        <div class="song-info-item">
          <div class="song-info-title">{{songTitle}}</div>
          <div class="song-info-artist">{{artistName}}</div>
        </div>
        <div class="buttons playlist">
          <button><i class="bi bi-play-fill play-btn"></i></button>
        </div>
      </div>
      `;
      songTemplate = templateReplace(songTemplate, albumSongInfo);
      this.songElements.unshift(songTemplate);
    });
  }

  getSongElements(){
    return this.songElements;
  }

  insertElements(container){
    this.getSongElements().forEach(entry => {
      container.insertAdjacentHTML("afterbegin", entry);
    });
  }

  resetQueue(){
    this.songQueue = null;
    this.songElements = new Array();
  }

  getCounter(){
    return this.counter;
  }

  nextCounter(){
    if(this.counter < this.queueContainer.childElementCount-1)
      this.counter+=1;
    console.log("Next Counter:" + this.counter);
  }

  prevCounter(){
    if(this.counter > 0)
      this.counter-=1;
    console.log("Previous Counter:" + this.counter);
  }

  triggerQueue(queueContainer){
    let songContainer = queueContainer.querySelectorAll('.playlist-item-flex.content-wrapper')[this.getCounter()];
    let playButton = songContainer.querySelector('.buttons.playlist button').firstChild;
    playButton.click();
  }
}

class audioItem {
  constructor(path, audio, time, playBtn, duration) {
    this.audioPath = path;
    this.audio = audio;
    this.currentTime = time;
    this.playButton = playBtn;
    this.duration = duration;
  }

  getAudio() {
    return this.audio;
  }

  getPlayButton() {
    return this.playButton;
  }

  getPath() {
    return this.audioPath;
  }

  getDuration() {
    return this.duration;
  }

  getCurrenTime() {
    return this.currentTime;
  }

  setCurrentTime(currentTime) {
    this.currentTime = currentTime;
  }

  setDuration(duration) {
    this.duration = duration;
  }

  getParent() {
    return this.playButton.closest(".content-wrapper");
  }
}

function musicManager() {
  const updatePlayControl = new CustomEvent("updatePlayControl");
  const updatePauseControl = new CustomEvent("updatePauseControl");
  const timeLineTrigger = new CustomEvent("timeLineTrigger");


  let timelineControlPanel = {
    musicDuration: null,
    musicTimePassed: null,
    playbackProgress: null,
    playbackProgressBg: null,
    playbackProgressBar: null,
    progressBarHandler: null,
    maxDuration: null,
    ratio: null,
    handlerMovement: null,

    init: function () {
      timelineControlPanel.playbackProgress = document.querySelector(
          ".music-playback-progress"
      );
      timelineControlPanel.playbackProgressBg = document.querySelector(
          ".timeline-progress-bg"
      );
      timelineControlPanel.playbackProgressBar = document.querySelector(
          ".timeline-progress-bar"
      );
      timelineControlPanel.progressBarHandler = document.querySelector(
          ".timeline-progress-handler"
      );
      timelineControlPanel.musicDuration = document.querySelector(
          ".music-time-duration"
      );
      timelineControlPanel.musicTimePassed =
          document.querySelector(".music-time-passed");
      timelineControlPanel.maxDuration = 0;
      timelineControlPanel.ratio = 0;
      timelineControlPanel.handlerMovement = 0;
    },

    initLogs: function () {
      console.log(timelineControlPanel.playbackProgress);
      console.log(timelineControlPanel.playbackProgressBg);
      console.log(timelineControlPanel.playbackProgressBar);
      console.log(timelineControlPanel.progressBarHandler);
      console.log(timelineControlPanel.musicDuration);
      console.log(timelineControlPanel.musicTimePassed);
      console.log(timelineControlPanel.maxDuration);
      console.log(timelineControlPanel.ratio);
      console.log(timelineControlPanel.handlerMovement);
    },

    songLogs: function () {
      // * Set the song duration in secs
      console.log(
          "Max Song Duration: " + timelineControlPanel.maxDuration + " secs"
      );
      timelineControlPanel.playbackProgress.setAttribute(
          "data-music-max",
          timelineControlPanel.maxDuration
      );
    },

    setMaxDuration: function (maxDuration) {
      timelineControlPanel.maxDuration = maxDuration;
    },

    setHandlerMovement: function (handlerMovement) {
      timelineControlPanel.handlerMovement = handlerMovement;
    },

    setMusicDuration: function () {
      let time = getTime(timelineControlPanel.maxDuration);
      timelineControlPanel.musicDuration.firstChild.innerHTML =
          formatTime(time);
    },

    resetHandler: function (){
      timelineControlPanel.handlerMovement = 0;
    },

    setRatio: function () {
      timelineControlPanel.ratio =
          parseFloat(
              getComputedStyle(
                  timelineControlPanel.playbackProgressBg
              ).getPropertyValue("width"),
              10
          ) / parseFloat(timelineControlPanel.maxDuration);
      console.log(timelineControlPanel.maxDuration);
      console.log(
          parseFloat(
              getComputedStyle(
                  timelineControlPanel.playbackProgressBg
              ).getPropertyValue("width"),
              10
          )
      );
      console.log(timelineControlPanel.ratio);
    },

    setMusicTimePassed: function (timePassed) {
      timelineControlPanel.musicTimePassed.firstChild.innerHTML = timePassed;
    },

    timeline: function () {
      // * Check the ratio
      currentAudio.getAudio().addEventListener("timeupdate", function (e) {
        let seconds = Math.floor(currentAudio.getAudio().currentTime);
        let time = getTime(seconds);
        timelineControlPanel.setMusicTimePassed(formatTime(time));
        timelineControlPanel.setHandlerMovement(
            seconds * timelineControlPanel.ratio
        );
        timelineControlPanel.progressBarHandler.style.left =
            timelineControlPanel.handlerMovement + "px";
        timelineControlPanel.playbackProgressBar.style.width =
            timelineControlPanel.handlerMovement + "px";
        currentAudio.setCurrentTime(currentAudio.getAudio().currentTime);
        e.stopImmediatePropagation();
      });
    },
  };

  let loadMetaData = function (audioTarget) {
    return new Promise((resolve, reject) => {
      audioTarget.addEventListener("loadedmetadata", function (e) {
        currentAudio.setDuration(audioTarget.duration);
      });
      setTimeout(() => resolve("Loaded Metadata..."), 30);
    });
  };

  let coogtuneContainer = document.querySelector(".music-manager-container");
  let queueContainer = document.querySelector(".queue-container");
  let playControl = document.querySelector(".icon :nth-child(2)");
  let currentAudio = new audioItem(null, null, null, null, null);
  let songQueue = new musicQueue(queueContainer);
  let wave = document.querySelector('.wave');
  let currentVolume = null;
  let playIcon = "bi-play-fill";
  let pauseIcon = "bi-pause-fill";

  timelineControlPanel.init();
  timelineControlPanel.initLogs();

  // * Volume Control
  volumeControlPanel();

  // * Next Song Control
  nextSong();

  // * Previous Song Control
  previousSong();

  // * TimeLine Trigger
  timelineControlPanel.playbackProgress.addEventListener("timeLineTrigger", function (e) {
        console.log("Triggered TimeLine Event");
        loadMetaData(currentAudio.getAudio()).then(function (res) {
          timelineControlPanel.setMaxDuration(currentAudio.getAudio().duration);
          timelineControlPanel.setMusicDuration();
          timelineControlPanel.setRatio();
          timelineControlPanel.songLogs();
          timelineControlPanel.timeline();
        });
      }
  );

  coogtuneContainer.addEventListener("click", function (evt) {
    let elem = evt.target;
    if (elem.classList.contains("play-btn")) {
      queueHelper(elem);
      playButtonControl();
    }
  });

  queueContainer.addEventListener("click", function (evt) {
    let elem = evt.target;
    if (elem.classList.contains("play-btn")) {
      playButtonControl();
    }
  });

  function playButtonControl() {
    let songContainer = queueContainer.querySelectorAll('.playlist-item-flex.content-wrapper')[songQueue.getCounter()];
    let playButton = songContainer.querySelector('.buttons.playlist button').firstChild;
    let musicState = isPlaying(playButton);
    let base_url = window.location.origin;
    let audioPath = songContainer.getAttribute('data-audio-path');

    let targetAudio = new Audio(base_url + audioPath.replace('.',''));
    if (musicState == "paused") {
      if (currentAudio.getPath() == null) {
        currentAudio = new audioItem(
            audioPath,
            targetAudio,
            targetAudio.currentTime,
            playButton,
            0
        );
        playControlPanel();
        masterPlaySongInfo();
      }
      else if(currentAudio.getPath() != audioPath){
        currentAudio.getAudio().pause();
        currentAudio.getPlayButton().classList.remove(pauseIcon);
        currentAudio.getPlayButton().classList.add(playIcon);
        currentAudio.getParent().setAttribute('data-music-state', 'paused');
        currentAudio = new audioItem(audioPath, targetAudio, targetAudio.currentTime, playButton, 0);
        timelineControlPanel.resetHandler();
        masterPlaySongInfo();
        wave.classList.toggle('active2');
      }
      masterPlaySongInfo();
      wave.classList.add('active2');
      currentAudio.getPlayButton().classList.remove(playIcon);
      currentAudio.getPlayButton().classList.add(pauseIcon);
      currentAudio.getParent().setAttribute("data-music-state", "playing");
      playControl.dispatchEvent(updatePlayControl);
      timelineControlPanel.playbackProgress.dispatchEvent(timeLineTrigger);
      currentAudio.getAudio().play();
      currentAudio.getAudio().volume = currentVolume;

    } else if (musicState == "playing") {
      currentAudio.getAudio().pause();
      wave.classList.remove('active2');
      currentAudio.getPlayButton().classList.remove(pauseIcon);
      currentAudio.getPlayButton().classList.add(playIcon);
      currentAudio.getParent().setAttribute("data-music-state", "paused");
      playControl.dispatchEvent(updatePauseControl);
    }

    songEnded(currentAudio.getAudio())

    function songEnded(audio){
      audio.addEventListener('ended', function () {
        songQueue.nextCounter();
        songQueue.triggerQueue(queueContainer);
      });
    }
  }

  // * Helper Functions
  function isAlbum(playButton){
    if(playButton.closest(".content-wrapper").getAttribute("data-audio-path"))
      return false;
    else
      return true;
  }

  function queueHelper(playButton){
    let parent = playButton.closest(".content-wrapper");
    if(isAlbum(playButton)){
      console.log('This is an album.');
      let songList = parent.getAttribute("data-album-songs");
      songQueue.setSongQueue(JSON.parse(songList));
      songQueue.createSongElements();
      songQueue.insertElements(queueContainer);
      songQueue.resetQueue();
    }
    else {
      console.log('This is a song.');
      let audio_path = parent.getAttribute("data-audio-path");
      let cover_path = playButton.closest(".content-wrapper").querySelector('.audio-cover img').src;
      let artist_name = playButton.closest(".content-wrapper").querySelector('.song-info-artist').innerHTML;
      let song_title = playButton.closest(".content-wrapper").querySelector('.song-info-title').innerHTML;
      let song_id = playButton.closest(".content-wrapper").getAttribute("data-song-id");

      const songObject = {
        songID : song_id,
        audioPath : audio_path,
        cover : cover_path,
        artistName: artist_name,
        songTitle: song_title,
      }
      songQueue.setSongQueue([songObject]);
      songQueue.createSongElements();
      songQueue.insertElements(queueContainer);
      songQueue.resetQueue();
    }
  }

  function isPlaying(playBtn) {
    let musicState = playBtn
        .closest(".content-wrapper")
        .getAttribute("data-music-state");
    return musicState == "playing"
        ? "playing"
        : musicState == "paused"
            ? "paused"
            : "error";
  }

  function masterPlaySongInfo(){
    let masterPlaycover = document.querySelector('.master-play-cover');
    let masterSongName = document.querySelector('.master-song-title');
    let masterArtistName = document.querySelector('.master-song-artist');
    masterPlaycover.setAttribute('src', currentAudio.getParent().querySelector('.audio-cover img').src);
    masterArtistName.innerHTML = currentAudio.getParent().querySelector('.song-info-artist').innerHTML;
    masterSongName.innerHTML = currentAudio.getParent().querySelector('.song-info-title').innerHTML;
  }

  function playControlPanel() {
    let controlPlayIcon = "bi-play-circle-fill";
    let controlPauseIcon = "bi-pause-circle-fill";
    let iconChild = playControl;

    function playToPause() {
      iconChild.classList.remove(controlPlayIcon);
      iconChild.classList.add(controlPauseIcon);
    }

    function pauseToPlay() {
      iconChild.classList.remove(controlPauseIcon);
      iconChild.classList.add(controlPlayIcon);
    }

    playControl.addEventListener("click", function () {
      let musicState = isPlaying(currentAudio.getPlayButton());
      console.log(musicState);
      if (musicState == "paused") {
        currentAudio.getAudio().play();
        wave.classList.toggle('active2');
        currentAudio.getPlayButton().classList.remove(playIcon);
        currentAudio.getPlayButton().classList.add(pauseIcon);
        currentAudio.getParent().setAttribute("data-music-state", "playing");
        playToPause();
      } else if (musicState == "playing") {
        currentAudio.getAudio().pause();
        wave.classList.toggle('active2');
        currentAudio.getPlayButton().classList.remove(pauseIcon);
        currentAudio.getPlayButton().classList.add(playIcon);
        currentAudio.setCurrentTime(currentAudio.getAudio().currentTime);
        currentAudio.getParent().setAttribute("data-music-state", "paused");
        pauseToPlay();
      }
    });

    playControl.addEventListener("updatePlayControl", function (e) {
      console.log("Logging Event: Update Play");
      playToPause();
      currentAudio.getPlayButton().classList.remove("bi-play-fill");
      currentAudio.getPlayButton().classList.add("bi-pause-fill");
    });

    playControl.addEventListener("updatePauseControl", function (e) {
      console.log("Logging Event: Update Pause");
      pauseToPlay();
      currentAudio.getPlayButton().classList.remove("bi-pause-fill");
      currentAudio.getPlayButton().classList.add("bi-play-fill");
    });

    playControl
  }

  function volumeControlPanel() {
    let masterVolume = document.querySelector(".master-volume");
    let volumeBackground = document.querySelector(".vol_bar");
    currentVolume = masterVolume.value * 0.01;
    console.log(masterVolume.value);

    masterVolume.setAttribute("value", masterVolume.value);
    volumeBackground.style.width = masterVolume.value + "%";

    masterVolume.addEventListener("input", function (evt) {
      let target = evt.target;
      volumeBackground.style.width = target.value + "%";
      masterVolume.setAttribute("value", target.value);
      if (currentAudio.getAudio() != null)
        currentAudio.getAudio().volume = target.value * 0.01;
        currentVolume = target.value * 0.01;
    });
  }

  function formatTime(time) {
    let finalTime = "";
    if (time.hours > 0) {
      finalTime += "" + time.hours + ":" + (time.minutes < 10 ? "0" : "");
    }
    finalTime += "" + time.minutes + ":" + (time.seconds < 10 ? "0" : "");
    finalTime += "" + time.seconds;

    return finalTime;
  }

  function getTime(time) {
    let hours = Math.floor(time / 3600);
    time = time - hours * 3600;
    let minutes = Math.floor(time / 60);
    let seconds = Math.floor(time - minutes * 60);
    return { hours, minutes, seconds };
  }

  function setSeeker(evt){
    const width = this.clientWidth;
    const clickX = evt.offsetX;
    const audioDuration = currentAudio.getDuration();
    let seekTarget = (clickX / width) * audioDuration;

    currentAudio.setCurrentTime(seekTarget);
  }

  function previousSong(){
    let prevContainer = document.querySelector('.bi.bi-skip-start-fill');

    prevContainer.addEventListener('click', function () {
      console.log('Previous song...')
      songQueue.prevCounter(queueContainer);
      songQueue.triggerQueue(queueContainer);
    });
  }

  function nextSong(){
    let nextContainer = document.querySelector('.bi.bi-skip-end-fill');

    nextContainer.addEventListener('click', function () {
      console.log('Next song...')
      songQueue.nextCounter(queueContainer);
      songQueue.triggerQueue(queueContainer);
    });
  }

  // * Implement click event for seeker on progress bar
  // * Implement songEnded ended event on currentAudio to go to next song
  // * Implement nextSong
  // * Implement previousSong
  // * Implement music queue
}

window.addEventListener("DOMContentLoaded", function (evt) {
  musicManager();
});
