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
