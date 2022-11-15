export function dateParse(dateString){
    return new Date(dateString).toDateString().split(' ').slice(1).join(' ');
}