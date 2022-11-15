export function ajaxGetHandler(router) {
    return fetch(router).then((response) => response.json());
}

export function ajaxPostHandler(router, data) {
    return fetch(router, {
        method: "POST",
        body: JSON.stringify(data),
        headers: {
            Accept: "application/json",
            "Content-type": "application/json",
        },
    }).then((response) => response.json());
}


export function ajaxPutHandler(router, data) {
    return fetch(router, {
        method: "PUT",
        body: JSON.stringify(data),
        headers: {
            Accept: "application/json",
            "Content-type": "application/json",
        },
    }).then((response) => response.json());
}