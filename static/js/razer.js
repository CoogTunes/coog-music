export function templateReplace(str, replaceMap){
    let re = new RegExp(Object.keys(replaceMap).join("|"), "gi");
    return str.replace(re, (matched) => {
        return replaceMap[matched];
    });
}