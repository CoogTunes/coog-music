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

let musicManagerTest = {
  KYqAYsujixMlfw7:
      "/media/artist/lil_baby/albums/its_only_me/songs/Lil_Baby-Real_Spill.mp3",
  zX3mtNS60AxOHLQ:
      "/media/artist/bad_bunny/albums/un_verano_sin_ti/songs/Bad_Bunny-Moscow_Mule.mp3",
  tv2j4ayyyUyUZI3:
      "/media/artist/beyonce/albums/renaissance/songs/Beyonce-Cozy.mp3",
  "8M06h9HcaWHS1Ya":
      "/media/artist/morgan_allen/albums/dangerous/songs/Morgan_Wallen-Dangerous.mp3",
  UyTInzO6jCD1j2X:
      "/media/artist/stray_kids/albums/maxident/songs/Stray_Kids-Super_Board.mp3",
  fBRiFzGIVST51kD:
      "/media/artist/the_weeknd/albums/the_highlights/songs/The_Weeknd_&_Ariana_Grande-Save_Your_Tears.mp3",
  g7vtS4I1sZDY099:
      "/media/artist/harry_style/albums/harrys_house/songs/Harry_Styles-As_It_Was.mp3",
  tZM5ee4bCLQEmaM:
      "/media/artist/quavo_takeoff/albums/only_built_for_infinity_links/songs/Quavo_&_Takeoff-2.30.mp3",
  ChLf0BDpj8egdsm:
      "/media/artist/zach_bryan/albums/american_heartbreak/songs/Zach_Bryan-Something_In_The_Orange.mp3",
  pKixQhRYiWwFRod:
      "/media/artist/g_herbo/albums/survivors_remorse/songs/G_Herbo-4_Minutes_of_Hell.mp3",
  fX8LYvcQR3PziKl:
      "/media/artist/charlie_puth/albums/charlie/songs/Charlie_Puth-Left_And_Right_ft._Jungkook_of_BTS.mp3",
  "1jnOsMvleMc1zUa":
      "/media/artist/rod_wave/albums/beautiful_mind/songs/Rod_Wave-Yungen_ft._Jack_Harlow.mp3",
};

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
      setTimeout(() => resolve("Loaded Metadata..."), 100);
    });
  };

  let coogtuneContainer = document.querySelector(".music-manager-container");
  let playControl = document.querySelector(".icon :nth-child(2)");
  let currentAudio = new audioItem(null, null, null, null, null);
  let wave = document.querySelector('.wave');
  let currentVolume = null;

  timelineControlPanel.init();
  timelineControlPanel.initLogs();

  // * Volume Control
  volumeControlPanel();

  // * TimeLine Trigger
  timelineControlPanel.playbackProgress.addEventListener(
      "timeLineTrigger",
      function (e) {
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
      playButtonControl(elem);
    }
  });

  function playButtonControl(elem) {
    let playButton = elem;
    let musicState = isPlaying(playButton);
    let base_url = window.location.origin;
    let audioPath = playButton
        .closest(".content-wrapper")
        .getAttribute("data-audio-path");
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
        currentAudio.getPlayButton().classList.remove("bi-pause-fill");
        currentAudio.getPlayButton().classList.add("bi-play-fill");
        currentAudio.getParent().setAttribute('data-music-state', 'paused');
        currentAudio = new audioItem(audioPath, targetAudio, targetAudio.currentTime, playButton, 0);
        timelineControlPanel.resetHandler();
        masterPlaySongInfo();
        wave.classList.toggle('active2');
      }
      masterPlaySongInfo();
      wave.classList.toggle('active2');
      currentAudio.getPlayButton().classList.remove("bi-play-fill");
      currentAudio.getPlayButton().classList.add("bi-pause-fill");
      currentAudio.getParent().setAttribute("data-music-state", "playing");
      playControl.dispatchEvent(updatePlayControl);
      timelineControlPanel.playbackProgress.dispatchEvent(timeLineTrigger);
      currentAudio.getAudio().play();
      currentAudio.getAudio().volume = currentVolume;

    } else if (musicState == "playing") {
      currentAudio.getAudio().pause();
      wave.classList.toggle('active2');
      currentAudio.getPlayButton().classList.remove("bi-pause-fill");
      currentAudio.getPlayButton().classList.add("bi-play-fill");
      currentAudio.getParent().setAttribute("data-music-state", "paused");
      playControl.dispatchEvent(updatePauseControl);
    }
  }

  // * Helper Function
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
    let playIcon = "bi-play-fill";
    let pauseIcon = "bi-pause-fill";
    let iconChild = playControl;

    console.log(iconChild);

    function playToPause() {
      iconChild.classList.remove(playIcon);
      iconChild.classList.add(pauseIcon);
    }

    function pauseToPlay() {
      iconChild.classList.remove(pauseIcon);
      iconChild.classList.add(playIcon);
    }

    playControl.addEventListener("click", function () {
      let musicState = isPlaying(currentAudio.getPlayButton());
      console.log(musicState);
      if (musicState == "paused") {
        currentAudio.getAudio().play();
        wave.classList.toggle('active2');
        currentAudio.getPlayButton().classList.remove("bi-play-fill");
        currentAudio.getPlayButton().classList.add("bi-pause-fill");
        currentAudio.getParent().setAttribute("data-music-state", "playing");
        playToPause();
      } else if (musicState == "playing") {
        currentAudio.getAudio().pause();
        wave.classList.toggle('active2');
        currentAudio.getPlayButton().classList.remove("bi-pause-fill");
        currentAudio.getPlayButton().classList.add("bi-play-fill");
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

  // Time Formatter
  function formatTime(time) {
    let finalTime = "";
    if (time.hours > 0) {
      finalTime += "" + time.hours + ":" + (time.minutes < 10 ? "0" : "");
    }
    finalTime += "" + time.minutes + ":" + (time.seconds < 10 ? "0" : "");
    finalTime += "" + time.seconds;

    return finalTime;
  }

  // Time Split
  function getTime(time) {
    let hours = Math.floor(time / 3600);
    time = time - hours * 3600;
    let minutes = Math.floor(time / 60);
    let seconds = Math.floor(time - minutes * 60);
    return { hours, minutes, seconds };
  }
}

window.addEventListener("DOMContentLoaded", function (evt) {
  musicManager();
});
