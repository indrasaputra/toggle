import http from 'k6/http';
import { check } from 'k6';

const CHARACTERS = 'abcdefghijklmnopqrstuvwxyz0123456789'
const BASE_URL = 'http://host.docker.internal:8081/' // or change to 'http://localhost:8081/'
const KEY_PREFIX = 'loadtest-k6-'
const KEY_LENGTH = 10

const keys = []

export function createRandomKey(length) {
    var result           = KEY_PREFIX;
    var charactersLength = CHARACTERS.length;
    for ( var i = 0; i < length; i++ ) {
        result += CHARACTERS.charAt(Math.floor(Math.random() * charactersLength));
    }
    if (keys.length < 10) {
        keys.push(result)
    }
    return result;
}

export function createToggle() {
    var url = `${BASE_URL}v1/toggles`;
    var payload = JSON.stringify({
        key: createRandomKey(KEY_LENGTH),
    });
    var params = {
        headers: {
            'Content-Type': 'application/json',
        }
    };

    let resp = http.post(url, payload, params);
    check(resp, {
        'success create toggle': (resp) => resp.status === 200,
    })
}

export function enableToggle() {
    var key = getKey();
    var url = `${BASE_URL}v1/toggles/${key}/enable`;

    let resp = http.put(url);
    check(resp, {
        'success enable toggle': (resp) => resp.status === 200,
    })
}

export function disableToggle() {
    var key = getKey();
    var url = `${BASE_URL}v1/toggles/${key}/disable`;

    let resp = http.put(url);
    check(resp, {
        'success disable toggle': (resp) => resp.status === 200,
    })
}

export function getToggle() {
    var key = getKey();
    var url = `${BASE_URL}v1/toggles/${key}`;

    let resp = http.get(url);
    check(resp, {
        'success get toggle': (resp) => resp.status === 200,
    })
}

function getKey() {
    var n = keys.length;
    return keys[Math.floor(Math.random() * n)];
}