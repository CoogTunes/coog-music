export function isEmpty(obj) {
    return Object.keys(obj).length === 0;
}

export function songCount(songCount){
    songCount = Object.keys(songCount).length;
    let str = "";
    if(songCount > 1)
        str = songCount + " Songs";
    else
        str = songCount + " Song";

    return str;
}

// Time Formatter
export function formatTime(time) {
    let finalTime = "";
    if (time.hours > 0) {
        finalTime += "" + time.hours + ":" + (time.minutes < 10 ? "0" : "");
    }
    finalTime += "" + time.minutes + ":" + (time.seconds < 10 ? "0" : "");
    finalTime += "" + time.seconds;

    return finalTime;
}

// Time Split
export function getTime(time) {
    let hours = Math.floor(time / 3600);
    time = time - hours * 3600;
    let minutes = Math.floor(time / 60);
    let seconds = Math.floor(time - minutes * 60);
    return { hours, minutes, seconds };
}