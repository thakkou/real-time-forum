// import { socket } from '../core/websocket.js';

export async function render() {
    return `
        <h1>Chat</h1>
        <div id="messages"></div>
    `;
}

// document.addEventListener(
//     'ws:new_message',
//     e => {
//         console.log(e.detail);
//     }
// );