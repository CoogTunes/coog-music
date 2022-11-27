export function dateParse(dateString){
    return new Date(dateString.replace(/-/g, '\/').replace(/T.+/, '')).toDateString().split(' ').slice(1).join(' ');
}