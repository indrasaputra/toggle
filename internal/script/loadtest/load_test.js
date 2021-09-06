import { group } from 'k6';
import { createToggle, enableToggle, disableToggle, getToggle } from './helper.js';

export let options = {
    stages: [
        { duration: '1m', target: 100 },
        { duration: '5m', target: 100 },
        { duration: '2m', target: 0 },
    ],
    thresholds: {
        'group_duration{group:::create toggle}': ['avg<100', 'p(95)<175','p(99)<500'],
        'group_duration{group:::enable toggle}': ['avg<100', 'p(95)<175','p(99)<500'],
        'group_duration{group:::disable toggle}': ['avg<100', 'p(95)<175','p(99)<500'],
        'group_duration{group:::get toggle}': ['avg<60', 'p(95)<100','p(99)<500'],

        'http_req_failed': ['rate<0.01'],
    },
};

export default function () {
    group('create toggle', function() {
        createToggle()
    })
    group('enable toggle', function() {
        enableToggle()
    })
    group('disable toggle', function() {
        disableToggle()
    })
    group('get toggle', function() {
        getToggle()
    })
}
